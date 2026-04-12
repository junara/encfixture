package usecase_test

import (
	"image"
	"image/color"
	"io"

	"github.com/junara/encfixture/domain"
)

type mockFFmpeg struct {
	runCalled          bool
	runArgs            []string
	runErr             error
	runWithStdinCalled bool
	runWithStdinArgs   []string
	runWithStdinErr    error
	checkErr           error
}

func (m *mockFFmpeg) CheckAvailable() error {
	return m.checkErr
}

func (m *mockFFmpeg) Run(args ...string) error {
	m.runCalled = true
	m.runArgs = args

	return m.runErr
}

func (m *mockFFmpeg) RunWithStdin(stdin io.Reader, args ...string) error {
	m.runWithStdinCalled = true
	m.runWithStdinArgs = args

	// drain stdin to prevent pipe deadlock
	go func() {
		buf := make([]byte, 4096)
		for {
			_, err := stdin.Read(buf)
			if err != nil {
				return
			}
		}
	}()

	return m.runWithStdinErr
}

type mockRenderer struct {
	solidImageCalled      bool
	drawTextAtCalls       []drawTextAtCall
	drawTestPatternCalled bool
	writePNGCalled        bool
	writePNGPath          string
	writePNGErr           error
	parseColorResult      color.Color
	contrastColorResult   color.Color
}

type drawTextAtCall struct {
	text  string
	pos   domain.TextPosition
	scale int
}

func newMockRenderer() *mockRenderer {
	return &mockRenderer{
		parseColorResult:    color.Black,
		contrastColorResult: color.White,
	}
}

func (m *mockRenderer) SolidImage(width, height int, _ color.Color) *image.RGBA {
	m.solidImageCalled = true

	return image.NewRGBA(image.Rect(0, 0, width, height))
}

func (m *mockRenderer) DrawScaledText(_ *image.RGBA, _ string, _ color.Color, _ int) {}

func (m *mockRenderer) DrawScaledTextAt(_ *image.RGBA, text string, _ color.Color, scale int, pos domain.TextPosition) {
	m.drawTextAtCalls = append(m.drawTextAtCalls, drawTextAtCall{text: text, pos: pos, scale: scale})
}

func (m *mockRenderer) DrawTestPattern(_ *image.RGBA) {
	m.drawTestPatternCalled = true
}

func (m *mockRenderer) WritePNG(path string, _ *image.RGBA) error {
	m.writePNGCalled = true
	m.writePNGPath = path

	return m.writePNGErr
}

func (m *mockRenderer) ParseColor(_ string) color.Color {
	return m.parseColorResult
}

func (m *mockRenderer) ContrastColor(_ color.Color) color.Color {
	return m.contrastColorResult
}
