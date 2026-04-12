---
title: オーバーレイ
description: テキストオーバーレイの仕組みと使い方
---

## 配置

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

## 予約語

| 値 | 説明 |
|---|---|
| `frame` | フレーム番号（動画では動的にカウント、画像では `0`） |
| `timecode` | タイムコード `HH:MM:SS:FF`（動画では動的、画像では `00:00:00:00`） |
| `filename` | 出力ファイル名 |

予約語以外の文字列はそのまま任意テキストとして表示されます。

## 使用例

### 画像にファイル名を表示

```bash
encfixture image --center filename -o sample.png
```

### 動画にフレーム番号とタイムコードを表示

```bash
encfixture video --tl frame --tr timecode -d 10 -o counter.mp4
```

### 全位置にオーバーレイを配置

```bash
encfixture video --tl frame --tr timecode --bl filename --br "CLIP-001" --center "SAMPLE" -d 10 -o full.mp4
```

### カラーバー背景と組み合わせ

```bash
encfixture video -b test --tl frame --tr timecode --center "TEST" -d 5 -o colorbar.mp4
```

## テキストの拡大率

`--scale`（`-S`）フラグでテキストの大きさを変更できます。デフォルトは `4` です。

```bash
# 大きいテキスト
encfixture image --center "BIG" -S 8 -o big.png

# 小さいテキスト
encfixture image --center "small" -S 2 -o small.png
```
