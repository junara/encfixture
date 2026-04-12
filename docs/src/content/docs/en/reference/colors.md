---
title: Colors
description: Available color options
---

## Named colors

| Name | Color |
|---|---|
| `black` | Black |
| `white` | White |
| `red` | Red |
| `green` | Green |
| `blue` | Blue |
| `yellow` | Yellow |
| `cyan` | Cyan |
| `magenta` | Magenta |
| `gray` / `grey` | Gray |

## Hex color codes

Specify any color using `#RRGGBB` format.

```bash
encfixture image -c "#ff6600" -o orange.png
encfixture image -c "#333333" -o dark_gray.png
encfixture video -c "#0066ff" -d 5 -o blue.mp4
```

## Text color

Text color is automatically selected based on the background color (contrast auto-adjustment).

- Dark background → White text
- Light background → Black text
