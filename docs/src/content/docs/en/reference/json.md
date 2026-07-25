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

When a command fails, `--json` emits an error object to stdout (the human-readable message still goes to stderr) and the exit code is non-zero.

```bash
$ encfixture verify --json missing.mp4
{"status":"error","error":"verify failed: probe failed: ..."}
```

| Field | Type | Description |
|---|---|---|
| `status` | string | Always `"error"` |
| `error` | string | Error message |

### batch

The `batch` command emits an aggregate object containing per-job results and totals. See [Batch Processing](/encfixture/en/usage/batch/) for details.

```bash
$ encfixture batch --json jobs.json
{"results":[{"index":0,"type":"image","file":"a.png","status":"ok"}],"succeeded":1,"failed":0}
```
