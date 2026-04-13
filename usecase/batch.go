package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/junara/encfixture/domain"
)

// ErrUnknownJobType indicates a batch job has an unrecognized Type.
var ErrUnknownJobType = errors.New("unknown job type")

// ErrMissingJobConfig indicates a batch job has no config for its declared Type.
var ErrMissingJobConfig = errors.New("missing job config")

// BatchOptions controls concurrency and failure behavior for BatchUseCase.Generate.
type BatchOptions struct {
	// Parallel is the maximum number of jobs running concurrently. Values < 1
	// are treated as 1.
	Parallel int
	// FailFast cancels scheduling of new jobs after the first failure. Jobs
	// already running continue to completion; pending jobs return a context
	// cancellation error in their JobResult.
	FailFast bool
}

// BatchUseCase orchestrates multiple media generation jobs, delegating each to
// the matching single-media use case.
type BatchUseCase struct {
	image *ImageUseCase
	video *VideoUseCase
	audio *AudioUseCase
}

// NewBatchUseCase constructs a BatchUseCase wired to the given single-media use cases.
func NewBatchUseCase(image *ImageUseCase, video *VideoUseCase, audio *AudioUseCase) *BatchUseCase {
	return &BatchUseCase{image: image, video: video, audio: audio}
}

// Generate runs every job in batch, respecting opts.Parallel as a concurrency
// ceiling. Results are returned in the same order as batch.Jobs. Failures in
// individual jobs are reported via JobResult.Err; Generate itself does not
// return an error.
func (uc *BatchUseCase) Generate(ctx context.Context, batch domain.Batch, opts BatchOptions) []domain.JobResult {
	parallel := max(opts.Parallel, 1)

	results := make([]domain.JobResult, len(batch.Jobs))
	sem := make(chan struct{}, parallel)

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	for i, job := range batch.Jobs {
		wg.Go(func() {
			results[i] = uc.runJob(runCtx, sem, i, job, opts.FailFast, cancel)
		})
	}

	wg.Wait()

	return results
}

func (uc *BatchUseCase) runJob(
	ctx context.Context, sem chan struct{}, idx int, job domain.Job, failFast bool, cancel context.CancelFunc,
) domain.JobResult {
	result := domain.JobResult{Index: idx, Type: job.Type, Output: jobOutput(job), Err: nil}

	select {
	case <-ctx.Done():
		result.Err = fmt.Errorf("job skipped: %w", ctx.Err())

		return result
	case sem <- struct{}{}:
	}

	defer func() { <-sem }()

	err := uc.execute(job)
	if err != nil {
		result.Err = err

		if failFast {
			cancel()
		}
	}

	return result
}

func (uc *BatchUseCase) execute(job domain.Job) error {
	switch job.Type {
	case domain.JobTypeImage:
		if job.Image == nil {
			return fmt.Errorf("%w: image", ErrMissingJobConfig)
		}

		return uc.image.Generate(*job.Image)
	case domain.JobTypeVideo:
		if job.Video == nil {
			return fmt.Errorf("%w: video", ErrMissingJobConfig)
		}

		return uc.video.Generate(*job.Video)
	case domain.JobTypeAudio:
		if job.Audio == nil {
			return fmt.Errorf("%w: audio", ErrMissingJobConfig)
		}

		return uc.audio.Generate(*job.Audio)
	default:
		return fmt.Errorf("%w: %s", ErrUnknownJobType, job.Type)
	}
}

func jobOutput(job domain.Job) string {
	switch job.Type {
	case domain.JobTypeImage:
		if job.Image != nil {
			return job.Image.Output
		}
	case domain.JobTypeVideo:
		if job.Video != nil {
			return job.Video.Output
		}
	case domain.JobTypeAudio:
		if job.Audio != nil {
			return job.Audio.Output
		}
	}

	return ""
}
