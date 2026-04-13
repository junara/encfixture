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
| `status` | string | 常に `"ok"` |
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

### batch

`batch` コマンドはジョブごとの結果と集計を含む集約オブジェクトを出力します。詳細は [バッチ処理](/encfixture/ja/usage/batch/) を参照してください。

```bash
$ encfixture batch --json jobs.json
{"results":[{"index":0,"type":"image","file":"a.png","status":"ok"}],"succeeded":1,"failed":0}
```
