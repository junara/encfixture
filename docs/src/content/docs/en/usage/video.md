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
