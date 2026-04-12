package domain_test

import (
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
