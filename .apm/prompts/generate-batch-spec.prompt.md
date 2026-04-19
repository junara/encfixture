---
description: Interactively build an encfixture batch JSON spec from free-form requirements and run it with encfixture batch.
author: junara
input:
  - requirements
  - spec_path
  - parallel
---

# Build a batch spec from requirements and run it

You drive the full loop of authoring an `encfixture batch` JSON spec interactively and executing it.

## Inputs

- `${input:requirements}` — natural-language description of the desired assets (e.g. "1080p and 4K videos with frame overlay, plus a 440Hz sine wave for BGM")
- `${input:spec_path}` — destination for the spec JSON (default: `fixtures.json`)
- `${input:parallel}` — parallelism (default: unset, i.e. `NumCPU/2`)

## Steps

1. **Decompose the requirements**: read `${input:requirements}` and enumerate the required jobs. If anything is ambiguous (resolution, duration, whether audio is needed, etc.), ask the user once in a single consolidated question. Avoid back-and-forth.

2. **Build the spec** under these constraints:
   - Every job needs `type` (`image` / `video` / `audio`) and `output`
   - Lift shared values into `defaults`
   - Sample rate is `sampleRate` (camelCase) in JSON
   - Overlay slots are `tl` / `tr` / `center` / `bl` / `br`. Reserved keywords: `frame` / `timecode` / `filename`
   - Unknown fields error out, so field names must match CLI flag names exactly

3. **Save the spec** to `${input:spec_path}` (default `fixtures.json`).

4. **Choose parallelism**: use `1`–`2` if videos are mostly high-res, `NumCPU` if image-heavy, default otherwise. `${input:parallel}` takes precedence when provided.

5. **Execute**: run `encfixture batch ${input:spec_path} --json [--parallel N]`.

6. **Report**: summarize `succeeded` / `failed` from the result JSON. List output paths of successful jobs. Surface `error` messages for any failures.

## Notes

- When the user requests high parallelism, remind them once that ffmpeg is already multithreaded so video jobs do not scale linearly. Don't repeat this for low parallelism.
- Only suggest `--fail-fast` when the user's context signals CI usage.
