package domain_test

import (
	"errors"
	"testing"

	"github.com/junara/encfixture/domain"
)

func TestParseFileSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{name: "plain bytes", input: "1024", want: 1024, wantErr: false},
		{name: "explicit bytes", input: "512B", want: 512, wantErr: false},
		{name: "kilobytes", input: "10KB", want: 10240, wantErr: false},
		{name: "bare k suffix", input: "10k", want: 10240, wantErr: false},
		{name: "kibibytes spelling", input: "10KiB", want: 10240, wantErr: false},
		{name: "megabytes", input: "10MB", want: 10485760, wantErr: false},
		{name: "lowercase megabytes", input: "10m", want: 10485760, wantErr: false},
		{name: "fractional gigabytes", input: "1.5G", want: 1610612736, wantErr: false},
		{name: "surrounding spaces", input: " 2MB ", want: 2097152, wantErr: false},
		{name: "empty", input: "", want: 0, wantErr: true},
		{name: "zero", input: "0", want: 0, wantErr: true},
		{name: "negative", input: "-5MB", want: 0, wantErr: true},
		{name: "not a number", input: "big", want: 0, wantErr: true},
		{name: "unknown suffix", input: "10TB", want: 0, wantErr: true},
		{name: "rounds down to zero", input: "0.4B", want: 0, wantErr: true},
		{name: "beyond int64", input: "9999999999G", want: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.ParseFileSize(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseFileSize(%q) error = nil, want error", tt.input)
				}

				if !errors.Is(err, domain.ErrInvalidTargetSize) {
					t.Errorf("ParseFileSize(%q) error = %v, want ErrInvalidTargetSize", tt.input, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("ParseFileSize(%q) error = %v", tt.input, err)
			}

			if got != tt.want {
				t.Errorf("ParseFileSize(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestBitrateForSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		bytes    int64
		duration float64
		reserved int
		want     int
	}{
		// 1MB over 10s is 830472 bits/s after the 1% container margin.
		{name: "no reservation", bytes: 1048576, duration: 10, reserved: 0, want: 830472},
		{name: "audio reserved", bytes: 1048576, duration: 10, reserved: 128000, want: 702472},
		{name: "reservation exceeds budget", bytes: 1024, duration: 10, reserved: 128000, want: 0},
		{name: "zero duration", bytes: 1048576, duration: 0, reserved: 0, want: 0},
		{name: "zero size", bytes: 0, duration: 10, reserved: 0, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := domain.BitrateForSize(tt.bytes, tt.duration, tt.reserved); got != tt.want {
				t.Errorf("BitrateForSize(%d, %v, %d) = %d, want %d", tt.bytes, tt.duration, tt.reserved, got, tt.want)
			}
		})
	}
}

func TestFormatFileSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{name: "bytes", bytes: 512, want: "512B"},
		{name: "kilobytes", bytes: 10240, want: "10.00KB"},
		{name: "megabytes", bytes: 1048576, want: "1.00MB"},
		{name: "gigabytes", bytes: 1610612736, want: "1.50GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := domain.FormatFileSize(tt.bytes); got != tt.want {
				t.Errorf("FormatFileSize(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}
