package usecase

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/junara/encfixture/domain"
)

// imageFPS is the nominal frame rate used when rendering a single still image,
// so animated background patterns resolve to their frame-zero state.
const imageFPS = 1

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

	writeErr := uc.renderer.WriteImage(cfg.Output, img, cfg.Quality)
	if writeErr != nil {
		return fmt.Errorf("write image failed: %w", writeErr)
	}

	return nil
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
