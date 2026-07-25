package domain_test

import (
	"errors"
	"testing"

	"github.com/junara/encfixture/domain"
)

func TestParseExpectation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    domain.Expectation
		wantErr bool
	}{
		{
			name:  "string field",
			input: "codec=h264",
			want:  domain.Expectation{Field: "codec", Value: "h264", Tolerance: 0},
		},
		{
			name:  "integer field",
			input: "width=1920",
			want:  domain.Expectation{Field: "width", Value: "1920", Tolerance: 0},
		},
		{
			name:  "kebab-case alias",
			input: "pix-fmt=yuv420p",
			want:  domain.Expectation{Field: "pixFmt", Value: "yuv420p", Tolerance: 0},
		},
		{
			name:  "case-insensitive key",
			input: "PixFmt=yuv420p",
			want:  domain.Expectation{Field: "pixFmt", Value: "yuv420p", Tolerance: 0},
		},
		{
			name:  "lowercase alias",
			input: "audiocodec=aac",
			want:  domain.Expectation{Field: "audioCodec", Value: "aac", Tolerance: 0},
		},
		{
			name:  "surrounding whitespace",
			input: " codec = h264 ",
			want:  domain.Expectation{Field: "codec", Value: "h264", Tolerance: 0},
		},
		{
			name:  "sample-rate alias",
			input: "sample-rate=48000",
			want:  domain.Expectation{Field: "sampleRate", Value: "48000", Tolerance: 0},
		},
		{
			name:  "duration with ascii tolerance",
			input: "duration=5+-0.2",
			want:  domain.Expectation{Field: "duration", Value: "5", Tolerance: 0.2},
		},
		{
			name:  "duration with unicode tolerance",
			input: "duration=5±0.2",
			want:  domain.Expectation{Field: "duration", Value: "5", Tolerance: 0.2},
		},
		{
			name:    "unknown field",
			input:   "nope=1",
			wantErr: true,
		},
		{
			name:    "missing value",
			input:   "codec=",
			wantErr: true,
		},
		{
			name:    "missing separator",
			input:   "codec",
			wantErr: true,
		},
		{
			name:    "non-numeric width",
			input:   "width=wide",
			wantErr: true,
		},
		{
			name:    "negative tolerance",
			input:   "duration=5+--1",
			wantErr: true,
		},
		{
			name:  "string field keeps tolerance-like text verbatim",
			input: "codec=h264+-0.5",
			want:  domain.Expectation{Field: "codec", Value: "h264+-0.5", Tolerance: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.ParseExpectation(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseExpectation(%q) error = nil, want error", tt.input)
				}

				if !errors.Is(err, domain.ErrInvalidExpectation) {
					t.Errorf("error = %v, want ErrInvalidExpectation", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("ParseExpectation(%q) error = %v", tt.input, err)
			}

			if got != tt.want {
				t.Errorf("ParseExpectation(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func testMediaInfo() domain.MediaInfo {
	return domain.MediaInfo{
		Format: domain.FormatInfo{FormatName: "mov,mp4", Duration: "5.016000", Size: "1000", BitRate: "1"},
		Streams: []domain.StreamInfo{
			{Index: 0, Type: domain.StreamTypeVideo, Codec: "h264", Width: 1920, Height: 1080, FPS: "29.970", PixFmt: "yuv420p"},
			{Index: 1, Type: domain.StreamTypeAudio, Codec: "aac", SampleRate: "48000", Channels: 2},
		},
	}
}

func TestEvaluateExpectations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		exp      domain.Expectation
		wantPass bool
	}{
		{name: "codec match", exp: domain.Expectation{Field: "codec", Value: "h264"}, wantPass: true},
		{name: "codec mismatch", exp: domain.Expectation{Field: "codec", Value: "hevc"}, wantPass: false},
		{name: "width match", exp: domain.Expectation{Field: "width", Value: "1920"}, wantPass: true},
		{name: "height mismatch", exp: domain.Expectation{Field: "height", Value: "720"}, wantPass: false},
		{name: "fps ntsc match", exp: domain.Expectation{Field: "fps", Value: "29.97"}, wantPass: true},
		{name: "fps integer mismatch", exp: domain.Expectation{Field: "fps", Value: "30"}, wantPass: false},
		{name: "duration within default tolerance", exp: domain.Expectation{Field: "duration", Value: "5"}, wantPass: true},
		{name: "duration outside default tolerance", exp: domain.Expectation{Field: "duration", Value: "4.5"}, wantPass: false},
		{name: "duration explicit tolerance", exp: domain.Expectation{Field: "duration", Value: "4.5", Tolerance: 0.6}, wantPass: true},
		{name: "pixFmt match", exp: domain.Expectation{Field: "pixFmt", Value: "yuv420p"}, wantPass: true},
		{name: "audioCodec match", exp: domain.Expectation{Field: "audioCodec", Value: "aac"}, wantPass: true},
		{name: "sampleRate match", exp: domain.Expectation{Field: "sampleRate", Value: "48000"}, wantPass: true},
		{name: "channels mismatch", exp: domain.Expectation{Field: "channels", Value: "6"}, wantPass: false},
	}

	info := testMediaInfo()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			results := domain.EvaluateExpectations(info, []domain.Expectation{tt.exp})
			if len(results) != 1 {
				t.Fatalf("len(results) = %d, want 1", len(results))
			}

			if results[0].Pass != tt.wantPass {
				t.Errorf("%s: Pass = %v, want %v (expected %s, actual %s)",
					tt.name, results[0].Pass, tt.wantPass, results[0].Expected, results[0].Actual)
			}
		})
	}
}

func TestEvaluateExpectations_MissingStreams(t *testing.T) {
	t.Parallel()

	info := domain.MediaInfo{
		Format:  domain.FormatInfo{FormatName: "wav", Duration: "3.0", Size: "10", BitRate: "1"},
		Streams: []domain.StreamInfo{{Index: 0, Type: domain.StreamTypeAudio, Codec: "pcm_s16le", SampleRate: "48000", Channels: 2}},
	}

	results := domain.EvaluateExpectations(info, []domain.Expectation{
		{Field: "codec", Value: "h264", Tolerance: 0},
	})

	if results[0].Pass {
		t.Error("codec expectation should fail without a video stream")
	}

	if results[0].Actual != "(no video stream)" {
		t.Errorf("Actual = %q, want %q", results[0].Actual, "(no video stream)")
	}
}

func TestEvaluateExpectations_ExpectedDisplayShowsTolerance(t *testing.T) {
	t.Parallel()

	results := domain.EvaluateExpectations(testMediaInfo(), []domain.Expectation{
		{Field: "duration", Value: "9", Tolerance: 0.5},
	})

	if results[0].Expected != "9±0.5" {
		t.Errorf("Expected = %q, want %q", results[0].Expected, "9±0.5")
	}
}
