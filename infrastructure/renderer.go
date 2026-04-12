package infrastructure

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/junara/encfixture/domain"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	centerDivisor      = 2
	textMargin         = 4
	secondsPerHour     = 3600
	secondsPerMinute   = 60
	luminanceThreshold = 0.5
	maxColorComponentF = 65535.0
)

// ImageRenderer provides image rendering and manipulation operations.
type ImageRenderer struct{}

// NewImageRenderer creates a new ImageRenderer instance.
func NewImageRenderer() *ImageRenderer {
	return &ImageRenderer{}
}

// SolidImage creates a new RGBA image filled with a solid color.
func (r *ImageRenderer) SolidImage(width, height int, clr color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for yPos := range height {
		for xPos := range width {
			img.Set(xPos, yPos, clr)
		}
	}

	return img
}

// DrawScaledText draws centered text on the image at the given scale factor.
func (r *ImageRenderer) DrawScaledText(img *image.RGBA, text string, col color.Color, scale int) {
	if scale < 1 {
		scale = 1
	}

	face := basicfont.Face7x13
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	scaledW := imgWidth / scale
	scaledH := imgHeight / scale

	small := image.NewRGBA(image.Rect(0, 0, scaledW, scaledH))

	lines := splitLines(text, scaledW, face)
	lineHeight := face.Metrics().Height.Ceil()
	totalHeight := lineHeight * len(lines)
	startY := (scaledH-totalHeight)/centerDivisor + face.Metrics().Ascent.Ceil()

	for idx, line := range lines {
		textWidth := measureText(line, face)
		xPos := (scaledW - textWidth) / centerDivisor
		yPos := startY + idx*lineHeight

		drawer := &font.Drawer{
			Dst:  small,
			Src:  image.NewUniform(col),
			Face: face,
			Dot:  fixed.P(xPos, yPos),
		}
		drawer.DrawString(line)
	}

	for yPos := range imgHeight {
		for xPos := range imgWidth {
			srcX := xPos / scale
			srcY := yPos / scale

			if srcX < scaledW && srcY < scaledH {
				pixel := small.At(srcX, srcY)

				_, _, _, alpha := pixel.RGBA()
				if alpha > 0 {
					img.Set(xPos, yPos, pixel)
				}
			}
		}
	}
}

// DrawScaledTextAt draws text at a specific position on the image at the given scale factor.
func (r *ImageRenderer) DrawScaledTextAt(
	img *image.RGBA, text string, col color.Color, scale int, pos domain.TextPosition,
) {
	if scale < 1 {
		scale = 1
	}

	face := basicfont.Face7x13
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	scaledW := imgWidth / scale
	scaledH := imgHeight / scale

	small := image.NewRGBA(image.Rect(0, 0, scaledW, scaledH))

	lines := splitLines(text, scaledW-textMargin*centerDivisor, face)
	lineHeight := face.Metrics().Height.Ceil()
	totalHeight := lineHeight * len(lines)

	for idx, line := range lines {
		textWidth := measureText(line, face)
		xPos, yPos := calcTextPosition(pos, scaledW, scaledH, textWidth, totalHeight, lineHeight, idx, face)

		drawer := &font.Drawer{
			Dst:  small,
			Src:  image.NewUniform(col),
			Face: face,
			Dot:  fixed.P(xPos, yPos),
		}
		drawer.DrawString(line)
	}

	for yPos := range imgHeight {
		for xPos := range imgWidth {
			srcX := xPos / scale
			srcY := yPos / scale

			if srcX < scaledW && srcY < scaledH {
				pixel := small.At(srcX, srcY)

				_, _, _, alpha := pixel.RGBA()
				if alpha > 0 {
					img.Set(xPos, yPos, pixel)
				}
			}
		}
	}
}

func calcTextPosition(
	pos domain.TextPosition,
	scaledW, scaledH, textWidth, totalHeight, lineHeight, lineIdx int,
	face font.Face,
) (int, int) {
	ascent := face.Metrics().Ascent.Ceil()

	var xPos, yPos int

	switch pos {
	case domain.PositionTopLeft:
		xPos = textMargin
		yPos = textMargin + ascent + lineIdx*lineHeight
	case domain.PositionTopRight:
		xPos = scaledW - textWidth - textMargin
		yPos = textMargin + ascent + lineIdx*lineHeight
	case domain.PositionCenter:
		xPos = (scaledW - textWidth) / centerDivisor
		yPos = (scaledH-totalHeight)/centerDivisor + ascent + lineIdx*lineHeight
	case domain.PositionBottomLeft:
		xPos = textMargin
		yPos = scaledH - totalHeight - textMargin + ascent + lineIdx*lineHeight
	case domain.PositionBottomRight:
		xPos = scaledW - textWidth - textMargin
		yPos = scaledH - totalHeight - textMargin + ascent + lineIdx*lineHeight
	default:
		xPos = (scaledW - textWidth) / centerDivisor
		yPos = (scaledH-totalHeight)/centerDivisor + ascent + lineIdx*lineHeight
	}

	return xPos, yPos
}

// DrawTestPattern draws SMPTE-style color bars on the image.
func (r *ImageRenderer) DrawTestPattern(img *image.RGBA) {
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	colors := []color.Color{
		color.RGBA{R: 192, G: 192, B: 192, A: 255},
		color.RGBA{R: 192, G: 192, B: 0, A: 255},
		color.RGBA{R: 0, G: 192, B: 192, A: 255},
		color.RGBA{R: 0, G: 192, B: 0, A: 255},
		color.RGBA{R: 192, G: 0, B: 192, A: 255},
		color.RGBA{R: 192, G: 0, B: 0, A: 255},
		color.RGBA{R: 0, G: 0, B: 192, A: 255},
	}

	barWidth := imgWidth / len(colors)

	for idx, clr := range colors {
		for yPos := range imgHeight {
			for xPos := idx * barWidth; xPos < (idx+1)*barWidth && xPos < imgWidth; xPos++ {
				img.Set(xPos, yPos, clr)
			}
		}
	}
}

// WritePNG writes an RGBA image to a PNG file at the given path.
func (r *ImageRenderer) WritePNG(path string, img *image.RGBA) error {
	file, err := os.Create(path) //nolint:gosec // path is provided by the user via CLI flags
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	encodeErr := png.Encode(file, img)
	if encodeErr != nil {
		return fmt.Errorf("PNG encoding failed: %w", encodeErr)
	}

	return nil
}

// ParseColor converts a color name or hex string to a color.Color value.
func (r *ImageRenderer) ParseColor(name string) color.Color {
	switch name {
	case "black":
		return color.Black
	case "white":
		return color.White
	case "red":
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	case "green":
		return color.RGBA{R: 0, G: 255, B: 0, A: 255}
	case "blue":
		return color.RGBA{R: 0, G: 0, B: 255, A: 255}
	case "yellow":
		return color.RGBA{R: 255, G: 255, B: 0, A: 255}
	case "cyan":
		return color.RGBA{R: 0, G: 255, B: 255, A: 255}
	case "magenta":
		return color.RGBA{R: 255, G: 0, B: 255, A: 255}
	case "gray", "grey":
		return color.RGBA{R: 128, G: 128, B: 128, A: 255}
	default:
		var red, green, blue uint8

		_, scanErr := fmt.Sscanf(name, "#%02x%02x%02x", &red, &green, &blue)
		if scanErr == nil {
			return color.RGBA{R: red, G: green, B: blue, A: 255}
		}

		return color.Black
	}
}

// ContrastColor returns black or white depending on the luminance of the background color.
func (r *ImageRenderer) ContrastColor(bg color.Color) color.Color {
	rv, gv, bv, _ := bg.RGBA()
	luminance := (0.299*float64(rv) + 0.587*float64(gv) + 0.114*float64(bv)) / maxColorComponentF

	if luminance > luminanceThreshold {
		return color.Black
	}

	return color.White
}

// FormatTimecode converts a frame number and FPS into HH:MM:SS:FF timecode format.
func (r *ImageRenderer) FormatTimecode(frameNum, fps int) string {
	totalSeconds := frameNum / fps
	frames := frameNum % fps
	hours := totalSeconds / secondsPerHour
	minutes := (totalSeconds % secondsPerHour) / secondsPerMinute
	seconds := totalSeconds % secondsPerMinute

	return fmt.Sprintf("%02d:%02d:%02d:%02d", hours, minutes, seconds, frames)
}

func measureText(text string, face font.Face) int {
	advance := fixed.Int26_6(0)

	for _, char := range text {
		glyphAdvance, ok := face.GlyphAdvance(char)
		if ok {
			advance += glyphAdvance
		}
	}

	return advance.Ceil()
}

func splitLines(text string, maxWidth int, face font.Face) []string {
	if measureText(text, face) <= maxWidth {
		return []string{text}
	}

	var lines []string

	current := ""

	for _, char := range text {
		next := current + string(char)

		if measureText(next, face) > maxWidth {
			if current != "" {
				lines = append(lines, current)
			}

			current = string(char)
		} else {
			current = next
		}
	}

	if current != "" {
		lines = append(lines, current)
	}

	return lines
}
