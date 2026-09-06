---
title: Target File Size
description: Generating image, video, and audio files of an approximate size with --target-size
---

`--target-size` asks for an output of a given size, for tests where the size itself is what matters: upload limits, storage quotas, transfer timings.

Sizes take an optional `K`/`M`/`G` suffix and are binary — `1KB` is 1024 bytes, `10MB` is 10485760 — and `10M`, `10MB` and `10MiB` all mean the same thing. A plain number is a byte count.

## Basic

```bash
# 20MB video: the bitrate comes from size / duration
encfixture video --target-size 20MB -d 30 -o sized.mp4

# 2MB audio
encfixture audio -t sine -d 60 --target-size 2MB -o sized.mp3

# 500KB JPEG: the quality that fits is found by bisection
encfixture image -b test --target-size 500KB -o sized.jpg
```

## How close it lands

The size is **approximate**. Video and audio reach it by deriving a bitrate rather than padding the file, so how close the result lands depends on the encoder.

### Video

Setting `--target-size` on an `.mp4`, `.mov`, `.mkv` or `.m4v` output selects h264 and pins the encoder to CBR, which lands within a few percent even on a black frame:

```bash
$ encfixture video --target-size 5MB -d 10 -o sized.mp4 --json
{"status":"ok","file":"sized.mp4","type":"video","width":1920,"height":1080,"fps":30,"duration":"10","size":5144863,"targetSize":5242880}
```

Other codecs (`hevc`, `vp9`, `av1`, `prores`) only get bitrate limits, and content that compresses to almost nothing falls well short of them:

```bash
$ encfixture video --codec hevc --target-size 5MB -d 10 -o small.mp4
Created: small.mp4
warning: output is 61.71KB, 99% off the requested 5.00MB; to get closer, raise --width/--height or --fps, use a busier --bg (test, gradient, moving), or pick --codec h264, which tracks a target bitrate closest
```

### Audio

The bitrate is set for compressed formats. `.mp3` honors it closely; AAC and Opus collapse a silent track to a few kbit/s however much bitrate they are given, so pair `--target-size` with `-t noise` or `-t sine` for those.

PCM and lossless outputs (`.wav`, `.aiff`, `.au`, `.flac`) have no bitrate dial and are rejected. Since their size follows the sample rate, the error names the rate that would reach it:

```bash
$ encfixture audio --target-size 1MB -d 10 -o out.wav
encfixture: audio generation failed: target size not supported for this format: .wav is not bitrate-controlled (16-bit PCM reaches 1.00MB at --sample-rate 26214); write a compressed format such as .mp3, .m4a or .opus
```

### Images

A JPEG is sized by picking the quality that fits, so `--target-size` needs `.jpg`/`.jpeg` output — PNG is lossless and has no such dial. The ceiling is whatever quality 100 produces at the given resolution; raise `--width`/`--height` or pick a busier `--bg` for larger files.

## Reported size

Every generation result reports the size actually produced, and warns when it lands more than 10% off the request. With `--json` the size is machine-readable, so CI can assert on it:

| Field | Type | Description |
|---|---|---|
| `size` | int | Size of the generated file, in bytes |
| `targetSize` | int | The requested size in bytes (only with `--target-size`) |
| `warning` | string | Present when the result is more than 10% off the target |

## Conflicting flags

`--target-size` sets the same dial as `--bitrate` and `--crf` (video) and `--quality` (image), so combining them is a usage error:

```bash
$ encfixture video --target-size 5MB --bitrate 1M -o out.mp4
encfixture: conflicting flags: --target-size and --bitrate cannot be combined
```

## Batch

Batch jobs take the same value as a `targetSize` string, per job or in `defaults`. A job that sets `targetSize` alongside `bitrate` or `crf` fails with `target_size_conflict` rather than letting one of them silently win:

```json
{
  "defaults": {"targetSize": "1MB", "duration": "5"},
  "jobs": [
    {"type": "video", "output": "a.mp4"},
    {"type": "audio", "output": "b.m4a", "audio": "noise", "targetSize": "200KB"},
    {"type": "image", "output": "c.jpg", "bg": "gradient", "targetSize": "100KB"}
  ]
}
```

## Errors

| Code | Meaning |
|---|---|
| `invalid_target_size` | Malformed `--target-size` value |
| `target_size_unsupported` | Format has no size dial (PNG, WAV, FLAC, …) |
| `target_size_too_small` | The size leaves less bitrate than an encoder can use |
| `target_size_conflict` | A batch job set `targetSize` together with `bitrate` or `crf` |
