---
title: 検査（verify）
description: encfixture verify コマンドでメディアファイルを検査する
---

## 概要

`verify` は既存のメディアファイルを `ffprobe` で検査し、コンテナと各ストリームの情報を表示します。生成したファイル（や任意のファイル）が期待どおりのコーデック・解像度・fps・尺かを確認でき、CI やエージェントからの「生成 → 検証」フローに向いています。

```bash
encfixture verify test.mp4
```

## 出力

```
$ encfixture verify test.mp4
File:     test.mp4
Format:   mov,mp4,m4a,3gp,3g2,mj2
Duration: 2.000000s
Size:     28934 bytes
Stream 0: video  h264  1920x1080  30fps  yuv420p
Stream 1: audio  aac  48000Hz  2ch
```

## JSON 出力

`--json` を付けると機械可読の結果を返します。

```bash
$ encfixture verify --json test.mp4
{"file":"test.mp4","format":{"formatName":"mov,mp4,...","duration":"2.000000","size":"28934","bitRate":"115736"},"streams":[{"index":0,"type":"video","codec":"h264","width":1920,"height":1080,"fps":"30","pixFmt":"yuv420p"},{"index":1,"type":"audio","codec":"aac","sampleRate":"48000","channels":2}]}
```

### フィールド

| フィールド | 説明 |
|---|---|
| `format.formatName` | ffprobe が報告するコンテナ形式 |
| `format.duration` | 尺（秒） |
| `format.size` | ファイルサイズ（バイト） |
| `format.bitRate` | 全体のビットレート |
| `streams[].type` | `video` または `audio` |
| `streams[].codec` | コーデック名（例: `h264`, `aac`） |
| `streams[].width` / `height` | 動画の解像度 |
| `streams[].fps` | フレームレート（`30`, `29.970` など） |
| `streams[].pixFmt` | ピクセルフォーマット |
| `streams[].sampleRate` / `channels` | 音声のサンプルレートとチャンネル数 |

## 必要環境

`verify` には ffmpeg に同梱の `ffprobe` が必要です。ファイルが存在しない・読み取れない場合は ffprobe のエラーメッセージを表示し、非ゼロで終了します。
