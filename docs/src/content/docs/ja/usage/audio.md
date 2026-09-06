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

# ファイル容量を指定（容量と長さからビットレートを逆算）
encfixture audio -t sine -d 60 --target-size 2MB -o sized.mp3

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
| `--target-size` | | | 出力ファイル容量の目安。例: `2MB`（圧縮形式のみ。[ファイル容量の指定](/encfixture/ja/usage/target-size/)を参照） |
| `--no-clobber` | | | 既存の出力ファイルがある場合は上書きせずエラーにする |
