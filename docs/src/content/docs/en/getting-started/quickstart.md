---
title: Quick Start
description: Basic usage of encfixture
---

## Generate your first files

### Image

```bash
encfixture image -o test.png
```

Generates a black 1920x1080 PNG image.

### Video

```bash
encfixture video --tl frame --tr timecode -d 5 -o test.mp4
```

Generates a 5-second video with frame number and timecode overlays.

### Audio

```bash
encfixture audio -t sine -f 1000 -d 3 -o beep.wav
```

Generates a 3-second WAV file with a 1000Hz sine wave.

## Overlay system

Place text freely at 5 positions on images and videos.

```
┌──────────────────────────────┐
│ --tl              --tr       │
│                              │
│          --center            │
│                              │
│ --bl              --br       │
└──────────────────────────────┘
```

Each position accepts:

| Value | Description |
|---|---|
| `frame` | Frame number (dynamic in video, `0` in image) |
| `timecode` | Timecode `HH:MM:SS:FF` |
| `filename` | Output filename |
| Any other string | Displayed as-is (arbitrary text) |

### Example: all overlay positions

```bash
encfixture video --tl frame --tr timecode --bl filename --br "CLIP-001" --center "SAMPLE" -d 10 -o full.mp4
```
