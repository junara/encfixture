package usecase

import (
	"errors"
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
	// proResPixelFormat is the default pixel format for ProRes output (prores_ks).
	proResPixelFormat = "yuv422p10le"
	// syncFlashSeconds is the nominal duration of each A/V sync marker. It is
	// quantized to whole frames so the beep and the visual flash cover exactly
	// the same span, starting together so drift is visible.
	syncFlashSeconds = 0.08
	// defaultSyncInterval is the fallback gap between sync markers, in seconds.
	defaultSyncInterval = 1.0
	// syncBeepAmplitude keeps the beep below full scale to avoid clipping.
	syncBeepAmplitude = "0.5"
)

// ErrUnknownVideoCodec indicates an unrecognized video codec was specified.
var ErrUnknownVideoCodec = errors.New("unknown video codec")

func encoderName(codec domain.VideoCodec) (string, error) {
	switch codec {
	case domain.CodecH264:
		return "libx264", nil
	case domain.CodecHEVC:
		return "libx265", nil
	case domain.CodecVP9:
		return "libvpx-vp9", nil
	case domain.CodecAV1:
		return "libaom-av1", nil
	case domain.CodecProRes:
		return "prores_ks", nil
	default:
		return "", fmt.Errorf("%w: %s (supported: h264, hevc, vp9, av1, prores)", ErrUnknownVideoCodec, codec)
	}
}

// buildEncodeArgs assembles the ffmpeg output options for codec, quality and
// pixel format. defaultPixFmt is used when neither the config nor the codec
// dictates one; pass "" to leave the pixel format to ffmpeg.
func buildEncodeArgs(cfg domain.VideoConfig, ext, defaultPixFmt string) ([]string, error) {
	var args []string

	codec := cfg.Codec
	if codec == "" && ext == webmExtension {
		codec = domain.CodecVP9
	}

	if codec != "" {
		encoder, err := encoderName(codec)
		if err != nil {
			return nil, err
		}

		args = append(args, "-c:v", encoder)
	}

	if ext == webmExtension {
		args = append(args, "-c:a", "libopus")
	}

	if cfg.CRF != "" {
		args = append(args, "-crf", cfg.CRF)
	}

	if cfg.Bitrate != "" {
		args = append(args, "-b:v", cfg.Bitrate)
	}

	pixFmt := cfg.PixFmt
	if pixFmt == "" && codec == domain.CodecProRes {
		pixFmt = proResPixelFormat
	}

	if pixFmt == "" {
		pixFmt = defaultPixFmt
	}

	if pixFmt != "" {
		args = append(args, "-pix_fmt", pixFmt)
	}

	return args, nil
}

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

	if cfg.Background == domain.BackgroundTest || cfg.Overlay.HasContent() || cfg.Sync {
		return uc.generateWithFrames(cfg, ext)
	}

	return uc.generateSimple(cfg, ext)
}

// syncInterval returns the configured marker gap, falling back to the default
// when the value is unset or non-positive.
func syncInterval(cfg domain.VideoConfig) float64 {
	if cfg.SyncInterval > 0 {
		return cfg.SyncInterval
	}

	return defaultSyncInterval
}

// syncFrameWindow returns the marker period and flash length in frames. The
// visual flash and the audio beep are both derived from these frame counts, so
// they stay aligned even when interval*fps is not a whole number. Each is at
// least one frame, since a marker cannot be shorter than a single frame.
func syncFrameWindow(cfg domain.VideoConfig) (int, int) {
	fps := float64(cfg.FPS)
	intervalFrames := max(int(syncInterval(cfg)*fps), 1)
	flashFrames := max(int(syncFlashSeconds*fps), 1)

	return intervalFrames, flashFrames
}

// isSyncFlashFrame reports whether the given frame falls inside a sync marker's
// flash window (the first flashFrames of each interval).
func isSyncFlashFrame(cfg domain.VideoConfig, frameIdx int) bool {
	if cfg.FPS <= 0 {
		return false
	}

	intervalFrames, flashFrames := syncFrameWindow(cfg)

	return frameIdx%intervalFrames < flashFrames
}

// syncAudioWindow returns the beep period and length in seconds, matching the
// frame-quantized video flash so audio and video stay aligned. It falls back to
// the raw seconds when FPS is unusable (no frames are produced in that case).
func syncAudioWindow(cfg domain.VideoConfig) (float64, float64) {
	if cfg.FPS <= 0 {
		return syncInterval(cfg), syncFlashSeconds
	}

	intervalFrames, flashFrames := syncFrameWindow(cfg)
	fps := float64(cfg.FPS)

	return float64(intervalFrames) / fps, float64(flashFrames) / fps
}

func (uc *VideoUseCase) generateSimple(cfg domain.VideoConfig, ext string) error {
	videoFilter := fmt.Sprintf("color=c=%s:s=%dx%d:d=%s:r=%d", cfg.Color, cfg.Width, cfg.Height, cfg.Duration, cfg.FPS)
	audioFilter := uc.buildAudioFilter(cfg)

	encodeArgs, encodeErr := buildEncodeArgs(cfg, ext, "")
	if encodeErr != nil {
		return encodeErr
	}

	args := append([]string{
		"-y",
		"-f", "lavfi", "-i", videoFilter,
		"-f", "lavfi", "-i", audioFilter,
		"-t", cfg.Duration,
		"-shortest",
	}, encodeArgs...)
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

	encodeArgs, encodeErr := buildEncodeArgs(cfg, ext, "yuv420p")
	if encodeErr != nil {
		return encodeErr
	}

	totalFrames := int(durationSec * float64(cfg.FPS))
	bgColor := uc.renderer.ParseColor(cfg.Color)
	textColor := uc.renderer.ContrastColor(bgColor)
	pipeReader, pipeWriter := io.Pipe()
	errCh := make(chan error, 1)

	go uc.writeFrames(pipeWriter, errCh, cfg, totalFrames, bgColor, textColor)

	ffmpegErr := uc.runFFmpegWithPipe(pipeReader, cfg, encodeArgs)
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

func (uc *VideoUseCase) runFFmpegWithPipe(pipeReader *io.PipeReader, cfg domain.VideoConfig, encodeArgs []string) error {
	args := append([]string{
		"-y",
		"-f", "rawvideo",
		"-pixel_format", "rgba",
		"-video_size", fmt.Sprintf("%dx%d", cfg.Width, cfg.Height),
		"-framerate", strconv.Itoa(cfg.FPS),
		"-i", "pipe:0",
		"-f", "lavfi", "-i", uc.buildAudioFilter(cfg),
		"-t", cfg.Duration,
		"-shortest",
	}, encodeArgs...)
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
	frameBg, frameText := bgColor, textColor

	flash := cfg.Sync && isSyncFlashFrame(cfg, frameIdx)
	if flash {
		frameBg = color.White
		frameText = uc.renderer.ContrastColor(frameBg)
	}

	img := uc.renderer.SolidImage(cfg.Width, cfg.Height, frameBg)

	if cfg.Background == domain.BackgroundTest && !flash {
		uc.renderer.DrawTestPattern(img)
	}

	for _, entry := range cfg.Overlay.Entries() {
		text := resolveOverlayContent(entry.Content, frameIdx, cfg.FPS, cfg.Output)
		uc.renderer.DrawScaledTextAt(img, text, frameText, cfg.Scale, entry.Position)
	}

	return img
}

func (uc *VideoUseCase) buildAudioFilter(cfg domain.VideoConfig) string {
	if cfg.Sync {
		return buildSyncAudioFilter(cfg)
	}

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

// buildSyncAudioFilter builds a lavfi source that emits a short beep at the
// start of every interval, matching the visual flash produced by renderFrame.
// The comma inside the mod()/lt() calls is escaped so the filtergraph parser
// keeps it as a function argument rather than a filter separator.
func buildSyncAudioFilter(cfg domain.VideoConfig) string {
	intervalSec, flashSec := syncAudioWindow(cfg)

	expr := fmt.Sprintf(
		"%s*sin(2*PI*%.0f*t)*lt(mod(t\\,%g)\\,%g)",
		syncBeepAmplitude, cfg.Frequency, intervalSec, flashSec,
	)

	channels := max(cfg.Channels, 1)
	perChannel := make([]string, channels)

	for i := range perChannel {
		perChannel[i] = expr
	}

	return fmt.Sprintf(
		"aevalsrc=exprs=%s:s=%d:c=%s:d=%s",
		strings.Join(perChannel, "|"), cfg.SampleRate, channelLayout(cfg.Channels), cfg.Duration,
	)
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
