package agent

import (
	"sync"
)

// ThinkingLevel represents the reasoning depth
type ThinkingLevel string

const (
	ThinkOff    ThinkingLevel = "off"
	ThinkLow    ThinkingLevel = "low"
	ThinkMedium ThinkingLevel = "medium"
	ThinkHigh   ThinkingLevel = "high"
)

// SessionSettings holds per-session configuration
type SessionSettings struct {
	ThinkingLevel ThinkingLevel
	Verbose       bool
}

// SessionStore manages session settings
type SessionStore struct {
	settings map[string]*SessionSettings
	mu       sync.RWMutex
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	return &SessionStore{
		settings: make(map[string]*SessionSettings),
	}
}

// Get returns settings for a session, creating defaults if needed
func (s *SessionStore) Get(key string) *SessionSettings {
	s.mu.RLock()
	settings, ok := s.settings[key]
	s.mu.RUnlock()

	if ok {
		return settings
	}

	// Create default settings
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring write lock
	if settings, ok = s.settings[key]; ok {
		return settings
	}

	settings = &SessionSettings{
		ThinkingLevel: ThinkMedium,
		Verbose:       false,
	}
	s.settings[key] = settings
	return settings
}

// SetThinkingLevel sets the thinking level for a session
func (s *SessionStore) SetThinkingLevel(key string, level ThinkingLevel) {
	settings := s.Get(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	settings.ThinkingLevel = level
}

// SetVerbose sets verbose mode for a session
func (s *SessionStore) SetVerbose(key string, verbose bool) {
	settings := s.Get(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	settings.Verbose = verbose
}

// Clear removes settings for a session
func (s *SessionStore) Clear(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.settings, key)
}

// ThinkingBudgetTokens returns the token budget for Claude extended thinking.
// Returns 0 for non-thinking levels.
func ThinkingBudgetTokens(level ThinkingLevel) int {
	switch level {
	case ThinkLow:
		return 1024
	case ThinkMedium:
		return 4096
	case ThinkHigh:
		return 16384
	default:
		return 0
	}
}

// ThinkingPrompt returns the thinking instruction based on level
func ThinkingPrompt(level ThinkingLevel) string {
	switch level {
	case ThinkOff:
		return ""
	case ThinkLow:
		return "\n\n## Thinking Mode: Low\nThink briefly before responding. Consider the main points."
	case ThinkMedium:
		return "\n\n## Thinking Mode: Medium\nThink through the problem step by step. Consider multiple approaches."
	case ThinkHigh:
		return "\n\n## Thinking Mode: High\nThink deeply and thoroughly. Analyze from multiple angles, consider edge cases, and reason carefully before responding."
	default:
		return ""
	}
}
