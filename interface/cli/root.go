// Package cli provides the command-line interface for the encfixture tool.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

var jsonOutput bool

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
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "encfixture: %v\n", err)
		os.Exit(1)
	}
}
