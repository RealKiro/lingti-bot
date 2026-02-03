package voice

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// Speaker handles text-to-speech output
type Speaker struct {
	provider Provider
	voice    string
	speed    float64
}

// SpeakerConfig holds speaker configuration
type SpeakerConfig struct {
	Provider string  // "system", "openai", "elevenlabs"
	APIKey   string  // API key for cloud providers
	Voice    string  // Voice name/ID
	Speed    float64 // Speech rate (1.0 = normal)
}

// NewSpeaker creates a new speaker
func NewSpeaker(cfg SpeakerConfig) (*Speaker, error) {
	var provider Provider
	var err error

	switch cfg.Provider {
	case "openai":
		provider, err = NewOpenAIProvider(cfg.APIKey)
	case "elevenlabs":
		provider, err = NewElevenLabsProvider(cfg.APIKey)
	case "system", "":
		provider = NewSystemProvider()
	default:
		return nil, fmt.Errorf("unknown voice provider: %s", cfg.Provider)
	}

	if err != nil {
		return nil, err
	}

	speed := cfg.Speed
	if speed == 0 {
		speed = 1.0
	}

	return &Speaker{
		provider: provider,
		voice:    cfg.Voice,
		speed:    speed,
	}, nil
}

// Speak converts text to speech and plays it
func (s *Speaker) Speak(ctx context.Context, text string) error {
	// For system provider on macOS, use say directly for efficiency
	if s.provider.Name() == "system" && runtime.GOOS == "darwin" {
		return s.speakDirect(ctx, text)
	}

	// Get audio from provider
	audio, err := s.provider.TextToSpeech(ctx, text, TTSOptions{
		Voice: s.voice,
		Speed: s.speed,
	})
	if err != nil {
		return err
	}

	// Play the audio
	return playAudioData(audio)
}

// speakDirect uses system TTS directly without intermediate files (macOS)
func (s *Speaker) speakDirect(ctx context.Context, text string) error {
	args := []string{}

	if s.voice != "" {
		args = append(args, "-v", s.voice)
	}

	if s.speed != 1.0 {
		rate := int(175 * s.speed) // 175 wpm is default
		args = append(args, "-r", fmt.Sprintf("%d", rate))
	}

	args = append(args, text)

	cmd := exec.CommandContext(ctx, "say", args...)
	return cmd.Run()
}

// SpeakAsync speaks text asynchronously
func (s *Speaker) SpeakAsync(ctx context.Context, text string) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Speak(ctx, text)
		close(errCh)
	}()
	return errCh
}

// ProviderName returns the name of the underlying provider
func (s *Speaker) ProviderName() string {
	return s.provider.Name()
}

// ListVoices lists available voices (system provider only)
func (s *Speaker) ListVoices() ([]string, error) {
	if s.provider.Name() != "system" {
		return nil, fmt.Errorf("voice listing only supported for system provider")
	}

	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("say", "-v", "?")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		// Parse output - each line starts with voice name
		var voices []string
		lines := splitLines(string(output))
		for _, line := range lines {
			if len(line) > 0 {
				// Voice name is first word
				parts := splitFirst(line, " ")
				if len(parts) > 0 {
					voices = append(voices, parts[0])
				}
			}
		}
		return voices, nil

	default:
		return nil, fmt.Errorf("voice listing not supported on %s", runtime.GOOS)
	}
}

// playAudioData plays audio data through the system speaker
func playAudioData(audio []byte) error {
	// Determine format from magic bytes
	format := "wav"
	if len(audio) > 3 && audio[0] == 0x49 && audio[1] == 0x44 && audio[2] == 0x33 {
		format = "mp3"
	} else if len(audio) > 4 && audio[0] == 0xff && (audio[1]&0xf0) == 0xf0 {
		format = "mp3"
	}

	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("speak-%d.%s", time.Now().UnixNano(), format))
	defer os.Remove(tmpFile)

	if err := os.WriteFile(tmpFile, audio, 0644); err != nil {
		return fmt.Errorf("failed to write audio file: %w", err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("afplay", tmpFile)
	case "linux":
		// Try multiple players
		if _, err := exec.LookPath("aplay"); err == nil && format == "wav" {
			cmd = exec.Command("aplay", "-q", tmpFile)
		} else if _, err := exec.LookPath("paplay"); err == nil {
			cmd = exec.Command("paplay", tmpFile)
		} else if _, err := exec.LookPath("mpv"); err == nil {
			cmd = exec.Command("mpv", "--no-video", "--really-quiet", tmpFile)
		} else if _, err := exec.LookPath("ffplay"); err == nil {
			cmd = exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", tmpFile)
		} else {
			return fmt.Errorf("no audio player found (install pulseaudio, mpv, or ffmpeg)")
		}
	case "windows":
		// Use PowerShell to play audio
		cmd = exec.Command("powershell", "-c",
			fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", tmpFile))
	default:
		return fmt.Errorf("audio playback not supported on %s", runtime.GOOS)
	}

	return cmd.Run()
}

// Helper functions

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func splitFirst(s string, sep string) []string {
	for i := 0; i < len(s)-len(sep)+1; i++ {
		if s[i:i+len(sep)] == sep {
			return []string{s[:i], s[i+len(sep):]}
		}
	}
	return []string{s}
}
