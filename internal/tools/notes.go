package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// NotesListFolders lists all note folders (macOS)
func NotesListFolders(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	script := `
		tell application "Notes"
			set output to ""
			repeat with f in folders
				set folderName to name of f
				set noteCount to count of notes of f
				set output to output & folderName & " (" & noteCount & " notes)" & linefeed
			end repeat
			return output
		end tell
	`

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list folders: %v", err)), nil
	}

	return mcp.NewToolResultText("Note Folders:\n" + string(output)), nil
}

// NotesListNotes lists notes in a folder (macOS)
func NotesListNotes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	folder := "Notes"
	if f, ok := req.Params.Arguments["folder"].(string); ok && f != "" {
		folder = f
	}

	limit := 20
	if l, ok := req.Params.Arguments["limit"].(float64); ok && l > 0 {
		limit = int(l)
	}

	script := fmt.Sprintf(`
		tell application "Notes"
			set output to ""
			set noteCount to 0
			tell folder "%s"
				repeat with n in notes
					if noteCount ≥ %d then exit repeat
					set noteName to name of n
					set noteDate to modification date of n
					set output to output & noteName & " | Modified: " & (noteDate as string) & linefeed
					set noteCount to noteCount + 1
				end repeat
			end tell
			return output
		end tell
	`, escapeAppleScript(folder), limit)

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list notes: %v", err)), nil
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No notes found in folder '%s'", folder)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Notes in %s:\n%s", folder, output)), nil
}

// NotesRead reads a note's content (macOS)
func NotesRead(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, ok := req.Params.Arguments["title"].(string)
	if !ok || title == "" {
		return mcp.NewToolResultError("title is required"), nil
	}

	script := fmt.Sprintf(`
		tell application "Notes"
			repeat with f in folders
				repeat with n in notes of f
					if name of n is "%s" then
						return plaintext of n
					end if
				end repeat
			end repeat
			return "NotFound"
		end tell
	`, escapeAppleScript(title))

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to read note: %v", err)), nil
	}

	result := strings.TrimSpace(string(output))
	if result == "NotFound" {
		return mcp.NewToolResultText(fmt.Sprintf("Note '%s' not found", title)), nil
	}

	return mcp.NewToolResultText(result), nil
}

// NotesCreate creates a new note (macOS)
func NotesCreate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, ok := req.Params.Arguments["title"].(string)
	if !ok || title == "" {
		return mcp.NewToolResultError("title is required"), nil
	}

	body := ""
	if b, ok := req.Params.Arguments["body"].(string); ok {
		body = b
	}

	folder := "Notes"
	if f, ok := req.Params.Arguments["folder"].(string); ok && f != "" {
		folder = f
	}

	// Create HTML content (Notes app uses HTML internally)
	content := fmt.Sprintf("<h1>%s</h1>", escapeHTML(title))
	if body != "" {
		content += "<br>" + escapeHTML(body)
	}

	script := fmt.Sprintf(`
		tell application "Notes"
			tell folder "%s"
				make new note with properties {body:"%s"}
			end tell
		end tell
		return "OK"
	`, escapeAppleScript(folder), escapeAppleScript(content))

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create note: %v - %s", err, output)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created note: %s", title)), nil
}

// NotesSearch searches notes by keyword (macOS)
func NotesSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keyword, ok := req.Params.Arguments["keyword"].(string)
	if !ok || keyword == "" {
		return mcp.NewToolResultError("keyword is required"), nil
	}

	script := fmt.Sprintf(`
		set searchTerm to "%s"
		set output to ""
		tell application "Notes"
			repeat with f in folders
				set folderName to name of f
				repeat with n in notes of f
					set noteName to name of n
					set noteBody to plaintext of n
					if noteName contains searchTerm or noteBody contains searchTerm then
						set output to output & noteName & " [" & folderName & "]" & linefeed
					end if
				end repeat
			end repeat
		end tell
		return output
	`, escapeAppleScript(keyword))

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to search notes: %v", err)), nil
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No notes found matching '%s'", keyword)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Notes matching '%s':\n%s", keyword, output)), nil
}

// NotesDelete deletes a note (macOS)
func NotesDelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, ok := req.Params.Arguments["title"].(string)
	if !ok || title == "" {
		return mcp.NewToolResultError("title is required"), nil
	}

	script := fmt.Sprintf(`
		tell application "Notes"
			repeat with f in folders
				repeat with n in notes of f
					if name of n is "%s" then
						delete n
						return "Deleted"
					end if
				end repeat
			end repeat
			return "NotFound"
		end tell
	`, escapeAppleScript(title))

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to delete note: %v - %s", err, output)), nil
	}

	result := strings.TrimSpace(string(output))
	if result == "NotFound" {
		return mcp.NewToolResultText(fmt.Sprintf("Note '%s' not found", title)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted note: %s", title)), nil
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
