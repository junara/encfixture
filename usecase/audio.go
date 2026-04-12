package usecase

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/junara/encfixture/domain"
)

// ErrUnknownAudioType indicates an unrecognized audio type was specified.
var ErrUnknownAudioType = errors.New("unknown audio type")

// AudioUseCase handles audio file generation.
type AudioUseCase struct {
	ffmpeg FFmpegExecutor
}

// NewAudioUseCase creates a new AudioUseCase with the given ffmpeg executor.
func NewAudioUseCase(ffmpeg FFmpegExecutor) *AudioUseCase {
	return &AudioUseCase{ffmpeg: ffmpeg}
}

// Generate creates an audio file based on the provided configuration.
func (uc *AudioUseCase) Generate(cfg domain.AudioConfig) error {
	err := uc.ffmpeg.CheckAvailable()
	if err != nil {
		return fmt.Errorf("ffmpeg availability check failed: %w", err)
	}

	layout := channelLayout(cfg.Channels)

	var audioFilter string

	switch cfg.Type {
	case domain.AudioSilence:
		audioFilter = fmt.Sprintf("anullsrc=r=%d:cl=%s:d=%s", cfg.SampleRate, layout, cfg.Duration)
	case domain.AudioSine:
		audioFilter = fmt.Sprintf("sine=frequency=%.0f:sample_rate=%d:d=%s", cfg.Frequency, cfg.SampleRate, cfg.Duration)
	case domain.AudioNoise:
		audioFilter = fmt.Sprintf("anoisesrc=d=%s:c=white:r=%d:a=0.5", cfg.Duration, cfg.SampleRate)
	case domain.AudioTone:
		audioFilter = fmt.Sprintf("sine=frequency=%.0f:sample_rate=%d:d=%s", cfg.Frequency, cfg.SampleRate, cfg.Duration)
	default:
		return fmt.Errorf("%w: %s", ErrUnknownAudioType, cfg.Type)
	}

	args := []string{
		"-y",
		"-f", "lavfi",
		"-i", audioFilter,
		"-t", cfg.Duration,
	}

	if cfg.Channels > 1 && cfg.Type != domain.AudioSilence {
		args = append(args, "-ac", strconv.Itoa(cfg.Channels))
	}

	args = append(args, cfg.Output)

	runErr := uc.ffmpeg.Run(args...)
	if runErr != nil {
		return fmt.Errorf("ffmpeg audio generation failed: %w", runErr)
	}

	return nil
}
