package cli

import (
	"fmt"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/infrastructure"
	"github.com/junara/encfixture/usecase"

	"github.com/spf13/cobra"
)

var audioCmd = &cobra.Command{
	Use:   "audio",
	Short: "Generate a dummy audio file",
	Long:  "Generate a dummy audio file with various types (silence, sine, noise, tone).",
	Example: `  # Silent WAV (10s, stereo, 48kHz)
  encfixture audio -o silence.wav

  # 1000Hz sine wave
  encfixture audio -t sine -f 1000 -d 5 -o beep.wav

  # White noise
  encfixture audio -t noise -d 3 -o noise.mp3`,
	RunE: runAudio,
}

func init() {
	rootCmd.AddCommand(audioCmd)

	audioCmd.Flags().StringP("type", "t", "silence", "Audio type: silence, sine, noise, tone")
	audioCmd.Flags().StringP("duration", "d", "10", "Duration in seconds")
	audioCmd.Flags().IntP("sample-rate", "s", 48000, "Sample rate in Hz")
	audioCmd.Flags().IntP("channels", "C", 2, "Number of audio channels")
	audioCmd.Flags().Float64P("frequency", "f", 440.0, "Tone frequency in Hz (for sine/tone type)")
	audioCmd.Flags().StringP("output", "o", "output.wav", "Output file path (any format supported by ffmpeg)")
}

func runAudio(cmd *cobra.Command, _ []string) error {
	audioType, _ := cmd.Flags().GetString("type")
	duration, _ := cmd.Flags().GetString("duration")
	sampleRate, _ := cmd.Flags().GetInt("sample-rate")
	channels, _ := cmd.Flags().GetInt("channels")
	frequency, _ := cmd.Flags().GetFloat64("frequency")
	output, _ := cmd.Flags().GetString("output")

	cfg := domain.AudioConfig{
		Type:       domain.AudioType(audioType),
		Duration:   duration,
		SampleRate: sampleRate,
		Channels:   channels,
		Frequency:  frequency,
		Output:     output,
	}

	ffmpeg := infrastructure.NewFFmpeg()
	uc := usecase.NewAudioUseCase(ffmpeg)

	if err := uc.Generate(cfg); err != nil {
		return fmt.Errorf("audio generation failed: %w", err)
	}

	printResult(result{
		Status:   "ok",
		File:     output,
		Type:     "audio",
		Duration: duration,
	})

	return nil
}
