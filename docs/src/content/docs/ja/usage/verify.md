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

## アサーション（`--expect`）

`--expect key=value` を 1 つ以上渡すと、verify がワンショットの合否判定になります。各期待値を検査結果と突き合わせ、1 つでも失敗すれば非ゼロ（エラーコード `verify_failed`）で終了します。

```bash
encfixture verify test.mp4 --expect codec=h264 --expect width=1920 --expect duration=5
```

```
$ encfixture verify test.mp4 --expect codec=hevc --expect duration=5
File:     test.mp4
...
FAIL  codec: expected hevc, actual h264
PASS  duration = 5.016000
encfixture: verification failed: 1 of 2 expectations failed
```

| キー | 検査対象 |
|---|---|
| `codec`, `width`, `height`, `fps`, `pixFmt` | 最初の映像ストリーム |
| `audioCodec`, `sampleRate`, `channels` | 最初の音声ストリーム |
| `duration` | コンテナの尺 |

キーは大文字小文字を区別せず、ケバブケースの別名（`pix-fmt`, `audio-codec`, `sample-rate`）も使えます。数値キーは `5+-0.2`（または `5±0.2`）の許容誤差サフィックスを受け付けます。デフォルトでは `duration` が ±0.1 秒（コンテナの丸め分）、`fps` が ±0.001（`fps=29.97` が ffprobe の `29.970` にマッチ）で、それ以外は完全一致です。対象ストリームが無い場合は `(no video stream)` / `(no audio stream)` を実測値として失敗します。

## JSON 出力

`--json` を付けると機械可読の結果を返します。

```bash
$ encfixture verify --json test.mp4
{"status":"ok","file":"test.mp4","format":{"formatName":"mov,mp4,...","duration":"2.000000","size":"28934","bitRate":"115736"},"streams":[{"index":0,"type":"video","codec":"h264","width":1920,"height":1080,"fps":"30","pixFmt":"yuv420p"},{"index":1,"type":"audio","codec":"aac","sampleRate":"48000","channels":2}]}
```

`--expect` を付けた場合は `checks` 配列が追加され、`status` は `ok` / `failed` になります:

```bash
$ encfixture verify --json test.mp4 --expect codec=hevc
{"status":"failed","file":"test.mp4","format":{...},"streams":[...],"checks":[{"field":"codec","expected":"hevc","actual":"h264","pass":false}]}
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
