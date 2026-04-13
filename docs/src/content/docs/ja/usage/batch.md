---
title: バッチ処理
description: encfixture batch コマンドで複数ファイルを一括生成
---

## 概要

`batch` サブコマンドは、JSON ファイルに定義した複数のジョブを一括で実行します。CI や網羅テスト、解像度違いのサンプル一括生成などに便利です。

```bash
encfixture batch jobs.json
```

## JSON スキーマ

トップレベルのオブジェクトは任意の `defaults` と必須の `jobs` を持ちます。`defaults` の値は各ジョブで上書き可能です。

```json
{
  "defaults": {
    "width": 1920,
    "height": 1080
  },
  "jobs": [
    { "type": "video", "output": "clip.mp4", "duration": "5", "tl": "frame", "tr": "timecode" },
    { "type": "image", "output": "thumb.png", "bg": "test" },
    { "type": "audio", "output": "beep.wav", "audio": "sine", "frequency": 1000 }
  ]
}
```

各ジョブは `type` と `output` が必須で、残りのフィールドは対応するサブコマンドのフラグと同じ意味です。未知のフィールドはエラーになるため typo の早期発見に役立ちます。

## ジョブフィールド

| フィールド | 対応フラグ | 型 | 適用タイプ |
|---|---|---|---|
| `type` | — | `"image"` / `"video"` / `"audio"` | 必須 |
| `output` | `--output` | string | 必須 |
| `width` | `--width` | int | image, video |
| `height` | `--height` | int | image, video |
| `fps` | `--fps` | int | video |
| `duration` | `--duration` | string | video, audio |
| `bg` | `--bg` | string | image, video |
| `color` | `--color` | string | image, video |
| `tl` / `tr` / `center` / `bl` / `br` | `--tl` ほか | string | image, video |
| `scale` | `--scale` | int | image, video |
| `audio` | `--audio` / `--type` | string | video, audio |
| `sampleRate` | `--sample-rate` | int | video, audio |
| `channels` | `--channels` | int | video, audio |
| `frequency` | `--frequency` | float | video, audio |

## CLI フラグ

| フラグ | 短縮 | デフォルト | 説明 |
|---|---|---|---|
| `--parallel` | `-p` | `NumCPU/2`（最低 1） | 同時実行ジョブ数の上限 |
| `--fail-fast` | — | false | 最初の失敗以降、未着手のジョブをスキップ |
| `--json` | — | false | 結果を JSON で標準出力に出力 |

### 並列度の目安

ffmpeg は 1 プロセス内部で複数スレッドを使う CPU バウンドな処理なので、やみくもに並列度を上げるとかえって遅くなります。以下を目安にしてください。

| ジョブの傾向 | 推奨値 | 理由 |
|---|---|---|
| 動画（高解像度・長尺）中心 | `1` 〜 `2` | 1 本の ffmpeg が既に全コアを使う。並列しても取り合いになる |
| 動画（低解像度・短尺）中心 | `NumCPU/2`（デフォルト） | エンコード時間が短く、並列化の恩恵がある |
| 画像のみ | `NumCPU` | Go 側のレンダリング主体で ffmpeg を呼ばないため高並列が有効 |
| 音声のみ | `NumCPU/2` 〜 `NumCPU` | 処理が軽く I/O 待ちが支配的 |
| 混在（CI） | デフォルト | 中庸な値で十分。必要なら実測して調整 |

目安を超えた値を指定しても害はなく、内部のセマフォで上限が守られるだけです。最適値はマシン・コーデック・解像度に依存するため、時間が重要な用途では実測してください。

## 使用例

### 並列度を指定して実行

```bash
encfixture batch jobs.json --parallel 4
```

### 最初の失敗で残りをスキップ

```bash
encfixture batch jobs.json --fail-fast
```

### CI 向けに JSON で結果を受け取る

```bash
encfixture batch jobs.json --json
```

出力:

```json
{
  "results": [
    { "index": 0, "type": "video", "file": "clip.mp4", "status": "ok" },
    { "index": 1, "type": "image", "file": "thumb.png", "status": "ok" },
    { "index": 2, "type": "audio", "file": "beep.wav", "status": "error", "error": "..." }
  ],
  "succeeded": 2,
  "failed": 1
}
```

## 挙動

- ジョブの実行順序は保証されません（並列実行のため）。結果配列は入力順で返ります。
- `--fail-fast` 指定時、実行中のジョブは最後まで完走します。未着手のジョブは `job skipped: context canceled` エラーとして結果に記録されます。
- 失敗が 1 件でもあると終了コードは 1 になります。
