package logger

import "testing"

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input   string
		want    Level
		wantErr bool
	}{
		{"trace", LevelTrace, false},
		{"debug", LevelDebug, false},
		{"info", LevelInfo, false},
		{"warn", LevelWarn, false},
		{"warning", LevelWarn, false},
		{"error", LevelError, false},
		{"fatal", LevelFatal, false},
		{"panic", LevelPanic, false},
		{"INFO", LevelInfo, false},
		{"bogus", LevelInfo, true},
	}
	for _, tt := range tests {
		got, err := ParseLevel(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseLevel(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
