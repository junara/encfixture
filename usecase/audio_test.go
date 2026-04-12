package usecase_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/usecase"
)

func TestAudioUseCase_Generate_Silence(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	uc := usecase.NewAudioUseCase(ffmpeg)

	cfg := domain.AudioConfig{
		Type:       domain.AudioSilence,
		Duration:   "10",
		SampleRate: 48000,
		Channels:   2,
		Frequency:  440,
		Output:     "test.wav",
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !ffmpeg.runCalled {
		t.Fatal("ffmpeg.Run was not called")
	}

	args := strings.Join(ffmpeg.runArgs, " ")
	if !strings.Contains(args, "anullsrc") {
		t.Errorf("expected anullsrc filter, got args: %s", args)
	}

	if !strings.Contains(args, "test.wav") {
		t.Errorf("expected output file in args, got: %s", args)
	}
}

func TestAudioUseCase_Generate_Sine(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	uc := usecase.NewAudioUseCase(ffmpeg)

	cfg := domain.AudioConfig{
		Type:       domain.AudioSine,
		Duration:   "5",
		SampleRate: 48000,
		Channels:   2,
		Frequency:  1000,
		Output:     "sine.wav",
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	args := strings.Join(ffmpeg.runArgs, " ")
	if !strings.Contains(args, "sine=frequency=1000") {
		t.Errorf("expected sine filter with 1000Hz, got args: %s", args)
	}
}

func TestAudioUseCase_Generate_Noise(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	uc := usecase.NewAudioUseCase(ffmpeg)

	cfg := domain.AudioConfig{
		Type:       domain.AudioNoise,
		Duration:   "3",
		SampleRate: 48000,
		Channels:   2,
		Output:     "noise.wav",
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	args := strings.Join(ffmpeg.runArgs, " ")
	if !strings.Contains(args, "anoisesrc") {
		t.Errorf("expected anoisesrc filter, got args: %s", args)
	}
}

func TestAudioUseCase_Generate_UnknownType(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	uc := usecase.NewAudioUseCase(ffmpeg)

	cfg := domain.AudioConfig{
		Type:     domain.AudioType("invalid"),
		Duration: "5",
		Output:   "test.wav",
	}

	err := uc.Generate(cfg)
	if err == nil {
		t.Fatal("Generate() should return error for unknown type")
	}

	if !errors.Is(err, usecase.ErrUnknownAudioType) {
		t.Errorf("expected ErrUnknownAudioType, got: %v", err)
	}
}

func TestAudioUseCase_Generate_FFmpegUnavailable(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{checkErr: errors.New("not found")}
	uc := usecase.NewAudioUseCase(ffmpeg)

	cfg := domain.AudioConfig{
		Type:     domain.AudioSilence,
		Duration: "5",
		Output:   "test.wav",
	}

	err := uc.Generate(cfg)
	if err == nil {
		t.Fatal("Generate() should return error when ffmpeg unavailable")
	}
}

func TestAudioUseCase_Generate_FFmpegRunError(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{runErr: errors.New("ffmpeg failed")}
	uc := usecase.NewAudioUseCase(ffmpeg)

	cfg := domain.AudioConfig{
		Type:       domain.AudioSilence,
		Duration:   "5",
		SampleRate: 48000,
		Channels:   2,
		Output:     "test.wav",
	}

	err := uc.Generate(cfg)
	if err == nil {
		t.Fatal("Generate() should return error when ffmpeg fails")
	}
}

func TestAudioUseCase_Generate_VariousFormats(t *testing.T) {
	t.Parallel()

	formats := []string{".wav", ".flac", ".aac", ".mp3", ".ogg", ".opus", ".wma"}

	for _, ext := range formats {
		t.Run(ext, func(t *testing.T) {
			t.Parallel()

			ffmpeg := &mockFFmpeg{}
			uc := usecase.NewAudioUseCase(ffmpeg)

			cfg := domain.AudioConfig{
				Type:       domain.AudioSilence,
				Duration:   "1",
				SampleRate: 48000,
				Channels:   2,
				Output:     "test" + ext,
			}

			err := uc.Generate(cfg)
			if err != nil {
				t.Errorf("Generate() error for %s: %v", ext, err)
			}
		})
	}
}

func TestAudioUseCase_Generate_MonoChannels(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	uc := usecase.NewAudioUseCase(ffmpeg)

	cfg := domain.AudioConfig{
		Type:       domain.AudioSine,
		Duration:   "1",
		SampleRate: 44100,
		Channels:   1,
		Frequency:  440,
		Output:     "mono.wav",
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// mono sine should not add -ac flag (channels == 1 is not > 1)
	args := strings.Join(ffmpeg.runArgs, " ")
	if strings.Contains(args, "-ac") {
		t.Errorf("mono sine should not have -ac flag, got: %s", args)
	}
}
