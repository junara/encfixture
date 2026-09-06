package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	Status     string `json:"status"`
	File       string `json:"file"`
	Type       string `json:"type"`
	Width      int    `json:"width,omitzero"`
	Height     int    `json:"height,omitzero"`
	FPS        int    `json:"fps,omitzero"`
	Duration   string `json:"duration,omitempty"`
	Size       int64  `json:"size,omitzero"`
	TargetSize int64  `json:"targetSize,omitzero"`
	Warning    string `json:"warning,omitempty"`
}

// sizeTolerance is how far the produced file may drift from --target-size
// before the result says so. Encoders only approximate a requested size, and
// content that compresses to almost nothing can miss it by a wide margin.
const sizeTolerance = 0.1

// fileSize returns the size of path in bytes, or 0 when it cannot be read.
func fileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return info.Size()
}

// targetSizeWarning reports how far the generated file landed from the
// requested size, with the knobs that move it in the right direction. It
// returns an empty string when no target was set or the drift is within
// sizeTolerance.
func targetSizeWarning(kind string, actual, target int64) string {
	if target <= 0 || actual <= 0 {
		return ""
	}

	drift := math.Abs(float64(actual-target)) / float64(target)
	if drift <= sizeTolerance {
		return ""
	}

	return fmt.Sprintf(
		"output is %s, %.0f%% off the requested %s; to get closer, %s",
		domain.FormatFileSize(actual), drift*100, domain.FormatFileSize(target),
		targetSizeAdvice(kind, actual < target),
	)
}

// targetSizeAdvice names the flags that move a given media type toward its
// target size: each type is steered by a different set, and an encoder that
// cannot fill the requested bitrate needs more detail to encode, not a
// different rate.
func targetSizeAdvice(kind string, undershot bool) string {
	switch {
	case kind == "audio" && undershot:
		return "use -t noise or sine (silence collapses in AAC and Opus), or write .mp3, which honors the bitrate"
	case kind == "audio":
		return "shorten -d, or lower --sample-rate/--channels"
	case kind == "image" && undershot:
		return "raise --width/--height or use a busier --bg (test, gradient, moving); JPEG quality 100 is the ceiling"
	case kind == "image":
		return "lower --width/--height, or use a plainer --bg"
	case undershot:
		return "raise --width/--height or --fps, use a busier --bg (test, gradient, moving)," +
			" or pick --codec h264, which tracks a target bitrate closest"
	default:
		return "lower --width/--height or --fps, use a plainer --bg, or shorten -d"
	}
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
	if uErr, ok := errors.AsType[*usageError](err); ok {
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
	case errors.Is(err, domain.ErrInvalidTargetSize):
		return "invalid_target_size", "Pass a positive size with an optional K/M/G suffix, e.g. --target-size 500KB or 10MB."
	case errors.Is(err, domain.ErrTargetSizeUnsupported):
		return "target_size_unsupported", "--target-size needs a format with a bitrate or quality dial: .jpg for images, a compressed codec for audio."
	case errors.Is(err, domain.ErrTargetSizeConflict):
		return "target_size_conflict", "Set either --target-size or --bitrate/--crf (in a batch job: targetSize or bitrate/crf), not both."
	case errors.Is(err, domain.ErrTargetSizeTooSmall):
		return "target_size_too_small", "Raise --target-size or shorten the duration so the streams fit in the budget."
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

	if r.Warning != "" {
		fmt.Fprintf(os.Stderr, "warning: %s\n", r.Warning)
	}
}
