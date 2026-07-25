package usecase

import (
	"fmt"

	"github.com/junara/encfixture/domain"
)

// DoctorUseCase checks the environment encfixture depends on: the ffmpeg and
// ffprobe binaries and the encoders needed by the selectable codecs.
type DoctorUseCase struct {
	inspector Inspector
}

// NewDoctorUseCase creates a new DoctorUseCase with the given inspector.
func NewDoctorUseCase(inspector Inspector) *DoctorUseCase {
	return &DoctorUseCase{inspector: inspector}
}

// selectableCodecs lists every codec accepted by --codec, in display order.
func selectableCodecs() []domain.VideoCodec {
	return []domain.VideoCodec{
		domain.CodecH264,
		domain.CodecHEVC,
		domain.CodecVP9,
		domain.CodecAV1,
		domain.CodecProRes,
	}
}

// audioEncoders lists the audio encoders encfixture relies on: aac is the
// default for mp4/mov output, libopus is used for webm output.
func audioEncoders() []domain.EncoderStatus {
	return []domain.EncoderStatus{
		{Codec: "aac", Encoder: "aac", Available: false},
		{Codec: "opus", Encoder: "libopus", Available: false},
	}
}

// Report inspects the environment and returns the aggregated result. Encoder
// availability is only probed when ffmpeg itself is present; without ffmpeg
// every encoder is reported as unavailable.
func (uc *DoctorUseCase) Report() (domain.DoctorReport, error) {
	var empty domain.DoctorReport

	report := domain.DoctorReport{
		FFmpeg:        uc.inspector.ToolStatus("ffmpeg"),
		FFprobe:       uc.inspector.ToolStatus("ffprobe"),
		VideoEncoders: nil,
		AudioEncoders: nil,
	}

	available := map[string]bool{}

	if report.FFmpeg.Available {
		found, err := uc.inspector.Encoders()
		if err != nil {
			return empty, fmt.Errorf("list ffmpeg encoders: %w", err)
		}

		available = found
	}

	codecs := selectableCodecs()
	report.VideoEncoders = make([]domain.EncoderStatus, 0, len(codecs))

	for _, codec := range codecs {
		encoder, err := encoderName(codec)
		if err != nil {
			return empty, fmt.Errorf("resolve encoder for %s: %w", codec, err)
		}

		report.VideoEncoders = append(report.VideoEncoders, domain.EncoderStatus{
			Codec:     string(codec),
			Encoder:   encoder,
			Available: available[encoder],
		})
	}

	report.AudioEncoders = audioEncoders()
	for i := range report.AudioEncoders {
		report.AudioEncoders[i].Available = available[report.AudioEncoders[i].Encoder]
	}

	return report, nil
}
