// Package infrastructure provides implementations for ffmpeg execution and image rendering.
package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// ErrFFmpegNotFound indicates ffmpeg is not installed or not in PATH.
var ErrFFmpegNotFound = errors.New("ffmpeg not found in PATH: please install ffmpeg")

// ErrFFmpegFailed indicates an ffmpeg invocation exited with an error.
var ErrFFmpegFailed = errors.New("ffmpeg execution failed")

// ErrEncoderNotAvailable indicates the requested encoder is not compiled into
// the local ffmpeg build.
var ErrEncoderNotAvailable = errors.New("encoder not available in this ffmpeg build")

// FFmpeg provides methods for executing ffmpeg commands.
type FFmpeg struct {
	// Verbose streams ffmpeg's own log and encoding progress to stderr
	// instead of suppressing everything but errors.
	Verbose bool
}

// NewFFmpeg creates a new FFmpeg instance.
func NewFFmpeg() *FFmpeg {
	return &FFmpeg{Verbose: false}
}

// CheckAvailable verifies that ffmpeg is available in the system PATH.
func (f *FFmpeg) CheckAvailable() error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("%w", ErrFFmpegNotFound)
	}

	return nil
}

// Run executes ffmpeg with the given arguments.
func (f *FFmpeg) Run(args ...string) error {
	return f.run(nil, args)
}

// RunWithStdin executes ffmpeg with the given arguments and stdin reader.
func (f *FFmpeg) RunWithStdin(stdin io.Reader, args ...string) error {
	return f.run(stdin, args)
}

// run executes ffmpeg, keeping its output quiet unless it fails or Verbose is
// set. Non-verbose runs surface only ffmpeg's real error lines; verbose runs
// stream the banner-free log and encoding progress to stderr as they happen.
func (f *FFmpeg) run(stdin io.Reader, args []string) error {
	if f.Verbose {
		args = append([]string{"-hide_banner", "-stats"}, args...)
	} else {
		args = append([]string{"-hide_banner", "-loglevel", "error"}, args...)
	}

	cmd := exec.CommandContext(context.Background(), "ffmpeg", args...) //nolint:gosec // ffmpeg requires dynamic arguments
	cmd.Stdin = stdin

	if f.Verbose {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrFFmpegFailed, err)
		}

		return nil
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		fmt.Fprintln(os.Stderr, msg)

		if line, found := lineContaining(msg, "Unknown encoder"); found {
			return fmt.Errorf("%w: %s", ErrEncoderNotAvailable, line)
		}

		return fmt.Errorf("%w: %w", ErrFFmpegFailed, err)
	}

	return nil
}

// lineContaining returns the first line of msg containing substr, so a typed
// error can carry ffmpeg's own diagnostic without the rest of the log.
func lineContaining(msg, substr string) (string, bool) {
	for line := range strings.SplitSeq(msg, "\n") {
		if strings.Contains(line, substr) {
			return strings.TrimSpace(line), true
		}
	}

	return "", false
}
