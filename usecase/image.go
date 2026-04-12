package usecase

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/junara/encfixture/domain"
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
	bgColor := uc.renderer.ParseColor(cfg.Color)
	textColor := uc.renderer.ContrastColor(bgColor)
	img := uc.renderer.SolidImage(cfg.Width, cfg.Height, bgColor)

	if cfg.Background == domain.BackgroundTest {
		uc.renderer.DrawTestPattern(img)
	}

	for _, entry := range cfg.Overlay.Entries() {
		text := resolveOverlayContent(entry.Content, 0, 1, cfg.Output)
		uc.renderer.DrawScaledTextAt(img, text, textColor, cfg.Scale, entry.Position)
	}

	writeErr := uc.renderer.WritePNG(cfg.Output, img)
	if writeErr != nil {
		return fmt.Errorf("write PNG failed: %w", writeErr)
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
