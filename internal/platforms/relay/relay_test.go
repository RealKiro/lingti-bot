package relay

import "testing"

func TestWechatMediaType(t *testing.T) {
	tests := []struct {
		path      string
		mediaType string
		want      string
	}{
		{"photo.jpg", "", "image"},
		{"photo.PNG", "", "image"},
		{"photo.gif", "", "image"},
		{"audio.mp3", "", "voice"},
		{"audio.amr", "", "voice"},
		{"clip.mp4", "", "video"},
		{"doc.pdf", "", ""},
		{"file.txt", "", ""},
		{"any.bin", "image", "image"},
		{"any.bin", "voice", "voice"},
		{"any.bin", "video", "video"},
		{"any.bin", "thumb", "thumb"},
	}
	for _, tt := range tests {
		got := wechatMediaType(tt.path, tt.mediaType)
		if got != tt.want {
			t.Errorf("wechatMediaType(%q, %q) = %q, want %q", tt.path, tt.mediaType, got, tt.want)
		}
	}
}
