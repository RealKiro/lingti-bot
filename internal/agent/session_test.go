package agent

import "testing"

func TestSessionStore_Defaults(t *testing.T) {
	s := NewSessionStore()
	settings := s.Get("new-key")
	if settings.ThinkingLevel != ThinkMedium {
		t.Errorf("default thinking level: got %q, want %q", settings.ThinkingLevel, ThinkMedium)
	}
	if settings.Verbose {
		t.Error("default verbose should be false")
	}
}

func TestSessionStore_ThinkingLevel(t *testing.T) {
	s := NewSessionStore()
	s.SetThinkingLevel("k1", ThinkHigh)
	if s.Get("k1").ThinkingLevel != ThinkHigh {
		t.Errorf("expected ThinkHigh")
	}
}

func TestSessionStore_Verbose(t *testing.T) {
	s := NewSessionStore()
	s.SetVerbose("k1", true)
	if !s.Get("k1").Verbose {
		t.Error("expected verbose true")
	}
}

func TestSessionStore_Clear(t *testing.T) {
	s := NewSessionStore()
	s.SetThinkingLevel("k1", ThinkHigh)
	s.Clear("k1")
	if s.Get("k1").ThinkingLevel != ThinkMedium {
		t.Error("expected defaults after clear")
	}
}

func TestThinkingBudgetTokens(t *testing.T) {
	tests := []struct {
		level ThinkingLevel
		want  int
	}{
		{ThinkOff, 0},
		{ThinkLow, 1024},
		{ThinkMedium, 4096},
		{ThinkHigh, 16384},
	}
	for _, tt := range tests {
		if got := ThinkingBudgetTokens(tt.level); got != tt.want {
			t.Errorf("ThinkingBudgetTokens(%q) = %d, want %d", tt.level, got, tt.want)
		}
	}
}

func TestThinkingPrompt(t *testing.T) {
	if ThinkingPrompt(ThinkOff) != "" {
		t.Error("ThinkOff should return empty prompt")
	}
	for _, level := range []ThinkingLevel{ThinkLow, ThinkMedium, ThinkHigh} {
		if ThinkingPrompt(level) == "" {
			t.Errorf("ThinkingPrompt(%q) should not be empty", level)
		}
	}
}
