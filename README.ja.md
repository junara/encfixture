# encfixture

日本語 | [English](README.md)

ffmpeg エンコードテスト用のダミー素材（画像・動画・音声）を生成する Go CLI ツールです。

## 必要条件

- Go 1.22+
- ffmpeg（動画・音声の生成に必要）

## インストール

```bash
go install github.com/junara/encfixture@latest
```

## 使い方

### グローバルフラグ

| フラグ | 説明 |
|---|---|
| `--json` | 結果を JSON で出力 |
| `--version` | バージョンを表示 |

### オーバーレイの仕組み

画像・動画では、5つの位置にテキストを自由に配置できます。

```
┌──────────────────────────────┐
│ --tl              --tr       │
│                              │
│          --center            │
│                              │
│ --bl              --br       │
└──────────────────────────────┘
```

各位置には以下を指定できます：

| 値 | 説明 |
|---|---|
| `frame` | フレーム番号（動画では動的にカウント、画像では `0`） |
| `timecode` | タイムコード `HH:MM:SS:FF`（動画では動的、画像では `00:00:00:00`） |
| `filename` | 出力ファイル名 |
| その他の文字列 | そのまま表示（任意テキスト） |

### 画像の生成

```bash
# 単色画像（黒、1920x1080）
encfixture image -o test.png

# 青色の画像
encfixture image -c blue -o blue.png

# カラーバーのテストパターン
encfixture image -b test -o colorbar.png

# 全位置にオーバーレイを配置
encfixture image --tl frame --tr timecode --bl filename --br "ID-001" --center "TEST" -o info.png

# ファイル名だけ表示
encfixture image --center filename -o sample.png

# 解像度を指定
encfixture image -W 3840 -H 2160 -c white -o 4k.png

# 16進カラーコードで色を指定
encfixture image -c "#ff6600" -o orange.png

# カラーバー + テキストオーバーレイ
encfixture image -b test --center "SAMPLE" -o test_with_text.png

# JSON 出力
encfixture image --json --center "TEST" -o test.png
```

#### image フラグ

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

### 動画の生成

```bash
# 黒い無音動画（10秒、1080p、30fps）
encfixture video -o test.mp4

# フレーム番号 + タイムコードを表示
encfixture video --tl frame --tr timecode -d 5 -o counter.mp4

# 全位置にオーバーレイを配置
encfixture video --tl frame --tr timecode --bl filename --br "CLIP-001" --center "SAMPLE" -d 10 -o full.mp4

# ファイル名だけ中央に表示
encfixture video --center filename -d 5 -o sample.mov

# カラーバー背景 + オーバーレイ
encfixture video -b test --tl frame --tr timecode -d 5 -o colorbar.mp4

# サイン波音声付きの動画
encfixture video -c blue -a sine --frequency 1000 --center "BEEP" -o beep.mp4

# WebM 形式
encfixture video --tl frame -d 5 -o test.webm

# 解像度と FPS を指定
encfixture video -W 3840 -H 2160 -r 60 -d 10 --tl frame -o 4k60.mp4

# JSON 出力
encfixture video --json --tl frame -d 3 -o test.mp4
```

#### video フラグ

| フラグ | 短縮 | デフォルト | 説明 |
|---|---|---|---|
| `--width` | `-W` | 1920 | 動画の幅（px） |
| `--height` | `-H` | 1080 | 動画の高さ（px） |
| `--fps` | `-r` | 30 | フレームレート |
| `--duration` | `-d` | 10 | 長さ（秒） |
| `--bg` | `-b` | solid | 背景タイプ: solid, test |
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

### 音声の生成

```bash
# 無音の WAV（10秒）
encfixture audio -o silence.wav

# サイン波（440Hz）
encfixture audio -t sine -d 5 -o sine.wav

# 1000Hz のトーン
encfixture audio -t tone -f 1000 -d 3 -o tone.wav

# ホワイトノイズ
encfixture audio -t noise -d 5 -o noise.wav

# モノラル、44100Hz
encfixture audio -t silence -C 1 -s 44100 -o mono.wav

# MP3 形式
encfixture audio -t sine -d 5 -o beep.mp3

# FLAC 形式
encfixture audio -t silence -d 10 -o silence.flac

# JSON 出力
encfixture audio --json -t sine -d 3 -o beep.wav
```

#### audio フラグ

| フラグ | 短縮 | デフォルト | 説明 |
|---|---|---|---|
| `--type` | `-t` | silence | 音声タイプ: silence, sine, noise, tone |
| `--duration` | `-d` | 10 | 長さ（秒） |
| `--sample-rate` | `-s` | 48000 | サンプルレート（Hz） |
| `--channels` | `-C` | 2 | チャンネル数 |
| `--frequency` | `-f` | 440 | 周波数（Hz） |
| `--output` | `-o` | output.wav | 出力ファイルパス（ffmpeg 対応の任意フォーマット） |

### JSON 出力

`--json` フラグを付けると、結果を JSON で標準出力に出力します。

```bash
$ encfixture video --json --tl frame --tr timecode -d 5 -o test.mp4
{"status":"ok","file":"test.mp4","type":"video","width":1920,"height":1080,"fps":30,"duration":"5"}
```

## 対応色

名前で指定: `black`, `white`, `red`, `green`, `blue`, `yellow`, `cyan`, `magenta`, `gray`

16進カラーコードで指定: `#ff6600`, `#333333` など

## アーキテクチャ

クリーンアーキテクチャを採用しています。

```
encfixture/
├── main.go
├── domain/              # エンティティ・値オブジェクト
├── usecase/             # アプリケーションロジック・ポートインターフェース
├── infrastructure/      # ffmpeg実行・画像レンダリング実装
└── interface/cli/       # CLIアダプター（cobra）
```

## Claude Code での活用

encfixture は `--json` 出力と分かりやすい `--help` に対応しており、Claude Code の Bash ツールから簡単にダミーメディアファイルを生成できます。

### 例: Claude Code にテストフィクスチャの生成を依頼

```
> 5秒の720pテスト動画をフレームカウンターとタイムコード付きで作って

Claude が実行:
  encfixture video --json --tl frame --tr timecode -W 1280 -H 720 -d 5 -o test_720p.mp4
```

```
> カラーバー、青色単色、ファイル名表示の3種類のテスト画像を作って

Claude が実行:
  encfixture image --json -b test -o colorbar.png
  encfixture image --json -c blue -o blue.png
  encfixture image --json --center filename -o sample.png
```

```
> 3秒の1000Hzビープ音を作って

Claude が実行:
  encfixture audio --json -t sine -f 1000 -d 3 -o beep.wav
```

### CLAUDE.md への統合

プロジェクトの `CLAUDE.md` に以下を追加すると、Claude Code が encfixture を認識して使えるようになります:

```markdown
## Tools

- `encfixture` でダミーメディアファイル（画像・動画・音声）を生成できます。ffmpeg エンコードテスト用。
  - 常に `--json` フラグを使って構造化出力を取得してください。
  - `encfixture <subcommand> --help` で利用可能なフラグを確認できます。
```

## 開発

```bash
# リポジトリのクローン
git clone https://github.com/junara/encfixture.git
cd encfixture

# ビルド
go build -o encfixture .

# 実行
./encfixture --help

# lint（全リンター有効）
go tool golangci-lint run ./...

# テスト
go test ./...
```

## ライセンス

MIT
