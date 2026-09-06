---
title: ファイル容量の指定
description: --target-size で画像・動画・音声をおおよその容量で生成する
---

アップロード上限・ストレージ枠・転送時間など、ファイル容量そのものが試験対象になるケース向けに、`--target-size` で出力容量を指定できます。

値は `K`/`M`/`G` のサフィックスを付けられ、2進接頭辞として解釈します（`1KB` = 1024 バイト、`10MB` = 10485760 バイト。`10M` / `10MB` / `10MiB` はすべて同じ意味です）。サフィックスなしの数値はバイト数です。

## 基本

```bash
# 20MB の動画（容量 ÷ 長さからビットレートを算出）
encfixture video --target-size 20MB -d 30 -o sized.mp4

# 2MB の音声
encfixture audio -t sine -d 60 --target-size 2MB -o sized.mp3

# 500KB の JPEG（収まる品質を二分探索で決定）
encfixture image -b test --target-size 500KB -o sized.jpg
```

## 精度について

容量は **近似値** です。動画・音声はパディングではなくビットレートの逆算で合わせるため、どこまで近づくかはエンコーダに依存します。

### 動画

`.mp4` / `.mov` / `.mkv` / `.m4v` に `--target-size` を指定するとコーデックは h264 になり、CBR で固定します。黒一色の映像でも数%以内に収まります。

```bash
$ encfixture video --target-size 5MB -d 10 -o sized.mp4 --json
{"status":"ok","file":"sized.mp4","type":"video","width":1920,"height":1080,"fps":30,"duration":"10","size":5144863,"targetSize":5242880}
```

他のコーデック（`hevc`, `vp9`, `av1`, `prores`）はビットレート上限のみが効くため、圧縮が効きすぎる映像では大幅に下振れします。

```bash
$ encfixture video --codec hevc --target-size 5MB -d 10 -o small.mp4
Created: small.mp4
warning: output is 61.71KB, 99% off the requested 5.00MB; to get closer, raise --width/--height or --fps, use a busier --bg (test, gradient, moving), or pick --codec h264, which tracks a target bitrate closest
```

### 音声

圧縮形式にビットレートを設定します。`.mp3` はほぼ指定どおりになりますが、AAC や Opus は無音をどれだけビットレートを与えても数 kbit/s まで潰すため、`-t noise` や `-t sine` と組み合わせてください。

PCM / ロスレス（`.wav`, `.aiff`, `.au`, `.flac`）はビットレートで制御できないためエラーになります。容量はサンプルレートで決まるので、その容量になるサンプルレートをエラーメッセージで案内します。

```bash
$ encfixture audio --target-size 1MB -d 10 -o out.wav
encfixture: audio generation failed: target size not supported for this format: .wav is not bitrate-controlled (16-bit PCM reaches 1.00MB at --sample-rate 26214); write a compressed format such as .mp3, .m4a or .opus
```

### 画像

JPEG は収まる品質を選ぶ方式なので `.jpg` / `.jpeg` 出力が必要です（PNG はロスレスで、容量を動かすつまみがありません）。上限はその解像度における品質 100 の容量です。もっと大きくしたい場合は `--width`/`--height` を上げるか、`--bg` を情報量の多いパターンにしてください。

## 生成結果の容量

生成結果には実際の容量が含まれ、指定から 10% 以上ずれた場合は警告します。`--json` なら機械可読なので、CI でそのまま検証できます。

| フィールド | 型 | 説明 |
|---|---|---|
| `size` | int | 生成されたファイルの容量（バイト） |
| `targetSize` | int | 指定した容量（バイト。`--target-size` 指定時のみ） |
| `warning` | string | 指定から 10% 以上ずれた場合に付与 |

## 併用できないフラグ

`--target-size` は `--bitrate` / `--crf`（video）、`--quality`（image）と同じつまみを操作するため、併用するとエラーになります。

```bash
$ encfixture video --target-size 5MB --bitrate 1M -o out.mp4
encfixture: conflicting flags: --target-size and --bitrate cannot be combined
```

## バッチ

バッチジョブでは `targetSize` 文字列として指定します（ジョブ単位でも `defaults` でも指定できます）。同じジョブで `bitrate` や `crf` を併記した場合は、どちらかが暗黙に勝つのではなく `target_size_conflict` エラーになります。

```json
{
  "defaults": {"targetSize": "1MB", "duration": "5"},
  "jobs": [
    {"type": "video", "output": "a.mp4"},
    {"type": "audio", "output": "b.m4a", "audio": "noise", "targetSize": "200KB"},
    {"type": "image", "output": "c.jpg", "bg": "gradient", "targetSize": "100KB"}
  ]
}
```

## エラーコード

| コード | 意味 |
|---|---|
| `invalid_target_size` | `--target-size` の値が不正 |
| `target_size_unsupported` | 容量を制御できない形式（PNG, WAV, FLAC など） |
| `target_size_too_small` | 容量が小さすぎてエンコーダが扱えるビットレートを下回る |
| `target_size_conflict` | バッチジョブで `targetSize` と `bitrate` / `crf` を同時指定 |
