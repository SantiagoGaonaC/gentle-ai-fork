# Kiro IDE

← [Back to README](../README.md)

---

This document explains how gentle-ai integrates with **Kiro IDE** and what is installed in your local Kiro configuration.

## Overview

gentle-ai supports Kiro as a **native-subagent** platform (`kiro-ide`).

When configured, gentle-ai installs:

| Artifact | Path |
|----------|------|
| Steering file | `~/.kiro/steering/gentle-ai.md` |
| Native SDD agents | `~/.kiro/agents/sdd-{phase}.md` *(10 files)* |
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

Kiro runs with **native sub-agent delegation** via `~/.kiro/agents/`.

The orchestrator stays in the steering file and coordinates phase execution, while each phase runs in its dedicated Kiro agent file:

```
sdd-init → sdd-explore → sdd-propose → sdd-spec → sdd-design → sdd-tasks → sdd-apply → sdd-verify → sdd-archive (+ sdd-onboard)
```

This follows the same SDD architecture used in gentle-ai: orchestrator coordinates, phase agents execute, Engram persists artifacts across phases.

**Approval gates** remain required before `apply` and `archive`.

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

## Steering File Format

The steering file written by gentle-ai uses the following frontmatter:

```yaml
---
inclusion: always
---
```

`inclusion: always` ensures Kiro loads this context in every conversation automatically, regardless of workspace or file type.

## Native Agent Frontmatter

Kiro SDD phase agents are generated with YAML frontmatter including:

- `name`
- `description`
- `tools`
- `model`
- `includeMcpJson: true`

The `model` value is injected during sync from Claude alias assignments (`opus|sonnet|haiku`) to Kiro-native model IDs.

---

## Config Paths by Platform

### macOS

| Artifact | Path |
|----------|------|
| Global config dir | `~/Library/Application Support/Kiro/User` |
| Steering file | `~/.kiro/steering/gentle-ai.md` |
| Skills dir | `~/.kiro/skills/` |
| Settings path | `~/Library/Application Support/Kiro/User/settings.json` |
| MCP config | `~/.kiro/settings/mcp.json` |

### Windows

| Artifact | Path |
|----------|------|
| Global config dir | `%APPDATA%\kiro\User` |
| Steering file | `%USERPROFILE%\.kiro\steering\gentle-ai.md` |
| Skills dir | `%USERPROFILE%\.kiro\skills\` |
| Settings path | `%APPDATA%\kiro\User\settings.json` |
| MCP config | `%USERPROFILE%\.kiro\settings\mcp.json` |

### Linux (XDG)

| Artifact | Path |
|----------|------|
| Global config dir | `$XDG_CONFIG_HOME/kiro/user` *(fallback: `~/.config/kiro/user`)* |
| Steering file | `~/.kiro/steering/gentle-ai.md` |
| Skills dir | `~/.kiro/skills/` |
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
| Delegation model | Full (native subagents) |
| Auto-install | ❌ No — manual install required |
