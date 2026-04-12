package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

type result struct {
	Status   string `json:"status"`
	File     string `json:"file"`
	Type     string `json:"type"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	FPS      int    `json:"fps,omitempty"`
	Duration string `json:"duration,omitempty"`
}

func printResult(r result) {
	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)

		if err := enc.Encode(r); err != nil {
			fmt.Fprintf(os.Stderr, "json encode error: %v\n", err)
		}

		return
	}

	fmt.Fprintf(os.Stderr, "Created: %s\n", r.File)
}
