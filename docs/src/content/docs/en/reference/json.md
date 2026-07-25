---
title: JSON Output
description: Structured output with the --json flag
---

## Usage

Add `--json` to output results as JSON to stdout.

```bash
encfixture video --json --tl frame --tr timecode -d 5 -o test.mp4
```

Output:

```json
{"status":"ok","file":"test.mp4","type":"video","width":1920,"height":1080,"fps":30,"duration":"5"}
```

## Response fields

| Field | Type | Description |
|---|---|---|
| `status` | string | `"ok"` on success, `"error"` on failure |
| `file` | string | Output file path |
| `type` | string | `"image"`, `"video"`, or `"audio"` |
| `width` | int | Width (image/video only) |
| `height` | int | Height (image/video only) |
| `fps` | int | Frames per second (video only) |
| `duration` | string | Duration (video/audio only) |

## Examples

### image

```bash
$ encfixture image --json --center "TEST" -o test.png
{"status":"ok","file":"test.png","type":"image","width":1920,"height":1080}
```

### video

```bash
$ encfixture video --json --tl frame -d 3 -o test.mp4
{"status":"ok","file":"test.mp4","type":"video","width":1920,"height":1080,"fps":30,"duration":"3"}
```

### audio

```bash
$ encfixture audio --json -t sine -d 3 -o beep.wav
{"status":"ok","file":"beep.wav","type":"audio","duration":"3"}
```

### Errors

When a command fails, `--json` emits a structured error object to stdout (the human-readable message still goes to stderr) and the exit code is non-zero. The `code` is stable and machine-readable — branch on it instead of parsing the message.

```bash
$ encfixture video --json --codec av1 -o test.mp4
{"status":"error","code":"encoder_not_available","error":"...Unknown encoder 'libaom-av1'","hint":"Run 'encfixture doctor' to list the encoders your ffmpeg build supports, then pick an available --codec."}
```

| Field | Type | Description |
|---|---|---|
| `status` | string | Always `"error"` |
| `code` | string | Stable machine-readable error code (below) |
| `error` | string | Error message |
| `hint` | string | How to recover (omitted when there is nothing useful to say) |

| Code | Meaning |
|---|---|
| `usage` | Wrong flags or arguments |
| `ffmpeg_not_found` / `ffprobe_not_found` | Tool missing from PATH |
| `encoder_not_available` | Requested encoder not in this ffmpeg build |
| `unknown_codec` / `unknown_background` | Invalid `--codec` / `--bg` value |
| `invalid_duration` / `invalid_bitrate` | Malformed `-d` / `--bitrate` value |
| `invalid_expectation` | Malformed `--expect` assertion |
| `verify_failed` | One or more `--expect` assertions failed |
| `output_exists` | Refused to overwrite under `--no-clobber` |
| `probe_failed` | ffprobe could not read the file |
| `ffmpeg_failed` | ffmpeg exited with an unclassified error |
| `env_unhealthy` | `doctor` found ffmpeg/ffprobe missing |
| `error` | Anything else |

Without `--json`, the same hint is printed to stderr as a `hint:` line after the error message.

### batch

The `batch` command emits an aggregate object containing per-job results and totals. See [Batch Processing](/encfixture/en/usage/batch/) for details.

```bash
$ encfixture batch --json jobs.json
{"results":[{"index":0,"type":"image","file":"a.png","status":"ok"}],"succeeded":1,"failed":0}
```
