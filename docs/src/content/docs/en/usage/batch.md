---
title: Batch Processing
description: Generate multiple files at once with the encfixture batch command
---

## Overview

The `batch` subcommand runs multiple jobs defined in a JSON file. Useful for CI, exhaustive testing, or generating samples across different resolutions.

```bash
encfixture batch jobs.json
```

## JSON schema

The top-level object has an optional `defaults` and a required `jobs` array. Values in `defaults` can be overridden by individual jobs.

```json
{
  "defaults": {
    "width": 1920,
    "height": 1080
  },
  "jobs": [
    { "type": "video", "output": "clip.mp4", "duration": "5", "tl": "frame", "tr": "timecode" },
    { "type": "image", "output": "thumb.png", "bg": "test" },
    { "type": "audio", "output": "beep.wav", "audio": "sine", "frequency": 1000 }
  ]
}
```

Each job requires `type` and `output`. Remaining fields map directly to the flags of the corresponding subcommand. Unknown fields are rejected, which helps catch typos early.

## Job fields

| Field | Flag equivalent | Type | Applies to |
|---|---|---|---|
| `type` | — | `"image"` / `"video"` / `"audio"` | required |
| `output` | `--output` | string | required |
| `width` | `--width` | int | image, video |
| `height` | `--height` | int | image, video |
| `fps` | `--fps` | int | video |
| `duration` | `--duration` | string | video, audio |
| `bg` | `--bg` | string | image, video |
| `color` | `--color` | string | image, video |
| `tl` / `tr` / `center` / `bl` / `br` | `--tl` etc. | string | image, video |
| `scale` | `--scale` | int | image, video |
| `audio` | `--audio` / `--type` | string | video, audio |
| `sampleRate` | `--sample-rate` | int | video, audio |
| `channels` | `--channels` | int | video, audio |
| `frequency` | `--frequency` | float | video, audio |

## CLI flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--parallel` | `-p` | `NumCPU/2` (min 1) | Max concurrent jobs |
| `--fail-fast` | — | false | Skip pending jobs after the first failure |
| `--json` | — | false | Emit structured results to stdout |

### Choosing a parallelism value

ffmpeg is CPU-bound and already uses multiple threads inside a single process, so cranking `--parallel` up doesn't linearly scale throughput — it can make things slower due to contention.

| Batch profile | Suggested value | Reason |
|---|---|---|
| Video-heavy (high-res / long) | `1`–`2` | One ffmpeg already saturates the CPU; parallel jobs just compete |
| Video-heavy (low-res / short) | `NumCPU/2` (default) | Encode time is short; parallelism pays off |
| Image only | `NumCPU` | Pure Go rendering, no ffmpeg — high concurrency is fine |
| Audio only | `NumCPU/2`–`NumCPU` | Lightweight; often I/O bound |
| Mixed (CI) | default | Reasonable middle ground; measure if it matters |

Going above the suggested value is harmless (the internal semaphore still caps it), but optimal values depend on machine, codec, and resolution. For time-critical use, benchmark on your target machine.

## Examples

### Run with a specific concurrency limit

```bash
encfixture batch jobs.json --parallel 4
```

### Stop scheduling after the first failure

```bash
encfixture batch jobs.json --fail-fast
```

### Machine-readable output for CI

```bash
encfixture batch jobs.json --json
```

Output:

```json
{
  "results": [
    { "index": 0, "type": "video", "file": "clip.mp4", "status": "ok" },
    { "index": 1, "type": "image", "file": "thumb.png", "status": "ok" },
    { "index": 2, "type": "audio", "file": "beep.wav", "status": "error", "error": "..." }
  ],
  "succeeded": 2,
  "failed": 1
}
```

## Behavior

- Jobs may run in any order (they execute concurrently). The result array preserves input order.
- With `--fail-fast`, in-flight jobs always run to completion. Pending jobs are recorded with a `job skipped: context canceled` error.
- Exit code is 1 if any job fails.
