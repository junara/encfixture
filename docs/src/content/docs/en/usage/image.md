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

# JPEG output with quality
encfixture image -c blue --center "SAMPLE" -q 75 -o sample.jpg
```

## Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--width` | `-W` | 1920 | Image width (px) |
| `--height` | `-H` | 1080 | Image height (px) |
| `--bg` | `-b` | solid | Background type: solid, test, gradient, moving |
| `--color` | `-c` | black | Background color (name or #hex) |
| `--tl` | | | Top-left content |
| `--tr` | | | Top-right content |
| `--center` | | | Center content |
| `--bl` | | | Bottom-left content |
| `--br` | | | Bottom-right content |
| `--scale` | `-S` | 4 | Text scale factor |
| `--output` | `-o` | output.png | Output file path (.png, .jpg, .jpeg) |
| `--quality` | `-q` | 90 | JPEG quality 1-100 (for .jpg/.jpeg output) |
| `--no-clobber` | | | Fail if the output file already exists instead of overwriting |

## Output formats

The encoder is selected from the file extension: `.png` produces PNG, `.jpg`/`.jpeg` produces JPEG. Other extensions are rejected with an error.
