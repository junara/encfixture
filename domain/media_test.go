package domain_test

import (
	"errors"
	"testing"

	"github.com/junara/encfixture/domain"
)

func TestOverlay_HasContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		overlay domain.Overlay
		want    bool
	}{
		{
			name:    "empty overlay",
			overlay: domain.Overlay{},
			want:    false,
		},
		{
			name:    "top-left only",
			overlay: domain.Overlay{TopLeft: "frame"},
			want:    true,
		},
		{
			name:    "center only",
			overlay: domain.Overlay{Center: "text"},
			want:    true,
		},
		{
			name:    "all positions",
			overlay: domain.Overlay{TopLeft: "a", TopRight: "b", Center: "c", BottomLeft: "d", BottomRight: "e"},
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.overlay.HasContent(); got != tt.want {
				t.Errorf("HasContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOverlay_HasDynamicContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		overlay domain.Overlay
		want    bool
	}{
		{
			name:    "empty overlay",
			overlay: domain.Overlay{},
			want:    false,
		},
		{
			name:    "static text only",
			overlay: domain.Overlay{Center: "hello"},
			want:    false,
		},
		{
			name:    "filename keyword is not dynamic",
			overlay: domain.Overlay{BottomLeft: "filename"},
			want:    false,
		},
		{
			name:    "frame keyword",
			overlay: domain.Overlay{TopLeft: "frame"},
			want:    true,
		},
		{
			name:    "timecode keyword",
			overlay: domain.Overlay{TopRight: "timecode"},
			want:    true,
		},
		{
			name:    "mixed static and dynamic",
			overlay: domain.Overlay{TopLeft: "frame", Center: "hello"},
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.overlay.HasDynamicContent(); got != tt.want {
				t.Errorf("HasDynamicContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOverlay_All(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		overlay domain.Overlay
		want    int
	}{
		{
			name:    "empty",
			overlay: domain.Overlay{},
			want:    0,
		},
		{
			name:    "one position",
			overlay: domain.Overlay{Center: "text"},
			want:    1,
		},
		{
			name:    "all positions",
			overlay: domain.Overlay{TopLeft: "a", TopRight: "b", Center: "c", BottomLeft: "d", BottomRight: "e"},
			want:    5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.overlay.All()
			if len(got) != tt.want {
				t.Errorf("All() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}

func TestOverlay_Entries(t *testing.T) {
	t.Parallel()

	overlay := domain.Overlay{
		TopLeft:     "frame",
		TopRight:    "timecode",
		BottomLeft:  "filename",
		BottomRight: "ID-001",
	}

	entries := overlay.Entries()

	if len(entries) != 4 {
		t.Fatalf("Entries() returned %d items, want 4", len(entries))
	}

	expected := []struct {
		pos     domain.TextPosition
		content string
	}{
		{domain.PositionTopLeft, "frame"},
		{domain.PositionTopRight, "timecode"},
		{domain.PositionBottomLeft, "filename"},
		{domain.PositionBottomRight, "ID-001"},
	}

	for i, exp := range expected {
		if entries[i].Position != exp.pos {
			t.Errorf("Entries()[%d].Position = %v, want %v", i, entries[i].Position, exp.pos)
		}

		if entries[i].Content != exp.content {
			t.Errorf("Entries()[%d].Content = %q, want %q", i, entries[i].Content, exp.content)
		}
	}
}

func TestValidateDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		duration string
		wantErr  bool
	}{
		{name: "integer seconds", duration: "10", wantErr: false},
		{name: "fractional seconds", duration: "0.5", wantErr: false},
		{name: "not a number", duration: "abc", wantErr: true},
		{name: "empty", duration: "", wantErr: true},
		{name: "zero", duration: "0", wantErr: true},
		{name: "negative", duration: "-3", wantErr: true},
		{name: "infinity", duration: "Inf", wantErr: true},
		{name: "nan", duration: "NaN", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := domain.ValidateDuration(tt.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDuration(%q) error = %v, wantErr %v", tt.duration, err, tt.wantErr)
			}

			if err != nil && !errors.Is(err, domain.ErrInvalidDuration) {
				t.Errorf("ValidateDuration(%q) error = %v, want ErrInvalidDuration", tt.duration, err)
			}
		})
	}
}

func TestValidateBitrate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		bitrate string
		wantErr bool
	}{
		{name: "empty means encoder default", bitrate: "", wantErr: false},
		{name: "kilobits lowercase", bitrate: "800k", wantErr: false},
		{name: "megabits uppercase", bitrate: "5M", wantErr: false},
		{name: "plain number", bitrate: "5000000", wantErr: false},
		{name: "fractional with suffix", bitrate: "1.5M", wantErr: false},
		{name: "not a number", bitrate: "fast", wantErr: true},
		{name: "suffix only", bitrate: "M", wantErr: true},
		{name: "zero", bitrate: "0", wantErr: true},
		{name: "negative", bitrate: "-5M", wantErr: true},
		{name: "unknown suffix", bitrate: "5G", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := domain.ValidateBitrate(tt.bitrate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBitrate(%q) error = %v, wantErr %v", tt.bitrate, err, tt.wantErr)
			}

			if err != nil && !errors.Is(err, domain.ErrInvalidBitrate) {
				t.Errorf("ValidateBitrate(%q) error = %v, want ErrInvalidBitrate", tt.bitrate, err)
			}
		})
	}
}
