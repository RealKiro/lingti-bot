package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// WeatherCurrent gets current weather for a location using wttr.in
func WeatherCurrent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	location := ""
	if l, ok := req.Params.Arguments["location"].(string); ok {
		location = l
	}

	// URL encode the location
	encodedLoc := url.QueryEscape(location)
	if encodedLoc == "" {
		encodedLoc = "" // Empty means auto-detect
	}

	// Use wttr.in with compact format
	apiURL := fmt.Sprintf("https://wttr.in/%s?format=%%l:+%%c+%%C+%%t+%%h+%%w", encodedLoc)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get weather: %v", err)), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to read response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(body)), nil
}

// WeatherForecast gets weather forecast for a location
func WeatherForecast(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	location := ""
	if l, ok := req.Params.Arguments["location"].(string); ok {
		location = l
	}

	days := 3
	if d, ok := req.Params.Arguments["days"].(float64); ok && d > 0 && d <= 3 {
		days = int(d)
	}

	// URL encode the location
	encodedLoc := url.QueryEscape(location)

	// Use wttr.in with text format (limited days)
	apiURL := fmt.Sprintf("https://wttr.in/%s?%d&format=v2", encodedLoc, days)

	client := &http.Client{Timeout: 10 * time.Second}
	req2, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	req2.Header.Set("User-Agent", "curl/7.0") // wttr.in needs this for text output

	resp, err := client.Do(req2)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get forecast: %v", err)), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to read response: %v", err)), nil
	}

	// Clean ANSI codes for text output
	result := strings.ReplaceAll(string(body), "\x1b[", "")

	return mcp.NewToolResultText(result), nil
}
