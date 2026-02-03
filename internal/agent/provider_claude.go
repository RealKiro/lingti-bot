package agent

import (
	"context"
	"fmt"

	"github.com/liushuangls/go-anthropic/v2"
)

// ClaudeProvider implements the Provider interface for Claude/Anthropic
type ClaudeProvider struct {
	client *anthropic.Client
	model  string
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

	var client *anthropic.Client
	if cfg.BaseURL != "" {
		client = anthropic.NewClient(cfg.APIKey, anthropic.WithBaseURL(cfg.BaseURL))
	} else {
		client = anthropic.NewClient(cfg.APIKey)
	}

	return &ClaudeProvider{
		client: client,
		model:  cfg.Model,
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

	// Call Anthropic API
	resp, err := p.client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model:     anthropic.Model(p.model),
		MaxTokens: maxTokens,
		System:    req.SystemPrompt,
		Messages:  messages,
		Tools:     tools,
	})
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
				content = append(content, anthropic.NewToolUseMessageContent(tc.ID, tc.Name, tc.Input))
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
