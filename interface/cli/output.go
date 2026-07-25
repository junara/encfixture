package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// errOutputExists reports a refusal to overwrite an existing file under --no-clobber.
var errOutputExists = errors.New("output file already exists")

// checkClobber refuses to overwrite path when noClobber is set. By default
// generation commands overwrite existing files, matching their fixture
// regeneration purpose; --no-clobber opts into failing instead.
func checkClobber(path string, noClobber bool) error {
	if !noClobber {
		return nil
	}

	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("%w: %s (remove it or drop --no-clobber to overwrite)", errOutputExists, path)
	}

	return nil
}

type result struct {
	Status   string `json:"status"`
	File     string `json:"file"`
	Type     string `json:"type"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	FPS      int    `json:"fps,omitempty"`
	Duration string `json:"duration,omitempty"`
}

type errorOutput struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// jsonWritten tracks whether a JSON document was already emitted to stdout, so
// Execute does not append a second one after e.g. a batch summary with failures.
var jsonWritten bool

// writeJSON encodes v to stdout as a single JSON document.
func writeJSON(v any) {
	jsonWritten = true

	if err := json.NewEncoder(os.Stdout).Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "json encode error: %v\n", err)
	}
}

// writeJSONError reports err on stdout in JSON form when --json is active and
// nothing has been emitted yet, so machine consumers never get an empty stdout.
func writeJSONError(err error) {
	if !jsonOutput || jsonWritten {
		return
	}

	writeJSON(errorOutput{Status: "error", Error: err.Error()})
}

func printResult(r result) {
	if jsonOutput {
		writeJSON(r)

		return
	}

	fmt.Fprintf(os.Stderr, "Created: %s\n", r.File)
}
