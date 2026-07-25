package cli

import (
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
	Short: "Inspect a media file via ffprobe, optionally asserting expected properties",
	Long: `Inspect an existing media file with ffprobe and print its container and
per-stream properties.

With --expect key=value assertions, each expectation is checked against the
probed properties and the command exits non-zero if any of them fail.
Supported keys: codec, width, height, fps, pixFmt, duration, audioCodec,
sampleRate, channels. Numeric keys accept a tolerance suffix ("5+-0.2");
duration defaults to a 0.1s tolerance.`,
	Example: `  # Human-readable summary
  encfixture verify test.mp4

  # Machine-readable JSON
  encfixture verify --json test.mp4

  # Assert properties; exits 1 if any expectation fails
  encfixture verify test.mp4 --expect codec=h264 --expect width=1920 --expect duration=5

  # Duration with explicit tolerance
  encfixture verify --json test.mp4 --expect "duration=5+-0.5"`,
	Args: exactArgsWithHelpHint(1),
	RunE: runVerify,
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringArray("expect", nil, "Assertion key=value (repeatable); see --help for supported keys")
}

func runVerify(cmd *cobra.Command, args []string) error {
	path := args[0]

	expectValues, _ := cmd.Flags().GetStringArray("expect")

	exps, parseErr := parseExpectations(expectValues)
	if parseErr != nil {
		return parseErr
	}

	prober := infrastructure.NewFFprobe()
	uc := usecase.NewVerifyUseCase(prober)

	info, checks, err := uc.VerifyWithExpectations(path, exps)
	if err != nil {
		return fmt.Errorf("verify failed: %w", err)
	}

	printMediaInfo(path, info, checks)

	if failed := countFailed(checks); failed > 0 {
		return fmt.Errorf("%w: %d of %d expectations failed", errVerifyFailed, failed, len(checks))
	}

	return nil
}

// parseExpectations parses each --expect value. Malformed values surface as
// domain.ErrInvalidExpectation, whose classified hint lists the valid keys.
func parseExpectations(values []string) ([]domain.Expectation, error) {
	exps := make([]domain.Expectation, 0, len(values))

	for _, value := range values {
		exp, err := domain.ParseExpectation(value)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		exps = append(exps, exp)
	}

	return exps, nil
}

func countFailed(checks []domain.CheckResult) int {
	failed := 0

	for _, check := range checks {
		if !check.Pass {
			failed++
		}
	}

	return failed
}

type verifyResult struct {
	Status  string               `json:"status"`
	File    string               `json:"file"`
	Format  domain.FormatInfo    `json:"format"`
	Streams []domain.StreamInfo  `json:"streams"`
	Checks  []domain.CheckResult `json:"checks,omitempty"`
}

func printMediaInfo(path string, info domain.MediaInfo, checks []domain.CheckResult) {
	if jsonOutput {
		status := "ok"
		if countFailed(checks) > 0 {
			status = "failed"
		}

		writeJSON(verifyResult{Status: status, File: path, Format: info.Format, Streams: info.Streams, Checks: checks})

		return
	}

	lines := []string{"File:     " + path}
	lines = appendField(lines, "Format:   ", info.Format.FormatName, "")
	lines = appendField(lines, "Duration: ", info.Format.Duration, "s")
	lines = appendField(lines, "Size:     ", info.Format.Size, " bytes")

	for _, stream := range info.Streams {
		lines = append(lines, formatStreamLine(stream))
	}

	for _, check := range checks {
		lines = append(lines, formatCheckLine(check))
	}

	_, err := os.Stdout.WriteString(strings.Join(lines, "\n") + "\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
	}
}

func formatCheckLine(check domain.CheckResult) string {
	if check.Pass {
		return fmt.Sprintf("PASS  %s = %s", check.Field, check.Actual)
	}

	return fmt.Sprintf("FAIL  %s: expected %s, actual %s", check.Field, check.Expected, check.Actual)
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
