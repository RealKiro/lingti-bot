package mcpclient

import "testing"

func TestIsMCPTool(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"mcp_chrome_snapshot", true},
		{"mcp_", true},
		{"browser_click", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsMCPTool(tt.name); got != tt.want {
			t.Errorf("IsMCPTool(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"chrome-devtools", "chrome_devtools"},
		{"My Tool", "my_tool"},
		{"UPPER", "upper"},
		{"already_ok", "already_ok"},
	}
	for _, tt := range tests {
		if got := sanitizeName(tt.input); got != tt.want {
			t.Errorf("sanitizeName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
