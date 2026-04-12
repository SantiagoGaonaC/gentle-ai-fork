<!-- ⚠️ READ BEFORE SUBMITTING
  Every PR must be linked to an issue that has the "status:approved" label.
  PRs without a linked approved issue will be automatically rejected by CI.
  See CONTRIBUTING.md for the full contribution workflow.
-->

## 🔗 Linked Issue

Closes #166

---

## 🏷️ PR Type

What kind of change does this PR introduce?

- [ ] `type:bug` — Bug fix (non-breaking change that fixes an issue)
- [x] `type:feature` — New feature (non-breaking change that adds functionality)
- [ ] `type:docs` — Documentation only
- [ ] `type:refactor` — Code refactoring (no functional changes)
- [ ] `type:chore` — Build, CI, or tooling changes
- [ ] `type:breaking-change` — Breaking change (fix or feature that changes existing behavior)

---

## 📝 Summary

This PR adds first-class Kiro IDE support aligned with the current gentle-ai SDD architecture:

- Native SDD subagents in `~/.kiro/agents/` for all phases
- Steering strategy integration in `~/.kiro/steering/gentle-ai.md`
- Dedicated Kiro model configuration flow in TUI (`Configure Kiro models`) with independent per-phase assignments
- End-to-end sync/inject plumbing so selected Kiro models update `model:` in generated Kiro subagent frontmatter
- Windows compatibility fixes for skills path resolution and cross-platform path guidance in Kiro SDD subagent prompts

Credit: this extends the initial Kiro groundwork by @calvimor.

---

## 📂 Changes

| File / Area | What Changed |
|-------------|-------------|
| `internal/agents/kiro/adapter.go` | Added native subagent support methods; fixed `SkillsDir()` to `~/.kiro/skills/` |
| `internal/model/kiro_model.go` | Added/updated alias→Kiro model ID mapping (`opus/sonnet/haiku`) |
| `internal/components/sdd/inject.go` | Added Kiro model placeholder resolution and Kiro-specific assignment support |
| `internal/assets/kiro/agents/*.md` | Added/updated 10 Kiro SDD subagent templates |
| `internal/assets/kiro/sdd-orchestrator.md` | Updated to delegation model language |
| `internal/tui/screens/model_config.go` | Added `Configure Kiro models` option (ordered after OpenCode) |
| `internal/tui/screens/kiro_model_picker.go` | New dedicated Kiro model assignment screen/state |
| `internal/tui/model.go` | Added Kiro picker flow, navigation, and sync override wiring |
| `internal/model/selection.go` | Added `KiroModelAssignments` to Selection/SyncOverrides |
| `internal/cli/run.go`, `internal/cli/sync.go`, `internal/app/app.go` | Plumbed Kiro assignments through run/sync/app override pipeline |
| `internal/assets/skills/_shared/SKILL.md` | Added metadata placeholder for Kiro skill scanner compatibility |
| `docs/kiro.md`, `docs/agents.md`, `README.md` | Updated docs to reflect native Kiro subagents and behavior |
| `testdata/golden/sdd-kiro-*.golden` | Regenerated Kiro SDD goldens |
| `*_test.go` (model/agents/components/tui/app) | Added/updated tests for adapter, injection, TUI flow, and path/model behavior |

---

## 🧪 Test Plan

**Unit Tests**
```bash
go test ./...
```

**E2E Tests** (Docker required)
```bash
cd e2e && ./docker-test.sh
```

- [x] Unit tests pass (`go test ./...`)
- [x] E2E tests pass (`cd e2e && ./docker-test.sh`)
- [x] Manually tested locally

Additional manual checks:
- Verified Kiro generation under `~/.kiro/agents/`, `~/.kiro/steering/`, and `~/.kiro/skills/`
- Verified `Configure Models -> Configure Kiro models` updates `model:` in generated Kiro subagent files
- Verified Windows MCP file-read path behavior with `%USERPROFILE%\\.kiro\\skills\\...`

---

## 🤖 Automated Checks

The following checks run automatically on this PR:

| Check | Status | Description |
|-------|--------|-------------|
| Check Issue Reference | ⏳ | PR body must contain `Closes/Fixes/Resolves #N` |
| Check Issue Has `status:approved` | ⏳ | Linked issue must have been approved before work began |
| Check PR Has `type:*` Label | ⏳ | Exactly one `type:*` label must be applied |
| Unit Tests | ⏳ | `go test ./...` must pass |
| E2E Tests | ⏳ | `cd e2e && ./docker-test.sh` must pass |

---

## ✅ Contributor Checklist

- [x] PR is linked to an issue with `status:approved`
- [ ] I have added the appropriate `type:*` label to this PR
- [x] Unit tests pass (`go test ./...`)
- [x] E2E tests pass (`cd e2e && ./docker-test.sh`)
- [x] I have updated documentation if necessary
- [x] My commits follow [Conventional Commits](https://www.conventionalcommits.org/) format
- [x] My commits do not include `Co-Authored-By` trailers

---

## 💬 Notes for Reviewers

- Kiro assignment flow is intentionally separate from Claude/OpenCode in TUI and persisted independently (`KiroModelAssignments`).
- Kiro injection prioritizes `KiroModelAssignments` with backward-compatible fallback to legacy Claude assignments when unset.
- Kiro subagent prompts now include explicit Windows/macOS/Linux skill paths to avoid MCP access errors on Windows.
