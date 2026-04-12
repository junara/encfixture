// Package infrastructure provides implementations for ffmpeg execution and image rendering.
package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// ErrFFmpegNotFound indicates ffmpeg is not installed or not in PATH.
var ErrFFmpegNotFound = errors.New("ffmpeg not found in PATH: please install ffmpeg")

// FFmpeg provides methods for executing ffmpeg commands.
type FFmpeg struct{}

// NewFFmpeg creates a new FFmpeg instance.
func NewFFmpeg() *FFmpeg {
	return &FFmpeg{}
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
	cmd := exec.CommandContext(context.Background(), "ffmpeg", args...) //nolint:gosec // ffmpeg requires dynamic arguments

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(out))

		return fmt.Errorf("ffmpeg execution failed: %w", err)
	}

	return nil
}

// RunWithStdin executes ffmpeg with the given arguments and stdin reader.
func (f *FFmpeg) RunWithStdin(stdin io.Reader, args ...string) error {
	cmd := exec.CommandContext(context.Background(), "ffmpeg", args...) //nolint:gosec // ffmpeg requires dynamic arguments
	cmd.Stdin = stdin

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(out))

		return fmt.Errorf("ffmpeg execution failed: %w", err)
	}

	return nil
}
