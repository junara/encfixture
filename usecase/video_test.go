package usecase_test

import (
	"errors"
	"image/color"
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

func TestVideoUseCase_Generate_CodecSelection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		codec  domain.VideoCodec
		expect string
	}{
		{domain.CodecH264, "libx264"},
		{domain.CodecHEVC, "libx265"},
		{domain.CodecVP9, "libvpx-vp9"},
		{domain.CodecAV1, "libaom-av1"},
		{domain.CodecProRes, "prores_ks"},
	}

	for _, tt := range tests {
		t.Run(string(tt.codec), func(t *testing.T) {
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
				Output:     "test.mkv",
				Audio:      domain.AudioSilence,
				SampleRate: 48000,
				Channels:   2,
				Codec:      tt.codec,
			}

			err := uc.Generate(cfg)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			args := strings.Join(ffmpeg.runArgs, " ")
			if !strings.Contains(args, "-c:v "+tt.expect) {
				t.Errorf("expected encoder %q in args, got: %s", tt.expect, args)
			}
		})
	}
}

func TestVideoUseCase_Generate_UnknownCodec(t *testing.T) {
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
		Audio:      domain.AudioSilence,
		SampleRate: 48000,
		Channels:   2,
		Codec:      "mpeg99",
	}

	err := uc.Generate(cfg)
	if !errors.Is(err, usecase.ErrUnknownVideoCodec) {
		t.Fatalf("Generate() error = %v, want ErrUnknownVideoCodec", err)
	}

	if ffmpeg.runCalled {
		t.Error("ffmpeg.Run should not be called for unknown codec")
	}
}

func TestVideoUseCase_Generate_EncodeOptions(t *testing.T) {
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
		Audio:      domain.AudioSilence,
		SampleRate: 48000,
		Channels:   2,
		Codec:      domain.CodecH264,
		CRF:        "23",
		Bitrate:    "5M",
		PixFmt:     "yuv444p",
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	args := strings.Join(ffmpeg.runArgs, " ")

	for _, want := range []string{"-crf 23", "-b:v 5M", "-pix_fmt yuv444p"} {
		if !strings.Contains(args, want) {
			t.Errorf("expected %q in args, got: %s", want, args)
		}
	}
}

func TestVideoUseCase_Generate_ProResDefaultPixFmt(t *testing.T) {
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
		Output:     "test.mov",
		Audio:      domain.AudioSilence,
		SampleRate: 48000,
		Channels:   2,
		Codec:      domain.CodecProRes,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	args := strings.Join(ffmpeg.runArgs, " ")
	if !strings.Contains(args, "-pix_fmt yuv422p10le") {
		t.Errorf("expected ProRes default pix_fmt yuv422p10le, got: %s", args)
	}
}

func TestVideoUseCase_Generate_SyncUsesFramesAndBeep(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	cfg := domain.VideoConfig{
		Width:        64,
		Height:       48,
		FPS:          30,
		Duration:     "2",
		Background:   "solid",
		Color:        "black",
		Scale:        4,
		Output:       "sync.mp4",
		Audio:        domain.AudioSilence, // overridden by sync
		SampleRate:   48000,
		Channels:     2,
		Frequency:    440,
		Sync:         true,
		SyncInterval: 1.0,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !ffmpeg.runWithStdinCalled {
		t.Fatal("sync mode should render frames (RunWithStdin), not use the simple path")
	}

	args := strings.Join(ffmpeg.runWithStdinArgs, " ")
	if !strings.Contains(args, "aevalsrc=exprs=") {
		t.Errorf("expected aevalsrc beep source, got: %s", args)
	}

	// Sync overrides --audio silence: no anullsrc.
	if strings.Contains(args, "anullsrc") {
		t.Errorf("sync should override silence audio, got: %s", args)
	}

	// Stereo → one beep expression per channel joined by '|'.
	if got := strings.Count(args, "sin(2*PI"); got != 2 {
		t.Errorf("expected 2 channel expressions, got %d in: %s", got, args)
	}
}

func TestVideoUseCase_Generate_SyncFlashFrames(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	cfg := domain.VideoConfig{
		Width:        8,
		Height:       8,
		FPS:          30,
		Duration:     "2",
		Background:   "solid",
		Color:        "black",
		Scale:        4,
		Output:       "sync.mp4",
		Audio:        domain.AudioSilence,
		SampleRate:   48000,
		Channels:     2,
		Frequency:    440,
		Sync:         true,
		SyncInterval: 1.0,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// 2s * 30fps = 60 frames; flash window = first int(0.08*30)=2 frames of each
	// 30-frame interval.
	if len(renderer.solidImageColors) != 60 {
		t.Fatalf("rendered %d frames, want 60", len(renderer.solidImageColors))
	}

	assertFrameColor(t, "frame 0 (flash)", renderer.solidImageColors[0], color.White)
	assertFrameColor(t, "frame 1 (flash)", renderer.solidImageColors[1], color.White)
	assertFrameColor(t, "frame 2 (normal)", renderer.solidImageColors[2], color.Black)
	assertFrameColor(t, "frame 15 (normal)", renderer.solidImageColors[15], color.Black)
	assertFrameColor(t, "frame 30 (flash)", renderer.solidImageColors[30], color.White)
}

func TestVideoUseCase_Generate_SyncCustomInterval(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	cfg := domain.VideoConfig{
		Width:        8,
		Height:       8,
		FPS:          30,
		Duration:     "1",
		Background:   "solid",
		Color:        "black",
		Scale:        4,
		Output:       "sync.mp4",
		Audio:        domain.AudioSilence,
		SampleRate:   48000,
		Channels:     2,
		Frequency:    440,
		Sync:         true,
		SyncInterval: 0.5,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	args := strings.Join(ffmpeg.runWithStdinArgs, " ")
	if !strings.Contains(args, `mod(t\,0.5)`) {
		t.Errorf("expected 0.5s interval in beep filter, got: %s", args)
	}

	// 30fps, 0.5s interval → flash every 15 frames.
	assertFrameColor(t, "frame 0 (flash)", renderer.solidImageColors[0], color.White)
	assertFrameColor(t, "frame 15 (flash)", renderer.solidImageColors[15], color.White)
	assertFrameColor(t, "frame 7 (normal)", renderer.solidImageColors[7], color.Black)
}

func TestVideoUseCase_Generate_SyncQuantizedAlignment(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	// interval*fps is not an integer (0.33*30 = 9.9 → 9 frames). The audio beep
	// must use the same frame-quantized 9/30 = 0.3s period as the video flash,
	// otherwise the two drift apart over the clip.
	cfg := domain.VideoConfig{
		Width:        8,
		Height:       8,
		FPS:          30,
		Duration:     "1",
		Background:   "solid",
		Color:        "black",
		Scale:        4,
		Output:       "sync.mp4",
		Audio:        domain.AudioSilence,
		SampleRate:   48000,
		Channels:     2,
		Frequency:    440,
		Sync:         true,
		SyncInterval: 0.33,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	args := strings.Join(ffmpeg.runWithStdinArgs, " ")
	if !strings.Contains(args, `mod(t\,0.3)`) {
		t.Errorf("audio beep must use the frame-quantized 0.3s period, got: %s", args)
	}

	// Flash starts on frame 9 (=0.3s), not on frame 8.
	assertFrameColor(t, "frame 8 (normal)", renderer.solidImageColors[8], color.Black)
	assertFrameColor(t, "frame 9 (flash)", renderer.solidImageColors[9], color.White)
}

func TestVideoUseCase_Generate_SyncMono(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	cfg := domain.VideoConfig{
		Width:        8,
		Height:       8,
		FPS:          30,
		Duration:     "1",
		Background:   "solid",
		Color:        "black",
		Scale:        4,
		Output:       "sync.mp4",
		Audio:        domain.AudioSilence,
		SampleRate:   48000,
		Channels:     1,
		Frequency:    440,
		Sync:         true,
		SyncInterval: 1.0,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	args := strings.Join(ffmpeg.runWithStdinArgs, " ")
	if got := strings.Count(args, "sin(2*PI"); got != 1 {
		t.Errorf("mono should emit 1 channel expression, got %d in: %s", got, args)
	}

	if !strings.Contains(args, "c=mono") {
		t.Errorf("expected c=mono layout, got: %s", args)
	}
}

func assertFrameColor(t *testing.T, label string, got, want color.Color) {
	t.Helper()

	gotR, gotG, gotB, _ := got.RGBA()
	wantR, wantG, wantB, _ := want.RGBA()

	if gotR != wantR || gotG != wantG || gotB != wantB {
		t.Errorf("%s: color = (%d,%d,%d), want (%d,%d,%d)", label, gotR, gotG, gotB, wantR, wantG, wantB)
	}
}

func TestVideoUseCase_Generate_MovingBackgrounds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		background string
		wantGrad   int
		wantBox    int
	}{
		{domain.BackgroundGradient, 30, 0},
		{domain.BackgroundMoving, 0, 30},
	}

	for _, tt := range tests {
		t.Run(tt.background, func(t *testing.T) {
			t.Parallel()

			ffmpeg := &mockFFmpeg{}
			renderer := newMockRenderer()
			uc := usecase.NewVideoUseCase(ffmpeg, renderer)

			cfg := domain.VideoConfig{
				Width:      64,
				Height:     48,
				FPS:        30,
				Duration:   "1",
				Background: tt.background,
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

			// Non-solid backgrounds render frames rather than using the lavfi path.
			if !ffmpeg.runWithStdinCalled {
				t.Error("expected frame rendering (RunWithStdin) for animated background")
			}

			// 1s * 30fps = 30 frames, each drawing the pattern once.
			if renderer.drawGradientCalls != tt.wantGrad {
				t.Errorf("DrawScrollingGradient called %d times, want %d", renderer.drawGradientCalls, tt.wantGrad)
			}

			if renderer.drawMovingBoxCalls != tt.wantBox {
				t.Errorf("DrawMovingBox called %d times, want %d", renderer.drawMovingBoxCalls, tt.wantBox)
			}
		})
	}
}

func TestVideoUseCase_Generate_EmptyBackgroundIsSolid(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	cfg := domain.VideoConfig{
		Width:      64,
		Height:     48,
		FPS:        30,
		Duration:   "1",
		Background: "", // unspecified → solid
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

	// An empty background is treated as solid: the simple lavfi path, no frames.
	if !ffmpeg.runCalled {
		t.Error("empty background should use the solid (simple) path")
	}

	if ffmpeg.runWithStdinCalled {
		t.Error("empty background should not render frames")
	}
}

func TestVideoUseCase_Generate_UnknownBackground(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	cfg := domain.VideoConfig{
		Width:      64,
		Height:     48,
		FPS:        30,
		Duration:   "1",
		Background: "sparkles",
		Color:      "black",
		Scale:      4,
		Output:     "test.mp4",
		Audio:      domain.AudioSilence,
		SampleRate: 48000,
		Channels:   2,
	}

	err := uc.Generate(cfg)
	if !errors.Is(err, usecase.ErrUnknownBackground) {
		t.Fatalf("Generate() error = %v, want ErrUnknownBackground", err)
	}

	if ffmpeg.runCalled || ffmpeg.runWithStdinCalled {
		t.Error("ffmpeg should not run for an unknown background")
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
