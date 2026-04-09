# Kiro IDE

← [Back to README](../README.md)

---

This document explains how gentle-ai integrates with **Kiro IDE** and what is installed in your local Kiro configuration.

## Overview

gentle-ai supports Kiro as a **solo-agent** platform (`kiro-ide`).

When configured, gentle-ai installs:

| Artifact | Path |
|----------|------|
| System instructions | `<GlobalConfigDir>/prompts/gentle-ai.instructions.md` |
| Skills directory | `<GlobalConfigDir>/skills/` |
| MCP config | `~/.kiro/settings/mcp.json` *(separate root — see note below)* |

> **Auto-install not supported.** Kiro must be installed manually before running gentle-ai.
> Download from: [kiro.dev/downloads](https://kiro.dev/downloads)

---

## Detection

gentle-ai detects Kiro by resolving `kiro` from `PATH`:

```
exec.LookPath("kiro")
```

If `kiro` is not on `PATH`, Kiro will **not** be auto-detected during install. Make sure the Kiro binary is accessible in your shell environment before running `gentle-ai install`.

---

## SDD Execution Model

Kiro runs as a **solo-agent** platform — there is no custom sub-agent delegation.

All SDD phases run **inline in the same conversation**:

```
propose → spec → design → tasks → apply → verify → archive
```

The orchestrator IS the executor. Cross-phase persistence is handled by Engram, which lets each phase retrieve prior artifacts from memory rather than relying on conversation history alone.

**Approval gates** are required before the `apply` and `archive` phases — the orchestrator will pause and ask for confirmation before proceeding.

---

## Native Kiro Specs Integration

Kiro has a built-in spec workflow that gentle-ai leverages. For medium and large changes, the orchestrator will use native Kiro artifacts at:

```
.kiro/specs/<feature>/
├── requirements.md
├── design.md
└── tasks.md
```

**Steering files** at `.kiro/steering/*.md` provide persistent workspace context across sessions — treat them like always-on system context for your project conventions, architecture decisions, and team rules.

**Size classification** routes tasks through Small / Medium / Large paths to decide planning depth:

| Size | Approach |
|------|----------|
| Small | Inline — no formal SDD phases |
| Medium | Kiro native specs (`.kiro/specs/`) + Engram |
| Large | Full SDD cycle: explore → propose → spec → design → tasks → apply → verify → archive |

---

## Config Paths by Platform

### macOS

| Artifact | Path |
|----------|------|
| Global config dir | `~/Library/Application Support/Kiro/User` |
| System prompt file | `~/Library/Application Support/Kiro/User/prompts/gentle-ai.instructions.md` |
| Skills dir | `~/Library/Application Support/Kiro/User/skills/` |
| Settings path | `~/Library/Application Support/Kiro/User/settings.json` |
| MCP config | `~/.kiro/settings/mcp.json` |

### Windows

| Artifact | Path |
|----------|------|
| Global config dir | `%APPDATA%\kiro\User` |
| System prompt file | `%APPDATA%\kiro\User\prompts\gentle-ai.instructions.md` |
| Skills dir | `%APPDATA%\kiro\User\skills\` |
| Settings path | `%APPDATA%\kiro\User\settings.json` |
| MCP config | `%USERPROFILE%\.kiro\settings\mcp.json` |

### Linux (XDG)

| Artifact | Path |
|----------|------|
| Global config dir | `$XDG_CONFIG_HOME/kiro/user` *(fallback: `~/.config/kiro/user`)* |
| System prompt file | `$XDG_CONFIG_HOME/kiro/user/prompts/gentle-ai.instructions.md` |
| Skills dir | `$XDG_CONFIG_HOME/kiro/user/skills/` |
| Settings path | `$XDG_CONFIG_HOME/kiro/user/settings.json` |
| MCP config | `~/.kiro/settings/mcp.json` |

---

## ⚠️ MCP Path Separation

Kiro stores user config and MCP config in **different root directories**.

- **User config** (prompts, skills, settings) lives under the platform-native Kiro User dir (`~/Library/Application Support/Kiro/User`, `%APPDATA%\kiro\User`, or `$XDG_CONFIG_HOME/kiro/user`)
- **MCP config** is always at `~/.kiro/settings/mcp.json` (or `%USERPROFILE%\.kiro\settings\mcp.json` on Windows), regardless of platform

If MCP tools are not loading, check the **separate `.kiro/settings/` root** rather than the Kiro User config dir.

---

## Capability Snapshot

| Capability | Status |
|------------|--------|
| Skills | ✅ Yes |
| System prompt | ✅ Yes |
| MCP | ✅ Yes |
| Output styles | ❌ No |
| Slash commands | ❌ No |
| Delegation model | Solo-agent |
| Auto-install | ❌ No — manual install required |
