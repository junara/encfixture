---
title: Installation
description: How to install encfixture
---

## Requirements

- ffmpeg

## Homebrew (macOS / Linux)

```bash
brew install junara/tap/encfixture
```

## Go

```bash
go install github.com/junara/encfixture@latest
```

## Binary

Download the binary for your OS from [Releases](https://github.com/junara/encfixture/releases).

## Install the skill (for AI coding agents)

Separately from the CLI, encfixture ships an [Agent Skill](https://agentskills.io) (`SKILL.md`). Installing it lets agents such as Claude Code use encfixture with its flags, overlay keywords, and batch schema already in hand. The skill is documentation, not a binary — install the CLI with one of the methods above first.

### As a Claude Code plugin (recommended)

The repository doubles as a Claude Code plugin marketplace. Inside Claude Code:

```
/plugin marketplace add junara/encfixture
/plugin install encfixture@encfixture
```

Invoke it explicitly with `/encfixture:encfixture` (agents also load it automatically when relevant).

### With `gh skill` (Claude Code / Copilot / Cursor / Codex / Gemini CLI)

```bash
gh skill install junara/encfixture encfixture --agent claude-code
```

### With `apm` (skill + prompts)

```bash
apm install junara/encfixture
```

For install locations, version pinning, updates, and uninstall, see [Agent Skill](/encfixture/en/integration/agent-skill/). Pick a single method — installing with more than one loads the skill twice.
