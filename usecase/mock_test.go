package usecase_test

import (
	"image"
	"image/color"
	"io"
	"sync"

	"github.com/junara/encfixture/domain"
)

// mockFFmpeg records calls. Its mutex keeps the batch use case's concurrent job
// execution race-free when several jobs share one mock.
type mockFFmpeg struct {
	mu                 sync.Mutex
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
	m.mu.Lock()
	m.runCalled = true
	m.runArgs = args
	m.mu.Unlock()

	return m.runErr
}

func (m *mockFFmpeg) RunWithStdin(stdin io.Reader, args ...string) error {
	m.mu.Lock()
	m.runWithStdinCalled = true
	m.runWithStdinArgs = args
	m.mu.Unlock()

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
	mu                    sync.Mutex
	solidImageCalled      bool
	solidImageColors      []color.Color
	drawTextAtCalls       []drawTextAtCall
	drawTestPatternCalled bool
	drawGradientCalls     int
	drawMovingBoxCalls    int
	writeImageCalled      bool
	writeImagePath        string
	writeImageQuality     int
	writeImageErr         error
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

func (m *mockRenderer) SolidImage(width, height int, clr color.Color) *image.RGBA {
	m.mu.Lock()
	m.solidImageCalled = true
	m.solidImageColors = append(m.solidImageColors, clr)
	m.mu.Unlock()

	return image.NewRGBA(image.Rect(0, 0, width, height))
}

func (m *mockRenderer) DrawScaledText(_ *image.RGBA, _ string, _ color.Color, _ int) {}

func (m *mockRenderer) DrawScaledTextAt(_ *image.RGBA, text string, _ color.Color, scale int, pos domain.TextPosition) {
	m.mu.Lock()
	m.drawTextAtCalls = append(m.drawTextAtCalls, drawTextAtCall{text: text, pos: pos, scale: scale})
	m.mu.Unlock()
}

func (m *mockRenderer) DrawTestPattern(_ *image.RGBA) {
	m.mu.Lock()
	m.drawTestPatternCalled = true
	m.mu.Unlock()
}

func (m *mockRenderer) DrawScrollingGradient(_ *image.RGBA, _, _ int) {
	m.mu.Lock()
	m.drawGradientCalls++
	m.mu.Unlock()
}

func (m *mockRenderer) DrawMovingBox(_ *image.RGBA, _, _ int) {
	m.mu.Lock()
	m.drawMovingBoxCalls++
	m.mu.Unlock()
}

func (m *mockRenderer) WriteImage(path string, _ *image.RGBA, quality int) error {
	m.mu.Lock()
	m.writeImageCalled = true
	m.writeImagePath = path
	m.writeImageQuality = quality
	m.mu.Unlock()

	return m.writeImageErr
}

func (m *mockRenderer) ParseColor(_ string) color.Color {
	return m.parseColorResult
}

func (m *mockRenderer) ContrastColor(_ color.Color) color.Color {
	return m.contrastColorResult
}
