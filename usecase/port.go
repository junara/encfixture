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

// Prober inspects the properties of an existing media file.
type Prober interface {
	CheckAvailable() error
	Probe(path string) (domain.MediaInfo, error)
}

// Inspector reports the availability of external tools and of the encoders
// compiled into the local ffmpeg build.
type Inspector interface {
	ToolStatus(name string) domain.ToolStatus
	Encoders() (map[string]bool, error)
}

// Renderer defines the interface for image rendering operations.
type Renderer interface {
	SolidImage(width, height int, c color.Color) *image.RGBA
	DrawScaledText(img *image.RGBA, text string, col color.Color, scale int)
	DrawScaledTextAt(img *image.RGBA, text string, col color.Color, scale int, pos domain.TextPosition)
	DrawTestPattern(img *image.RGBA)
	DrawScrollingGradient(img *image.RGBA, frameIdx, fps int)
	DrawMovingBox(img *image.RGBA, frameIdx, fps int)
	WriteImage(path string, img *image.RGBA, quality int) error
	ParseColor(name string) color.Color
	ContrastColor(bg color.Color) color.Color
}
