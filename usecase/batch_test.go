package usecase_test

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/usecase"
)

type countingFFmpeg struct {
	inFlight *atomic.Int32
	peak     *atomic.Int32
}

func (c *countingFFmpeg) CheckAvailable() error {
	return nil
}

func (c *countingFFmpeg) Run(_ ...string) error {
	n := c.inFlight.Add(1)
	defer c.inFlight.Add(-1)

	for {
		cur := c.peak.Load()
		if n <= cur || c.peak.CompareAndSwap(cur, n) {
			break
		}
	}

	time.Sleep(5 * time.Millisecond)

	return nil
}

func (c *countingFFmpeg) RunWithStdin(_ io.Reader, _ ...string) error {
	return c.Run()
}

func makeBatch() domain.Batch {
	img := domain.ImageConfig{
		Width: 64, Height: 64, Background: "solid", Color: "black",
		Overlay: domain.Overlay{TopLeft: "", TopRight: "", Center: "", BottomLeft: "", BottomRight: ""},
		Scale:   1, Output: "a.png",
	}
	vid := domain.VideoConfig{
		Width: 64, Height: 64, FPS: 30, Duration: "1",
		Background: "solid", Color: "black",
		Overlay: domain.Overlay{TopLeft: "", TopRight: "", Center: "", BottomLeft: "", BottomRight: ""},
		Scale:   1, Output: "b.mp4",
		Audio: domain.AudioSilence, SampleRate: 48000, Channels: 2, Frequency: 440,
	}
	aud := domain.AudioConfig{
		Type: domain.AudioSilence, Duration: "1", SampleRate: 48000, Channels: 2,
		Frequency: 440, Output: "c.wav",
	}

	return domain.Batch{Jobs: []domain.Job{
		{Type: domain.JobTypeImage, Image: &img, Video: nil, Audio: nil},
		{Type: domain.JobTypeVideo, Image: nil, Video: &vid, Audio: nil},
		{Type: domain.JobTypeAudio, Image: nil, Video: nil, Audio: &aud},
	}}
}

func TestBatchUseCase_Generate_AllSucceed(t *testing.T) {
	t.Parallel()

	ffmpeg := &mockFFmpeg{}
	renderer := newMockRenderer()
	uc := usecase.NewBatchUseCase(
		usecase.NewImageUseCase(renderer),
		usecase.NewVideoUseCase(ffmpeg, renderer),
		usecase.NewAudioUseCase(ffmpeg),
	)

	results := uc.Generate(context.Background(), makeBatch(), usecase.BatchOptions{Parallel: 2, FailFast: false})

	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}

	for _, r := range results {
		if r.Err != nil {
			t.Errorf("job %d (%s) failed: %v", r.Index, r.Type, r.Err)
		}
	}

	wantOutputs := []string{"a.png", "b.mp4", "c.wav"}
	for i, want := range wantOutputs {
		if results[i].Output != want {
			t.Errorf("results[%d].Output = %q, want %q", i, results[i].Output, want)
		}
	}
}

func TestBatchUseCase_Generate_FailFast(t *testing.T) {
	t.Parallel()

	failingErr := errors.New("boom")
	ffmpeg := &mockFFmpeg{runErr: failingErr, runWithStdinErr: failingErr}
	renderer := newMockRenderer()
	renderer.writePNGErr = failingErr
	uc := usecase.NewBatchUseCase(
		usecase.NewImageUseCase(renderer),
		usecase.NewVideoUseCase(ffmpeg, renderer),
		usecase.NewAudioUseCase(ffmpeg),
	)

	results := uc.Generate(context.Background(), makeBatch(), usecase.BatchOptions{Parallel: 1, FailFast: true})

	if results[0].Err == nil {
		t.Fatal("expected first job to fail")
	}

	later := 0

	for _, r := range results[1:] {
		if r.Err != nil {
			later++
		}
	}

	if later == 0 {
		t.Error("expected at least one later job to be skipped/failed under fail-fast")
	}
}

func TestBatchUseCase_Generate_ParallelBounded(t *testing.T) {
	t.Parallel()

	renderer := newMockRenderer()

	var inFlight atomic.Int32

	var peak atomic.Int32

	ffmpeg := &countingFFmpeg{inFlight: &inFlight, peak: &peak}

	jobs := make([]domain.Job, 0, 6)

	for range 6 {
		aud := domain.AudioConfig{
			Type: domain.AudioSilence, Duration: "1", SampleRate: 48000, Channels: 2,
			Frequency: 440, Output: "x.wav",
		}
		jobs = append(jobs, domain.Job{Type: domain.JobTypeAudio, Image: nil, Video: nil, Audio: &aud})
	}

	uc := usecase.NewBatchUseCase(
		usecase.NewImageUseCase(renderer),
		usecase.NewVideoUseCase(ffmpeg, renderer),
		usecase.NewAudioUseCase(ffmpeg),
	)

	_ = uc.Generate(context.Background(), domain.Batch{Jobs: jobs}, usecase.BatchOptions{Parallel: 2, FailFast: false})

	if got := peak.Load(); got > 2 {
		t.Errorf("peak concurrency = %d, want <= 2", got)
	}
}
