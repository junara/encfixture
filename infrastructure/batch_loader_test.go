package infrastructure_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/infrastructure"
)

func writeJSON(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "batch.json")

	err := os.WriteFile(path, []byte(content), 0o600)
	if err != nil {
		t.Fatalf("write file: %v", err)
	}

	return path
}

func TestLoadBatch_DefaultsAppliedAndOverridden(t *testing.T) {
	t.Parallel()

	path := writeJSON(t, `{
  "defaults": {"width": 1280, "height": 720, "color": "blue"},
  "jobs": [
    {"type": "image", "output": "a.png"},
    {"type": "image", "output": "b.png", "color": "red"}
  ]
}`)

	batch, err := infrastructure.LoadBatch(path)
	if err != nil {
		t.Fatalf("LoadBatch: %v", err)
	}

	if len(batch.Jobs) != 2 {
		t.Fatalf("got %d jobs, want 2", len(batch.Jobs))
	}

	first := batch.Jobs[0].Image
	if first == nil {
		t.Fatal("job[0].Image is nil")
	}

	if first.Width != 1280 || first.Height != 720 || first.Color != "blue" {
		t.Errorf("defaults not applied: %+v", first)
	}

	second := batch.Jobs[1].Image
	if second == nil {
		t.Fatal("job[1].Image is nil")
	}

	if second.Color != "red" {
		t.Errorf("job override lost: color=%q", second.Color)
	}
}

func TestLoadBatch_AllTypes(t *testing.T) {
	t.Parallel()

	path := writeJSON(t, `{
  "jobs": [
    {"type": "image", "output": "i.png"},
    {"type": "video", "output": "v.mp4", "fps": 60},
    {"type": "audio", "output": "a.wav", "audio": "sine", "frequency": 1000}
  ]
}`)

	batch, err := infrastructure.LoadBatch(path)
	if err != nil {
		t.Fatalf("LoadBatch: %v", err)
	}

	if batch.Jobs[0].Type != domain.JobTypeImage || batch.Jobs[0].Image == nil {
		t.Error("image job not populated")
	}

	if batch.Jobs[1].Type != domain.JobTypeVideo || batch.Jobs[1].Video == nil || batch.Jobs[1].Video.FPS != 60 {
		t.Error("video job not populated correctly")
	}

	aud := batch.Jobs[2].Audio
	if aud == nil || aud.Type != domain.AudioSine || aud.Frequency != 1000 {
		t.Errorf("audio job not populated: %+v", aud)
	}
}

func TestLoadBatch_RejectsUnknownField(t *testing.T) {
	t.Parallel()

	path := writeJSON(t, `{"jobs":[{"type":"image","output":"x.png","bogus":1}]}`)

	_, err := infrastructure.LoadBatch(path)
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestLoadBatch_MissingOutput(t *testing.T) {
	t.Parallel()

	path := writeJSON(t, `{"jobs":[{"type":"image"}]}`)

	_, err := infrastructure.LoadBatch(path)
	if !errors.Is(err, infrastructure.ErrBatchInvalid) {
		t.Fatalf("expected ErrBatchInvalid, got %v", err)
	}
}

func TestLoadBatch_UnknownType(t *testing.T) {
	t.Parallel()

	path := writeJSON(t, `{"jobs":[{"type":"movie","output":"x.mov"}]}`)

	_, err := infrastructure.LoadBatch(path)
	if !errors.Is(err, infrastructure.ErrBatchInvalid) {
		t.Fatalf("expected ErrBatchInvalid, got %v", err)
	}
}
