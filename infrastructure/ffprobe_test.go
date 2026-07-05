package infrastructure_test

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/infrastructure"
)

func requireBinary(t *testing.T, name string) {
	t.Helper()

	_, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("%s not installed", name)
	}
}

func TestFFprobe_Probe_Video(t *testing.T) {
	t.Parallel()
	requireBinary(t, "ffmpeg")
	requireBinary(t, "ffprobe")

	out := filepath.Join(t.TempDir(), "v.mp4")
	cmd := exec.CommandContext(context.Background(), "ffmpeg", "-y",
		"-f", "lavfi", "-i", "color=c=black:s=64x48:d=1:r=30", "-pix_fmt", "yuv420p", out)

	runErr := cmd.Run()
	if runErr != nil {
		t.Fatalf("failed to create fixture: %v", runErr)
	}

	info, err := infrastructure.NewFFprobe().Probe(out)
	if err != nil {
		t.Fatalf("Probe() error = %v", err)
	}

	if len(info.Streams) == 0 {
		t.Fatal("no streams reported")
	}

	video := info.Streams[0]
	if video.Type != domain.StreamTypeVideo {
		t.Errorf("stream type = %q, want video", video.Type)
	}

	if video.Width != 64 || video.Height != 48 {
		t.Errorf("resolution = %dx%d, want 64x48", video.Width, video.Height)
	}

	if video.FPS != "30" {
		t.Errorf("fps = %q, want 30", video.FPS)
	}

	if video.Codec == "" {
		t.Error("codec is empty")
	}
}

func TestFFprobe_Probe_MissingFile(t *testing.T) {
	t.Parallel()
	requireBinary(t, "ffprobe")

	_, err := infrastructure.NewFFprobe().Probe(filepath.Join(t.TempDir(), "nope.mp4"))
	if err == nil {
		t.Fatal("Probe() should error on a missing file")
	}
}

func TestFFprobe_CheckAvailable(t *testing.T) {
	t.Parallel()
	requireBinary(t, "ffprobe")

	err := infrastructure.NewFFprobe().CheckAvailable()
	if err != nil {
		t.Errorf("CheckAvailable() error = %v", err)
	}
}
