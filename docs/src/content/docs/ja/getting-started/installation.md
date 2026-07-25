---
title: インストール
description: encfixture のインストール方法
---

## 必要条件

- ffmpeg

## Homebrew（macOS / Linux）

```bash
brew install junara/tap/encfixture
```

## Go

```bash
go install github.com/junara/encfixture@latest
```

## バイナリ

[Releases](https://github.com/junara/encfixture/releases) からお使いの OS に合ったバイナリをダウンロードしてください。

## Skill のインストール（AI コーディングエージェント向け）

encfixture は CLI とは別に [Agent Skill](https://agentskills.io)（`SKILL.md`）を同梱しています。導入しておくと、Claude Code などのエージェントがフラグ・オーバーレイ予約語・batch スキーマを把握した状態で encfixture を使えます。CLI 本体とは別物なので、上記のいずれかで CLI をインストールしたうえで追加してください。

### Claude Code プラグインとして（推奨）

このリポジトリは Claude Code のプラグインマーケットプレイスとしても機能します。Claude Code 内で:

```
/plugin marketplace add junara/encfixture
/plugin install encfixture@encfixture
```

`/encfixture:encfixture` で明示的に起動できます（関連するタスクでは自動で読み込まれます）。

### `gh skill` で（Claude Code / Copilot / Cursor / Codex / Gemini CLI）

```bash
gh skill install junara/encfixture encfixture --agent claude-code
```

### `apm` で（Skill + Prompt をまとめて）

```bash
apm install junara/encfixture
```

インストール先ディレクトリ、バージョン固定、更新・アンインストールなどの詳細は [Agent Skill](/encfixture/ja/integration/agent-skill/) を参照してください。複数の方法を併用するとスキルが二重に読み込まれるため、いずれか 1 つを選んでください。
