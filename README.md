# encfixture

[日本語](README.ja.md) | English | [Documentation](https://junara.github.io/encfixture/)

A Go CLI tool for generating dummy media files (image, video, audio) for ffmpeg encoding tests.

## Requirements

- ffmpeg

## Installation

### Homebrew (macOS / Linux)

```bash
brew install junara/tap/encfixture
```

### Go

```bash
go install github.com/junara/encfixture@latest
```

### Binary

Download from [Releases](https://github.com/junara/encfixture/releases).

## Usage

### Global Flags

| Flag | Description |
|---|---|
| `--json` | Output results as JSON |
| `--version` | Show version |

### Overlay System

Place text freely at 5 positions on images and videos.

```
┌──────────────────────────────┐
│ --tl              --tr       │
│                              │
│          --center            │
│                              │
│ --bl              --br       │
└──────────────────────────────┘
```

Each position accepts:

| Value | Description |
|---|---|
| `frame` | Frame number (dynamic in video, `0` in image) |
| `timecode` | Timecode `HH:MM:SS:FF` (dynamic in video, `00:00:00:00` in image) |
| `filename` | Output filename |
| Any other string | Displayed as-is (arbitrary text) |

### Image Generation

```bash
# Solid black image (1920x1080)
encfixture image -o test.png

# Blue background with centered text
encfixture image -c blue --center "SAMPLE" -o sample.png

# Color bar test pattern
encfixture image -b test -o colorbar.png

# All overlay positions
encfixture image --tl frame --tr timecode --bl filename --br "ID-001" --center "TEST" -o info.png

# Custom resolution
encfixture image -W 3840 -H 2160 -c white -o 4k.png

# Hex color code
encfixture image -c "#ff6600" -o orange.png

# JSON output
encfixture image --json --center "TEST" -o test.png
```

#### image Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--width` | `-W` | 1920 | Image width (px) |
| `--height` | `-H` | 1080 | Image height (px) |
| `--bg` | `-b` | solid | Background type: solid, test |
| `--color` | `-c` | black | Background color (name or #hex) |
| `--tl` | | | Top-left content |
| `--tr` | | | Top-right content |
| `--center` | | | Center content |
| `--bl` | | | Bottom-left content |
| `--br` | | | Bottom-right content |
| `--scale` | `-S` | 4 | Text scale factor |
| `--output` | `-o` | output.png | Output file path |

### Video Generation

```bash
# Black silent video (10s, 1080p, 30fps)
encfixture video -o test.mp4

# Frame counter + timecode
encfixture video --tl frame --tr timecode -d 5 -o counter.mp4

# All overlay positions
encfixture video --tl frame --tr timecode --bl filename --br "CLIP-001" --center "SAMPLE" -d 10 -o full.mp4

# Color bar background + overlays
encfixture video -b test --tl frame --tr timecode -d 5 -o colorbar.mp4

# With sine wave audio
encfixture video -c blue -a sine --frequency 1000 --center "BEEP" -o beep.mp4

# WebM format
encfixture video --tl frame -d 5 -o test.webm

# Custom resolution and FPS
encfixture video -W 3840 -H 2160 -r 60 -d 10 --tl frame -o 4k60.mp4

# JSON output
encfixture video --json --tl frame -d 3 -o test.mp4
```

#### video Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--width` | `-W` | 1920 | Video width (px) |
| `--height` | `-H` | 1080 | Video height (px) |
| `--fps` | `-r` | 30 | Frames per second |
| `--duration` | `-d` | 10 | Duration (seconds) |
| `--bg` | `-b` | solid | Background type: solid, test |
| `--color` | `-c` | black | Background color (name or #hex) |
| `--tl` | | | Top-left content |
| `--tr` | | | Top-right content |
| `--center` | | | Center content |
| `--bl` | | | Bottom-left content |
| `--br` | | | Bottom-right content |
| `--scale` | `-S` | 4 | Text scale factor |
| `--output` | `-o` | output.mp4 | Output file path (any format supported by ffmpeg) |
| `--audio` | `-a` | silence | Audio type: silence, sine, noise, tone |
| `--sample-rate` | `-s` | 48000 | Audio sample rate |
| `--channels` | `-C` | 2 | Audio channels |
| `--frequency` | | 440 | Tone frequency (Hz) |

### Audio Generation

```bash
# Silent WAV (10s)
encfixture audio -o silence.wav

# Sine wave (440Hz)
encfixture audio -t sine -d 5 -o sine.wav

# 1000Hz tone
encfixture audio -t tone -f 1000 -d 3 -o tone.wav

# White noise
encfixture audio -t noise -d 5 -o noise.wav

# Mono, 44100Hz
encfixture audio -t silence -C 1 -s 44100 -o mono.wav

# MP3 format
encfixture audio -t sine -d 5 -o beep.mp3

# JSON output
encfixture audio --json -t sine -d 3 -o beep.wav
```

#### audio Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--type` | `-t` | silence | Audio type: silence, sine, noise, tone |
| `--duration` | `-d` | 10 | Duration (seconds) |
| `--sample-rate` | `-s` | 48000 | Sample rate (Hz) |
| `--channels` | `-C` | 2 | Number of channels |
| `--frequency` | `-f` | 440 | Frequency (Hz) |
| `--output` | `-o` | output.wav | Output file path (any format supported by ffmpeg) |

### JSON Output

Add `--json` to output results as JSON to stdout.

```bash
$ encfixture video --json --tl frame --tr timecode -d 5 -o test.mp4
{"status":"ok","file":"test.mp4","type":"video","width":1920,"height":1080,"fps":30,"duration":"5"}
```

## Supported Colors

By name: `black`, `white`, `red`, `green`, `blue`, `yellow`, `cyan`, `magenta`, `gray`

By hex code: `#ff6600`, `#333333`, etc.

## Architecture

Built with clean architecture.

```
encfixture/
├── main.go
├── domain/              # Entities and value objects
├── usecase/             # Application logic and port interfaces
├── infrastructure/      # ffmpeg execution and image rendering
└── interface/cli/       # CLI adapter (cobra)
```

## Usage with Claude Code

encfixture supports `--json` output and descriptive `--help`, making it easy for Claude Code to generate dummy media files via Bash tool calls.

### Example: Ask Claude Code to generate test fixtures

```
> Generate a 5-second 720p test video with frame counter and timecode overlay

Claude runs:
  encfixture video --json --tl frame --tr timecode -W 1280 -H 720 -d 5 -o test_720p.mp4
```

```
> Create a set of test images: color bar, solid blue, and one with the filename displayed

Claude runs:
  encfixture image --json -b test -o colorbar.png
  encfixture image --json -c blue -o blue.png
  encfixture image --json --center filename -o sample.png
```

```
> Generate a 3-second 1000Hz beep audio file

Claude runs:
  encfixture audio --json -t sine -f 1000 -d 3 -o beep.wav
```

### CLAUDE.md integration

Add the following to your project's `CLAUDE.md` to let Claude Code know about encfixture:

```markdown
## Tools

- `encfixture` is available for generating dummy media files (image, video, audio) for ffmpeg encoding tests.
  - Always use `--json` flag for structured output.
  - Use `encfixture <subcommand> --help` for available flags.
```

## Development

```bash
# Clone
git clone https://github.com/junara/encfixture.git
cd encfixture

# Build
go build -o encfixture .

# Run
./encfixture --help

# Lint (all linters enabled)
go tool golangci-lint run ./...

# Test
go test ./...
```

## License

MIT
