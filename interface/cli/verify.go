package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/infrastructure"
	"github.com/junara/encfixture/usecase"

	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify <file>",
	Short: "Inspect a media file's codec, resolution, fps, and duration via ffprobe",
	Long:  "Inspect an existing media file with ffprobe and print its container and per-stream properties.",
	Example: `  # Human-readable summary
  encfixture verify test.mp4

  # Machine-readable JSON
  encfixture verify --json test.mp4`,
	Args: exactArgsWithHelpHint(1),
	RunE: runVerify,
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}

func runVerify(_ *cobra.Command, args []string) error {
	path := args[0]

	prober := infrastructure.NewFFprobe()
	uc := usecase.NewVerifyUseCase(prober)

	info, err := uc.Verify(path)
	if err != nil {
		return fmt.Errorf("verify failed: %w", err)
	}

	printMediaInfo(path, info)

	return nil
}

type verifyResult struct {
	File    string              `json:"file"`
	Format  domain.FormatInfo   `json:"format"`
	Streams []domain.StreamInfo `json:"streams"`
}

func printMediaInfo(path string, info domain.MediaInfo) {
	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)

		encErr := enc.Encode(verifyResult{File: path, Format: info.Format, Streams: info.Streams})
		if encErr != nil {
			fmt.Fprintf(os.Stderr, "json encode error: %v\n", encErr)
		}

		return
	}

	lines := []string{"File:     " + path}
	lines = appendField(lines, "Format:   ", info.Format.FormatName, "")
	lines = appendField(lines, "Duration: ", info.Format.Duration, "s")
	lines = appendField(lines, "Size:     ", info.Format.Size, " bytes")

	for _, stream := range info.Streams {
		lines = append(lines, formatStreamLine(stream))
	}

	_, err := os.Stdout.WriteString(strings.Join(lines, "\n") + "\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
	}
}

// appendField adds "<label><value><suffix>" only when value is non-empty, so an
// unknown property (e.g. a still image's duration) is omitted rather than shown
// as an empty stub.
func appendField(lines []string, label, value, suffix string) []string {
	if value == "" {
		return lines
	}

	return append(lines, label+value+suffix)
}

func formatStreamLine(s domain.StreamInfo) string {
	switch s.Type {
	case domain.StreamTypeVideo:
		parts := []string{fmt.Sprintf("Stream %d: video", s.Index), s.Codec, fmt.Sprintf("%dx%d", s.Width, s.Height)}
		if s.FPS != "" {
			parts = append(parts, s.FPS+"fps")
		}

		if s.PixFmt != "" {
			parts = append(parts, s.PixFmt)
		}

		return strings.Join(parts, "  ")
	case domain.StreamTypeAudio:
		parts := []string{fmt.Sprintf("Stream %d: audio", s.Index), s.Codec}
		if s.SampleRate != "" {
			parts = append(parts, s.SampleRate+"Hz")
		}

		if s.Channels > 0 {
			parts = append(parts, strconv.Itoa(s.Channels)+"ch")
		}

		return strings.Join(parts, "  ")
	default:
		return fmt.Sprintf("Stream %d: %s  %s", s.Index, s.Type, s.Codec)
	}
}
