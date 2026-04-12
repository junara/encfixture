// Package domain defines the core entities and value objects for media generation.
package domain

// TextPosition represents where text is drawn on the image.
type TextPosition int

// TextPosition constants define the available text placement positions.
const (
	PositionTopLeft     TextPosition = iota // top-left corner
	PositionTopRight                        // top-right corner
	PositionCenter                          // center of the image
	PositionBottomLeft                      // bottom-left corner
	PositionBottomRight                     // bottom-right corner
)

// OverlayKeyword constants define reserved keywords for overlay content.
const (
	KeywordFrame    = "frame"
	KeywordTimecode = "timecode"
	KeywordFilename = "filename"
)

// BackgroundTest is the background type for color bar test patterns.
const BackgroundTest = "test"

// AudioType represents the type of audio to generate.
type AudioType string

// AudioType constants define the available audio generation modes.
const (
	AudioSilence AudioType = "silence"
	AudioSine    AudioType = "sine"
	AudioNoise   AudioType = "noise"
	AudioTone    AudioType = "tone"
)

// Overlay holds the content for each display position.
// Values can be a keyword (frame, timecode, filename) or arbitrary text.
// An empty string means nothing is displayed at that position.
type Overlay struct {
	TopLeft     string
	TopRight    string
	Center      string
	BottomLeft  string
	BottomRight string
}

// HasContent returns true if any position has content to display.
func (o Overlay) HasContent() bool {
	return o.TopLeft != "" || o.TopRight != "" || o.Center != "" ||
		o.BottomLeft != "" || o.BottomRight != ""
}

// HasDynamicContent returns true if any position uses frame-dependent keywords.
func (o Overlay) HasDynamicContent() bool {
	for _, content := range o.All() {
		if content == KeywordFrame || content == KeywordTimecode {
			return true
		}
	}

	return false
}

// All returns all non-empty overlay contents as position-content pairs.
func (o Overlay) All() []string {
	var result []string

	if o.TopLeft != "" {
		result = append(result, o.TopLeft)
	}

	if o.TopRight != "" {
		result = append(result, o.TopRight)
	}

	if o.Center != "" {
		result = append(result, o.Center)
	}

	if o.BottomLeft != "" {
		result = append(result, o.BottomLeft)
	}

	if o.BottomRight != "" {
		result = append(result, o.BottomRight)
	}

	return result
}

// Entries returns all non-empty overlay entries with their positions.
func (o Overlay) Entries() []OverlayEntry {
	var entries []OverlayEntry

	if o.TopLeft != "" {
		entries = append(entries, OverlayEntry{Position: PositionTopLeft, Content: o.TopLeft})
	}

	if o.TopRight != "" {
		entries = append(entries, OverlayEntry{Position: PositionTopRight, Content: o.TopRight})
	}

	if o.Center != "" {
		entries = append(entries, OverlayEntry{Position: PositionCenter, Content: o.Center})
	}

	if o.BottomLeft != "" {
		entries = append(entries, OverlayEntry{Position: PositionBottomLeft, Content: o.BottomLeft})
	}

	if o.BottomRight != "" {
		entries = append(entries, OverlayEntry{Position: PositionBottomRight, Content: o.BottomRight})
	}

	return entries
}

// OverlayEntry pairs a text position with its content.
type OverlayEntry struct {
	Position TextPosition
	Content  string
}

// ImageConfig holds the configuration for image generation.
type ImageConfig struct {
	Width      int
	Height     int
	Background string // "solid" or "test"
	Color      string
	Overlay    Overlay
	Scale      int
	Output     string
}

// VideoConfig holds the configuration for video generation.
type VideoConfig struct {
	Width      int
	Height     int
	FPS        int
	Duration   string
	Background string // "solid" or "test"
	Color      string
	Overlay    Overlay
	Scale      int
	Output     string
	Audio      AudioType
	SampleRate int
	Channels   int
	Frequency  float64
}

// AudioConfig holds the configuration for audio generation.
type AudioConfig struct {
	Type       AudioType
	Duration   string
	SampleRate int
	Channels   int
	Frequency  float64
	Output     string
}
