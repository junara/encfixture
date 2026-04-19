---
description: Generate a single test video with frame / timecode / filename overlays using encfixture. Resolution, fps, duration, and audio type are configurable.
author: junara
input:
  - output
  - width
  - height
  - fps
  - duration
  - audio_type
---

# Generate a single test video with overlays

You are generating one test video that has frame number and timecode burned in, via `encfixture video`. Useful for verifying frame alignment and splice points in encoded streams.

## Prerequisites

- `encfixture` must be on PATH
- `ffmpeg` must be on PATH

## Parameters

| Input | Default | Notes |
|---|---|---|
| `${input:output}` | `test_overlay.mp4` | Extension selects the output format (`.mp4` / `.webm` / `.mkv` / …) |
| `${input:width}` | `1920` | |
| `${input:height}` | `1080` | |
| `${input:fps}` | `30` | |
| `${input:duration}` | `10` | Seconds |
| `${input:audio_type}` | `silence` | One of `silence` / `sine` / `noise` / `tone` |

Use the default whenever an input is omitted.

## Command

Always apply these overlays:

- Top-left: `frame` (frame number)
- Top-right: `timecode` (`HH:MM:SS:FF`)
- Bottom-left: `filename` (output filename)

Execution example:

```bash
encfixture video \
  --tl frame --tr timecode --bl filename \
  -W ${input:width} -H ${input:height} -r ${input:fps} -d ${input:duration} \
  -a ${input:audio_type} \
  -o ${input:output} \
  --json
```

## Reporting

Parse the `--json` output and report `file` / `width` / `height` / `fps` / `duration` in 1–2 lines. On error, surface the `error` field verbatim.
