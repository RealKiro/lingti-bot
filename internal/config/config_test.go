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
