// Package domain defines the core entities and value objects for media generation.
package domain

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

// ErrInvalidDuration indicates a duration that is not a positive number of seconds.
var ErrInvalidDuration = errors.New("invalid duration")

// ErrInvalidBitrate indicates a malformed bitrate value.
var ErrInvalidBitrate = errors.New("invalid bitrate")

// ValidateDuration checks that s parses as a positive number of seconds, so a
// typo fails fast with a clear message instead of deep inside ffmpeg.
func ValidateDuration(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 {
		return fmt.Errorf("%w: %q (must be a positive number of seconds)", ErrInvalidDuration, s)
	}

	return nil
}

// ValidateBitrate checks that s looks like an ffmpeg bitrate such as "800k" or
// "5M". An empty value is valid and means the encoder default.
func ValidateBitrate(s string) error {
	if s == "" {
		return nil
	}

	num := s
	switch s[len(s)-1] {
	case 'k', 'K', 'm', 'M':
		num = s[:len(s)-1]
	}

	v, err := strconv.ParseFloat(num, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 {
		return fmt.Errorf("%w: %q (expected a number with an optional k/M suffix, e.g. 800k or 5M)", ErrInvalidBitrate, s)
	}

	return nil
}

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

// Background type constants define the available background patterns.
const (
	BackgroundSolid    = "solid"    // flat fill of the chosen color
	BackgroundTest     = "test"     // SMPTE-style color bars
	BackgroundGradient = "gradient" // horizontally scrolling color gradient
	BackgroundMoving   = "moving"   // a box translating across the frame
)

// IsValidBackground reports whether s is a supported background type.
func IsValidBackground(s string) bool {
	switch s {
	case BackgroundSolid, BackgroundTest, BackgroundGradient, BackgroundMoving:
		return true
	default:
		return false
	}
}

// VideoCodec represents the video encoder selection.
// An empty value means the container's default codec is used.
type VideoCodec string

// VideoCodec constants define the selectable video codecs.
const (
	CodecH264   VideoCodec = "h264"
	CodecHEVC   VideoCodec = "hevc"
	CodecVP9    VideoCodec = "vp9"
	CodecAV1    VideoCodec = "av1"
	CodecProRes VideoCodec = "prores"
)

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
	Quality    int // JPEG quality (1-100); ignored for PNG output
}

// VideoConfig holds the configuration for video generation.
type VideoConfig struct {
	Width        int
	Height       int
	FPS          int
	Duration     string
	Background   string // "solid" or "test"
	Color        string
	Overlay      Overlay
	Scale        int
	Output       string
	Audio        AudioType
	SampleRate   int
	Channels     int
	Frequency    float64
	Codec        VideoCodec // empty means container default
	CRF          string     // constant rate factor; empty means encoder default
	Bitrate      string     // video bitrate (e.g. "5M"); empty means encoder default
	PixFmt       string     // pixel format (e.g. "yuv420p"); empty means codec-dependent default
	Sync         bool       // generate an A/V sync test pattern (periodic beep + visual flash)
	SyncInterval float64    // seconds between sync markers; <=0 means the default
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
