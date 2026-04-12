---
title: Audio Generation
description: How to use the encfixture audio command
---

## Basic

```bash
encfixture audio -o silence.wav
```

## Examples

```bash
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

# FLAC format
encfixture audio -t silence -d 10 -o silence.flac
```

## Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--type` | `-t` | silence | Audio type: silence, sine, noise, tone |
| `--duration` | `-d` | 10 | Duration (seconds) |
| `--sample-rate` | `-s` | 48000 | Sample rate (Hz) |
| `--channels` | `-C` | 2 | Number of channels |
| `--frequency` | `-f` | 440 | Frequency (Hz) |
| `--output` | `-o` | output.wav | Output file path (any format supported by ffmpeg) |
