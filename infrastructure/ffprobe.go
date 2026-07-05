package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/junara/encfixture/domain"
)

// ErrFFprobeNotFound indicates ffprobe is not installed or not in PATH.
var ErrFFprobeNotFound = errors.New("ffprobe not found in PATH: please install ffmpeg")

// FFprobe inspects media files by shelling out to ffprobe.
type FFprobe struct{}

// NewFFprobe creates a new FFprobe instance.
func NewFFprobe() *FFprobe {
	return &FFprobe{}
}

// CheckAvailable verifies that ffprobe is available in the system PATH.
func (f *FFprobe) CheckAvailable() error {
	_, err := exec.LookPath("ffprobe")
	if err != nil {
		return fmt.Errorf("%w", ErrFFprobeNotFound)
	}

	return nil
}

// Probe runs ffprobe on the file at path and returns its normalized properties.
func (f *FFprobe) Probe(path string) (domain.MediaInfo, error) {
	var empty domain.MediaInfo

	cmd := exec.CommandContext(context.Background(), "ffprobe", //nolint:gosec // path is user-supplied by design
		"-v", "error", "-print_format", "json", "-show_format", "-show_streams", path)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return empty, fmt.Errorf("ffprobe execution failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	var parsed ffprobeOutput

	unmarshalErr := json.Unmarshal(out, &parsed)
	if unmarshalErr != nil {
		return empty, fmt.Errorf("parse ffprobe output: %w", unmarshalErr)
	}

	return toMediaInfo(parsed), nil
}

type ffprobeOutput struct {
	Streams []ffprobeStream `json:"streams"`
	Format  ffprobeFormat   `json:"format"`
}

type ffprobeStream struct {
	Index      int    `json:"index"`
	CodecType  string `json:"codec_type"`
	CodecName  string `json:"codec_name"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	RFrameRate string `json:"r_frame_rate"`
	PixFmt     string `json:"pix_fmt"`
	SampleRate string `json:"sample_rate"`
	Channels   int    `json:"channels"`
}

type ffprobeFormat struct {
	FormatName string `json:"format_name"`
	Duration   string `json:"duration"`
	Size       string `json:"size"`
	BitRate    string `json:"bit_rate"`
}

func toMediaInfo(out ffprobeOutput) domain.MediaInfo {
	streams := make([]domain.StreamInfo, 0, len(out.Streams))
	for i := range out.Streams {
		streams = append(streams, toStreamInfo(out.Streams[i]))
	}

	return domain.MediaInfo{
		Format: domain.FormatInfo{
			FormatName: out.Format.FormatName,
			Duration:   out.Format.Duration,
			Size:       out.Format.Size,
			BitRate:    out.Format.BitRate,
		},
		Streams: streams,
	}
}

func toStreamInfo(s ffprobeStream) domain.StreamInfo {
	info := domain.StreamInfo{
		Index:      s.Index,
		Type:       s.CodecType,
		Codec:      s.CodecName,
		Width:      0,
		Height:     0,
		FPS:        "",
		PixFmt:     "",
		SampleRate: "",
		Channels:   0,
	}

	switch s.CodecType {
	case domain.StreamTypeVideo:
		info.Width = s.Width
		info.Height = s.Height
		info.FPS = formatFrameRate(s.RFrameRate)
		info.PixFmt = s.PixFmt
	case domain.StreamTypeAudio:
		info.SampleRate = s.SampleRate
		info.Channels = s.Channels
	}

	return info
}

// formatFrameRate turns ffprobe's "num/den" rate into a readable value, e.g.
// "30/1" -> "30" and "30000/1001" -> "29.970". It returns "" for missing or
// zero rates (still images report "0/0").
func formatFrameRate(rate string) string {
	num, den, ok := parseRatio(rate)
	if !ok || den == 0 || num == 0 {
		return ""
	}

	if num%den == 0 {
		return strconv.Itoa(num / den)
	}

	return fmt.Sprintf("%.3f", float64(num)/float64(den))
}

func parseRatio(rate string) (int, int, bool) {
	parts := strings.SplitN(rate, "/", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}

	num, numErr := strconv.Atoi(parts[0])
	den, denErr := strconv.Atoi(parts[1])

	if numErr != nil || denErr != nil {
		return 0, 0, false
	}

	return num, den, true
}
