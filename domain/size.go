package domain

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ErrInvalidTargetSize indicates a malformed target file size value.
var ErrInvalidTargetSize = errors.New("invalid target size")

// ErrTargetSizeUnsupported indicates a target size was requested for an output
// format whose size no encoder setting can steer.
var ErrTargetSizeUnsupported = errors.New("target size not supported for this format")

// ErrTargetSizeTooSmall indicates the requested size leaves no bits for the
// stream once the fixed overheads are subtracted.
var ErrTargetSizeTooSmall = errors.New("target size too small")

// ErrTargetSizeConflict indicates a target size was combined with a manual
// rate-control setting that steers the same dial.
var ErrTargetSizeConflict = errors.New("target size conflicts with manual rate control")

const (
	bytesPerKiB = 1024
	bytesPerMiB = bytesPerKiB * bytesPerKiB
	bytesPerGiB = bytesPerMiB * bytesPerKiB
	bitsPerByte = 8
	// containerOverheadRatio is the share of the target size handed to the
	// media streams; the rest covers container headers and muxing overhead,
	// which an encoder bitrate does not account for.
	containerOverheadRatio = 0.99
)

// ParseFileSize parses a human-readable file size such as "10MB", "500k",
// "1.5G" or a plain byte count into bytes. Suffixes are binary (1K is 1024
// bytes), and the "B"/"iB" spellings mean the same thing: 10K, 10KB and 10KiB
// are all 10240 bytes.
func ParseFileSize(s string) (int64, error) {
	num := strings.ToLower(strings.TrimSpace(s))
	num = strings.TrimSuffix(strings.TrimSuffix(num, "ib"), "b")

	multiplier := int64(1)

	for _, unit := range []struct {
		suffix string
		size   int64
	}{{suffix: "k", size: bytesPerKiB}, {suffix: "m", size: bytesPerMiB}, {suffix: "g", size: bytesPerGiB}} {
		if rest, ok := strings.CutSuffix(num, unit.suffix); ok {
			num, multiplier = rest, unit.size

			break
		}
	}

	value, err := strconv.ParseFloat(num, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value <= 0 {
		return 0, fmt.Errorf(
			"%w: %q (expected a positive number with an optional K/M/G suffix, e.g. 500KB or 10MB)",
			ErrInvalidTargetSize, s,
		)
	}

	total := value * float64(multiplier)
	if total >= math.MaxInt64 {
		return 0, fmt.Errorf("%w: %q (larger than any file can be)", ErrInvalidTargetSize, s)
	}

	size := int64(total)
	if size <= 0 {
		return 0, fmt.Errorf("%w: %q (rounds down to zero bytes)", ErrInvalidTargetSize, s)
	}

	return size, nil
}

// BitrateForSize returns the average bitrate, in bits per second, that fills
// targetBytes over durationSec once reservedBitrate (bits per second spent on
// the other streams, typically audio) and a margin for container overhead are
// subtracted. It returns 0 when nothing is left for the stream.
func BitrateForSize(targetBytes int64, durationSec float64, reservedBitrate int) int {
	if targetBytes <= 0 || durationSec <= 0 {
		return 0
	}

	available := float64(targetBytes) * bitsPerByte * containerOverheadRatio / durationSec

	bitrate := int(available) - reservedBitrate
	if bitrate <= 0 {
		return 0
	}

	return bitrate
}

// FormatFileSize renders a byte count in the same binary units ParseFileSize
// accepts, so messages and flag values speak the same language.
func FormatFileSize(size int64) string {
	switch {
	case size >= bytesPerGiB:
		return fmt.Sprintf("%.2fGB", float64(size)/bytesPerGiB)
	case size >= bytesPerMiB:
		return fmt.Sprintf("%.2fMB", float64(size)/bytesPerMiB)
	case size >= bytesPerKiB:
		return fmt.Sprintf("%.2fKB", float64(size)/bytesPerKiB)
	default:
		return strconv.FormatInt(size, 10) + "B"
	}
}
