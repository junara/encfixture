package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/infrastructure"
	"github.com/junara/encfixture/usecase"
)

// statusError is the JSON status value shared by every failure shape.
const statusError = "error"

// errOutputExists reports a refusal to overwrite an existing file under --no-clobber.
var errOutputExists = errors.New("output file already exists")

// errVerifyFailed reports that one or more --expect assertions did not hold.
var errVerifyFailed = errors.New("verification failed")

// checkClobber refuses to overwrite path when noClobber is set. By default
// generation commands overwrite existing files, matching their fixture
// regeneration purpose; --no-clobber opts into failing instead.
func checkClobber(path string, noClobber bool) error {
	if !noClobber {
		return nil
	}

	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("%w: %s (remove it or drop --no-clobber to overwrite)", errOutputExists, path)
	}

	return nil
}

type result struct {
	Status   string `json:"status"`
	File     string `json:"file"`
	Type     string `json:"type"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	FPS      int    `json:"fps,omitempty"`
	Duration string `json:"duration,omitempty"`
}

type errorOutput struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Error  string `json:"error"`
	Hint   string `json:"hint,omitempty"`
}

// classifyError maps err to a stable machine-readable code plus a recovery
// hint, so agents and scripts can branch on the failure without parsing the
// human-oriented message. Unrecognized errors get the generic "error" code.
func classifyError(err error) (string, string) {
	var uErr *usageError
	if errors.As(err, &uErr) {
		return "usage", fmt.Sprintf("Run '%s --help' for usage.", uErr.path)
	}

	switch {
	case errors.Is(err, infrastructure.ErrFFmpegNotFound):
		return "ffmpeg_not_found", "Install ffmpeg (e.g. 'brew install ffmpeg') and ensure it is in PATH. Run 'encfixture doctor' to check the environment."
	case errors.Is(err, infrastructure.ErrFFprobeNotFound):
		return "ffprobe_not_found", "ffprobe ships with ffmpeg. Install ffmpeg (e.g. 'brew install ffmpeg') and ensure ffprobe is in PATH."
	case errors.Is(err, infrastructure.ErrEncoderNotAvailable):
		return "encoder_not_available", "Run 'encfixture doctor' to list the encoders your ffmpeg build supports, then pick an available --codec."
	case errors.Is(err, usecase.ErrUnknownVideoCodec):
		return "unknown_codec", "Supported codecs: h264, hevc, vp9, av1, prores."
	case errors.Is(err, usecase.ErrUnknownBackground):
		return "unknown_background", "Supported backgrounds: solid, test, gradient, moving."
	case errors.Is(err, domain.ErrInvalidDuration):
		return "invalid_duration", "Pass a positive number of seconds, e.g. -d 5 or -d 2.5."
	case errors.Is(err, domain.ErrInvalidBitrate):
		return "invalid_bitrate", "Pass a number with an optional k/M/G suffix, e.g. --bitrate 800k or 5M."
	case errors.Is(err, domain.ErrInvalidExpectation):
		return "invalid_expectation", "Use --expect key=value with keys: codec, width, height, fps, pixFmt, duration, audioCodec, sampleRate, channels."
	case errors.Is(err, errOutputExists):
		return "output_exists", "Remove the existing file or drop --no-clobber to overwrite it."
	case errors.Is(err, errVerifyFailed):
		return "verify_failed", "Compare the expected and actual values in the reported checks."
	case errors.Is(err, errEnvUnhealthy):
		return "env_unhealthy", "Install ffmpeg (includes ffprobe), e.g. 'brew install ffmpeg'."
	case errors.Is(err, infrastructure.ErrFFprobeFailed):
		return "probe_failed", "Check that the file exists and is a readable media file."
	case errors.Is(err, infrastructure.ErrFFmpegFailed):
		return "ffmpeg_failed", "Re-run with --verbose to see the full ffmpeg log."
	default:
		return statusError, ""
	}
}

// jsonWritten tracks whether a JSON document was already emitted to stdout, so
// Execute does not append a second one after e.g. a batch summary with failures.
var jsonWritten bool

// writeJSON encodes v to stdout as a single JSON document.
func writeJSON(v any) {
	jsonWritten = true

	if err := json.NewEncoder(os.Stdout).Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "json encode error: %v\n", err)
	}
}

// writeJSONError reports err on stdout in JSON form when --json is active and
// nothing has been emitted yet, so machine consumers never get an empty stdout.
func writeJSONError(err error, code, hint string) {
	if !jsonOutput || jsonWritten {
		return
	}

	writeJSON(errorOutput{Status: statusError, Code: code, Error: err.Error(), Hint: hint})
}

// printProgress emits a status line to stderr when it is a terminal, so
// interactive users see activity during long encodes without adding noise to
// logs or pipelines.
func printProgress(msg string) {
	info, err := os.Stderr.Stat()
	if err != nil || info.Mode()&os.ModeCharDevice == 0 {
		return
	}

	fmt.Fprintln(os.Stderr, msg)
}

func printResult(r result) {
	if jsonOutput {
		writeJSON(r)

		return
	}

	fmt.Fprintf(os.Stderr, "Created: %s\n", r.File)
}
