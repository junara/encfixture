package cli

import (
	"fmt"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/infrastructure"
	"github.com/junara/encfixture/usecase"

	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate a dummy image file",
	Long:  "Generate a dummy image file with overlays at specified positions.",
	Example: `  # Black solid image (1920x1080)
  encfixture image -o test.png

  # Blue background with centered text
  encfixture image -c blue --center "SAMPLE" -o sample.png

  # Color bar test pattern
  encfixture image -b test -o colorbar.png

  # All overlay positions
  encfixture image --tl frame --tr timecode --bl filename --br "ID-001" --center "TEST" -o info.png`,
	RunE: runImage,
}

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.Flags().IntP("width", "W", 1920, "Image width in pixels")
	imageCmd.Flags().IntP("height", "H", 1080, "Image height in pixels")
	imageCmd.Flags().StringP("bg", "b", "solid", "Background type: solid, test")
	imageCmd.Flags().StringP("color", "c", "black", "Background color (name or #hex)")
	imageCmd.Flags().String("tl", "", "Top-left content (frame, timecode, filename, or text)")
	imageCmd.Flags().String("tr", "", "Top-right content (frame, timecode, filename, or text)")
	imageCmd.Flags().String("center", "", "Center content (frame, timecode, filename, or text)")
	imageCmd.Flags().String("bl", "", "Bottom-left content (frame, timecode, filename, or text)")
	imageCmd.Flags().String("br", "", "Bottom-right content (frame, timecode, filename, or text)")
	imageCmd.Flags().IntP("scale", "S", 4, "Text scale factor")
	imageCmd.Flags().StringP("output", "o", "output.png", "Output file path")
}

func runImage(cmd *cobra.Command, _ []string) error {
	width, _ := cmd.Flags().GetInt("width")
	height, _ := cmd.Flags().GetInt("height")
	bg, _ := cmd.Flags().GetString("bg")
	colorName, _ := cmd.Flags().GetString("color")
	tl, _ := cmd.Flags().GetString("tl")
	tr, _ := cmd.Flags().GetString("tr")
	center, _ := cmd.Flags().GetString("center")
	bl, _ := cmd.Flags().GetString("bl")
	br, _ := cmd.Flags().GetString("br")
	scale, _ := cmd.Flags().GetInt("scale")
	output, _ := cmd.Flags().GetString("output")

	cfg := domain.ImageConfig{
		Width:      width,
		Height:     height,
		Background: bg,
		Color:      colorName,
		Overlay: domain.Overlay{
			TopLeft:     tl,
			TopRight:    tr,
			Center:      center,
			BottomLeft:  bl,
			BottomRight: br,
		},
		Scale:  scale,
		Output: output,
	}

	renderer := infrastructure.NewImageRenderer()
	uc := usecase.NewImageUseCase(renderer)

	if err := uc.Generate(cfg); err != nil {
		return fmt.Errorf("image generation failed: %w", err)
	}

	printResult(result{
		Status: "ok",
		File:   output,
		Type:   "image",
		Width:  width,
		Height: height,
	})

	return nil
}
