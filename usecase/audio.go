package usecase

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/junara/encfixture/domain"
)

// pcmBytesPerSample is the width of one 16-bit PCM sample, used to suggest a
// sample rate that would land an uncompressed file on a requested size.
const pcmBytesPerSample = 2

// minAudioBitrate is the lowest bitrate an audio encoder will accept; below it
// the requested size cannot be met by any encoder setting.
const minAudioBitrate = 8000

// isUncompressedAudio reports whether ext names a format whose size no bitrate
// option can steer: PCM size follows sample rate, channels and duration, and
// lossless codecs follow the content.
func isUncompressedAudio(ext string) bool {
	switch ext {
	case ".wav", ".aif", ".aiff", ".au", ".pcm", ".raw", ".flac":
		return true
	default:
		return false
	}
}

// ErrUnknownAudioType indicates an unrecognized audio type was specified.
var ErrUnknownAudioType = errors.New("unknown audio type")

// targetSizeAudioArgs derives the audio bitrate that aims the encode at
// cfg.TargetSize. Only compressed formats have such a dial; for PCM and
// lossless output it reports the sample rate that would reach the size
// instead, since that is the knob which actually moves it.
func targetSizeAudioArgs(cfg domain.AudioConfig) ([]string, error) {
	ext := strings.ToLower(filepath.Ext(cfg.Output))
	if isUncompressedAudio(ext) {
		return nil, fmt.Errorf(
			"%w: %s is not bitrate-controlled%s; write a compressed format such as .mp3, .m4a or .opus",
			domain.ErrTargetSizeUnsupported, ext, sampleRateHint(cfg, ext),
		)
	}

	durationSec, err := strconv.ParseFloat(cfg.Duration, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", domain.ErrInvalidDuration, cfg.Duration)
	}

	bitrate := domain.BitrateForSize(cfg.TargetSize, durationSec, 0)
	if bitrate < minAudioBitrate {
		return nil, fmt.Errorf(
			"%w: %s over %ss works out to %d bps, below the %d bps an encoder can use;"+
				" raise --target-size or shorten -d",
			domain.ErrTargetSizeTooSmall, domain.FormatFileSize(cfg.TargetSize), cfg.Duration,
			bitrate, minAudioBitrate,
		)
	}

	return []string{"-b:a", strconv.Itoa(bitrate)}, nil
}

// sampleRateHint returns the " (try --sample-rate N)" fragment for PCM output,
// where the sample rate is what sizes the file. Lossless codecs get no hint:
// their size follows the content, not a setting.
func sampleRateHint(cfg domain.AudioConfig, ext string) string {
	if ext == ".flac" || cfg.Channels < 1 {
		return ""
	}

	durationSec, err := strconv.ParseFloat(cfg.Duration, 64)
	if err != nil || durationSec <= 0 {
		return ""
	}

	rate := int(float64(cfg.TargetSize) / (durationSec * float64(cfg.Channels) * pcmBytesPerSample))
	if rate <= 0 {
		return ""
	}

	return fmt.Sprintf(" (16-bit PCM reaches %s at --sample-rate %d)", domain.FormatFileSize(cfg.TargetSize), rate)
}

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
	durErr := domain.ValidateDuration(cfg.Duration)
	if durErr != nil {
		return fmt.Errorf("%w", durErr)
	}

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
		"-f", lavfiFormat,
		"-i", audioFilter,
		"-t", cfg.Duration,
	}

	if cfg.TargetSize > 0 {
		sizeArgs, sizeErr := targetSizeAudioArgs(cfg)
		if sizeErr != nil {
			return sizeErr
		}

		args = append(args, sizeArgs...)
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
