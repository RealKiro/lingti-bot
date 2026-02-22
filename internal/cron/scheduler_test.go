package cron

import "testing"

func TestNormalizeCron(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"43 * * * *", "0 43 * * * *"},
		{"0 9 * * 1-5", "0 0 9 * * 1-5"},
		{"0 43 * * * *", "0 43 * * * *"},   // already 6-field
		{"30 0 9 * * 1-5", "30 0 9 * * 1-5"}, // already 6-field
	}
	for _, tt := range tests {
		got := normalizeCron(tt.input)
		if got != tt.want {
			t.Errorf("normalizeCron(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
