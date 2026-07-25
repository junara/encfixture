package usecase_test

import (
	"errors"
	"testing"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/usecase"
)

type mockInspector struct {
	statuses    map[string]domain.ToolStatus
	encoders    map[string]bool
	encodersErr error
}

func (m *mockInspector) ToolStatus(name string) domain.ToolStatus {
	return m.statuses[name]
}

func (m *mockInspector) Encoders() (map[string]bool, error) {
	return m.encoders, m.encodersErr
}

func TestDoctorUseCase_Report_AllAvailable(t *testing.T) {
	t.Parallel()

	inspector := &mockInspector{
		statuses: map[string]domain.ToolStatus{
			"ffmpeg":  {Name: "ffmpeg", Available: true, Version: "7.1", Path: "/usr/bin/ffmpeg"},
			"ffprobe": {Name: "ffprobe", Available: true, Version: "7.1", Path: "/usr/bin/ffprobe"},
		},
		encoders: map[string]bool{
			"libx264": true, "libx265": true, "libvpx-vp9": true,
			"libaom-av1": true, "prores_ks": true, "aac": true, "libopus": true,
		},
	}

	report, err := usecase.NewDoctorUseCase(inspector).Report()
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	if !report.Healthy() {
		t.Error("Healthy() = false, want true")
	}

	if len(report.VideoEncoders) != 5 {
		t.Fatalf("len(VideoEncoders) = %d, want 5", len(report.VideoEncoders))
	}

	for _, enc := range report.VideoEncoders {
		if !enc.Available {
			t.Errorf("encoder %s reported unavailable", enc.Encoder)
		}
	}

	for _, enc := range report.AudioEncoders {
		if !enc.Available {
			t.Errorf("audio encoder %s reported unavailable", enc.Encoder)
		}
	}
}

func TestDoctorUseCase_Report_MissingEncoder(t *testing.T) {
	t.Parallel()

	inspector := &mockInspector{
		statuses: map[string]domain.ToolStatus{
			"ffmpeg":  {Name: "ffmpeg", Available: true},
			"ffprobe": {Name: "ffprobe", Available: true},
		},
		encoders: map[string]bool{"libx264": true},
	}

	report, err := usecase.NewDoctorUseCase(inspector).Report()
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	if !report.Healthy() {
		t.Error("Healthy() = false, want true (missing encoders are not fatal)")
	}

	got := map[string]bool{}
	for _, enc := range report.VideoEncoders {
		got[enc.Codec] = enc.Available
	}

	if !got["h264"] {
		t.Error("h264 should be available")
	}

	if got["hevc"] || got["av1"] {
		t.Error("hevc/av1 should be unavailable")
	}
}

func TestDoctorUseCase_Report_FFmpegMissing(t *testing.T) {
	t.Parallel()

	inspector := &mockInspector{
		statuses: map[string]domain.ToolStatus{
			"ffmpeg":  {Name: "ffmpeg", Available: false},
			"ffprobe": {Name: "ffprobe", Available: true},
		},
		// Encoders must not be consulted when ffmpeg is absent.
		encodersErr: errors.New("should not be called"),
	}

	report, err := usecase.NewDoctorUseCase(inspector).Report()
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	if report.Healthy() {
		t.Error("Healthy() = true, want false")
	}

	for _, enc := range report.VideoEncoders {
		if enc.Available {
			t.Errorf("encoder %s should be unavailable without ffmpeg", enc.Encoder)
		}
	}
}

func TestDoctorUseCase_Report_EncodersError(t *testing.T) {
	t.Parallel()

	inspector := &mockInspector{
		statuses: map[string]domain.ToolStatus{
			"ffmpeg":  {Name: "ffmpeg", Available: true},
			"ffprobe": {Name: "ffprobe", Available: true},
		},
		encodersErr: errors.New("boom"),
	}

	_, err := usecase.NewDoctorUseCase(inspector).Report()
	if err == nil {
		t.Fatal("Report() error = nil, want error")
	}
}
