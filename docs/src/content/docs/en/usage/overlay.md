---
title: Overlay
description: How the text overlay system works
---

## Layout

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

## Keywords

| Value | Description |
|---|---|
| `frame` | Frame number (dynamic in video, `0` in image) |
| `timecode` | Timecode `HH:MM:SS:FF` (dynamic in video, `00:00:00:00` in image) |
| `filename` | Output filename |

Any other string is displayed as-is (arbitrary text).

## Examples

### Display filename on image

```bash
encfixture image --center filename -o sample.png
```

### Frame counter + timecode on video

```bash
encfixture video --tl frame --tr timecode -d 10 -o counter.mp4
```

### All overlay positions

```bash
encfixture video --tl frame --tr timecode --bl filename --br "CLIP-001" --center "SAMPLE" -d 10 -o full.mp4
```

### Combined with color bar background

```bash
encfixture video -b test --tl frame --tr timecode --center "TEST" -d 5 -o colorbar.mp4
```

## Text scale

Use `--scale` (`-S`) to change the text size. Default is `4`.

```bash
# Large text
encfixture image --center "BIG" -S 8 -o big.png

# Small text
encfixture image --center "small" -S 2 -o small.png
```
