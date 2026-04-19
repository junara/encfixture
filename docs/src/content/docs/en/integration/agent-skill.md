---
title: Agent Skill
description: Install encfixture as an Agent Skill for Claude Code, Copilot, Cursor, and other AI coding agents
---

encfixture ships an [Agent Skill](https://agentskills.io) (`SKILL.md`) and an [APM](https://microsoft.github.io/apm/) package. Installing it into your project gives AI coding agents pre-loaded knowledge of encfixture's flags, overlay keywords, batch schema, and common recipes — so you don't have to explain them every time.

## What's in the package

- **Skill** (`skills/encfixture/SKILL.md`) — the on-demand usage guide loaded when an agent decides encfixture is relevant.
- **Prompts** (via APM only):
  - `generate-test-fixtures` — build a standard encode-test bundle (1080p30 / 4K60 / webm / mkv / sine audio / thumbnail).
  - `generate-video-with-overlay` — one video with frame / timecode / filename overlays, parameters configurable.
  - `generate-batch-spec` — interactively assemble a `fixtures.json` from free-form requirements, then run it.

## Install with `gh skill` (GitHub CLI)

Works for Claude Code, GitHub Copilot, Cursor, Codex, and Gemini CLI.

```bash
# Claude Code (project scope — default)
gh skill install junara/encfixture encfixture --agent claude-code

# GitHub Copilot at user scope
gh skill install junara/encfixture encfixture --agent github-copilot --scope user

# Cursor
gh skill install junara/encfixture encfixture --agent cursor

# Pin to a release tag
gh skill install junara/encfixture encfixture@v1.0.0 --agent claude-code
```

After install, the skill lives in the agent's native directory:

| Agent | Project scope | User scope |
|---|---|---|
| Claude Code | `.claude/skills/encfixture/` | `~/.claude/skills/encfixture/` |
| GitHub Copilot | `.agents/skills/encfixture/` | `~/.copilot/skills/encfixture/` |
| Cursor | `.agents/skills/encfixture/` | `~/.cursor/skills/encfixture/` |
| Codex | `.agents/skills/encfixture/` | `~/.codex/skills/encfixture/` |
| Gemini CLI | `.agents/skills/encfixture/` | `~/.gemini/skills/encfixture/` |

Update with:

```bash
gh skill update --all
```

`gh skill` installs only the skill (no prompts). Use APM if you want the prompts as well.

## Install with `apm`

Gets skill + prompts in one step.

```bash
# From a project that consumes encfixture
apm install junara/encfixture

# Pin to a release
apm install junara/encfixture#v1.0.0
```

APM deploys to every agent directory present in the project (`.claude/`, `.github/`, `.cursor/`, `.opencode/`). Track dependencies in `apm.yml`:

```yaml
dependencies:
  apm:
    - junara/encfixture#v1.0.0
```

Add `apm_modules/` to `.gitignore` and commit `apm.lock.yaml` for reproducible installs.

## Prerequisites

The skill documents how to call the `encfixture` CLI — it does not include the binary. Install the CLI separately (see [Installation](/encfixture/en/getting-started/installation/)). `ffmpeg` must also be on `PATH`.
