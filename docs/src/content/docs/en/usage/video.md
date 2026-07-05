---
title: Video Generation
description: How to use the encfixture video command
---

## Basic

```bash
encfixture video -o test.mp4
```

## Examples

```bash
# Frame counter + timecode
encfixture video --tl frame --tr timecode -d 5 -o counter.mp4

# All overlay positions
encfixture video --tl frame --tr timecode --bl filename --br "CLIP-001" --center "SAMPLE" -d 10 -o full.mp4

# Color bar background + overlays
encfixture video -b test --tl frame --tr timecode -d 5 -o colorbar.mp4

# With sine wave audio
encfixture video -c blue -a sine --frequency 1000 --center "BEEP" -o beep.mp4

# WebM format
encfixture video --tl frame -d 5 -o test.webm

# Custom resolution and FPS
encfixture video -W 3840 -H 2160 -r 60 -d 10 --tl frame -o 4k60.mp4

# Specific codec and quality (HEVC, CRF 28)
encfixture video --codec hevc --crf 28 --tl frame -d 5 -o hevc.mp4

# ProRes for editing workflows (10bit 4:2:2 by default)
encfixture video --codec prores -d 5 -o prores.mov

# Fixed bitrate and pixel format
encfixture video --codec h264 --bitrate 5M --pix-fmt yuv420p -d 5 -o cbr.mp4

# A/V sync test pattern (beep + white flash every second)
encfixture video --sync --tr timecode -d 10 -o sync.mp4

# A/V sync with a 0.5s marker interval
encfixture video --sync --sync-interval 0.5 -d 5 -o sync_fast.mp4
```

## Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--width` | `-W` | 1920 | Video width (px) |
| `--height` | `-H` | 1080 | Video height (px) |
| `--fps` | `-r` | 30 | Frames per second |
| `--duration` | `-d` | 10 | Duration (seconds) |
| `--bg` | `-b` | solid | Background type: solid, test |
| `--color` | `-c` | black | Background color (name or #hex) |
| `--tl` | | | Top-left content |
| `--tr` | | | Top-right content |
| `--center` | | | Center content |
| `--bl` | | | Bottom-left content |
| `--br` | | | Bottom-right content |
| `--scale` | `-S` | 4 | Text scale factor |
| `--output` | `-o` | output.mp4 | Output file path (any format supported by ffmpeg) |
| `--audio` | `-a` | silence | Audio type: silence, sine, noise, tone |
| `--sample-rate` | `-s` | 48000 | Audio sample rate |
| `--channels` | `-C` | 2 | Audio channels |
| `--frequency` | | 440 | Tone frequency (Hz) |
| `--codec` | | | Video codec: h264, hevc, vp9, av1, prores (default: container default) |
| `--crf` | | | Constant rate factor for h264/hevc/vp9/av1 (default: encoder default) |
| `--bitrate` | | | Video bitrate, e.g. `5M` or `800k` (default: encoder default) |
| `--pix-fmt` | | | Pixel format, e.g. `yuv420p`, `yuv422p10le` (default: codec-dependent) |
| `--sync` | | false | A/V sync test pattern: a beep and a white flash fire together at each marker |
| `--sync-interval` | | 1.0 | Seconds between sync markers |

## A/V sync test pattern

`--sync` generates a pattern for detecting audio/video drift. At the start of every interval (`--sync-interval`, default 1s), a white full-frame flash and a beep fire together, each lasting ~0.08s. If the beep and flash appear at different times during playback, the file (or the pipeline that produced it) has A/V desync.

- The beep pitch comes from `--frequency` (default 440Hz).
- Overlays such as `--tr timecode` still render, so you can identify which marker drifted.
- `--sync` overrides `--audio` (the beep track replaces the chosen audio source).

```bash
encfixture video --sync --tr timecode -d 10 -o sync.mp4
```

## Codec selection

`--codec` maps to the following ffmpeg encoders:

| Value | Encoder | Notes |
|---|---|---|
| `h264` | libx264 | |
| `hevc` | libx265 | |
| `vp9` | libvpx-vp9 | default for `.webm` output |
| `av1` | libaom-av1 | slow; use short durations |
| `prores` | prores_ks | defaults to `yuv422p10le` pixel format |

Without `--codec`, the container's default codec is used (e.g. H.264 for `.mp4`, VP9 for `.webm`).
