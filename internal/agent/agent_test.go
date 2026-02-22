package agent

import (
	"strings"
	"testing"

	"github.com/pltanton/lingti-bot/internal/router"
)

func TestCreateProvider_ValidProviders(t *testing.T) {
	providers := []string{"claude", "deepseek", "kimi", "moonshot", "qwen", "qianwen", "tongyi", "openai", "zhipu", "gemini"}
	for _, p := range providers {
		_, err := createProvider(Config{Provider: p, APIKey: "test-key"})
		if err != nil {
			t.Errorf("createProvider(%q) failed: %v", p, err)
		}
	}
}

func TestCreateProvider_Aliases(t *testing.T) {
	aliases := []string{"glm", "gpt", "chatgpt", "google", "xai"}
	for _, alias := range aliases {
		_, err := createProvider(Config{Provider: alias, APIKey: "test-key"})
		if err != nil {
			t.Errorf("createProvider(%q) failed: %v", alias, err)
		}
	}
}

func TestCreateProvider_EmptyDefaultsClaude(t *testing.T) {
	p, err := createProvider(Config{Provider: "", APIKey: "test-key"})
	if err != nil {
		t.Fatalf("empty provider should default to claude: %v", err)
	}
	if p.Name() != "claude" {
		t.Errorf("expected claude, got %s", p.Name())
	}
}

func TestCreateProvider_Unknown(t *testing.T) {
	_, err := createProvider(Config{Provider: "nonexistent", APIKey: "test-key"})
	if err == nil {
		t.Error("expected error for unknown provider")
	}
	if !strings.Contains(err.Error(), "unknown provider") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCreateProvider_OllamaNoKey(t *testing.T) {
	_, err := createProvider(Config{Provider: "ollama"})
	if err != nil {
		t.Errorf("ollama should not require API key: %v", err)
	}
}

func TestHandleBuiltinCommand(t *testing.T) {
	agent, err := New(Config{Provider: "claude", APIKey: "test-key"})
	if err != nil {
		t.Fatalf("failed to create agent: %v", err)
	}

	tests := []struct {
		text    string
		handled bool
	}{
		{"/help", true},
		{"/new", true},
		{"/status", true},
		{"/whoami", true},
		{"/model", true},
		{"/tools", true},
		{"/think high", true},
		{"/verbose on", true},
		{"hello world", false},
	}
	for _, tt := range tests {
		msg := router.Message{Text: tt.text, Platform: "test", ChannelID: "c1", UserID: "u1", Username: "tester"}
		_, handled := agent.handleBuiltinCommand(msg)
		if handled != tt.handled {
			t.Errorf("handleBuiltinCommand(%q): got handled=%v, want %v", tt.text, handled, tt.handled)
		}
	}
}
