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

## Assertions (`--expect`)

Pass one or more `--expect key=value` assertions to turn verify into a one-shot pass/fail gate. Each expectation is checked against the probed properties; the command exits non-zero (error code `verify_failed`) if any fail.

```bash
encfixture verify test.mp4 --expect codec=h264 --expect width=1920 --expect duration=5
```

```
$ encfixture verify test.mp4 --expect codec=hevc --expect duration=5
File:     test.mp4
...
FAIL  codec: expected hevc, actual h264
PASS  duration = 5.016000
encfixture: verification failed: 1 of 2 expectations failed
```

| Key | Checks |
|---|---|
| `codec`, `width`, `height`, `fps`, `pixFmt` | First video stream |
| `audioCodec`, `sampleRate`, `channels` | First audio stream |
| `duration` | Container duration |

Keys are case-insensitive, and kebab-case aliases (`pix-fmt`, `audio-codec`, `sample-rate`) are also accepted. Numeric keys take an optional tolerance suffix written `5+-0.2` (or `5±0.2`). By default `duration` allows ±0.1s (container rounding) and `fps` ±0.001 (so `fps=29.97` matches ffprobe's `29.970`); other keys compare exactly. A missing stream fails the check with `(no video stream)` / `(no audio stream)` as the actual value.

## JSON output

Add `--json` for a machine-readable result.

```bash
$ encfixture verify --json test.mp4
{"status":"ok","file":"test.mp4","format":{"formatName":"mov,mp4,...","duration":"2.000000","size":"28934","bitRate":"115736"},"streams":[{"index":0,"type":"video","codec":"h264","width":1920,"height":1080,"fps":"30","pixFmt":"yuv420p"},{"index":1,"type":"audio","codec":"aac","sampleRate":"48000","channels":2}]}
```

With `--expect`, the result adds a `checks` array and `status` becomes `ok` or `failed`:

```bash
$ encfixture verify --json test.mp4 --expect codec=hevc
{"status":"failed","file":"test.mp4","format":{...},"streams":[...],"checks":[{"field":"codec","expected":"hevc","actual":"h264","pass":false}]}
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
