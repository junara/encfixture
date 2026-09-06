// Package cli provides the command-line interface for the encfixture tool.
package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/junara/encfixture/domain"
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

// errConflictingFlags reports two flags that steer the same setting in
// different ways being set together.
var errConflictingFlags = errors.New("conflicting flags")

// exclusiveFlags fails when more than one of names was set, so conflicting
// intent is reported instead of one flag silently winning.
func exclusiveFlags(cmd *cobra.Command, names ...string) error {
	var set []string

	for _, name := range names {
		if cmd.Flags().Changed(name) {
			set = append(set, name)
		}
	}

	if len(set) > 1 {
		return &usageError{
			err:  fmt.Errorf("%w: --%s cannot be combined", errConflictingFlags, strings.Join(set, " and --")),
			path: cmd.CommandPath(),
		}
	}

	return nil
}

// targetSizeFlag reads --target-size as a byte count. An unset flag yields 0,
// meaning the encoder settings decide the size as before.
func targetSizeFlag(cmd *cobra.Command) (int64, error) {
	raw, _ := cmd.Flags().GetString("target-size")
	if raw == "" {
		return 0, nil
	}

	size, err := domain.ParseFileSize(raw)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return size, nil
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
		code, hint := classifyError(err)
		writeJSONError(err, code, hint)
		fmt.Fprintf(os.Stderr, "encfixture: %v\n", err)

		if hint != "" {
			fmt.Fprintf(os.Stderr, "hint: %s\n", hint)
		}

		os.Exit(1)
	}
}
