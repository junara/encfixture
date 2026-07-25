package domain

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ErrInvalidExpectation indicates a malformed --expect assertion.
var ErrInvalidExpectation = errors.New("invalid expectation")

// Expectation field name constants, matching the verify JSON output keys.
const (
	ExpectCodec      = "codec"
	ExpectWidth      = "width"
	ExpectHeight     = "height"
	ExpectFPS        = "fps"
	ExpectPixFmt     = "pixFmt"
	ExpectDuration   = "duration"
	ExpectAudioCodec = "audioCodec"
	ExpectSampleRate = "sampleRate"
	ExpectChannels   = "channels"
)

const (
	// defaultDurationTolerance absorbs the container-level rounding between the
	// requested and the actual duration of generated files, in seconds.
	defaultDurationTolerance = 0.1
	// defaultRateTolerance absorbs float formatting noise when comparing frame
	// rates such as 29.97 against ffprobe's "29.970".
	defaultRateTolerance = 0.001
)

// Expectation is one assertion against a probed media file.
type Expectation struct {
	Field     string
	Value     string
	Tolerance float64 // numeric fields only; 0 means the field's default
}

// CheckResult is the outcome of evaluating one Expectation.
type CheckResult struct {
	Field    string `json:"field"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Pass     bool   `json:"pass"`
}

// expectationFields maps accepted key spellings to the canonical field name.
// Keys are matched lowercased, which makes the lookup case-insensitive and
// covers both the canonical camelCase and the kebab-case used by CLI flags.
func expectationFields() map[string]string {
	return map[string]string{
		"codec":       ExpectCodec,
		"width":       ExpectWidth,
		"height":      ExpectHeight,
		"fps":         ExpectFPS,
		"pixfmt":      ExpectPixFmt,
		"pix-fmt":     ExpectPixFmt,
		"duration":    ExpectDuration,
		"audiocodec":  ExpectAudioCodec,
		"audio-codec": ExpectAudioCodec,
		"samplerate":  ExpectSampleRate,
		"sample-rate": ExpectSampleRate,
		"channels":    ExpectChannels,
	}
}

// numericExpectation reports whether the field compares numerically, which
// also makes a "value±tolerance" spelling valid for it.
func numericExpectation(field string) bool {
	switch field {
	case ExpectWidth, ExpectHeight, ExpectFPS, ExpectDuration, ExpectSampleRate, ExpectChannels:
		return true
	default:
		return false
	}
}

// ParseExpectation parses a "key=value" assertion. Keys are case-insensitive
// and surrounding whitespace is ignored. Numeric fields accept a tolerance
// suffix, written "5±0.2" or "5+-0.2".
func ParseExpectation(s string) (Expectation, error) {
	var empty Expectation

	key, value, found := strings.Cut(s, "=")
	key, value = strings.TrimSpace(key), strings.TrimSpace(value)

	if !found || key == "" || value == "" {
		return empty, fmt.Errorf("%w: %q (want key=value)", ErrInvalidExpectation, s)
	}

	field, known := expectationFields()[strings.ToLower(key)]
	if !known {
		return empty, fmt.Errorf(
			"%w: unknown field %q (supported: codec, width, height, fps, pixFmt, duration, audioCodec, sampleRate, channels)",
			ErrInvalidExpectation, key,
		)
	}

	exp := Expectation{Field: field, Value: value, Tolerance: 0}

	if numericExpectation(field) {
		parsed, err := parseNumericValue(field, value)
		if err != nil {
			return empty, err
		}

		exp = parsed
	}

	return exp, nil
}

// parseNumericValue validates a numeric expectation value and splits off an
// optional "±tol" / "+-tol" tolerance suffix.
func parseNumericValue(field, value string) (Expectation, error) {
	var empty Expectation

	number, tolerance := value, ""

	for _, sep := range []string{"±", "+-"} {
		if n, t, found := strings.Cut(value, sep); found {
			number, tolerance = n, t

			break
		}
	}

	exp := Expectation{Field: field, Value: number, Tolerance: 0}

	if !validFloat(number) {
		return empty, fmt.Errorf("%w: %s=%q (want a number)", ErrInvalidExpectation, field, value)
	}

	if tolerance != "" {
		tol, err := strconv.ParseFloat(tolerance, 64)
		if err != nil || math.IsNaN(tol) || math.IsInf(tol, 0) || tol < 0 {
			return empty, fmt.Errorf("%w: %s=%q (want a non-negative tolerance)", ErrInvalidExpectation, field, value)
		}

		exp.Tolerance = tol
	}

	return exp, nil
}

func validFloat(s string) bool {
	v, err := strconv.ParseFloat(s, 64)

	return err == nil && !math.IsNaN(v) && !math.IsInf(v, 0)
}

// EvaluateExpectations checks each expectation against info and returns one
// result per expectation, in input order.
func EvaluateExpectations(info MediaInfo, exps []Expectation) []CheckResult {
	results := make([]CheckResult, 0, len(exps))
	for _, exp := range exps {
		results = append(results, evaluateExpectation(info, exp))
	}

	return results
}

func evaluateExpectation(info MediaInfo, exp Expectation) CheckResult {
	actual, ok := actualValue(info, exp.Field)
	if !ok {
		return CheckResult{Field: exp.Field, Expected: expectedDisplay(exp), Actual: actual, Pass: false}
	}

	return CheckResult{
		Field:    exp.Field,
		Expected: expectedDisplay(exp),
		Actual:   actual,
		Pass:     expectationHolds(exp, actual),
	}
}

// expectedDisplay renders the expected value with its tolerance when one
// applies, so a failed check shows the full assertion.
func expectedDisplay(exp Expectation) string {
	if exp.Tolerance > 0 {
		return fmt.Sprintf("%s±%g", exp.Value, exp.Tolerance)
	}

	return exp.Value
}

// actualValue extracts the probed value for field. The second return value is
// false when the file lacks the stream or property the field refers to; the
// first then carries the reason for display.
func actualValue(info MediaInfo, field string) (string, bool) {
	switch field {
	case ExpectCodec, ExpectWidth, ExpectHeight, ExpectFPS, ExpectPixFmt:
		return videoStreamValue(info, field)
	case ExpectAudioCodec, ExpectSampleRate, ExpectChannels:
		return audioStreamValue(info, field)
	case ExpectDuration:
		if info.Format.Duration == "" {
			return "(unknown duration)", false
		}

		return info.Format.Duration, true
	default:
		return "(unsupported field)", false
	}
}

func videoStreamValue(info MediaInfo, field string) (string, bool) {
	stream, found := firstStream(info, StreamTypeVideo)
	if !found {
		return "(no video stream)", false
	}

	switch field {
	case ExpectCodec:
		return stream.Codec, true
	case ExpectWidth:
		return strconv.Itoa(stream.Width), true
	case ExpectHeight:
		return strconv.Itoa(stream.Height), true
	case ExpectFPS:
		if stream.FPS == "" {
			return "(unknown fps)", false
		}

		return stream.FPS, true
	default:
		return stream.PixFmt, true
	}
}

func audioStreamValue(info MediaInfo, field string) (string, bool) {
	stream, found := firstStream(info, StreamTypeAudio)
	if !found {
		return "(no audio stream)", false
	}

	switch field {
	case ExpectAudioCodec:
		return stream.Codec, true
	case ExpectSampleRate:
		return stream.SampleRate, true
	default:
		return strconv.Itoa(stream.Channels), true
	}
}

func firstStream(info MediaInfo, streamType string) (StreamInfo, bool) {
	var empty StreamInfo

	for _, stream := range info.Streams {
		if stream.Type == streamType {
			return stream, true
		}
	}

	return empty, false
}

// expectationHolds compares the actual value against the expectation. Numeric
// fields compare within tolerance; string fields compare exactly.
func expectationHolds(exp Expectation, actual string) bool {
	if !numericExpectation(exp.Field) {
		return exp.Value == actual
	}

	expected, expErr := strconv.ParseFloat(exp.Value, 64)
	got, gotErr := strconv.ParseFloat(actual, 64)

	if expErr != nil || gotErr != nil {
		return exp.Value == actual
	}

	return math.Abs(expected-got) <= toleranceFor(exp)
}

// toleranceFor returns the explicit tolerance, or the field's default: a small
// epsilon for rates, a container-rounding allowance for duration, and exact
// match for integer fields.
func toleranceFor(exp Expectation) float64 {
	if exp.Tolerance > 0 {
		return exp.Tolerance
	}

	switch exp.Field {
	case ExpectDuration:
		return defaultDurationTolerance
	case ExpectFPS:
		return defaultRateTolerance
	default:
		return 0
	}
}
