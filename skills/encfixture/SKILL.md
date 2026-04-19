---
name: encfixture
description: Guide to using `encfixture`, a CLI that generates dummy media files (image / video / audio) for ffmpeg encoding tests. Covers subcommands image / video / audio / batch, overlay placement, and JSON output.
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
```

Common to all commands: pass `--json` to print a structured result to stdout.

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
```

### image flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--width` | `-W` | 1920 | Width (px) |
| `--height` | `-H` | 1080 | Height (px) |
| `--bg` | `-b` | solid | Background: solid, test |
| `--color` | `-c` | black | Background color (name or #hex) |
| `--tl` | | | Top-left content |
| `--tr` | | | Top-right content |
| `--center` | | | Center content |
| `--bl` | | | Bottom-left content |
| `--br` | | | Bottom-right content |
| `--scale` | `-S` | 4 | Text scale factor |
| `--output` | `-o` | output.png | Output path |

## Generating videos

```bash
# Default (black, silent, 10s, 1080p, 30fps)
encfixture video -o output.mp4

# Frame counter + timecode
encfixture video --tl frame --tr timecode -d 5 -o counter.mp4

# Color bars + overlays
encfixture video -b test --tl frame --tr timecode --bl filename --br "CLIP-001" -d 10 -o full.mp4

# With sine-wave audio
encfixture video -a sine --frequency 1000 -d 3 -o beep.mp4

# Custom resolution / fps
encfixture video -W 3840 -H 2160 -r 60 -d 10 -o 4k60.mp4
```

### video flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--width` | `-W` | 1920 | Width (px) |
| `--height` | `-H` | 1080 | Height (px) |
| `--fps` | `-r` | 30 | Frame rate |
| `--duration` | `-d` | 10 | Length (s) |
| `--bg` | `-b` | solid | Background: solid, test |
| `--color` | `-c` | black | Background color |
| `--tl/--tr/--center/--bl/--br` | | | Overlays (same as image) |
| `--scale` | `-S` | 4 | Text scale factor |
| `--output` | `-o` | output.mp4 | Output path (any ffmpeg-supported format) |
| `--audio` | `-a` | silence | Audio: silence, sine, noise, tone |
| `--sample-rate` | `-s` | 48000 | Sample rate |
| `--channels` | `-C` | 2 | Channel count |
| `--frequency` | | 440 | Frequency (Hz) |

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
- `--sample-rate` becomes `sampleRate` (camelCase) in JSON.
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
