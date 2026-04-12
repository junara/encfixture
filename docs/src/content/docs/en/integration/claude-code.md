---
title: Claude Code
description: How to use encfixture with Claude Code
---

encfixture supports `--json` output and descriptive `--help`, making it easy for Claude Code to generate dummy media files via Bash tool calls.

## Examples

### Generate test video

```
> Create a 5-second 720p test video with frame counter and timecode

Claude runs:
  encfixture video --json --tl frame --tr timecode -W 1280 -H 720 -d 5 -o test_720p.mp4
```

### Generate multiple test images

```
> Create 3 test images: color bar, solid blue, and one with filename displayed

Claude runs:
  encfixture image --json -b test -o colorbar.png
  encfixture image --json -c blue -o blue.png
  encfixture image --json --center filename -o sample.png
```

### Generate test audio

```
> Create a 3-second 1000Hz beep

Claude runs:
  encfixture audio --json -t sine -f 1000 -d 3 -o beep.wav
```

## CLAUDE.md integration

Add the following to your project's `CLAUDE.md` to let Claude Code know about encfixture:

```markdown
## Tools

- `encfixture` is available for generating dummy media files (image, video, audio) for ffmpeg encoding tests.
  - Always use `--json` flag for structured output.
  - Use `encfixture <subcommand> --help` for available flags.
```
