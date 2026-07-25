// Package cli provides the command-line interface for the encfixture tool.
package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/junara/encfixture/infrastructure"

	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

var (
	jsonOutput bool
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "encfixture",
	Short: "Generate dummy media files for ffmpeg encoding tests",
	Long: `encfixture generates dummy image, video, and audio files using ffmpeg
for encoding test purposes.

Each position flag (--tl, --tr, --center, --bl, --br) accepts:
  frame     - frame number (dynamic in video)
  timecode  - HH:MM:SS:FF timecode (dynamic in video)
  filename  - output filename
  <text>    - any other string is displayed as-is`,
	Version: Version,
	// Execute prints the error itself; without this cobra would print it too,
	// showing every failure twice.
	SilenceErrors: true,
	// Show usage only for usage errors (bad flags/args, hinted below), not for
	// runtime failures where it would drown the actual error.
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Show ffmpeg log output and encoding progress")

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return &usageError{err: err, path: cmd.CommandPath()}
	})
}

// usageError marks errors caused by wrong flags or arguments, so Execute can
// point the user at --help without doing that for runtime failures.
type usageError struct {
	err  error
	path string // command path for the --help hint, e.g. "encfixture verify"
}

func (u *usageError) Error() string { return u.err.Error() }

func (u *usageError) Unwrap() error { return u.err }

// newFFmpeg builds the ffmpeg runner, honoring the global --verbose flag.
func newFFmpeg() *infrastructure.FFmpeg {
	ffmpeg := infrastructure.NewFFmpeg()
	ffmpeg.Verbose = verbose

	return ffmpeg
}

// exactArgsWithHelpHint validates like cobra.ExactArgs but marks failures as
// usage errors, since usage is no longer printed automatically on errors.
func exactArgsWithHelpHint(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(n)(cmd, args); err != nil {
			return &usageError{err: err, path: cmd.CommandPath()}
		}

		return nil
	}
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		writeJSONError(err)
		fmt.Fprintf(os.Stderr, "encfixture: %v\n", err)

		var uErr *usageError
		if errors.As(err, &uErr) {
			fmt.Fprintf(os.Stderr, "Run '%s --help' for usage.\n", uErr.path)
		}

		os.Exit(1)
	}
}
