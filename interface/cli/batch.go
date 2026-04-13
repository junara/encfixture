package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/junara/encfixture/domain"
	"github.com/junara/encfixture/infrastructure"
	"github.com/junara/encfixture/usecase"

	"github.com/spf13/cobra"
)

var errBatchHadFailures = errors.New("batch had failures")

var batchCmd = &cobra.Command{
	Use:   "batch <file.json>",
	Short: "Run multiple media generation jobs defined in a JSON file",
	Long: `Load a JSON batch file and run each job, optionally in parallel.

Each job is an object with "type" (image, video, or audio), "output", and any
flags the matching subcommand accepts (width, height, fps, duration, tl, tr,
center, bl, br, bg, color, scale, audio, sampleRate, channels, frequency).

A top-level "defaults" object applies to every job; individual jobs override it.`,
	Example: `  # Run all jobs with default parallelism (NumCPU/2)
  encfixture batch jobs.json

  # Cap concurrency and stop on first failure
  encfixture batch jobs.json --parallel 4 --fail-fast

  # Machine-readable output for CI
  encfixture batch jobs.json --json`,
	Args: cobra.ExactArgs(1),
	RunE: runBatch,
}

func init() {
	rootCmd.AddCommand(batchCmd)

	batchCmd.Flags().IntP("parallel", "p", 0, "Max concurrent jobs (default: NumCPU/2, minimum 1)")
	batchCmd.Flags().Bool("fail-fast", false, "Stop scheduling new jobs after the first failure")
}

func runBatch(cmd *cobra.Command, args []string) error {
	path := args[0]
	parallel, _ := cmd.Flags().GetInt("parallel")
	failFast, _ := cmd.Flags().GetBool("fail-fast")

	if parallel <= 0 {
		parallel = max(runtime.NumCPU()/2, 1)
	}

	batch, err := infrastructure.LoadBatch(path)
	if err != nil {
		return fmt.Errorf("load batch: %w", err)
	}

	ffmpeg := infrastructure.NewFFmpeg()
	renderer := infrastructure.NewImageRenderer()
	uc := usecase.NewBatchUseCase(
		usecase.NewImageUseCase(renderer),
		usecase.NewVideoUseCase(ffmpeg, renderer),
		usecase.NewAudioUseCase(ffmpeg),
	)

	results := uc.Generate(cmd.Context(), batch, usecase.BatchOptions{
		Parallel: parallel,
		FailFast: failFast,
	})

	succeeded, failed := printBatchResults(results)
	if failed > 0 {
		return fmt.Errorf("%w: %d failed, %d succeeded", errBatchHadFailures, failed, succeeded)
	}

	return nil
}

type batchJobResult struct {
	Index  int    `json:"index"`
	Type   string `json:"type"`
	File   string `json:"file"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type batchOutput struct {
	Results   []batchJobResult `json:"results"`
	Succeeded int              `json:"succeeded"`
	Failed    int              `json:"failed"`
}

func printBatchResults(results []domain.JobResult) (int, int) {
	succeeded, failed := 0, 0

	for _, r := range results {
		if r.Err != nil {
			failed++
		} else {
			succeeded++
		}
	}

	if jsonOutput {
		emitBatchJSON(results, succeeded, failed)

		return succeeded, failed
	}

	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "[%d] FAIL  %s  %v\n", r.Index, r.Output, r.Err)
		} else {
			fmt.Fprintf(os.Stderr, "[%d] ok    %s\n", r.Index, r.Output)
		}
	}

	fmt.Fprintf(os.Stderr, "%d succeeded, %d failed\n", succeeded, failed)

	return succeeded, failed
}

func emitBatchJSON(results []domain.JobResult, succeeded, failed int) {
	out := batchOutput{
		Results:   make([]batchJobResult, len(results)),
		Succeeded: succeeded,
		Failed:    failed,
	}

	for i, r := range results {
		status := "ok"
		errStr := ""

		if r.Err != nil {
			status = "error"
			errStr = r.Err.Error()
		}

		out.Results[i] = batchJobResult{
			Index:  r.Index,
			Type:   string(r.Type),
			File:   r.Output,
			Status: status,
			Error:  errStr,
		}
	}

	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "json encode error: %v\n", err)
	}
}
