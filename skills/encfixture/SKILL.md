---
name: encfixture
description: Guide to using `encfixture`, a CLI that generates dummy media files (image / video / audio) for ffmpeg encoding tests. Covers subcommands image / video / audio / batch / verify / doctor, overlay placement, background patterns (solid / test bars / scrolling gradient / moving box), video codec selection (h264 / hevc / vp9 / av1 / prores) with CRF / bitrate / pixel format, PNG / JPEG image output, A/V sync test patterns (periodic beep + visual flash), inspecting and asserting files with verify --expect, environment / encoder diagnosis with doctor, and JSON output with structured error codes.
---

# encfixture — dummy media file generation skill

A CLI tool that generates dummy assets (image / video / audio) for ffmpeg encoding tests.

## Prerequisites

- `encfixture` must be built (`go build -o encfixture .`)
- `ffmpeg` must be installed (required for video / audio generation)

## Command layout

```
encfixture image [flags]          # generate image
encfixture video [flags]          # generate video
encfixture audio [flags]          # generate audio
encfixture batch <file.json>      # run multiple jobs defined in JSON
encfixture verify <file>          # inspect a file via ffprobe, optionally asserting --expect key=value
encfixture doctor                 # check ffmpeg/ffprobe availability and encoder support
```

Common to all commands: pass `--json` to print a structured result to stdout (errors too — see "Structured errors" below), and `--verbose` to stream ffmpeg's log and encoding progress to stderr.

## Checking the environment first

Run `doctor` before generating when the environment is unknown — it reports whether ffmpeg/ffprobe are installed and which `--codec` values will actually work on this machine (e.g. Homebrew ffmpeg often lacks `libaom-av1`):

```bash
encfixture doctor --json
```

```json
{
  "status": "ok",
  "ffmpeg": { "name": "ffmpeg", "available": true, "version": "7.1", "path": "/opt/homebrew/bin/ffmpeg" },
  "ffprobe": { "name": "ffprobe", "available": true, "version": "7.1", "path": "/opt/homebrew/bin/ffprobe" },
  "videoEncoders": [
    { "codec": "h264", "encoder": "libx264", "available": true },
    { "codec": "av1", "encoder": "libaom-av1", "available": false }
  ],
  "audioEncoders": [
    { "codec": "aac", "encoder": "aac", "available": true },
    { "codec": "opus", "encoder": "libopus", "available": true }
  ]
}
```

Exit code is non-zero only when ffmpeg or ffprobe is missing (`"status": "error"`); unavailable encoders are informational — avoid those `--codec` values.

## Generating images

```bash
# Default (black, 1920x1080)
encfixture image -o output.png

# Solid background color
encfixture image -c blue -o blue.png
encfixture image -c "#ff6600" -o orange.png

# Color bars
encfixture image -b test -o colorbar.png

# Text overlay at each slot
encfixture image --tl frame --tr timecode --bl filename --br "ID-001" --center "TEST" -o info.png

# Custom resolution
encfixture image -W 3840 -H 2160 -o 4k.png

# JPEG output with quality (extension picks the format: .png / .jpg / .jpeg)
encfixture image -c blue -q 75 -o sample.jpg
```

### image flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--width` | `-W` | 1920 | Width (px) |
| `--height` | `-H` | 1080 | Height (px) |
| `--bg` | `-b` | solid | Background: solid, test, gradient, moving |
| `--color` | `-c` | black | Background color (name or #hex) |
| `--tl` | | | Top-left content |
| `--tr` | | | Top-right content |
| `--center` | | | Center content |
| `--bl` | | | Bottom-left content |
| `--br` | | | Bottom-right content |
| `--scale` | `-S` | 4 | Text scale factor |
| `--output` | `-o` | output.png | Output path (.png, .jpg, .jpeg) |
| `--quality` | `-q` | 90 | JPEG quality 1-100 (.jpg/.jpeg only) |
| `--no-clobber` | | | Fail if the output file already exists instead of overwriting |

## Generating videos

```bash
# Default (black, silent, 10s, 1080p, 30fps)
encfixture video -o output.mp4

# Frame counter + timecode
encfixture video --tl frame --tr timecode -d 5 -o counter.mp4

# Color bars + overlays
encfixture video -b test --tl frame --tr timecode --bl filename --br "CLIP-001" -d 10 -o full.mp4

# Moving backgrounds (motion for codec compression tests)
encfixture video -b gradient -d 5 -o gradient.mp4
encfixture video -b moving --tr timecode -d 5 -o moving.mp4

# With sine-wave audio
encfixture video -a sine --frequency 1000 -d 3 -o beep.mp4

# Custom resolution / fps
encfixture video -W 3840 -H 2160 -r 60 -d 10 -o 4k60.mp4

# Specific codec / quality
encfixture video --codec hevc --crf 28 -d 5 -o hevc.mp4
encfixture video --codec prores -d 5 -o prores.mov
encfixture video --codec h264 --bitrate 5M --pix-fmt yuv420p -d 5 -o cbr.mp4

# A/V sync test pattern (beep + white flash every second; overrides --audio)
encfixture video --sync --tr timecode -d 10 -o sync.mp4
encfixture video --sync --sync-interval 0.5 -d 5 -o sync_fast.mp4
```

### video flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--width` | `-W` | 1920 | Width (px) |
| `--height` | `-H` | 1080 | Height (px) |
| `--fps` | `-r` | 30 | Frame rate |
| `--duration` | `-d` | 10 | Length (s) |
| `--bg` | `-b` | solid | Background: solid, test, gradient, moving |
| `--color` | `-c` | black | Background color |
| `--tl/--tr/--center/--bl/--br` | | | Overlays (same as image) |
| `--scale` | `-S` | 4 | Text scale factor |
| `--output` | `-o` | output.mp4 | Output path (any ffmpeg-supported format) |
| `--audio` | `-a` | silence | Audio: silence, sine, noise, tone |
| `--sample-rate` | `-s` | 48000 | Sample rate |
| `--channels` | `-C` | 2 | Channel count |
| `--frequency` | | 440 | Frequency (Hz) |
| `--codec` | | | Video codec: h264, hevc, vp9, av1, prores (default: container default) |
| `--crf` | | | CRF for h264/hevc/vp9/av1 (default: encoder default) |
| `--bitrate` | | | Video bitrate, e.g. `5M` (default: encoder default) |
| `--pix-fmt` | | | Pixel format, e.g. `yuv420p` (default: codec-dependent; prores uses `yuv422p10le`) |
| `--sync` | | false | A/V sync test pattern: beep + white flash together at each marker (overrides `--audio`) |
| `--sync-interval` | | 1.0 | Seconds between sync markers |
| `--no-clobber` | | | Fail if the output file already exists instead of overwriting |

## Generating audio

```bash
# Silent WAV (10s)
encfixture audio -o silence.wav

# Sine wave
encfixture audio -t sine -f 1000 -d 5 -o beep.wav

# White noise
encfixture audio -t noise -d 3 -o noise.mp3
```

### audio flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--type` | `-t` | silence | Type: silence, sine, noise, tone |
| `--duration` | `-d` | 10 | Length (s) |
| `--sample-rate` | `-s` | 48000 | Sample rate |
| `--channels` | `-C` | 2 | Channel count |
| `--frequency` | `-f` | 440 | Frequency (Hz) |
| `--output` | `-o` | output.wav | Output path |
| `--no-clobber` | | | Fail if the output file already exists instead of overwriting |

## Verifying files

Inspect an existing file's codec, resolution, fps, and duration via ffprobe.
Handy for a "generate then check" loop.

```bash
encfixture verify out.mp4            # human-readable
encfixture verify --json out.mp4     # machine-readable
```

`--json` shape:

```json
{
  "status": "ok",
  "file": "out.mp4",
  "format": { "formatName": "mov,mp4,...", "duration": "2.000000", "size": "28934", "bitRate": "115736" },
  "streams": [
    { "index": 0, "type": "video", "codec": "h264", "width": 1920, "height": 1080, "fps": "30", "pixFmt": "yuv420p" },
    { "index": 1, "type": "audio", "codec": "aac", "sampleRate": "48000", "channels": 2 }
  ]
}
```

Requires `ffprobe` (ships with ffmpeg). Missing/unreadable files exit non-zero
with ffprobe's error message.

### Asserting expectations (`--expect`)

Turn verify into a one-shot pass/fail gate — prefer this over reading the JSON
and comparing yourself. Each `--expect key=value` is checked against the probed
properties; the command exits non-zero if any fail.

```bash
encfixture verify out.mp4 --expect codec=h264 --expect width=1920 --expect duration=5
encfixture verify --json out.mp4 --expect codec=h264 --expect audioCodec=aac
```

- Supported keys: `codec`, `width`, `height`, `fps`, `pixFmt`, `duration`,
  `audioCodec`, `sampleRate`, `channels`. Keys are case-insensitive and
  kebab-case aliases (`pix-fmt`, `audio-codec`, `sample-rate`) also work.
- `codec`/`width`/`height`/`fps`/`pixFmt` check the first video stream;
  `audioCodec`/`sampleRate`/`channels` the first audio stream; `duration` the
  container. A missing stream fails the check with e.g. `"(no video stream)"`.
- Numeric keys accept a tolerance suffix: `duration=5+-0.2` (or `5±0.2`).
  Defaults: `duration` ±0.1s (container rounding), `fps` ±0.001 (so
  `fps=29.97` matches ffprobe's `29.970`), all other keys exact.

With `--json`, the output gains `"status": "ok" | "failed"` and a `checks` array:

```json
{ "status": "failed", "file": "out.mp4", "format": {}, "streams": [],
  "checks": [ { "field": "codec", "expected": "hevc", "actual": "h264", "pass": false } ] }
```

On failure the process exits 1 with error code `verify_failed`.

## Batch processing

Define multiple jobs in a JSON file and run them in one go. Unknown fields cause an error, so typos are caught early.

```bash
encfixture batch jobs.json [--parallel N] [--fail-fast] [--json]
```

### JSON schema

```json
{
  "defaults": { "width": 1920, "height": 1080 },
  "jobs": [
    { "type": "video", "output": "clip.mp4", "duration": "5", "tl": "frame", "tr": "timecode" },
    { "type": "image", "output": "thumb.png", "bg": "test" },
    { "type": "audio", "output": "beep.wav", "audio": "sine", "frequency": 1000 }
  ]
}
```

- `type` and `output` are required. Other fields mirror the corresponding subcommand flags.
- `--sample-rate` becomes `sampleRate` and `--pix-fmt` becomes `pixFmt` (camelCase) in JSON.
- Video jobs also accept `codec`, `crf`, `bitrate`, `pixFmt`, `sync`, `syncInterval`; image jobs accept `quality`.
- `defaults` applies to all jobs and can be overridden per job.

### batch flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--parallel` | `-p` | `NumCPU/2` (min 1) | Max concurrent jobs |
| `--fail-fast` | | false | Skip remaining jobs after the first failure |

### `--json` output

```json
{
  "results": [
    { "index": 0, "type": "image", "file": "a.png", "status": "ok" },
    { "index": 1, "type": "video", "file": "b.mp4", "status": "error", "error": "..." }
  ],
  "succeeded": 1,
  "failed": 1
}
```

Exit code is 1 if any job fails.

### Parallelism guidance

ffmpeg is already multithreaded per process, so high concurrency does not scale linearly. Image jobs (pure Go rendering) benefit the most from parallelism.

| Job mix | Recommended `--parallel` |
|---|---|
| Mostly video (high-res / long) | `1`–`2` |
| Mostly video (low-res / short) | `NumCPU/2` (default) |
| Images only | `NumCPU` |
| Audio only | `NumCPU/2`–`NumCPU` |

Reference benchmark (30 images at 1920×1080, 10-core machine): `--parallel 1` → 1.19s, `--parallel 10` → 0.24s (~5x). Saturation point is near the CPU core count.

## Overlay placement

Five slots accept free-form text.

```
┌──────────────────────────────┐
│ --tl              --tr       │
│                              │
│          --center            │
│                              │
│ --bl              --br       │
└──────────────────────────────┘
```

Reserved keywords:
- `frame` — frame number (dynamic per-frame in video, `0` in image)
- `timecode` — `HH:MM:SS:FF` (dynamic in video, `00:00:00:00` in image)
- `filename` — output filename
- anything else — rendered literally

## Supported colors

Names: `black`, `white`, `red`, `green`, `blue`, `yellow`, `cyan`, `magenta`, `gray`

Hex: `#ff6600`, `#333333`, etc.

## JSON output

```bash
$ encfixture video --json --tl frame -d 5 -o test.mp4
{"status":"ok","file":"test.mp4","type":"video","width":1920,"height":1080,"fps":30,"duration":"5"}
```

### Structured errors

On failure, `--json` emits an error object to stdout and exits non-zero. The
`code` is stable and machine-readable — branch on it instead of parsing the
message; `hint` says how to recover:

```bash
$ encfixture video --json --codec av1 -o test.mp4
{"status":"error","code":"encoder_not_available","error":"...Unknown encoder 'libaom-av1'","hint":"Run 'encfixture doctor' to list the encoders your ffmpeg build supports, then pick an available --codec."}
```

| `code` | meaning | typical recovery |
|---|---|---|
| `usage` | wrong flags/args | fix the command line |
| `ffmpeg_not_found` / `ffprobe_not_found` | tool missing from PATH | install ffmpeg |
| `encoder_not_available` | encoder not in this ffmpeg build | run `doctor`, pick an available `--codec` |
| `unknown_codec` / `unknown_background` | invalid `--codec` / `--bg` | use a listed value |
| `invalid_duration` / `invalid_bitrate` / `invalid_expectation` | malformed value | fix the value |
| `verify_failed` | `--expect` assertions failed | inspect `checks` in the output |
| `output_exists` | refused overwrite under `--no-clobber` | remove file or drop the flag |
| `probe_failed` | ffprobe could not read the file | check path / file integrity |
| `ffmpeg_failed` | unclassified ffmpeg error | re-run with `--verbose` |
| `env_unhealthy` | `doctor` found ffmpeg/ffprobe missing | install ffmpeg |
| `error` | anything else | read `error` |

Without `--json`, the same hint is printed to stderr as a `hint:` line.

## Common recipes

### Encode-test asset set (individual commands)

```bash
encfixture video --tl frame --tr timecode -d 30 -o test_1080p30.mp4
encfixture video --tl frame --tr timecode -W 3840 -H 2160 -r 60 -d 10 -o test_4k60.mp4
encfixture video --tl frame -d 5 -o test.webm
encfixture video --tl frame -d 5 -o test.mkv
encfixture audio -t sine -f 1000 -d 10 -o test_audio.wav
encfixture image --center "THUMBNAIL" -o test_thumb.png
```

### Same set via batch (reproducible, CI-friendly)

`fixtures.json`:

```json
{
  "defaults": { "tl": "frame", "tr": "timecode" },
  "jobs": [
    { "type": "video", "output": "test_1080p30.mp4", "duration": "30" },
    { "type": "video", "output": "test_4k60.mp4", "width": 3840, "height": 2160, "fps": 60, "duration": "10" },
    { "type": "video", "output": "test.webm", "duration": "5" },
    { "type": "video", "output": "test.mkv", "duration": "5" },
    { "type": "audio", "output": "test_audio.wav", "audio": "sine", "frequency": 1000, "duration": "10" },
    { "type": "image", "output": "test_thumb.png", "center": "THUMBNAIL" }
  ]
}
```

```bash
encfixture batch fixtures.json --json
```

### Codec coverage set (same content across codecs)

Generate one clip per codec to test a decoder / pipeline against multiple encoders.

`codecs.json`:

```json
{
  "defaults": { "type": "video", "duration": "5", "tl": "frame", "tr": "timecode" },
  "jobs": [
    { "output": "cov_h264.mp4", "codec": "h264", "crf": 23 },
    { "output": "cov_hevc.mp4", "codec": "hevc", "crf": 28 },
    { "output": "cov_vp9.webm", "codec": "vp9" },
    { "output": "cov_prores.mov", "codec": "prores" }
  ]
}
```

```bash
# Check encoder availability first (e.g. Homebrew ffmpeg often lacks libaom-av1)
encfixture doctor --json

encfixture batch codecs.json --json
```

Then assert each output got the codec / pixel format you asked for:

```bash
encfixture verify cov_hevc.mp4 --expect codec=hevc --expect pixFmt=yuv420p --expect "duration=5+-0.2"
```

### A/V sync check clip

Generate a clip with a beep + flash every second, then confirm the beep timing
lines up with the visual flash to detect audio/video drift.

```bash
encfixture video --sync --tr timecode -d 10 -o sync.mp4

# beep onsets should sit at 0s, 1s, 2s, ... (each lasting ~0.08s)
ffmpeg -v error -i sync.mp4 -af "silencedetect=noise=-40dB:d=0.03,ametadata=print:file=-" -f null -
```
