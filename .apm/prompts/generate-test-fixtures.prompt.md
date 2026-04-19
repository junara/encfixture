---
description: Generate a standard set of ffmpeg encoding-test fixtures with encfixture (1080p30 / 4K60 / webm / mkv videos, sine-wave audio, thumbnail image).
author: junara
input:
  - output_dir
---

# Generate a dummy-asset set for encoding tests

You are generating a standard bundle of dummy media used in ffmpeg encoding tests, via the `encfixture` CLI.

## Prerequisites

- `encfixture` must be on PATH (if missing, instruct the user to build it with `go build -o encfixture .`)
- `ffmpeg` must be on PATH

## Output directory

`${input:output_dir}` — if unspecified, use the current directory. Create it first if it does not exist.

## Assets to generate

Produce these 6 files inside the output directory:

1. `test_1080p30.mp4` — 1920×1080 / 30fps / 30s / overlays: frame (top-left), timecode (top-right)
2. `test_4k60.mp4` — 3840×2160 / 60fps / 10s / same overlays
3. `test.webm` — 1920×1080 / 30fps / 5s / frame overlay
4. `test.mkv` — 1920×1080 / 30fps / 5s / frame overlay
5. `test_audio.wav` — 1 kHz sine wave / 10s
6. `test_thumb.png` — 1920×1080 with "THUMBNAIL" centered

## Execution

For reproducibility and progress visibility, **use batch mode**. Build a `fixtures.json` matching the spec above and run `encfixture batch fixtures.json --json`, then summarize the returned result JSON.

If any job fails, surface its `error` message verbatim and add a one-line guess at the cause (e.g. ffmpeg not installed, output path permission, unsupported codec).
