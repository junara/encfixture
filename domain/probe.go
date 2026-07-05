package domain

// Stream type constants identify the kind of a media stream.
const (
	StreamTypeVideo = "video"
	StreamTypeAudio = "audio"
)

// MediaInfo describes an existing media file as reported by a prober.
type MediaInfo struct {
	Format  FormatInfo   `json:"format"`
	Streams []StreamInfo `json:"streams"`
}

// FormatInfo holds container-level properties of a media file.
type FormatInfo struct {
	FormatName string `json:"formatName"`
	Duration   string `json:"duration"`
	Size       string `json:"size"`
	BitRate    string `json:"bitRate"`
}

// StreamInfo holds the properties of a single media stream. Fields that do not
// apply to the stream's type are left at their zero value.
type StreamInfo struct {
	Index      int    `json:"index"`
	Type       string `json:"type"`
	Codec      string `json:"codec"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	FPS        string `json:"fps,omitempty"`
	PixFmt     string `json:"pixFmt,omitempty"`
	SampleRate string `json:"sampleRate,omitempty"`
	Channels   int    `json:"channels,omitempty"`
}
