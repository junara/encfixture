---
title: クイックスタート
description: encfixture の基本的な使い方
---

## 最初のファイルを生成する

### 画像

```bash
encfixture image -o test.png
```

黒い 1920x1080 の PNG 画像が生成されます。

### 動画

```bash
encfixture video --tl frame --tr timecode -d 5 -o test.mp4
```

フレーム番号とタイムコードが表示される 5 秒の動画が生成されます。

### 音声

```bash
encfixture audio -t sine -f 1000 -d 3 -o beep.wav
```

1000Hz のサイン波が 3 秒間鳴る WAV ファイルが生成されます。

## オーバーレイ

画像・動画では、5 つの位置にテキストを自由に配置できます。

```
┌──────────────────────────────┐
│ --tl              --tr       │
│                              │
│          --center            │
│                              │
│ --bl              --br       │
└──────────────────────────────┘
```

各位置には以下を指定できます：

| 値 | 説明 |
|---|---|
| `frame` | フレーム番号（動画では動的、画像では `0`） |
| `timecode` | タイムコード `HH:MM:SS:FF` |
| `filename` | 出力ファイル名 |
| その他 | そのまま表示（任意テキスト） |

### 全位置にオーバーレイを配置する例

```bash
encfixture video --tl frame --tr timecode --bl filename --br "CLIP-001" --center "SAMPLE" -d 10 -o full.mp4
```
