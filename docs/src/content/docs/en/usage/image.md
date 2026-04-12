---
title: Image Generation
description: How to use the encfixture image command
---

## Basic

```bash
encfixture image -o test.png
```

## Examples

```bash
# Blue background
encfixture image -c blue -o blue.png

# Color bar test pattern
encfixture image -b test -o colorbar.png

# All overlay positions
encfixture image --tl frame --tr timecode --bl filename --br "ID-001" --center "TEST" -o info.png

# Custom resolution
encfixture image -W 3840 -H 2160 -c white -o 4k.png

# Hex color code
encfixture image -c "#ff6600" -o orange.png

# Color bar + text overlay
encfixture image -b test --center "SAMPLE" -o test_with_text.png
```

## Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--width` | `-W` | 1920 | Image width (px) |
| `--height` | `-H` | 1080 | Image height (px) |
| `--bg` | `-b` | solid | Background type: solid, test |
| `--color` | `-c` | black | Background color (name or #hex) |
| `--tl` | | | Top-left content |
| `--tr` | | | Top-right content |
| `--center` | | | Center content |
| `--bl` | | | Bottom-left content |
| `--br` | | | Bottom-right content |
| `--scale` | `-S` | 4 | Text scale factor |
| `--output` | `-o` | output.png | Output file path |
