---
title: Verify
description: Inspect a media file with the encfixture verify command
---

## Overview

`verify` inspects an existing media file with `ffprobe` and prints its container and per-stream properties. Use it to confirm a generated file (or any file) has the codec, resolution, fps, and duration you expect — a natural "generate then check" step for CI or agents.

```bash
encfixture verify test.mp4
```

## Output

```
$ encfixture verify test.mp4
File:     test.mp4
Format:   mov,mp4,m4a,3gp,3g2,mj2
Duration: 2.000000s
Size:     28934 bytes
Stream 0: video  h264  1920x1080  30fps  yuv420p
Stream 1: audio  aac  48000Hz  2ch
```

## JSON output

Add `--json` for a machine-readable result.

```bash
$ encfixture verify --json test.mp4
{"file":"test.mp4","format":{"formatName":"mov,mp4,...","duration":"2.000000","size":"28934","bitRate":"115736"},"streams":[{"index":0,"type":"video","codec":"h264","width":1920,"height":1080,"fps":"30","pixFmt":"yuv420p"},{"index":1,"type":"audio","codec":"aac","sampleRate":"48000","channels":2}]}
```

### Fields

| Field | Description |
|---|---|
| `format.formatName` | Container format(s) reported by ffprobe |
| `format.duration` | Duration in seconds |
| `format.size` | File size in bytes |
| `format.bitRate` | Overall bit rate |
| `streams[].type` | `video` or `audio` |
| `streams[].codec` | Codec name (e.g. `h264`, `aac`) |
| `streams[].width` / `height` | Video resolution |
| `streams[].fps` | Frame rate (`30`, `29.970`, …) |
| `streams[].pixFmt` | Pixel format |
| `streams[].sampleRate` / `channels` | Audio sample rate and channel count |

## Requirements

`verify` requires `ffprobe`, which ships with ffmpeg. A missing file or unreadable input surfaces ffprobe's error message and exits non-zero.
