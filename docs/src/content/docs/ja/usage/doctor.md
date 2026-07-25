---
title: 環境診断（doctor）
description: encfixture doctor で ffmpeg/ffprobe とエンコーダの対応状況を診断する
---

## 概要

`doctor` は encfixture が依存する環境を診断します。`ffmpeg` / `ffprobe` の有無と、選択可能なエンコーダ(`--codec` の各値と、mp4/webm 出力に使う音声エンコーダ)のうちローカルの ffmpeg ビルドが対応しているものを一覧します。環境が不明なときは生成の前に実行してください。たとえば Homebrew の ffmpeg には `libaom-av1` が入っていないことがよくあります。

```bash
encfixture doctor
```

## 出力

```
$ encfixture doctor
ffmpeg:  OK 7.1 (/opt/homebrew/bin/ffmpeg)
ffprobe: OK 7.1 (/opt/homebrew/bin/ffprobe)
video encoders:
  h264     libx264      OK
  hevc     libx265      OK
  vp9      libvpx-vp9   OK
  av1      libaom-av1   MISSING
  prores   prores_ks    OK
audio encoders:
  aac      aac          OK
  opus     libopus      OK
```

## JSON 出力

```bash
$ encfixture doctor --json
{"status":"ok","ffmpeg":{"name":"ffmpeg","available":true,"version":"7.1","path":"/opt/homebrew/bin/ffmpeg"},"ffprobe":{"name":"ffprobe","available":true,"version":"7.1","path":"/opt/homebrew/bin/ffprobe"},"videoEncoders":[{"codec":"h264","encoder":"libx264","available":true},{"codec":"av1","encoder":"libaom-av1","available":false}],"audioEncoders":[{"codec":"aac","encoder":"aac","available":true},{"codec":"opus","encoder":"libopus","available":true}]}
```

| フィールド | 説明 |
|---|---|
| `status` | ffmpeg と ffprobe が揃っていれば `ok`、欠けていれば `error` |
| `ffmpeg` / `ffprobe` | `available`(有無)・`path`(解決されたパス)・`version` |
| `videoEncoders[]` | `--codec` の各値と、対応する ffmpeg エンコーダの利用可否 |
| `audioEncoders[]` | `aac`(mp4/mov のデフォルト)と `opus`(webm 出力用) |

## 終了コード

`doctor` が非ゼロで終了するのは `ffmpeg` か `ffprobe` が無いときだけです(エラーコード `env_unhealthy`)。エンコーダの欠落は情報提供のみで、必要とする生成コマンドが個別に `encoder_not_available` エラーコードで失敗し、その hint が `doctor` を案内します。
