package infrastructure

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/junara/encfixture/domain"
)

// versionToken is the index of the version number in a tool's banner line,
// e.g. "ffmpeg version 7.1 Copyright ...".
const versionToken = 2

// Inspector probes the local ffmpeg installation for the doctor command.
type Inspector struct{}

// NewInspector creates a new Inspector instance.
func NewInspector() *Inspector {
	return &Inspector{}
}

// ToolStatus reports whether the named tool is in PATH, along with its
// resolved path and version.
func (i *Inspector) ToolStatus(name string) domain.ToolStatus {
	status := domain.ToolStatus{Name: name, Available: false, Version: "", Path: ""}

	path, err := exec.LookPath(name)
	if err != nil {
		return status
	}

	status.Available = true
	status.Path = path
	status.Version = toolVersion(name)

	return status
}

// toolVersion extracts the version number from "<tool> -version". A parse
// failure returns "" rather than an error: the tool is present, we just could
// not name its version.
func toolVersion(name string) string {
	out, err := exec.CommandContext(context.Background(), name, "-version").Output()
	if err != nil {
		return ""
	}

	line, _, _ := strings.Cut(string(out), "\n")

	fields := strings.Fields(line)
	if len(fields) > versionToken && fields[versionToken-1] == "version" {
		return fields[versionToken]
	}

	return ""
}

// Encoders returns the set of encoder names supported by the local ffmpeg
// build, as listed by "ffmpeg -encoders".
func (i *Inspector) Encoders() (map[string]bool, error) {
	out, err := exec.CommandContext(context.Background(), "ffmpeg", "-hide_banner", "-encoders").Output()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg -encoders failed: %w", err)
	}

	encoders := make(map[string]bool)
	inList := false

	for line := range strings.SplitSeq(string(out), "\n") {
		// The encoder table starts after a "------" separator line; everything
		// before it is the flag legend.
		if !inList {
			inList = strings.HasPrefix(strings.TrimSpace(line), "------")

			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			encoders[fields[1]] = true
		}
	}

	return encoders, nil
}
