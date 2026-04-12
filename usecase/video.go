package usecase

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/junara/encfixture/domain"
)

const (
	// webmExtension is the file extension for WebM format.
	webmExtension    = ".webm"
	secondsPerHour   = 3600
	secondsPerMinute = 60
)

// VideoUseCase handles video file generation.
type VideoUseCase struct {
	ffmpeg   FFmpegExecutor
	renderer Renderer
}

// NewVideoUseCase creates a new VideoUseCase with the given ffmpeg executor and renderer.
func NewVideoUseCase(ffmpeg FFmpegExecutor, renderer Renderer) *VideoUseCase {
	return &VideoUseCase{ffmpeg: ffmpeg, renderer: renderer}
}

// Generate creates a video file based on the provided configuration.
func (uc *VideoUseCase) Generate(cfg domain.VideoConfig) error {
	err := uc.ffmpeg.CheckAvailable()
	if err != nil {
		return fmt.Errorf("ffmpeg availability check failed: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(cfg.Output))

	if cfg.Background == domain.BackgroundTest || cfg.Overlay.HasContent() {
		return uc.generateWithFrames(cfg, ext)
	}

	return uc.generateSimple(cfg, ext)
}

func (uc *VideoUseCase) generateSimple(cfg domain.VideoConfig, ext string) error {
	videoFilter := fmt.Sprintf("color=c=%s:s=%dx%d:d=%s:r=%d", cfg.Color, cfg.Width, cfg.Height, cfg.Duration, cfg.FPS)
	audioFilter := uc.buildAudioFilter(cfg)

	args := []string{
		"-y",
		"-f", "lavfi", "-i", videoFilter,
		"-f", "lavfi", "-i", audioFilter,
		"-t", cfg.Duration,
		"-shortest",
	}

	if ext == webmExtension {
		args = append(args, "-c:v", "libvpx-vp9", "-c:a", "libopus")
	}

	args = append(args, cfg.Output)

	runErr := uc.ffmpeg.Run(args...)
	if runErr != nil {
		return fmt.Errorf("ffmpeg video generation failed: %w", runErr)
	}

	return nil
}

func (uc *VideoUseCase) generateWithFrames(cfg domain.VideoConfig, ext string) error {
	durationSec, err := strconv.ParseFloat(cfg.Duration, 64)
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	totalFrames := int(durationSec * float64(cfg.FPS))
	bgColor := uc.renderer.ParseColor(cfg.Color)
	textColor := uc.renderer.ContrastColor(bgColor)
	pipeReader, pipeWriter := io.Pipe()
	errCh := make(chan error, 1)

	go uc.writeFrames(pipeWriter, errCh, cfg, totalFrames, bgColor, textColor)

	ffmpegErr := uc.runFFmpegWithPipe(pipeReader, cfg, ext)
	if ffmpegErr != nil {
		return fmt.Errorf("ffmpeg failed: %w", ffmpegErr)
	}

	frameErr := <-errCh
	if frameErr != nil {
		return fmt.Errorf("frame generation failed: %w", frameErr)
	}

	return nil
}

func (uc *VideoUseCase) writeFrames(
	pipeWriter *io.PipeWriter,
	errCh chan<- error,
	cfg domain.VideoConfig,
	totalFrames int,
	bgColor, textColor color.Color,
) {
	defer func() {
		closeErr := pipeWriter.Close()
		if closeErr != nil {
			errCh <- closeErr
		}
	}()

	for frameIdx := range totalFrames {
		frameImg := uc.renderFrame(cfg, frameIdx, bgColor, textColor)

		writeErr := writeRawRGBA(pipeWriter, frameImg)
		if writeErr != nil {
			errCh <- writeErr

			return
		}
	}

	errCh <- nil
}

func (uc *VideoUseCase) runFFmpegWithPipe(pipeReader *io.PipeReader, cfg domain.VideoConfig, ext string) error {
	args := []string{
		"-y",
		"-f", "rawvideo",
		"-pixel_format", "rgba",
		"-video_size", fmt.Sprintf("%dx%d", cfg.Width, cfg.Height),
		"-framerate", strconv.Itoa(cfg.FPS),
		"-i", "pipe:0",
		"-f", "lavfi", "-i", uc.buildAudioFilter(cfg),
		"-t", cfg.Duration,
		"-pix_fmt", "yuv420p",
		"-shortest",
	}

	if ext == webmExtension {
		args = append(args, "-c:v", "libvpx-vp9", "-c:a", "libopus")
	}

	args = append(args, cfg.Output)

	runErr := uc.ffmpeg.RunWithStdin(pipeReader, args...)
	if runErr != nil {
		return fmt.Errorf("ffmpeg pipe execution failed: %w", runErr)
	}

	return nil
}

func (uc *VideoUseCase) renderFrame(
	cfg domain.VideoConfig,
	frameIdx int,
	bgColor, textColor color.Color,
) *image.RGBA {
	img := uc.renderer.SolidImage(cfg.Width, cfg.Height, bgColor)

	if cfg.Background == domain.BackgroundTest {
		uc.renderer.DrawTestPattern(img)
	}

	for _, entry := range cfg.Overlay.Entries() {
		text := resolveOverlayContent(entry.Content, frameIdx, cfg.FPS, cfg.Output)
		uc.renderer.DrawScaledTextAt(img, text, textColor, cfg.Scale, entry.Position)
	}

	return img
}

func (uc *VideoUseCase) buildAudioFilter(cfg domain.VideoConfig) string {
	layout := channelLayout(cfg.Channels)

	switch cfg.Audio {
	case domain.AudioSilence:
		return fmt.Sprintf("anullsrc=r=%d:cl=%s:d=%s", cfg.SampleRate, layout, cfg.Duration)
	case domain.AudioSine:
		return fmt.Sprintf("sine=frequency=%.0f:sample_rate=%d:d=%s", cfg.Frequency, cfg.SampleRate, cfg.Duration)
	case domain.AudioNoise:
		return fmt.Sprintf("anoisesrc=d=%s:c=white:r=%d:a=0.5", cfg.Duration, cfg.SampleRate)
	case domain.AudioTone:
		return fmt.Sprintf("sine=frequency=%.0f:sample_rate=%d:d=%s", cfg.Frequency, cfg.SampleRate, cfg.Duration)
	default:
		return fmt.Sprintf("anullsrc=r=%d:cl=%s:d=%s", cfg.SampleRate, layout, cfg.Duration)
	}
}

func writeRawRGBA(writer io.Writer, img *image.RGBA) error {
	_, err := writer.Write(img.Pix)
	if err != nil {
		return fmt.Errorf("raw RGBA write failed: %w", err)
	}

	return nil
}

func channelLayout(channels int) string {
	switch channels {
	case 1:
		return "mono"
	case 2:
		return "stereo"
	case 6:
		return "5.1"
	case 8:
		return "7.1"
	default:
		return fmt.Sprintf("%d channels", channels)
	}
}
