---
title: Doctor
description: Check ffmpeg/ffprobe availability and encoder support with encfixture doctor
---

## Overview

`doctor` checks the environment encfixture depends on: whether `ffmpeg` and `ffprobe` are installed, and which of the selectable encoders (`--codec` values, plus the audio encoders used for mp4/webm output) the local ffmpeg build supports. Run it before generating when the environment is unknown — for example, Homebrew's ffmpeg often lacks `libaom-av1`.

```bash
encfixture doctor
```

## Output

```
$ encfixture doctor
ffmpeg:  OK 7.1 (/opt/homebrew/bin/ffmpeg)
ffprobe: OK 7.1 (/opt/homebrew/bin/ffprobe)
video encoders:
  h264     libx264      OK
  hevc     libx265      OK
  vp9      libvpx-vp9   OK
  av1      libaom-av1   MISSING
  prores   prores_ks    OK
audio encoders:
  aac      aac          OK
  opus     libopus      OK
```

## JSON output

```bash
$ encfixture doctor --json
{"status":"ok","ffmpeg":{"name":"ffmpeg","available":true,"version":"7.1","path":"/opt/homebrew/bin/ffmpeg"},"ffprobe":{"name":"ffprobe","available":true,"version":"7.1","path":"/opt/homebrew/bin/ffprobe"},"videoEncoders":[{"codec":"h264","encoder":"libx264","available":true},{"codec":"av1","encoder":"libaom-av1","available":false}],"audioEncoders":[{"codec":"aac","encoder":"aac","available":true},{"codec":"opus","encoder":"libopus","available":true}]}
```

| Field | Description |
|---|---|
| `status` | `ok` when ffmpeg and ffprobe are both present, `error` otherwise |
| `ffmpeg` / `ffprobe` | `available`, resolved `path`, and `version` |
| `videoEncoders[]` | One entry per `--codec` value with the ffmpeg encoder it maps to |
| `audioEncoders[]` | `aac` (mp4/mov default) and `opus` (webm output) |

## Exit code

`doctor` exits non-zero only when `ffmpeg` or `ffprobe` is missing (error code `env_unhealthy`). Unavailable encoders are informational: generation commands that need them fail individually with the `encoder_not_available` error code, whose hint points back to `doctor`.
