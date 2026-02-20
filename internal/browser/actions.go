package browser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

// waitStable waits for the page to stop changing, but caps the wait at maxWait.
// rod's WaitStable(interval) can block forever on pages with continuous animations or
// infinite scroll (e.g. Zhihu search results). This wrapper prevents that.
func waitStable(page *rod.Page, interval, maxWait time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), maxWait)
	defer cancel()
	_ = page.Context(ctx).WaitStable(interval)
}

// Click clicks the element identified by the given ref number.
// It scrolls the element into view, waits for it to be interactable, then clicks.
func Click(page *rod.Page, b *Browser, ref int) error {
	el, err := resolveRef(page, b, ref)
	if err != nil {
		return captureErrorScreenshot(page, b, "click_resolve", ref, err)
	}

	// Wrap all element operations in a bounded context so that clicking elements that
	// trigger heavy AJAX (e.g. Zhihu comment panels) never hangs indefinitely.
	opCtx, opCancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer opCancel()
	bel := el.Context(opCtx)

	// Scroll into view so the element is visible
	if err := bel.ScrollIntoView(); err != nil {
		_ = err // not fatal — element might already be in view
	}

	// Wait for element to be interactable (visible + not covered)
	if _, err := bel.Interactable(); err != nil {
		time.Sleep(300 * time.Millisecond)
		if _, err := bel.Interactable(); err != nil {
			return captureErrorScreenshot(page, b, "click_not_interactable", ref, fmt.Errorf("element [%d] not interactable: %w", ref, err))
		}
	}

	// Use JavaScript click() instead of rod's mouse simulation.
	// Rod's Click() can hang when the target triggers a page navigation or heavy AJAX
	// (e.g. Zhihu comment panels) because rod waits internally for load events.
	// JS element.click() fires the event and returns immediately — no waiting.
	if _, err := bel.Eval(`() => { this.click(); return true; }`); err != nil {
		// Fall back to rod mouse click if JS eval fails
		if err2 := bel.Click(proto.InputMouseButtonLeft, 1); err2 != nil {
			return captureErrorScreenshot(page, b, "click_failed", ref, fmt.Errorf("click failed: %w", err2))
		}
	}

	// Brief settle wait (capped so it never hangs)
	waitStable(page, 300*time.Millisecond, 2*time.Second)

	return nil
}

// Type inputs text into the element identified by the given ref number.
// It clicks the element first to ensure focus, then types.
func Type(page *rod.Page, b *Browser, ref int, text string, submit bool) error {
	el, err := resolveRef(page, b, ref)
	if err != nil {
		return captureErrorScreenshot(page, b, "type_resolve", ref, err)
	}

	opCtx, opCancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer opCancel()
	bel := el.Context(opCtx)

	// Scroll into view
	if err := bel.ScrollIntoView(); err != nil {
		_ = err
	}

	// Click to focus the element first — critical for search boxes, inputs, etc.
	if err := bel.Click(proto.InputMouseButtonLeft, 1); err != nil {
		// Try Focus as fallback
		if err := bel.Focus(); err != nil {
			return captureErrorScreenshot(page, b, "type_focus", ref, fmt.Errorf("failed to focus element [%d]: %w", ref, err))
		}
	}

	// Small delay to let focus animations/handlers run
	time.Sleep(200 * time.Millisecond)

	// Clear existing content
	if err := bel.SelectAllText(); err != nil {
		_ = err
	}

	// Input text
	if err := bel.Input(text); err != nil {
		return captureErrorScreenshot(page, b, "type_input", ref, fmt.Errorf("failed to type text: %w", err))
	}

	if submit {
		time.Sleep(100 * time.Millisecond)
		if err := bel.Type(input.Enter); err != nil {
			return captureErrorScreenshot(page, b, "type_submit", ref, fmt.Errorf("failed to press Enter: %w", err))
		}
		waitStable(page, 500*time.Millisecond, 3*time.Second)
	}

	return nil
}

// Press sends a keyboard key press to the page.
func Press(page *rod.Page, key string) error {
	k, ok := keyMap[key]
	if !ok {
		return fmt.Errorf("unknown key: %q (supported: Enter, Tab, Escape, Backspace, ArrowUp, ArrowDown, ArrowLeft, ArrowRight, Space, Delete, Home, End, PageUp, PageDown)", key)
	}
	if err := page.Keyboard.Type(k); err != nil {
		return err
	}
	// Wait for page to settle after key press
	waitStable(page, 300*time.Millisecond, 3*time.Second)
	return nil
}

// Hover moves the mouse over the element identified by the given ref number.
func Hover(page *rod.Page, b *Browser, ref int) error {
	el, err := resolveRef(page, b, ref)
	if err != nil {
		return err
	}
	if err := el.ScrollIntoView(); err != nil {
		_ = err
	}
	return el.Hover()
}

// ClickAll clicks every element matching the CSS selector with a delay between each.
// It scrolls down repeatedly to find new elements, stopping only when no more appear.
// Elements matching skipSelector are skipped (e.g. already-liked items).
// Returns the number of elements successfully clicked.
func ClickAll(page *rod.Page, selector string, delay time.Duration, skipSelector string) (int, error) {
	clicked := 0
	seen := map[string]bool{} // track by outerHTML to avoid re-processing

	for {
		elements, err := page.Elements(selector)
		if err != nil {
			return clicked, fmt.Errorf("failed to find elements matching %q: %w", selector, err)
		}

		newClicks := 0
		for _, el := range elements {
			// Deduplicate by object remote ID
			objID := fmt.Sprintf("%p", el)
			html, _ := el.HTML()
			key := html
			if key == "" {
				key = objID
			}
			if seen[key] {
				continue
			}
			seen[key] = true

			// Skip elements matching the skip selector (e.g. already active/liked)
			if skipSelector != "" {
				matched, _ := el.Eval(`(sel) => this.matches(sel) || this.querySelector(sel) !== null`, skipSelector)
				if matched != nil && matched.Value.Bool() {
					continue
				}
			}

			_ = el.ScrollIntoView()
			time.Sleep(200 * time.Millisecond)

			if err := el.Click(proto.InputMouseButtonLeft, 1); err != nil {
				continue
			}
			clicked++
			newClicks++

			if delay > 0 {
				time.Sleep(delay)
			}
		}

		// Scroll down to load more content
		_ = page.Mouse.Scroll(0, 800, 0)
		time.Sleep(1500 * time.Millisecond)
		waitStable(page, 500*time.Millisecond, 3*time.Second)

		// Check if new elements appeared after scrolling
		newElements, _ := page.Elements(selector)
		if len(newElements) <= len(elements) && newClicks == 0 {
			// No new elements and nothing new was clicked — we're done
			break
		}
	}

	return clicked, nil
}

// resolveRef looks up a ref number in the browser's ref map and returns the corresponding element.
func resolveRef(page *rod.Page, b *Browser, ref int) (*rod.Element, error) {
	entry, ok := b.GetRef(ref)
	if !ok {
		return nil, fmt.Errorf("ref %d not found in snapshot (run browser_snapshot first, or page may have changed)", ref)
	}

	if entry.BackendDOMNodeID == 0 {
		return nil, fmt.Errorf("ref %d has no backend DOM node (element may be virtual)", ref)
	}

	return resolveBackendNode(page, entry.BackendDOMNodeID)
}

// resolveBackendNode converts a BackendDOMNodeID to a rod Element.
func resolveBackendNode(page *rod.Page, backendNodeID proto.DOMBackendNodeID) (*rod.Element, error) {
	// Use a bounded context so CDP calls don't hang when the page is navigating.
	rCtx, rCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer rCancel()
	result, err := proto.DOMResolveNode{
		BackendNodeID: backendNodeID,
	}.Call(page.Context(rCtx))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve node (element may have been removed from page): %w", err)
	}

	if result.Object.ObjectID == "" {
		return nil, fmt.Errorf("resolved node has no object ID")
	}

	el, err := page.ElementFromObject(result.Object)
	if err != nil {
		return nil, fmt.Errorf("failed to create element from object: %w", err)
	}
	return el, nil
}

// keyMap maps human-readable key names to rod input keys.
var keyMap = map[string]input.Key{
	"Enter":      input.Enter,
	"Tab":        input.Tab,
	"Escape":     input.Escape,
	"Backspace":  input.Backspace,
	"Delete":     input.Delete,
	"ArrowUp":    input.ArrowUp,
	"ArrowDown":  input.ArrowDown,
	"ArrowLeft":  input.ArrowLeft,
	"ArrowRight": input.ArrowRight,
	"Space":      input.Space,
	"Home":       input.Home,
	"End":        input.End,
	"PageUp":     input.PageUp,
	"PageDown":   input.PageDown,
}

// captureErrorScreenshot captures a screenshot when an action fails (if debug mode is enabled).
// It returns the original error with additional context about the screenshot location.
func captureErrorScreenshot(page *rod.Page, b *Browser, action string, ref int, originalErr error) error {
	if !b.IsDebugMode() {
		return originalErr
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05.000")
	filename := fmt.Sprintf("error_%s_ref%d_%s.png", action, ref, timestamp)
	screenshotPath := filepath.Join(b.DebugDir(), filename)

	screenshot, err := page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
	})
	if err != nil {
		// Failed to capture screenshot, just return original error
		return originalErr
	}

	if err := os.WriteFile(screenshotPath, screenshot, 0644); err != nil {
		// Failed to save screenshot, just return original error
		return originalErr
	}

	return fmt.Errorf("%w (debug screenshot saved to: %s)", originalErr, screenshotPath)
}
