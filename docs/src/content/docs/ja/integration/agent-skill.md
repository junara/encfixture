---
title: Agent Skill
description: encfixture を Agent Skill として Claude Code / Copilot / Cursor などのコーディングエージェントにインストールする方法
---

encfixture は [Agent Skill](https://agentskills.io)（`SKILL.md`）と [APM](https://microsoft.github.io/apm/) パッケージを同梱しています。プロジェクトにインストールしておくと、AI コーディングエージェントが encfixture のフラグ・オーバーレイ予約語・batch スキーマ・典型的な使い方を毎回説明せずに扱えるようになります。

## パッケージの中身

- **Skill**（`skills/encfixture/SKILL.md`）— エージェントが encfixture を使うと判断したときに読み込まれる使い方ガイド。
- **Prompts**（APM 経由のみ）:
  - `generate-test-fixtures` — エンコードテスト用の標準素材セット（1080p30 / 4K60 / webm / mkv / サイン波音声 / サムネ）を生成。
  - `generate-video-with-overlay` — frame / timecode / filename オーバーレイ付きの動画 1 本を生成（パラメータ可変）。
  - `generate-batch-spec` — 自然言語の要件から `fixtures.json` を組み立てて実行。

## `gh skill` でインストール（GitHub CLI）

Claude Code / GitHub Copilot / Cursor / Codex / Gemini CLI に対応しています。

```bash
# Claude Code（プロジェクトスコープ、既定）
gh skill install junara/encfixture encfixture --agent claude-code

# GitHub Copilot をユーザースコープで
gh skill install junara/encfixture encfixture --agent github-copilot --scope user

# Cursor
gh skill install junara/encfixture encfixture --agent cursor

# リリースタグに固定
gh skill install junara/encfixture encfixture@v1.0.0 --agent claude-code
```

インストール後、スキルは各エージェントのネイティブディレクトリに配置されます:

| エージェント | プロジェクトスコープ | ユーザースコープ |
|---|---|---|
| Claude Code | `.claude/skills/encfixture/` | `~/.claude/skills/encfixture/` |
| GitHub Copilot | `.agents/skills/encfixture/` | `~/.copilot/skills/encfixture/` |
| Cursor | `.agents/skills/encfixture/` | `~/.cursor/skills/encfixture/` |
| Codex | `.agents/skills/encfixture/` | `~/.codex/skills/encfixture/` |
| Gemini CLI | `.agents/skills/encfixture/` | `~/.gemini/skills/encfixture/` |

アップデート:

```bash
gh skill update --all
```

`gh skill` はスキルのみをインストールします（prompt は含まれません）。prompt も欲しい場合は APM を使ってください。

## `apm` でインストール

スキルと prompt をまとめて取得できます。

```bash
# encfixture を使うプロジェクトの中で
apm install junara/encfixture

# リリースに固定
apm install junara/encfixture#v1.0.0
```

APM はプロジェクト内に存在するエージェント分だけ `.claude/` / `.github/` / `.cursor/` / `.opencode/` にデプロイします。依存は `apm.yml` で管理します:

```yaml
dependencies:
  apm:
    - junara/encfixture#v1.0.0
```

`apm_modules/` を `.gitignore` に追加し、`apm.lock.yaml` をコミットすると再現性が得られます。

## 前提

スキル本体は `encfixture` CLI の使い方を記述したドキュメントであり、バイナリは含まれません。CLI は別途インストールが必要です（[インストール](/encfixture/ja/getting-started/installation/)を参照）。`ffmpeg` も `PATH` に必要です。
