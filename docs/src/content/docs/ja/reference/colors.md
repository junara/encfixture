---
title: 対応色
description: 利用可能な色の指定方法
---

## 名前で指定

| 名前 | 色 |
|---|---|
| `black` | 黒 |
| `white` | 白 |
| `red` | 赤 |
| `green` | 緑 |
| `blue` | 青 |
| `yellow` | 黄 |
| `cyan` | シアン |
| `magenta` | マゼンタ |
| `gray` / `grey` | 灰色 |

## 16進カラーコードで指定

`#RRGGBB` 形式で任意の色を指定できます。

```bash
encfixture image -c "#ff6600" -o orange.png
encfixture image -c "#333333" -o dark_gray.png
encfixture video -c "#0066ff" -d 5 -o blue.mp4
```

## テキストの色

テキストの色は背景色に応じて自動的に白または黒が選択されます（コントラスト自動調整）。

- 暗い背景 → 白いテキスト
- 明るい背景 → 黒いテキスト
