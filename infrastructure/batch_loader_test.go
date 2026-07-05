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

func TestLoadBatch_EncodeFields(t *testing.T) {
	t.Parallel()

	path := writeJSON(t, `{
  "defaults": {"codec": "h264", "crf": 23},
  "jobs": [
    {"type": "video", "output": "v.mp4", "bitrate": "5M", "pixFmt": "yuv444p"},
    {"type": "video", "output": "w.mp4", "codec": "vp9"},
    {"type": "image", "output": "i.jpg", "quality": 75}
  ]
}`)

	batch, err := infrastructure.LoadBatch(path)
	if err != nil {
		t.Fatalf("LoadBatch: %v", err)
	}

	first := batch.Jobs[0].Video
	if first.Codec != domain.CodecH264 || first.CRF != "23" || first.Bitrate != "5M" || first.PixFmt != "yuv444p" {
		t.Errorf("encode fields not applied: %+v", first)
	}

	if batch.Jobs[1].Video.Codec != domain.CodecVP9 {
		t.Errorf("codec override lost: %q", batch.Jobs[1].Video.Codec)
	}

	if batch.Jobs[2].Image.Quality != 75 {
		t.Errorf("quality = %d, want 75", batch.Jobs[2].Image.Quality)
	}
}

func TestLoadBatch_SyncFields(t *testing.T) {
	t.Parallel()

	path := writeJSON(t, `{
  "jobs": [
    {"type": "video", "output": "s.mp4", "sync": true, "syncInterval": 0.5},
    {"type": "video", "output": "n.mp4"}
  ]
}`)

	batch, err := infrastructure.LoadBatch(path)
	if err != nil {
		t.Fatalf("LoadBatch: %v", err)
	}

	first := batch.Jobs[0].Video
	if !first.Sync || first.SyncInterval != 0.5 {
		t.Errorf("sync fields not applied: sync=%v interval=%v", first.Sync, first.SyncInterval)
	}

	if second := batch.Jobs[1].Video; second.Sync {
		t.Errorf("sync should default to false, got %v", second.Sync)
	}
}

func TestLoadBatch_QualityDefault(t *testing.T) {
	t.Parallel()

	path := writeJSON(t, `{"jobs":[{"type":"image","output":"i.jpg"}]}`)

	batch, err := infrastructure.LoadBatch(path)
	if err != nil {
		t.Fatalf("LoadBatch: %v", err)
	}

	if batch.Jobs[0].Image.Quality != 90 {
		t.Errorf("quality default = %d, want 90", batch.Jobs[0].Image.Quality)
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
