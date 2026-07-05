---
title: 動画の生成
description: encfixture video コマンドの使い方
---

## 基本

```bash
encfixture video -o test.mp4
```

## 使用例

```bash
# フレーム番号 + タイムコードを表示
encfixture video --tl frame --tr timecode -d 5 -o counter.mp4

# 全位置にオーバーレイを配置
encfixture video --tl frame --tr timecode --bl filename --br "CLIP-001" --center "SAMPLE" -d 10 -o full.mp4

# カラーバー背景 + オーバーレイ
encfixture video -b test --tl frame --tr timecode -d 5 -o colorbar.mp4

# 動きのある背景（コーデックの圧縮特性テスト用）
encfixture video -b gradient -d 5 -o gradient.mp4
encfixture video -b moving --tr timecode -d 5 -o moving.mp4

# サイン波音声付き
encfixture video -c blue -a sine --frequency 1000 --center "BEEP" -o beep.mp4

# WebM 形式
encfixture video --tl frame -d 5 -o test.webm

# 解像度と FPS を指定
encfixture video -W 3840 -H 2160 -r 60 -d 10 --tl frame -o 4k60.mp4

# コーデックと品質を指定（HEVC、CRF 28）
encfixture video --codec hevc --crf 28 --tl frame -d 5 -o hevc.mp4

# 編集ワークフロー向け ProRes（デフォルトで 10bit 4:2:2）
encfixture video --codec prores -d 5 -o prores.mov

# 固定ビットレートとピクセルフォーマットを指定
encfixture video --codec h264 --bitrate 5M --pix-fmt yuv420p -d 5 -o cbr.mp4

# A/V シンクテストパターン（毎秒ビープ音 + 白フラッシュ）
encfixture video --sync --tr timecode -d 10 -o sync.mp4

# マーカー間隔を 0.5 秒に指定
encfixture video --sync --sync-interval 0.5 -d 5 -o sync_fast.mp4
```

## フラグ

| フラグ | 短縮 | デフォルト | 説明 |
|---|---|---|---|
| `--width` | `-W` | 1920 | 動画の幅（px） |
| `--height` | `-H` | 1080 | 動画の高さ（px） |
| `--fps` | `-r` | 30 | フレームレート |
| `--duration` | `-d` | 10 | 長さ（秒） |
| `--bg` | `-b` | solid | 背景タイプ: solid, test, gradient, moving |
| `--color` | `-c` | black | 背景色（名前または #hex） |
| `--tl` | | | 左上に表示する内容 |
| `--tr` | | | 右上に表示する内容 |
| `--center` | | | 中央に表示する内容 |
| `--bl` | | | 左下に表示する内容 |
| `--br` | | | 右下に表示する内容 |
| `--scale` | `-S` | 4 | テキストの拡大率 |
| `--output` | `-o` | output.mp4 | 出力ファイルパス（ffmpeg 対応の任意フォーマット） |
| `--audio` | `-a` | silence | 音声タイプ: silence, sine, noise, tone |
| `--sample-rate` | `-s` | 48000 | 音声サンプルレート |
| `--channels` | `-C` | 2 | 音声チャンネル数 |
| `--frequency` | | 440 | 音声の周波数（Hz） |
| `--codec` | | | 動画コーデック: h264, hevc, vp9, av1, prores（デフォルト: コンテナ既定） |
| `--crf` | | | CRF 値（h264/hevc/vp9/av1 向け。デフォルト: エンコーダ既定） |
| `--bitrate` | | | 動画ビットレート（例: `5M`, `800k`。デフォルト: エンコーダ既定） |
| `--pix-fmt` | | | ピクセルフォーマット（例: `yuv420p`, `yuv422p10le`。デフォルト: コーデック依存） |
| `--sync` | | false | A/V シンクテストパターン: 各マーカーでビープ音と白フラッシュが同時に発生 |
| `--sync-interval` | | 1.0 | シンクマーカーの間隔（秒） |

## A/V シンクテストパターン

`--sync` は音ズレ検出用のパターンを生成します。各間隔（`--sync-interval`、既定1秒）の先頭で、画面全体の白フラッシュとビープ音が同時に約0.08秒間発生します。再生時にビープとフラッシュがずれて見えれば、そのファイル（またはそれを生成したパイプライン）に A/V ずれがあります。

- ビープの音程は `--frequency`（既定440Hz）で決まります。
- `--tr timecode` などのオーバーレイはそのまま描画されるので、どのマーカーでズレたか特定できます。
- `--sync` は `--audio` より優先されます（ビープ音が選択した音声ソースを置き換えます）。

```bash
encfixture video --sync --tr timecode -d 10 -o sync.mp4
```

## 背景タイプ

`--bg` でフレームを埋める内容を選びます。

| 値 | 説明 |
|---|---|
| `solid` | `--color` の単色塗り（既定）。ffmpeg の `color` ソースで生成。 |
| `test` | 静止した SMPTE 風カラーバー。 |
| `gradient` | 時間とともに横方向へスクロールする斜めの色勾配。 |
| `moving` | `--color` の背景上を往復移動する白いボックス。 |

`gradient` と `moving` は実際の動きを加えるため、コーデックの動き推定を働かせ、圧縮対象として意味のある内容になります（単色や静止フレームはほぼ無に圧縮される）。どちらも決定論的で、同じフラグなら常に同じフレームを生成します。フレーム単位で描画するため、オーバーレイや `--sync` と併用できます。

## コーデックの選択

`--codec` は以下の ffmpeg エンコーダにマッピングされます。

| 値 | エンコーダ | 備考 |
|---|---|---|
| `h264` | libx264 | |
| `hevc` | libx265 | |
| `vp9` | libvpx-vp9 | `.webm` 出力時のデフォルト |
| `av1` | libaom-av1 | 低速なので短い動画向け |
| `prores` | prores_ks | ピクセルフォーマットは `yuv422p10le` がデフォルト |

`--codec` 未指定時はコンテナの既定コーデックが使われます（`.mp4` なら H.264、`.webm` なら VP9）。
