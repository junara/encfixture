package infrastructure_test

import (
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/infrastructure"
)

func TestImageRenderer_SolidImage(t *testing.T) {
	t.Parallel()

	r := infrastructure.NewImageRenderer()
	img := r.SolidImage(100, 50, color.White)

	bounds := img.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 50 {
		t.Errorf("SolidImage size = %dx%d, want 100x50", bounds.Dx(), bounds.Dy())
	}

	pixel := img.At(0, 0)
	rv, gv, bv, _ := pixel.RGBA()

	if rv != 0xffff || gv != 0xffff || bv != 0xffff {
		t.Errorf("pixel color = (%d,%d,%d), want white", rv, gv, bv)
	}
}

func TestImageRenderer_ParseColor(t *testing.T) {
	t.Parallel()

	r := infrastructure.NewImageRenderer()

	tests := []struct {
		name  string
		input string
		wantR uint8
		wantG uint8
		wantB uint8
	}{
		{"black", "black", 0, 0, 0},
		{"white", "white", 255, 255, 255},
		{"red", "red", 255, 0, 0},
		{"green", "green", 0, 255, 0},
		{"blue", "blue", 0, 0, 255},
		{"yellow", "yellow", 255, 255, 0},
		{"cyan", "cyan", 0, 255, 255},
		{"magenta", "magenta", 255, 0, 255},
		{"gray", "gray", 128, 128, 128},
		{"grey", "grey", 128, 128, 128},
		{"hex orange", "#ff6600", 255, 102, 0},
		{"hex white", "#ffffff", 255, 255, 255},
		{"unknown defaults to black", "unknown", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := r.ParseColor(tt.input)
			rv, gv, bv, _ := c.RGBA()

			gotR := uint8(rv >> 8)
			gotG := uint8(gv >> 8)
			gotB := uint8(bv >> 8)

			if gotR != tt.wantR || gotG != tt.wantG || gotB != tt.wantB {
				t.Errorf("ParseColor(%q) = (%d,%d,%d), want (%d,%d,%d)",
					tt.input, gotR, gotG, gotB, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

func TestImageRenderer_ContrastColor(t *testing.T) {
	t.Parallel()

	r := infrastructure.NewImageRenderer()

	tests := []struct {
		name string
		bg   color.Color
		want color.Color
	}{
		{"black bg gives white text", color.Black, color.White},
		{"white bg gives black text", color.White, color.Black},
		{"dark blue gives white text", color.RGBA{R: 0, G: 0, B: 128, A: 255}, color.White},
		{"bright yellow gives black text", color.RGBA{R: 255, G: 255, B: 0, A: 255}, color.Black},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := r.ContrastColor(tt.bg)

			gotR, gotG, gotB, _ := got.RGBA()
			wantR, wantG, wantB, _ := tt.want.RGBA()

			if gotR != wantR || gotG != wantG || gotB != wantB {
				t.Errorf("ContrastColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageRenderer_FormatTimecode(t *testing.T) {
	t.Parallel()

	r := infrastructure.NewImageRenderer()

	tests := []struct {
		name     string
		frame    int
		fps      int
		expected string
	}{
		{"zero", 0, 30, "00:00:00:00"},
		{"one frame", 1, 30, "00:00:00:01"},
		{"one second", 30, 30, "00:00:01:00"},
		{"one minute", 1800, 30, "00:01:00:00"},
		{"one hour", 108000, 30, "01:00:00:00"},
		{"mixed", 3661*24 + 12, 24, "01:01:01:12"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := r.FormatTimecode(tt.frame, tt.fps)
			if got != tt.expected {
				t.Errorf("FormatTimecode(%d, %d) = %q, want %q", tt.frame, tt.fps, got, tt.expected)
			}
		})
	}
}

func TestImageRenderer_WritePNG(t *testing.T) {
	t.Parallel()

	r := infrastructure.NewImageRenderer()
	img := r.SolidImage(10, 10, color.Black)

	dir := t.TempDir()
	path := filepath.Join(dir, "test.png")

	err := r.WritePNG(path, img)
	if err != nil {
		t.Fatalf("WritePNG() error = %v", err)
	}

	info, statErr := os.Stat(path)
	if statErr != nil {
		t.Fatalf("output file not found: %v", statErr)
	}

	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

func TestImageRenderer_DrawTestPattern(t *testing.T) {
	t.Parallel()

	r := infrastructure.NewImageRenderer()
	img := r.SolidImage(140, 100, color.Black)

	r.DrawTestPattern(img)

	// first bar should be gray (192, 192, 192)
	pixel := img.At(0, 0)
	rv, gv, bv, _ := pixel.RGBA()

	gotR := uint8(rv >> 8)
	if gotR != 192 {
		t.Errorf("first bar pixel R = %d, want 192", gotR)
	}

	_ = gv
	_ = bv
}

func TestImageRenderer_DrawScaledTextAt(t *testing.T) {
	t.Parallel()

	r := infrastructure.NewImageRenderer()
	img := r.SolidImage(200, 200, color.Black)

	positions := []domain.TextPosition{
		domain.PositionTopLeft,
		domain.PositionTopRight,
		domain.PositionCenter,
		domain.PositionBottomLeft,
		domain.PositionBottomRight,
	}

	for _, pos := range positions {
		r.DrawScaledTextAt(img, "X", color.White, 1, pos)
	}

	// verify at least some non-black pixels exist (text was drawn)
	hasWhite := false

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			rv, _, _, _ := img.At(x, y).RGBA()

			if a > 0 && rv > 0 {
				hasWhite = true

				break
			}
		}

		if hasWhite {
			break
		}
	}

	if !hasWhite {
		t.Error("DrawScaledTextAt did not draw any visible pixels")
	}
}
