---
title: 音声の生成
description: encfixture audio コマンドの使い方
---

## 基本

```bash
encfixture audio -o silence.wav
```

## 使用例

```bash
# サイン波（440Hz）
encfixture audio -t sine -d 5 -o sine.wav

# 1000Hz のトーン
encfixture audio -t tone -f 1000 -d 3 -o tone.wav

# ホワイトノイズ
encfixture audio -t noise -d 5 -o noise.wav

# モノラル、44100Hz
encfixture audio -t silence -C 1 -s 44100 -o mono.wav

# MP3 形式
encfixture audio -t sine -d 5 -o beep.mp3

# FLAC 形式
encfixture audio -t silence -d 10 -o silence.flac
```

## フラグ

| フラグ | 短縮 | デフォルト | 説明 |
|---|---|---|---|
| `--type` | `-t` | silence | 音声タイプ: silence, sine, noise, tone |
| `--duration` | `-d` | 10 | 長さ（秒） |
| `--sample-rate` | `-s` | 48000 | サンプルレート（Hz） |
| `--channels` | `-C` | 2 | チャンネル数 |
| `--frequency` | `-f` | 440 | 周波数（Hz） |
| `--output` | `-o` | output.wav | 出力ファイルパス（ffmpeg 対応の任意フォーマット） |
