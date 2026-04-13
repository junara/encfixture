package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/junara/encfixture/domain"
)

// ErrBatchInvalid is returned when a batch file is structurally valid JSON but
// fails domain validation (missing required fields, unknown job type, etc.).
var ErrBatchInvalid = errors.New("invalid batch")

type jobDTO struct {
	Type       *string  `json:"type"`
	Width      *int     `json:"width"`
	Height     *int     `json:"height"`
	FPS        *int     `json:"fps"`
	Duration   *string  `json:"duration"`
	BG         *string  `json:"bg"`
	Color      *string  `json:"color"`
	TL         *string  `json:"tl"`
	TR         *string  `json:"tr"`
	Center     *string  `json:"center"`
	BL         *string  `json:"bl"`
	BR         *string  `json:"br"`
	Scale      *int     `json:"scale"`
	Output     *string  `json:"output"`
	Audio      *string  `json:"audio"`
	SampleRate *int     `json:"sampleRate"`
	Channels   *int     `json:"channels"`
	Frequency  *float64 `json:"frequency"`
}

type batchDTO struct {
	Defaults *jobDTO  `json:"defaults"`
	Jobs     []jobDTO `json:"jobs"`
}

// LoadBatch reads a JSON batch definition from path, merges per-job values with
// the file-level defaults, and returns a normalized domain.Batch.
func LoadBatch(path string) (domain.Batch, error) {
	file, err := os.Open(path) //nolint:gosec // path is user-supplied by design.
	if err != nil {
		return domain.Batch{Jobs: nil}, fmt.Errorf("open batch file: %w", err)
	}

	defer func() { _ = file.Close() }()

	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()

	var dto batchDTO

	decErr := dec.Decode(&dto)
	if decErr != nil {
		return domain.Batch{Jobs: nil}, fmt.Errorf("decode batch file: %w", decErr)
	}

	defaults := jobDTO{
		Type: nil, Width: nil, Height: nil, FPS: nil, Duration: nil,
		BG: nil, Color: nil, TL: nil, TR: nil, Center: nil, BL: nil, BR: nil,
		Scale: nil, Output: nil, Audio: nil, SampleRate: nil, Channels: nil, Frequency: nil,
	}
	if dto.Defaults != nil {
		defaults = *dto.Defaults
	}

	jobs := make([]domain.Job, 0, len(dto.Jobs))

	for i := range dto.Jobs {
		job, convErr := dtoToJob(&defaults, &dto.Jobs[i])
		if convErr != nil {
			return domain.Batch{Jobs: nil}, fmt.Errorf("job[%d]: %w", i, convErr)
		}

		jobs = append(jobs, job)
	}

	return domain.Batch{Jobs: jobs}, nil
}

func dtoToJob(defaults, job *jobDTO) (domain.Job, error) {
	typePtr := job.Type
	if typePtr == nil {
		typePtr = defaults.Type
	}

	if typePtr == nil {
		return domain.Job{Type: "", Image: nil, Video: nil, Audio: nil},
			fmt.Errorf("%w: missing type", ErrBatchInvalid)
	}

	output := pickString(job.Output, defaults.Output, "")
	if output == "" {
		return domain.Job{Type: "", Image: nil, Video: nil, Audio: nil},
			fmt.Errorf("%w: missing output", ErrBatchInvalid)
	}

	jobType := domain.JobType(*typePtr)

	switch jobType {
	case domain.JobTypeImage:
		img := buildImage(defaults, job, output)

		return domain.Job{Type: jobType, Image: &img, Video: nil, Audio: nil}, nil
	case domain.JobTypeVideo:
		vid := buildVideo(defaults, job, output)

		return domain.Job{Type: jobType, Image: nil, Video: &vid, Audio: nil}, nil
	case domain.JobTypeAudio:
		aud := buildAudio(defaults, job, output)

		return domain.Job{Type: jobType, Image: nil, Video: nil, Audio: &aud}, nil
	default:
		return domain.Job{Type: "", Image: nil, Video: nil, Audio: nil},
			fmt.Errorf("%w: unknown type %q", ErrBatchInvalid, *typePtr)
	}
}

func buildImage(defaults, job *jobDTO, output string) domain.ImageConfig {
	return domain.ImageConfig{
		Width:      pickInt(job.Width, defaults.Width, 1920),
		Height:     pickInt(job.Height, defaults.Height, 1080),
		Background: pickString(job.BG, defaults.BG, "solid"),
		Color:      pickString(job.Color, defaults.Color, "black"),
		Overlay:    buildOverlay(defaults, job),
		Scale:      pickInt(job.Scale, defaults.Scale, 4),
		Output:     output,
	}
}

func buildVideo(defaults, job *jobDTO, output string) domain.VideoConfig {
	return domain.VideoConfig{
		Width:      pickInt(job.Width, defaults.Width, 1920),
		Height:     pickInt(job.Height, defaults.Height, 1080),
		FPS:        pickInt(job.FPS, defaults.FPS, 30),
		Duration:   pickString(job.Duration, defaults.Duration, "10"),
		Background: pickString(job.BG, defaults.BG, "solid"),
		Color:      pickString(job.Color, defaults.Color, "black"),
		Overlay:    buildOverlay(defaults, job),
		Scale:      pickInt(job.Scale, defaults.Scale, 4),
		Output:     output,
		Audio:      domain.AudioType(pickString(job.Audio, defaults.Audio, string(domain.AudioSilence))),
		SampleRate: pickInt(job.SampleRate, defaults.SampleRate, 48000),
		Channels:   pickInt(job.Channels, defaults.Channels, 2),
		Frequency:  pickFloat(job.Frequency, defaults.Frequency, 440.0),
	}
}

func buildAudio(defaults, job *jobDTO, output string) domain.AudioConfig {
	return domain.AudioConfig{
		Type:       domain.AudioType(pickString(job.Audio, defaults.Audio, string(domain.AudioSilence))),
		Duration:   pickString(job.Duration, defaults.Duration, "10"),
		SampleRate: pickInt(job.SampleRate, defaults.SampleRate, 48000),
		Channels:   pickInt(job.Channels, defaults.Channels, 2),
		Frequency:  pickFloat(job.Frequency, defaults.Frequency, 440.0),
		Output:     output,
	}
}

func buildOverlay(defaults, job *jobDTO) domain.Overlay {
	return domain.Overlay{
		TopLeft:     pickString(job.TL, defaults.TL, ""),
		TopRight:    pickString(job.TR, defaults.TR, ""),
		Center:      pickString(job.Center, defaults.Center, ""),
		BottomLeft:  pickString(job.BL, defaults.BL, ""),
		BottomRight: pickString(job.BR, defaults.BR, ""),
	}
}

func pickString(primary, fallback *string, def string) string {
	if primary != nil {
		return *primary
	}

	if fallback != nil {
		return *fallback
	}

	return def
}

func pickInt(primary, fallback *int, def int) int {
	if primary != nil {
		return *primary
	}

	if fallback != nil {
		return *fallback
	}

	return def
}

func pickFloat(primary, fallback *float64, def float64) float64 {
	if primary != nil {
		return *primary
	}

	if fallback != nil {
		return *fallback
	}

	return def
}
