package config

import "testing"

func TestResolveAI_NoOverride(t *testing.T) {
	cfg := &Config{AI: AIConfig{Provider: "claude", APIKey: "key1", Model: "sonnet"}}
	resolved := cfg.ResolveAI("telegram", "chan1")
	if resolved.Provider != "claude" || resolved.APIKey != "key1" {
		t.Errorf("expected base config, got provider=%s key=%s", resolved.Provider, resolved.APIKey)
	}
}

func TestResolveAI_PlatformOverride(t *testing.T) {
	cfg := &Config{AI: AIConfig{
		Provider: "claude",
		APIKey:   "key1",
		Overrides: []AIOverride{
			{Platform: "telegram", Provider: "deepseek", APIKey: "key2"},
		},
	}}
	resolved := cfg.ResolveAI("telegram", "any-channel")
	if resolved.Provider != "deepseek" || resolved.APIKey != "key2" {
		t.Errorf("expected override, got provider=%s key=%s", resolved.Provider, resolved.APIKey)
	}
}

func TestResolveAI_ChannelOverride(t *testing.T) {
	cfg := &Config{AI: AIConfig{
		Provider: "claude",
		APIKey:   "key1",
		Overrides: []AIOverride{
			{Platform: "telegram", APIKey: "platform-key"},
			{Platform: "telegram", ChannelID: "special", APIKey: "channel-key"},
		},
	}}
	resolved := cfg.ResolveAI("telegram", "special")
	if resolved.APIKey != "channel-key" {
		t.Errorf("channel override should win, got key=%s", resolved.APIKey)
	}
}

func TestResolveProvider_ExactMatch(t *testing.T) {
	cfg := &Config{
		Providers: map[string]ProviderEntry{
			"my-kimi": {Provider: "kimi", APIKey: "ak-xxx", Model: "kimi-k2.5"},
		},
	}
	e, ok := cfg.ResolveProvider("my-kimi")
	if !ok || e.Provider != "kimi" || e.APIKey != "ak-xxx" {
		t.Errorf("expected exact match, got ok=%v entry=%+v", ok, e)
	}
}

func TestResolveProvider_ByProviderType(t *testing.T) {
	cfg := &Config{
		Providers: map[string]ProviderEntry{
			"my-kimi": {Provider: "kimi", APIKey: "ak-xxx", Model: "kimi-k2.5"},
		},
	}
	e, ok := cfg.ResolveProvider("kimi")
	if !ok || e.APIKey != "ak-xxx" {
		t.Errorf("expected match by provider type, got ok=%v entry=%+v", ok, e)
	}
}

func TestResolveProvider_BackwardCompat(t *testing.T) {
	cfg := &Config{
		AI: AIConfig{Provider: "deepseek", APIKey: "sk-xxx", Model: "deepseek-chat"},
	}
	e, ok := cfg.ResolveProvider("deepseek")
	if !ok || e.APIKey != "sk-xxx" || e.Model != "deepseek-chat" {
		t.Errorf("expected backward compat from ai: block, got ok=%v entry=%+v", ok, e)
	}
}

func TestResolveProvider_EmptyName(t *testing.T) {
	cfg := &Config{
		Relay: RelayConfig{Provider: "my-kimi"},
		Providers: map[string]ProviderEntry{
			"my-kimi": {Provider: "kimi", APIKey: "ak-xxx"},
		},
	}
	e, ok := cfg.ResolveProvider("")
	if !ok || e.Provider != "kimi" {
		t.Errorf("expected fallback to relay.provider, got ok=%v entry=%+v", ok, e)
	}
}

func TestResolveProvider_EmptyFallbackToAI(t *testing.T) {
	cfg := &Config{
		AI: AIConfig{Provider: "deepseek", APIKey: "sk-xxx"},
	}
	e, ok := cfg.ResolveProvider("")
	if !ok || e.Provider != "deepseek" {
		t.Errorf("expected fallback to ai.provider, got ok=%v entry=%+v", ok, e)
	}
}

func TestResolveProvider_DifferentProviderNoMap(t *testing.T) {
	cfg := &Config{
		AI: AIConfig{Provider: "deepseek", APIKey: "sk-xxx"},
	}
	e, ok := cfg.ResolveProvider("kimi")
	if !ok || e.Provider != "kimi" {
		t.Errorf("expected provider-only entry for CLI override, got ok=%v entry=%+v", ok, e)
	}
	if e.APIKey != "" {
		t.Errorf("expected empty api_key for different provider, got %s", e.APIKey)
	}
}

func TestApplyOverride_ProviderChange(t *testing.T) {
	base := AIConfig{Provider: "claude", BaseURL: "https://api.anthropic.com", Model: "sonnet"}
	result := applyOverride(base, AIOverride{Provider: "deepseek", APIKey: "dk"})
	if result.Provider != "deepseek" {
		t.Errorf("expected deepseek, got %s", result.Provider)
	}
	if result.BaseURL != "" {
		t.Errorf("base_url should be cleared on provider change, got %s", result.BaseURL)
	}
	if result.Model != "" {
		t.Errorf("model should be cleared on provider change, got %s", result.Model)
	}
}
