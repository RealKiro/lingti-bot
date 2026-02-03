package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// RemindersToday gets today's reminders (macOS)
func RemindersToday(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	script := `
		set output to ""
		tell application "Reminders"
			set allLists to every list
			repeat with reminderList in allLists
				set listName to name of reminderList
				set incompleteReminders to (every reminder of reminderList whose completed is false)
				repeat with r in incompleteReminders
					set rName to name of r
					set rDue to ""
					try
						set rDue to due date of r as string
					end try
					if rDue is not "" then
						set output to output & "☐ " & rName & " [" & listName & "] - Due: " & rDue & linefeed
					else
						set output to output & "☐ " & rName & " [" & listName & "]" & linefeed
					end if
				end repeat
			end repeat
		end tell
		return output
	`

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get reminders: %v", err)), nil
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		return mcp.NewToolResultText("No reminders found"), nil
	}

	return mcp.NewToolResultText("Reminders:\n" + string(output)), nil
}

// RemindersAdd creates a new reminder (macOS)
func RemindersAdd(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, ok := req.Params.Arguments["title"].(string)
	if !ok || title == "" {
		return mcp.NewToolResultError("title is required"), nil
	}

	list := "Reminders"
	if l, ok := req.Params.Arguments["list"].(string); ok && l != "" {
		list = l
	}

	notes := ""
	if n, ok := req.Params.Arguments["notes"].(string); ok {
		notes = n
	}

	dueDate := ""
	if d, ok := req.Params.Arguments["due"].(string); ok && d != "" {
		// Parse the due date
		t, err := time.Parse("2006-01-02 15:04", d)
		if err != nil {
			// Try date only
			t, err = time.Parse("2006-01-02", d)
			if err != nil {
				return mcp.NewToolResultError("invalid due date format, use YYYY-MM-DD or YYYY-MM-DD HH:MM"), nil
			}
		}
		dueDate = t.Format("January 2, 2006 at 3:04:05 PM")
	}

	var script string
	if dueDate != "" {
		script = fmt.Sprintf(`
			tell application "Reminders"
				tell list "%s"
					make new reminder with properties {name:"%s", body:"%s", due date:date "%s"}
				end tell
			end tell
			return "OK"
		`, escapeAppleScript(list), escapeAppleScript(title), escapeAppleScript(notes), dueDate)
	} else {
		script = fmt.Sprintf(`
			tell application "Reminders"
				tell list "%s"
					make new reminder with properties {name:"%s", body:"%s"}
				end tell
			end tell
			return "OK"
		`, escapeAppleScript(list), escapeAppleScript(title), escapeAppleScript(notes))
	}

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create reminder: %v - %s", err, output)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created reminder: %s", title)), nil
}

// RemindersComplete marks a reminder as complete (macOS)
func RemindersComplete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, ok := req.Params.Arguments["title"].(string)
	if !ok || title == "" {
		return mcp.NewToolResultError("title is required"), nil
	}

	script := fmt.Sprintf(`
		tell application "Reminders"
			repeat with reminderList in every list
				set matchingReminders to (every reminder of reminderList whose name is "%s" and completed is false)
				repeat with r in matchingReminders
					set completed of r to true
					return "Completed"
				end repeat
			end repeat
			return "NotFound"
		end tell
	`, escapeAppleScript(title))

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to complete reminder: %v - %s", err, output)), nil
	}

	result := strings.TrimSpace(string(output))
	if result == "NotFound" {
		return mcp.NewToolResultText(fmt.Sprintf("Reminder '%s' not found", title)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Completed reminder: %s", title)), nil
}

// RemindersDelete deletes a reminder (macOS)
func RemindersDelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, ok := req.Params.Arguments["title"].(string)
	if !ok || title == "" {
		return mcp.NewToolResultError("title is required"), nil
	}

	script := fmt.Sprintf(`
		tell application "Reminders"
			repeat with reminderList in every list
				set matchingReminders to (every reminder of reminderList whose name is "%s")
				repeat with r in matchingReminders
					delete r
					return "Deleted"
				end repeat
			end repeat
			return "NotFound"
		end tell
	`, escapeAppleScript(title))

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to delete reminder: %v - %s", err, output)), nil
	}

	result := strings.TrimSpace(string(output))
	if result == "NotFound" {
		return mcp.NewToolResultText(fmt.Sprintf("Reminder '%s' not found", title)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted reminder: %s", title)), nil
}

// RemindersListLists lists all reminder lists (macOS)
func RemindersListLists(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	script := `
		tell application "Reminders"
			set output to ""
			repeat with reminderList in every list
				set listName to name of reminderList
				set reminderCount to count of (every reminder of reminderList whose completed is false)
				set output to output & listName & " (" & reminderCount & " items)" & linefeed
			end repeat
			return output
		end tell
	`

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list reminder lists: %v", err)), nil
	}

	return mcp.NewToolResultText("Reminder Lists:\n" + string(output)), nil
}
