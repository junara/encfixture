package cli

import (
	"fmt"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/infrastructure"
	"github.com/junara/encfixture/usecase"

	"github.com/spf13/cobra"
)

var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Generate a dummy video file",
	Long:  "Generate a dummy video file with overlays at specified positions.",
	Example: `  # Black silent video (10s, 1080p, 30fps)
  encfixture video -o test.mp4

  # Frame counter + timecode overlay
  encfixture video --tl frame --tr timecode -d 5 -o counter.mp4

  # All overlays with color bar background
  encfixture video -b test --tl frame --tr timecode --bl filename --br "CLIP-001" --center "SAMPLE" -d 10 -o full.mp4

  # With sine wave audio
  encfixture video -a sine --frequency 1000 --center "BEEP" -d 3 -o beep.mp4`,
	RunE: runVideo,
}

func init() {
	rootCmd.AddCommand(videoCmd)

	videoCmd.Flags().IntP("width", "W", 1920, "Video width in pixels")
	videoCmd.Flags().IntP("height", "H", 1080, "Video height in pixels")
	videoCmd.Flags().IntP("fps", "r", 30, "Frames per second")
	videoCmd.Flags().StringP("duration", "d", "10", "Duration in seconds")
	videoCmd.Flags().StringP("bg", "b", "solid", "Background type: solid, test")
	videoCmd.Flags().StringP("color", "c", "black", "Background color (name or #hex)")
	videoCmd.Flags().String("tl", "", "Top-left content (frame, timecode, filename, or text)")
	videoCmd.Flags().String("tr", "", "Top-right content (frame, timecode, filename, or text)")
	videoCmd.Flags().String("center", "", "Center content (frame, timecode, filename, or text)")
	videoCmd.Flags().String("bl", "", "Bottom-left content (frame, timecode, filename, or text)")
	videoCmd.Flags().String("br", "", "Bottom-right content (frame, timecode, filename, or text)")
	videoCmd.Flags().IntP("scale", "S", 4, "Text scale factor")
	videoCmd.Flags().StringP("output", "o", "output.mp4", "Output file path (any format supported by ffmpeg)")
	videoCmd.Flags().StringP("audio", "a", "silence", "Audio type: silence, sine, noise, tone")
	videoCmd.Flags().IntP("sample-rate", "s", 48000, "Audio sample rate")
	videoCmd.Flags().IntP("channels", "C", 2, "Audio channels")
	videoCmd.Flags().Float64("frequency", 440.0, "Tone frequency in Hz (for sine/tone audio)")
}

func runVideo(cmd *cobra.Command, _ []string) error {
	width, _ := cmd.Flags().GetInt("width")
	height, _ := cmd.Flags().GetInt("height")
	fps, _ := cmd.Flags().GetInt("fps")
	duration, _ := cmd.Flags().GetString("duration")
	bg, _ := cmd.Flags().GetString("bg")
	colorName, _ := cmd.Flags().GetString("color")
	tl, _ := cmd.Flags().GetString("tl")
	tr, _ := cmd.Flags().GetString("tr")
	center, _ := cmd.Flags().GetString("center")
	bl, _ := cmd.Flags().GetString("bl")
	br, _ := cmd.Flags().GetString("br")
	scale, _ := cmd.Flags().GetInt("scale")
	output, _ := cmd.Flags().GetString("output")
	audio, _ := cmd.Flags().GetString("audio")
	sampleRate, _ := cmd.Flags().GetInt("sample-rate")
	channels, _ := cmd.Flags().GetInt("channels")
	frequency, _ := cmd.Flags().GetFloat64("frequency")

	cfg := domain.VideoConfig{
		Width:      width,
		Height:     height,
		FPS:        fps,
		Duration:   duration,
		Background: bg,
		Color:      colorName,
		Overlay: domain.Overlay{
			TopLeft:     tl,
			TopRight:    tr,
			Center:      center,
			BottomLeft:  bl,
			BottomRight: br,
		},
		Scale:      scale,
		Output:     output,
		Audio:      domain.AudioType(audio),
		SampleRate: sampleRate,
		Channels:   channels,
		Frequency:  frequency,
	}

	ffmpeg := infrastructure.NewFFmpeg()
	renderer := infrastructure.NewImageRenderer()
	uc := usecase.NewVideoUseCase(ffmpeg, renderer)

	if err := uc.Generate(cfg); err != nil {
		return fmt.Errorf("video generation failed: %w", err)
	}

	printResult(result{
		Status:   "ok",
		File:     output,
		Type:     "video",
		Width:    width,
		Height:   height,
		FPS:      fps,
		Duration: duration,
	})

	return nil
}
