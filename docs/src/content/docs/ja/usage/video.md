---
title: 動画の生成
description: encfixture video コマンドの使い方
---

## 基本

```bash
encfixture video -o test.mp4
```

## 使用例

```bash
# フレーム番号 + タイムコードを表示
encfixture video --tl frame --tr timecode -d 5 -o counter.mp4

# 全位置にオーバーレイを配置
encfixture video --tl frame --tr timecode --bl filename --br "CLIP-001" --center "SAMPLE" -d 10 -o full.mp4

# カラーバー背景 + オーバーレイ
encfixture video -b test --tl frame --tr timecode -d 5 -o colorbar.mp4

# サイン波音声付き
encfixture video -c blue -a sine --frequency 1000 --center "BEEP" -o beep.mp4

# WebM 形式
encfixture video --tl frame -d 5 -o test.webm

# 解像度と FPS を指定
encfixture video -W 3840 -H 2160 -r 60 -d 10 --tl frame -o 4k60.mp4
```

## フラグ

| フラグ | 短縮 | デフォルト | 説明 |
|---|---|---|---|
| `--width` | `-W` | 1920 | 動画の幅（px） |
| `--height` | `-H` | 1080 | 動画の高さ（px） |
| `--fps` | `-r` | 30 | フレームレート |
| `--duration` | `-d` | 10 | 長さ（秒） |
| `--bg` | `-b` | solid | 背景タイプ: solid, test |
| `--color` | `-c` | black | 背景色（名前または #hex） |
| `--tl` | | | 左上に表示する内容 |
| `--tr` | | | 右上に表示する内容 |
| `--center` | | | 中央に表示する内容 |
| `--bl` | | | 左下に表示する内容 |
| `--br` | | | 右下に表示する内容 |
| `--scale` | `-S` | 4 | テキストの拡大率 |
| `--output` | `-o` | output.mp4 | 出力ファイルパス（ffmpeg 対応の任意フォーマット） |
| `--audio` | `-a` | silence | 音声タイプ: silence, sine, noise, tone |
| `--sample-rate` | `-s` | 48000 | 音声サンプルレート |
| `--channels` | `-C` | 2 | 音声チャンネル数 |
| `--frequency` | | 440 | 音声の周波数（Hz） |
