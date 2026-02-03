package tools

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/mark3labs/mcp-go/mcp"
)

// ClipboardRead reads content from the clipboard
func ClipboardRead(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.CommandContext(ctx, "pbpaste")
	case "linux":
		// Try xclip first, then xsel
		cmd = exec.CommandContext(ctx, "xclip", "-selection", "clipboard", "-o")
	case "windows":
		cmd = exec.CommandContext(ctx, "powershell", "-command", "Get-Clipboard")
	default:
		return mcp.NewToolResultError(fmt.Sprintf("clipboard not supported on %s", runtime.GOOS)), nil
	}

	output, err := cmd.Output()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to read clipboard: %v", err)), nil
	}

	if len(output) == 0 {
		return mcp.NewToolResultText("Clipboard is empty"), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// ClipboardWrite writes content to the clipboard
func ClipboardWrite(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content, ok := req.Params.Arguments["content"].(string)
	if !ok {
		return mcp.NewToolResultError("content is required"), nil
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.CommandContext(ctx, "pbcopy")
	case "linux":
		cmd = exec.CommandContext(ctx, "xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.CommandContext(ctx, "powershell", "-command", fmt.Sprintf("Set-Clipboard -Value '%s'", content))
		// For Windows, we don't need stdin
		if err := cmd.Run(); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to write clipboard: %v", err)), nil
		}
		return mcp.NewToolResultText("Content copied to clipboard"), nil
	default:
		return mcp.NewToolResultError(fmt.Sprintf("clipboard not supported on %s", runtime.GOOS)), nil
	}

	// For macOS and Linux, pipe content to stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create stdin pipe: %v", err)), nil
	}

	if err := cmd.Start(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to start clipboard command: %v", err)), nil
	}

	_, err = stdin.Write([]byte(content))
	stdin.Close()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to write to clipboard: %v", err)), nil
	}

	if err := cmd.Wait(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("clipboard command failed: %v", err)), nil
	}

	return mcp.NewToolResultText("Content copied to clipboard"), nil
}
