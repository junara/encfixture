package cli

import (
	"encoding/json"
	"fmt"
	"os"
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
	Args: cobra.ExactArgs(1),
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

	lines := []string{
		"File:     " + path,
		"Format:   " + info.Format.FormatName,
		"Duration: " + info.Format.Duration + "s",
		"Size:     " + info.Format.Size + " bytes",
	}

	for _, stream := range info.Streams {
		lines = append(lines, formatStreamLine(stream))
	}

	_, err := os.Stdout.WriteString(strings.Join(lines, "\n") + "\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
	}
}

func formatStreamLine(s domain.StreamInfo) string {
	switch s.Type {
	case domain.StreamTypeVideo:
		return fmt.Sprintf("Stream %d: video  %s  %dx%d  %sfps  %s", s.Index, s.Codec, s.Width, s.Height, s.FPS, s.PixFmt)
	case domain.StreamTypeAudio:
		return fmt.Sprintf("Stream %d: audio  %s  %sHz  %dch", s.Index, s.Codec, s.SampleRate, s.Channels)
	default:
		return fmt.Sprintf("Stream %d: %s  %s", s.Index, s.Type, s.Codec)
	}
}
