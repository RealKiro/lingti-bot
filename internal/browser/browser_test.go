package browser

import "testing"

func TestParseScreenSize(t *testing.T) {
	tests := []struct {
		input      string
		wantW      int
		wantH      int
		wantOK     bool
	}{
		{"1024x768", 1024, 768, true},
		{"1920X1080", 1920, 1080, true},
		{"abc", 0, 0, false},
		{"100x", 0, 0, false},
		{"x200", 0, 0, false},
		{"0x768", 0, 0, false},
		{"-1x768", 0, 0, false},
		{"", 0, 0, false},
	}
	for _, tt := range tests {
		w, h, ok := parseScreenSize(tt.input)
		if ok != tt.wantOK || w != tt.wantW || h != tt.wantH {
			t.Errorf("parseScreenSize(%q) = (%d, %d, %v), want (%d, %d, %v)",
				tt.input, w, h, ok, tt.wantW, tt.wantH, tt.wantOK)
		}
	}
}
