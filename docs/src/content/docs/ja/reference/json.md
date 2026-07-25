---
title: JSON 出力
description: --json フラグによる構造化出力
---

## 使い方

`--json` フラグを付けると、結果を JSON で標準出力に出力します。

```bash
encfixture video --json --tl frame --tr timecode -d 5 -o test.mp4
```

出力:

```json
{"status":"ok","file":"test.mp4","type":"video","width":1920,"height":1080,"fps":30,"duration":"5"}
```

## レスポンスフィールド

| フィールド | 型 | 説明 |
|---|---|---|
| `status` | string | 成功時は `"ok"`、失敗時は `"error"` |
| `file` | string | 出力ファイルパス |
| `type` | string | `"image"`, `"video"`, `"audio"` |
| `width` | int | 幅（image/video のみ） |
| `height` | int | 高さ（image/video のみ） |
| `fps` | int | フレームレート（video のみ） |
| `duration` | string | 長さ（video/audio のみ） |

## 各コマンドの例

### image

```bash
$ encfixture image --json --center "TEST" -o test.png
{"status":"ok","file":"test.png","type":"image","width":1920,"height":1080}
```

### video

```bash
$ encfixture video --json --tl frame -d 3 -o test.mp4
{"status":"ok","file":"test.mp4","type":"video","width":1920,"height":1080,"fps":30,"duration":"3"}
```

### audio

```bash
$ encfixture audio --json -t sine -d 3 -o beep.wav
{"status":"ok","file":"beep.wav","type":"audio","duration":"3"}
```

### エラー

コマンドが失敗した場合、`--json` 指定時は構造化されたエラーオブジェクトを標準出力に出力します（人間向けメッセージは従来どおり標準エラー出力）。終了コードは非ゼロです。`code` は安定した機械可読の値なので、メッセージを解析せずに `code` で分岐してください。

```bash
$ encfixture video --json --codec av1 -o test.mp4
{"status":"error","code":"encoder_not_available","error":"...Unknown encoder 'libaom-av1'","hint":"Run 'encfixture doctor' to list the encoders your ffmpeg build supports, then pick an available --codec."}
```

| フィールド | 型 | 説明 |
|---|---|---|
| `status` | string | 常に `"error"` |
| `code` | string | 安定した機械可読のエラーコード（下表） |
| `error` | string | エラーメッセージ |
| `hint` | string | 復旧方法（有用な情報が無い場合は省略） |

| code | 意味 |
|---|---|
| `usage` | フラグ・引数の誤り |
| `ffmpeg_not_found` / `ffprobe_not_found` | ツールが PATH に無い |
| `encoder_not_available` | 要求したエンコーダがこの ffmpeg ビルドに無い |
| `unknown_codec` / `unknown_background` | `--codec` / `--bg` の値が不正 |
| `invalid_duration` / `invalid_bitrate` | `-d` / `--bitrate` の値が不正 |
| `invalid_expectation` | `--expect` の書式が不正 |
| `verify_failed` | `--expect` のアサーションが失敗 |
| `output_exists` | `--no-clobber` 下で上書きを拒否 |
| `probe_failed` | ffprobe がファイルを読めない |
| `ffmpeg_failed` | ffmpeg が分類外のエラーで終了 |
| `env_unhealthy` | `doctor` が ffmpeg/ffprobe の欠落を検出 |
| `error` | その他 |

`--json` なしの場合も、同じ hint がエラーメッセージの後に `hint:` 行として標準エラー出力に表示されます。

### batch

`batch` コマンドはジョブごとの結果と集計を含む集約オブジェクトを出力します。詳細は [バッチ処理](/encfixture/ja/usage/batch/) を参照してください。

```bash
$ encfixture batch --json jobs.json
{"results":[{"index":0,"type":"image","file":"a.png","status":"ok"}],"succeeded":1,"failed":0}
```
