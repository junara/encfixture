---
title: 画像の生成
description: encfixture image コマンドの使い方
---

## 基本

```bash
encfixture image -o test.png
```

## 使用例

```bash
# 青色の画像
encfixture image -c blue -o blue.png

# カラーバーのテストパターン
encfixture image -b test -o colorbar.png

# 全位置にオーバーレイを配置
encfixture image --tl frame --tr timecode --bl filename --br "ID-001" --center "TEST" -o info.png

# 解像度を指定
encfixture image -W 3840 -H 2160 -c white -o 4k.png

# 16進カラーコードで色を指定
encfixture image -c "#ff6600" -o orange.png

# カラーバー + テキストオーバーレイ
encfixture image -b test --center "SAMPLE" -o test_with_text.png
```

## フラグ

| フラグ | 短縮 | デフォルト | 説明 |
|---|---|---|---|
| `--width` | `-W` | 1920 | 画像の幅（px） |
| `--height` | `-H` | 1080 | 画像の高さ（px） |
| `--bg` | `-b` | solid | 背景タイプ: solid, test |
| `--color` | `-c` | black | 背景色（名前または #hex） |
| `--tl` | | | 左上に表示する内容 |
| `--tr` | | | 右上に表示する内容 |
| `--center` | | | 中央に表示する内容 |
| `--bl` | | | 左下に表示する内容 |
| `--br` | | | 右下に表示する内容 |
| `--scale` | `-S` | 4 | テキストの拡大率 |
| `--output` | `-o` | output.png | 出力ファイルパス |
