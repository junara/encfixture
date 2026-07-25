package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/infrastructure"
	"github.com/junara/encfixture/usecase"

	"github.com/spf13/cobra"
)

// errEnvUnhealthy reports that a required external tool is missing.
var errEnvUnhealthy = errors.New("environment is not ready")

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check ffmpeg/ffprobe availability and encoder support",
	Long: `Check the environment encfixture depends on: whether ffmpeg and ffprobe are
installed, and which of the selectable encoders (--codec values and the audio
encoders used for mp4/webm output) the local ffmpeg build supports.

Exits non-zero when ffmpeg or ffprobe is missing. Missing encoders are
reported but do not fail the check.`,
	Example: `  # Human-readable environment report
  encfixture doctor

  # Machine-readable JSON
  encfixture doctor --json`,
	Args: exactArgsWithHelpHint(0),
	RunE: runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(_ *cobra.Command, _ []string) error {
	uc := usecase.NewDoctorUseCase(infrastructure.NewInspector())

	report, err := uc.Report()
	if err != nil {
		return fmt.Errorf("doctor failed: %w", err)
	}

	printDoctorReport(report)

	if !report.Healthy() {
		missing := missingTools(report)

		return fmt.Errorf("%w: %s not found in PATH", errEnvUnhealthy, strings.Join(missing, ", "))
	}

	return nil
}

func missingTools(report domain.DoctorReport) []string {
	var missing []string

	if !report.FFmpeg.Available {
		missing = append(missing, "ffmpeg")
	}

	if !report.FFprobe.Available {
		missing = append(missing, "ffprobe")
	}

	return missing
}

type doctorResult struct {
	Status        string                 `json:"status"`
	FFmpeg        domain.ToolStatus      `json:"ffmpeg"`
	FFprobe       domain.ToolStatus      `json:"ffprobe"`
	VideoEncoders []domain.EncoderStatus `json:"videoEncoders"`
	AudioEncoders []domain.EncoderStatus `json:"audioEncoders"`
}

func printDoctorReport(report domain.DoctorReport) {
	if jsonOutput {
		status := "ok"
		if !report.Healthy() {
			status = statusError
		}

		writeJSON(doctorResult{
			Status:        status,
			FFmpeg:        report.FFmpeg,
			FFprobe:       report.FFprobe,
			VideoEncoders: report.VideoEncoders,
			AudioEncoders: report.AudioEncoders,
		})

		return
	}

	lines := []string{
		formatToolLine(report.FFmpeg),
		formatToolLine(report.FFprobe),
		"video encoders:",
	}

	for _, enc := range report.VideoEncoders {
		lines = append(lines, formatEncoderLine(enc))
	}

	lines = append(lines, "audio encoders:")
	for _, enc := range report.AudioEncoders {
		lines = append(lines, formatEncoderLine(enc))
	}

	if _, err := os.Stdout.WriteString(strings.Join(lines, "\n") + "\n"); err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
	}
}

func formatToolLine(tool domain.ToolStatus) string {
	if !tool.Available {
		return fmt.Sprintf("%-8s MISSING (install ffmpeg, e.g. 'brew install ffmpeg')", tool.Name+":")
	}

	line := fmt.Sprintf("%-8s OK", tool.Name+":")
	if tool.Version != "" {
		line += " " + tool.Version
	}

	if tool.Path != "" {
		line += " (" + tool.Path + ")"
	}

	return line
}

func formatEncoderLine(enc domain.EncoderStatus) string {
	state := "OK"
	if !enc.Available {
		state = "MISSING"
	}

	return fmt.Sprintf("  %-8s %-12s %s", enc.Codec, enc.Encoder, state)
}
