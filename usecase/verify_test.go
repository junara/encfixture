package usecase_test

import (
	"errors"
	"testing"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/usecase"
)

type mockProber struct {
	checkErr  error
	info      domain.MediaInfo
	probeErr  error
	probePath string
}

func (m *mockProber) CheckAvailable() error {
	return m.checkErr
}

func (m *mockProber) Probe(path string) (domain.MediaInfo, error) {
	m.probePath = path

	return m.info, m.probeErr
}

func TestVerifyUseCase_Verify_Success(t *testing.T) {
	t.Parallel()

	prober := &mockProber{
		info: domain.MediaInfo{
			Format: domain.FormatInfo{FormatName: "mp4", Duration: "2", Size: "100", BitRate: "50"},
			Streams: []domain.StreamInfo{
				{Index: 0, Type: "video", Codec: "h264", Width: 1920, Height: 1080, FPS: "30", PixFmt: "yuv420p"},
			},
		},
	}
	uc := usecase.NewVerifyUseCase(prober)

	info, err := uc.Verify("test.mp4")
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	if prober.probePath != "test.mp4" {
		t.Errorf("probed path = %q, want test.mp4", prober.probePath)
	}

	if info.Format.FormatName != "mp4" {
		t.Errorf("format = %q, want mp4", info.Format.FormatName)
	}

	if len(info.Streams) != 1 || info.Streams[0].Codec != "h264" {
		t.Errorf("streams not passed through: %+v", info.Streams)
	}
}

func TestVerifyUseCase_Verify_Unavailable(t *testing.T) {
	t.Parallel()

	prober := &mockProber{checkErr: errors.New("no ffprobe")}
	uc := usecase.NewVerifyUseCase(prober)

	_, err := uc.Verify("test.mp4")
	if err == nil {
		t.Fatal("Verify() should error when ffprobe is unavailable")
	}

	if prober.probePath != "" {
		t.Error("Probe should not be called when the availability check fails")
	}
}

func TestVerifyUseCase_Verify_ProbeError(t *testing.T) {
	t.Parallel()

	prober := &mockProber{probeErr: errors.New("corrupt file")}
	uc := usecase.NewVerifyUseCase(prober)

	_, err := uc.Verify("bad.mp4")
	if err == nil {
		t.Fatal("Verify() should propagate a probe error")
	}
}
