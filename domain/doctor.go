package domain

// ToolStatus describes whether an external tool is installed, and where.
type ToolStatus struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
	Version   string `json:"version,omitempty"`
	Path      string `json:"path,omitempty"`
}

// EncoderStatus reports whether a selectable codec's ffmpeg encoder is present
// in the local ffmpeg build.
type EncoderStatus struct {
	Codec     string `json:"codec"`
	Encoder   string `json:"encoder"`
	Available bool   `json:"available"`
}

// DoctorReport aggregates the environment checks encfixture depends on.
type DoctorReport struct {
	FFmpeg        ToolStatus      `json:"ffmpeg"`
	FFprobe       ToolStatus      `json:"ffprobe"`
	VideoEncoders []EncoderStatus `json:"videoEncoders"`
	AudioEncoders []EncoderStatus `json:"audioEncoders"`
}

// Healthy reports whether the core tools are present. Missing optional
// encoders do not make the environment unhealthy; commands that need them
// fail individually with a pointer back to this report.
func (r DoctorReport) Healthy() bool {
	return r.FFmpeg.Available && r.FFprobe.Available
}
