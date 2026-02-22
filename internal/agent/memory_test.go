package agent

import (
	"testing"
	"time"
)

func TestConversationKey(t *testing.T) {
	key := ConversationKey("slack", "C123", "U456")
	want := "slack:C123:U456"
	if key != want {
		t.Errorf("got %q, want %q", key, want)
	}
}

func TestMemory_AddAndGet(t *testing.T) {
	m := NewMemory(10, time.Hour)
	m.AddMessage("k1", Message{Role: "user", Content: "hello"})

	history := m.GetHistory("k1")
	if len(history) != 1 {
		t.Fatalf("expected 1 message, got %d", len(history))
	}
	if history[0].Content != "hello" {
		t.Errorf("got %q, want %q", history[0].Content, "hello")
	}
}

func TestMemory_MaxMessages(t *testing.T) {
	m := NewMemory(4, time.Hour)
	for i := range 6 {
		role := "user"
		if i%2 == 1 {
			role = "assistant"
		}
		m.AddMessage("k1", Message{Role: role, Content: "msg"})
	}
	history := m.GetHistory("k1")
	if len(history) > 4 {
		t.Errorf("expected at most 4 messages, got %d", len(history))
	}
}

func TestMemory_TTLExpiry(t *testing.T) {
	m := NewMemory(10, 50*time.Millisecond)
	m.AddMessage("k1", Message{Role: "user", Content: "hello"})

	time.Sleep(100 * time.Millisecond)

	history := m.GetHistory("k1")
	if len(history) != 0 {
		t.Errorf("expected empty history after TTL, got %d", len(history))
	}
}

func TestMemory_Clear(t *testing.T) {
	m := NewMemory(10, time.Hour)
	m.AddMessage("k1", Message{Role: "user", Content: "hello"})
	m.Clear("k1")

	if history := m.GetHistory("k1"); len(history) != 0 {
		t.Errorf("expected empty after clear, got %d", len(history))
	}
}

func TestMemory_AddExchange(t *testing.T) {
	m := NewMemory(10, time.Hour)
	m.AddExchange("k1",
		Message{Role: "user", Content: "hi"},
		Message{Role: "assistant", Content: "hello"},
	)

	history := m.GetHistory("k1")
	if len(history) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(history))
	}
	if history[0].Role != "user" || history[1].Role != "assistant" {
		t.Error("unexpected roles")
	}
}
