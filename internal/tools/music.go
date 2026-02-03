package tools

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// MusicPlay starts or resumes music playback
func MusicPlay(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	app := detectMusicApp()
	if app == "" {
		return mcp.NewToolResultError("no music app detected (Spotify or Apple Music)"), nil
	}

	if runtime.GOOS != "darwin" {
		return mcp.NewToolResultError("music control only supported on macOS"), nil
	}

	script := fmt.Sprintf(`tell application "%s" to play`, app)
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to play: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Playing on %s", app)), nil
}

// MusicPause pauses music playback
func MusicPause(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	app := detectMusicApp()
	if app == "" {
		return mcp.NewToolResultError("no music app detected"), nil
	}

	if runtime.GOOS != "darwin" {
		return mcp.NewToolResultError("music control only supported on macOS"), nil
	}

	script := fmt.Sprintf(`tell application "%s" to pause`, app)
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to pause: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Paused %s", app)), nil
}

// MusicNext skips to the next track
func MusicNext(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	app := detectMusicApp()
	if app == "" {
		return mcp.NewToolResultError("no music app detected"), nil
	}

	if runtime.GOOS != "darwin" {
		return mcp.NewToolResultError("music control only supported on macOS"), nil
	}

	script := fmt.Sprintf(`tell application "%s" to next track`, app)
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to skip: %v", err)), nil
	}

	return mcp.NewToolResultText("Skipped to next track"), nil
}

// MusicPrevious goes to the previous track
func MusicPrevious(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	app := detectMusicApp()
	if app == "" {
		return mcp.NewToolResultError("no music app detected"), nil
	}

	if runtime.GOOS != "darwin" {
		return mcp.NewToolResultError("music control only supported on macOS"), nil
	}

	script := fmt.Sprintf(`tell application "%s" to previous track`, app)
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to go back: %v", err)), nil
	}

	return mcp.NewToolResultText("Went to previous track"), nil
}

// MusicNowPlaying gets the currently playing track
func MusicNowPlaying(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	app := detectMusicApp()
	if app == "" {
		return mcp.NewToolResultError("no music app detected"), nil
	}

	if runtime.GOOS != "darwin" {
		return mcp.NewToolResultError("music control only supported on macOS"), nil
	}

	var script string
	if app == "Spotify" {
		script = `
			tell application "Spotify"
				if player state is playing then
					set trackName to name of current track
					set artistName to artist of current track
					set albumName to album of current track
					return trackName & " by " & artistName & " (" & albumName & ")"
				else
					return "Not playing"
				end if
			end tell
		`
	} else {
		// Apple Music
		script = `
			tell application "Music"
				if player state is playing then
					set trackName to name of current track
					set artistName to artist of current track
					set albumName to album of current track
					return trackName & " by " & artistName & " (" & albumName & ")"
				else
					return "Not playing"
				end if
			end tell
		`
	}

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get now playing: %v", err)), nil
	}

	return mcp.NewToolResultText(strings.TrimSpace(string(output))), nil
}

// MusicSetVolume sets the music volume
func MusicSetVolume(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	volume, ok := req.Params.Arguments["volume"].(float64)
	if !ok || volume < 0 || volume > 100 {
		return mcp.NewToolResultError("volume is required (0-100)"), nil
	}

	app := detectMusicApp()
	if app == "" {
		return mcp.NewToolResultError("no music app detected"), nil
	}

	if runtime.GOOS != "darwin" {
		return mcp.NewToolResultError("music control only supported on macOS"), nil
	}

	script := fmt.Sprintf(`tell application "%s" to set sound volume to %d`, app, int(volume))
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to set volume: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Volume set to %d%%", int(volume))), nil
}

// MusicSearch searches for and plays a track
func MusicSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := req.Params.Arguments["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("query is required"), nil
	}

	if runtime.GOOS != "darwin" {
		return mcp.NewToolResultError("music control only supported on macOS"), nil
	}

	// Try Spotify first (it has better search)
	script := fmt.Sprintf(`
		tell application "Spotify"
			activate
			-- Use Spotify URI to search
			set searchURI to "spotify:search:" & "%s"
			open location searchURI
		end tell
		return "Searching in Spotify"
	`, escapeAppleScript(query))

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to search: %v - %s", err, output)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Searching for '%s' in Spotify", query)), nil
}

// detectMusicApp detects which music app is available
func detectMusicApp() string {
	if runtime.GOOS != "darwin" {
		return ""
	}

	// Check if Spotify is running
	script := `tell application "System Events" to (name of processes) contains "Spotify"`
	cmd := exec.Command("osascript", "-e", script)
	output, _ := cmd.Output()
	if strings.TrimSpace(string(output)) == "true" {
		return "Spotify"
	}

	// Check if Music is running
	script = `tell application "System Events" to (name of processes) contains "Music"`
	cmd = exec.Command("osascript", "-e", script)
	output, _ = cmd.Output()
	if strings.TrimSpace(string(output)) == "true" {
		return "Music"
	}

	// Default to Spotify if neither is running
	return "Spotify"
}
