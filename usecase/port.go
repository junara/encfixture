// Package usecase implements the application logic for media file generation.
package usecase

import (
	"image"
	"image/color"
	"io"

	"github.com/junara/encfixture/domain"
)

// FFmpegExecutor defines the interface for running ffmpeg commands.
type FFmpegExecutor interface {
	Run(args ...string) error
	RunWithStdin(stdin io.Reader, args ...string) error
	CheckAvailable() error
}

// Renderer defines the interface for image rendering operations.
type Renderer interface {
	SolidImage(width, height int, c color.Color) *image.RGBA
	DrawScaledText(img *image.RGBA, text string, col color.Color, scale int)
	DrawScaledTextAt(img *image.RGBA, text string, col color.Color, scale int, pos domain.TextPosition)
	DrawTestPattern(img *image.RGBA)
	WritePNG(path string, img *image.RGBA) error
	ParseColor(name string) color.Color
	ContrastColor(bg color.Color) color.Color
}
