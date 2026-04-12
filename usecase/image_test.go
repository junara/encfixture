package usecase_test

import (
	"errors"
	"testing"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/usecase"
)

func TestImageUseCase_Generate_Solid(t *testing.T) {
	t.Parallel()

	renderer := newMockRenderer()
	uc := usecase.NewImageUseCase(renderer)

	cfg := domain.ImageConfig{
		Width:      640,
		Height:     480,
		Background: "solid",
		Color:      "black",
		Output:     "test.png",
		Scale:      4,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !renderer.solidImageCalled {
		t.Error("SolidImage was not called")
	}

	if !renderer.writePNGCalled {
		t.Error("WritePNG was not called")
	}

	if renderer.writePNGPath != "test.png" {
		t.Errorf("WritePNG path = %q, want %q", renderer.writePNGPath, "test.png")
	}
}

func TestImageUseCase_Generate_TestPattern(t *testing.T) {
	t.Parallel()

	renderer := newMockRenderer()
	uc := usecase.NewImageUseCase(renderer)

	cfg := domain.ImageConfig{
		Width:      640,
		Height:     480,
		Background: "test",
		Color:      "black",
		Output:     "test.png",
		Scale:      4,
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !renderer.drawTestPatternCalled {
		t.Error("DrawTestPattern was not called")
	}
}

func TestImageUseCase_Generate_AllOverlays(t *testing.T) {
	t.Parallel()

	renderer := newMockRenderer()
	uc := usecase.NewImageUseCase(renderer)

	cfg := domain.ImageConfig{
		Width:      640,
		Height:     480,
		Background: "solid",
		Color:      "black",
		Overlay: domain.Overlay{
			TopLeft:     "frame",
			TopRight:    "timecode",
			Center:      "HELLO",
			BottomLeft:  "filename",
			BottomRight: "ID-001",
		},
		Scale:  4,
		Output: "sample.png",
	}

	err := uc.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if len(renderer.drawTextAtCalls) != 5 {
		t.Fatalf("DrawScaledTextAt called %d times, want 5", len(renderer.drawTextAtCalls))
	}

	expected := []struct {
		text string
		pos  domain.TextPosition
	}{
		{"0", domain.PositionTopLeft},
		{"00:00:00:00", domain.PositionTopRight},
		{"HELLO", domain.PositionCenter},
		{"sample.png", domain.PositionBottomLeft},
		{"ID-001", domain.PositionBottomRight},
	}

	for i, exp := range expected {
		call := renderer.drawTextAtCalls[i]
		if call.text != exp.text {
			t.Errorf("call[%d].text = %q, want %q", i, call.text, exp.text)
		}

		if call.pos != exp.pos {
			t.Errorf("call[%d].pos = %v, want %v", i, call.pos, exp.pos)
		}
	}
}

func TestImageUseCase_Generate_WritePNGError(t *testing.T) {
	t.Parallel()

	renderer := newMockRenderer()
	renderer.writePNGErr = errors.New("disk full")
	uc := usecase.NewImageUseCase(renderer)

	cfg := domain.ImageConfig{
		Width:      640,
		Height:     480,
		Background: "solid",
		Color:      "black",
		Output:     "test.png",
		Scale:      4,
	}

	err := uc.Generate(cfg)
	if err == nil {
		t.Fatal("Generate() should return error")
	}
}
