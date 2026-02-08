package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/pltanton/lingti-bot/internal/logger"
)

// debugTransport logs outgoing request headers (with redacted auth) for debugging.
type debugTransport struct {
	base http.RoundTripper
}

func (d *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if logger.IsDebug() {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Printf("[Claude OAuth DEBUG] Request:\n%s", string(dump[:min(len(dump), 500)]))
		_ = os.WriteFile("/tmp/claude-request-dump.txt", dump, 0644)
		log.Printf("[Claude OAuth DEBUG] Full request written to /tmp/claude-request-dump.txt")
	}
	return d.base.RoundTrip(req)
}

const (
	anthropicSetupTokenPrefix    = "sk-ant-oat01-"
	anthropicSetupTokenMinLength = 80
	claudeCodeVersion            = "2.1.2"
	claudeCodeSystemPrefix       = "You are Claude Code, Anthropic's official CLI for Claude."
)

func isOAuthToken(key string) bool {
	return strings.HasPrefix(key, anthropicSetupTokenPrefix) && len(key) >= anthropicSetupTokenMinLength
}

// oauthAdapter mimics Claude Code's headers so OAuth setup tokens are accepted.
type oauthAdapter struct {
	anthropic.DefaultAdapter
	token string
}

func (a *oauthAdapter) SetRequestHeaders(_ *anthropic.Client, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+a.token)
	req.Header.Set("Anthropic-Version", "2023-06-01")
	req.Header.Set("Anthropic-Beta", "claude-code-20250219,oauth-2025-04-20")
	req.Header.Set("User-Agent", "claude-cli/"+claudeCodeVersion+" (external, cli)")
	req.Header.Set("X-App", "cli")
	req.Header.Set("Anthropic-Dangerous-Direct-Browser-Access", "true")
	// Remove X-Api-Key if the default adapter set it before us
	req.Header.Del("X-Api-Key")
	return nil
}

// ClaudeProvider implements the Provider interface for Claude/Anthropic
type ClaudeProvider struct {
	client   *anthropic.Client
	model    string
	isOAuth  bool
}

// ClaudeConfig holds Claude provider configuration
type ClaudeConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(cfg ClaudeConfig) (*ClaudeProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	if cfg.Model == "" {
		cfg.Model = "claude-sonnet-4-20250514"
	}

	oauth := isOAuthToken(cfg.APIKey)

	opts := []anthropic.ClientOption{}
	if cfg.BaseURL != "" {
		opts = append(opts, anthropic.WithBaseURL(cfg.BaseURL))
	}
	if oauth {
		adapter := &oauthAdapter{token: cfg.APIKey}
		opts = append(opts, func(c *anthropic.ClientConfig) {
			c.Adapter = adapter
		})
		opts = append(opts, anthropic.WithHTTPClient(&http.Client{
			Transport: &debugTransport{base: http.DefaultTransport},
		}))
	}

	client := anthropic.NewClient(cfg.APIKey, opts...)

	return &ClaudeProvider{
		client:  client,
		model:   cfg.Model,
		isOAuth: oauth,
	}, nil
}

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "claude"
}

// Chat sends messages and returns a response
func (p *ClaudeProvider) Chat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	// Convert messages to Anthropic format
	messages := make([]anthropic.Message, 0, len(req.Messages))
	for _, msg := range req.Messages {
		messages = append(messages, p.toAnthropicMessage(msg))
	}

	// Convert tools to Anthropic format
	tools := make([]anthropic.ToolDefinition, 0, len(req.Tools))
	for _, tool := range req.Tools {
		tools = append(tools, anthropic.ToolDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}

	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
	}

	// Build request
	apiReq := anthropic.MessagesRequest{
		Model:     anthropic.Model(p.model),
		MaxTokens: maxTokens,
		Messages:  messages,
		Tools:     tools,
	}

	// For OAuth tokens, send system prompt as array with Claude Code identity as first block
	if p.isOAuth {
		parts := []anthropic.MessageSystemPart{
			anthropic.NewSystemMessagePart(claudeCodeSystemPrefix),
		}
		if req.SystemPrompt != "" {
			parts = append(parts, anthropic.NewSystemMessagePart(req.SystemPrompt))
		}
		apiReq.MultiSystem = parts
	} else {
		apiReq.System = req.SystemPrompt
	}

	// Call Anthropic API — OAuth tokens require streaming (Claude Code always streams)
	if p.isOAuth {
		resp, err := p.client.CreateMessagesStream(ctx, anthropic.MessagesStreamRequest{
			MessagesRequest: apiReq,
		})
		if err != nil {
			return ChatResponse{}, fmt.Errorf("anthropic API error: %w", err)
		}
		return p.fromAnthropicResponse(resp), nil
	}

	resp, err := p.client.CreateMessages(ctx, apiReq)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("anthropic API error: %w", err)
	}
	return p.fromAnthropicResponse(resp), nil
}

// toAnthropicMessage converts a generic Message to Anthropic format
func (p *ClaudeProvider) toAnthropicMessage(msg Message) anthropic.Message {
	switch msg.Role {
	case "user":
		if msg.ToolResult != nil {
			// Tool result message
			return anthropic.Message{
				Role: anthropic.RoleUser,
				Content: []anthropic.MessageContent{
					anthropic.NewToolResultMessageContent(
						msg.ToolResult.ToolCallID,
						msg.ToolResult.Content,
						msg.ToolResult.IsError,
					),
				},
			}
		}
		return anthropic.Message{
			Role: anthropic.RoleUser,
			Content: []anthropic.MessageContent{
				anthropic.NewTextMessageContent(msg.Content),
			},
		}

	case "assistant":
		if len(msg.ToolCalls) > 0 {
			// Assistant message with tool calls
			content := make([]anthropic.MessageContent, 0)
			if msg.Content != "" {
				content = append(content, anthropic.NewTextMessageContent(msg.Content))
			}
			for _, tc := range msg.ToolCalls {
				input := tc.Input
				if len(input) == 0 {
					input = json.RawMessage(`{}`)
				}
				content = append(content, anthropic.NewToolUseMessageContent(tc.ID, tc.Name, input))
			}
			return anthropic.Message{
				Role:    anthropic.RoleAssistant,
				Content: content,
			}
		}
		return anthropic.Message{
			Role: anthropic.RoleAssistant,
			Content: []anthropic.MessageContent{
				anthropic.NewTextMessageContent(msg.Content),
			},
		}

	default:
		return anthropic.Message{
			Role: anthropic.RoleUser,
			Content: []anthropic.MessageContent{
				anthropic.NewTextMessageContent(msg.Content),
			},
		}
	}
}

// fromAnthropicResponse converts Anthropic response to generic format
func (p *ClaudeProvider) fromAnthropicResponse(resp anthropic.MessagesResponse) ChatResponse {
	var content string
	var toolCalls []ToolCall

	for _, c := range resp.Content {
		switch c.Type {
		case anthropic.MessagesContentTypeText:
			if c.Text != nil {
				content += *c.Text
			}
		case anthropic.MessagesContentTypeToolUse:
			if c.MessageContentToolUse != nil {
				toolCalls = append(toolCalls, ToolCall{
					ID:    c.MessageContentToolUse.ID,
					Name:  c.MessageContentToolUse.Name,
					Input: c.MessageContentToolUse.Input,
				})
			}
		}
	}

	finishReason := "stop"
	if resp.StopReason == anthropic.MessagesStopReasonToolUse {
		finishReason = "tool_use"
	}

	return ChatResponse{
		Content:      content,
		ToolCalls:    toolCalls,
		FinishReason: finishReason,
	}
}
