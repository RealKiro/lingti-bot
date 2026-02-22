package router

import (
	"context"
	"errors"
	"testing"
)

func TestFriendlyError(t *testing.T) {
	tests := []struct {
		err      string
		contains string
	}{
		{"overdue-payment blah", "充值"},
		{"invalid_api_key", "API Key 无效"},
		{"Rate limit exceeded", "频率超限"},
		{"only authorized for use with Claude Code", "Setup Token"},
		{"unexpected EOF", "连接中断"},
		{"EOF", "连接中断"},
		{"connection reset by peer", "连接中断"},
		{"some random error", "处理消息时出错"},
	}
	for _, tt := range tests {
		got := friendlyError(errors.New(tt.err))
		if !contains(got, tt.contains) {
			t.Errorf("friendlyError(%q) = %q, want containing %q", tt.err, got, tt.contains)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && searchString(s, sub)
}

func searchString(s, sub string) bool {
	for i := range len(s) - len(sub) + 1 {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestProgressContext(t *testing.T) {
	called := false
	fn := ProgressFunc(func(text string) { called = true })

	ctx := ContextWithProgress(context.Background(), fn)
	got := ProgressFromContext(ctx)
	if got == nil {
		t.Fatal("expected non-nil ProgressFunc")
	}
	got("test")
	if !called {
		t.Error("ProgressFunc was not called")
	}
}

func TestProgressContext_Nil(t *testing.T) {
	got := ProgressFromContext(context.Background())
	if got != nil {
		t.Error("expected nil ProgressFunc from plain context")
	}
}
