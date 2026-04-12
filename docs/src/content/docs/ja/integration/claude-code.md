---
title: Claude Code
description: Claude Code での encfixture の活用方法
---

encfixture は `--json` 出力と分かりやすい `--help` に対応しており、Claude Code の Bash ツールから簡単にダミーメディアファイルを生成できます。

## 使用例

### テスト動画の生成

```
> 5秒の720pテスト動画をフレームカウンターとタイムコード付きで作って

Claude が実行:
  encfixture video --json --tl frame --tr timecode -W 1280 -H 720 -d 5 -o test_720p.mp4
```

### テスト画像の一括生成

```
> カラーバー、青色単色、ファイル名表示の3種類のテスト画像を作って

Claude が実行:
  encfixture image --json -b test -o colorbar.png
  encfixture image --json -c blue -o blue.png
  encfixture image --json --center filename -o sample.png
```

### テスト音声の生成

```
> 3秒の1000Hzビープ音を作って

Claude が実行:
  encfixture audio --json -t sine -f 1000 -d 3 -o beep.wav
```

## CLAUDE.md への統合

プロジェクトの `CLAUDE.md` に以下を追加すると、Claude Code が encfixture を認識して使えるようになります:

```markdown
## Tools

- `encfixture` でダミーメディアファイル（画像・動画・音声）を生成できます。ffmpeg エンコードテスト用。
  - 常に `--json` フラグを使って構造化出力を取得してください。
  - `encfixture <subcommand> --help` で利用可能なフラグを確認できます。
```
