package usecase

import (
	"fmt"
	"image"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/junara/encfixture/domain"
)

const (
	// imageFPS is the nominal frame rate used when rendering a single still image,
	// so animated background patterns resolve to their frame-zero state.
	imageFPS = 1
	// minJPEGQuality and maxJPEGQuality bound the search for the quality that
	// lands closest to a requested target size.
	minJPEGQuality = 1
	maxJPEGQuality = 100
)

// ImageUseCase handles image file generation.
type ImageUseCase struct {
	renderer Renderer
}

// NewImageUseCase creates a new ImageUseCase with the given renderer.
func NewImageUseCase(renderer Renderer) *ImageUseCase {
	return &ImageUseCase{renderer: renderer}
}

// Generate creates an image file based on the provided configuration.
func (uc *ImageUseCase) Generate(cfg domain.ImageConfig) error {
	if cfg.Background == "" {
		cfg.Background = domain.BackgroundSolid
	}

	if !domain.IsValidBackground(cfg.Background) {
		return fmt.Errorf("%w: %s (supported: solid, test, gradient, moving)", ErrUnknownBackground, cfg.Background)
	}

	bgColor := uc.renderer.ParseColor(cfg.Color)
	textColor := uc.renderer.ContrastColor(bgColor)
	img := uc.renderer.SolidImage(cfg.Width, cfg.Height, bgColor)

	drawBackground(uc.renderer, img, cfg.Background, 0, imageFPS)

	for _, entry := range cfg.Overlay.Entries() {
		text := resolveOverlayContent(entry.Content, 0, imageFPS, cfg.Output)
		uc.renderer.DrawScaledTextAt(img, text, textColor, cfg.Scale, entry.Position)
	}

	quality := cfg.Quality

	if cfg.TargetSize > 0 {
		fitted, fitErr := uc.qualityForTargetSize(cfg, img)
		if fitErr != nil {
			return fitErr
		}

		quality = fitted
	}

	writeErr := uc.renderer.WriteImage(cfg.Output, img, quality)
	if writeErr != nil {
		return fmt.Errorf("write image failed: %w", writeErr)
	}

	return nil
}

// qualityForTargetSize returns the highest JPEG quality whose encoded size
// stays within cfg.TargetSize, found by bisecting the quality range. Only JPEG
// output has such a dial: PNG is lossless, so its size follows the pixels
// alone. The result is approximate - when even quality 1 overshoots, or
// quality 100 undershoots, the closest end of the range is returned and the
// caller reports the size actually produced.
func (uc *ImageUseCase) qualityForTargetSize(cfg domain.ImageConfig, img *image.RGBA) (int, error) {
	ext := strings.ToLower(filepath.Ext(cfg.Output))
	if ext != ".jpg" && ext != ".jpeg" {
		return 0, fmt.Errorf(
			"%w: %s has no size dial (PNG is lossless); write .jpg/.jpeg, or change --width/--height",
			domain.ErrTargetSizeUnsupported, ext,
		)
	}

	best := minJPEGQuality
	low, high := minJPEGQuality, maxJPEGQuality

	for low <= high {
		mid := (low + high) / 2

		encoded, err := uc.renderer.EncodeImage(cfg.Output, img, mid)
		if err != nil {
			return 0, fmt.Errorf("encode image for target size: %w", err)
		}

		if int64(len(encoded)) <= cfg.TargetSize {
			best = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return best, nil
}

func resolveOverlayContent(content string, frameIdx, fps int, output string) string {
	switch content {
	case domain.KeywordFrame:
		return strconv.Itoa(frameIdx)
	case domain.KeywordTimecode:
		totalSeconds := frameIdx / fps
		frames := frameIdx % fps
		hours := totalSeconds / secondsPerHour
		minutes := (totalSeconds % secondsPerHour) / secondsPerMinute
		seconds := totalSeconds % secondsPerMinute

		return fmt.Sprintf("%02d:%02d:%02d:%02d", hours, minutes, seconds, frames)
	case domain.KeywordFilename:
		return filepath.Base(output)
	default:
		return content
	}
}
