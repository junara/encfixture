package usecase_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/usecase"
)

func TestVideoUseCase_Generate_SimpleSolid(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	cfg := domain.VideoConfig{
		Width:      640,
		Height:     480,
		FPS:        30,
		Duration:   "5",
		Background: "solid",
		Color:      "black",
		Scale:      4,
		Output:     "test.mp4",
		Audio:      domain.AudioSilence,
		SampleRate: 48000,
		Channels:   2,
		Frequency:  440,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !ffmpeg.runCalled {
		t.Error("ffmpeg.Run was not called")
	}

	args := strings.Join(ffmpeg.runArgs, " ")
	if !strings.Contains(args, "color=c=black") {
		t.Errorf("expected color filter, got: %s", args)
	}

	if !strings.Contains(args, "test.mp4") {
		t.Errorf("expected output file, got: %s", args)
	}
}

func TestVideoUseCase_Generate_WithOverlays(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	cfg := domain.VideoConfig{
		Width:      64,
		Height:     48,
		FPS:        1,
		Duration:   "2",
		Background: "solid",
		Color:      "black",
		Overlay: domain.Overlay{
			TopLeft:  "frame",
			TopRight: "timecode",
		},
		Scale:      4,
		Output:     "test.mp4",
		Audio:      domain.AudioSilence,
		SampleRate: 48000,
		Channels:   2,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !ffmpeg.runWithStdinCalled {
		t.Error("ffmpeg.RunWithStdin was not called for overlay rendering")
	}

	if !renderer.solidImageCalled {
		t.Error("SolidImage was not called")
	}

	// 2 frames (1fps x 2s), each frame has 2 overlays
	if len(renderer.drawTextAtCalls) != 4 {
		t.Errorf("DrawScaledTextAt called %d times, want 4", len(renderer.drawTextAtCalls))
	}
}

func TestVideoUseCase_Generate_TestBackground(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	cfg := domain.VideoConfig{
		Width:      64,
		Height:     48,
		FPS:        1,
		Duration:   "1",
		Background: "test",
		Color:      "black",
		Scale:      4,
		Output:     "test.mp4",
		Audio:      domain.AudioSilence,
		SampleRate: 48000,
		Channels:   2,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !renderer.drawTestPatternCalled {
		t.Error("DrawTestPattern was not called")
	}
}

func TestVideoUseCase_Generate_VariousFormats(t *testing.T) {
	t.Parallel()

	formats := []string{".mp4", ".avi", ".mov", ".webm", ".mkv", ".ts", ".flv"}

	for _, ext := range formats {
		t.Run(ext, func(t *testing.T) {
			t.Parallel()

			ffmpeg := &mockFFmpeg{}
			renderer := newMockRenderer()
			uc := usecase.NewVideoUseCase(ffmpeg, renderer)

			cfg := domain.VideoConfig{
				Width:      640,
				Height:     480,
				FPS:        30,
				Duration:   "1",
				Background: "solid",
				Color:      "black",
				Scale:      4,
				Output:     "test" + ext,
				Audio:      domain.AudioSilence,
				SampleRate: 48000,
				Channels:   2,
			}

			err := uc.Generate(cfg)
			if err != nil {
				t.Errorf("Generate() error for %s: %v", ext, err)
			}
		})
	}
}

func TestVideoUseCase_Generate_WebmCodec(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	cfg := domain.VideoConfig{
		Width:      640,
		Height:     480,
		FPS:        30,
		Duration:   "1",
		Background: "solid",
		Color:      "black",
		Scale:      4,
		Output:     "test.webm",
		Audio:      domain.AudioSilence,
		SampleRate: 48000,
		Channels:   2,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	args := strings.Join(ffmpeg.runArgs, " ")
	if !strings.Contains(args, "libvpx-vp9") {
		t.Errorf("expected libvpx-vp9 codec for webm, got: %s", args)
	}

	if !strings.Contains(args, "libopus") {
		t.Errorf("expected libopus codec for webm, got: %s", args)
	}
}

func TestVideoUseCase_Generate_FFmpegUnavailable(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{checkErr: errors.New("not found")}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	cfg := domain.VideoConfig{
		Width:    640,
		Height:   480,
		FPS:      30,
		Duration: "5",
		Output:   "test.mp4",
	}

	err := uc.Generate(cfg)
	if err == nil {
		t.Fatal("Generate() should return error when ffmpeg unavailable")
	}
}

func TestVideoUseCase_Generate_AudioTypes(t *testing.T) {
	t.Parallel()

	audioTypes := []struct {
		audioType domain.AudioType
		expect    string
	}{
		{domain.AudioSilence, "anullsrc"},
		{domain.AudioSine, "sine=frequency"},
		{domain.AudioNoise, "anoisesrc"},
		{domain.AudioTone, "sine=frequency"},
	}

	for _, tt := range audioTypes {
		t.Run(string(tt.audioType), func(t *testing.T) {
			t.Parallel()

			ffmpeg := &mockFFmpeg{}
			renderer := newMockRenderer()
			uc := usecase.NewVideoUseCase(ffmpeg, renderer)

			cfg := domain.VideoConfig{
				Width:      640,
				Height:     480,
				FPS:        30,
				Duration:   "1",
				Background: "solid",
				Color:      "black",
				Scale:      4,
				Output:     "test.mp4",
				Audio:      tt.audioType,
				SampleRate: 48000,
				Channels:   2,
				Frequency:  440,
			}

			err := uc.Generate(cfg)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			args := strings.Join(ffmpeg.runArgs, " ")
			if !strings.Contains(args, tt.expect) {
				t.Errorf("expected %q in args for %s, got: %s", tt.expect, tt.audioType, args)
			}
		})
	}
}
