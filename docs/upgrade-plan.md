# Nebula API Upgrade Plan

This document tracks the staged `main-plus` upgrade path from `v0.11.8` to
`v1.0.0-rc.8`.

## Upgrade Gate

Before starting each upgrade stage:

1. Inspect the diff and expected conflicts for that exact `From -> To` range.
2. Read the upstream GitHub release notes for the target range and summarize:
   - Newly added or changed user-facing features.
   - Bug fixes and behavior corrections.
   - Operationally relevant migrations, defaults, or compatibility changes.
3. Write a detailed conflict and risk summary for the user.
4. Give every conflict a stable short ID, such as `C2-01`, where the first
   number is the plan step and the second number is the conflict sequence.
   Use those IDs in all follow-up discussion and resolution decisions.
5. Wait for explicit user confirmation on the conflict handling plan.
6. Only then apply the upgrade.

After a stage is completed, update the matching row from `[ ]` to `[x]` and add
the upgrade commit, verification result, or remaining issue in `Notes`.

## Plan

| Done | Step | From | To | Status | Notes |
|---|---:|---|---|---|---|
| [x] | 0 | `v0.11.8` | `v0.11.8` | Baseline confirmed | `main-plus` merge-base is `v0.11.8`. |
| [x] | 1 | `v0.11.8` | `v0.11.9` | Completed | Merged `v0.11.9`; removed obsolete fingerprint/Claude OAuth UI paths; `go test ./...` and `mise run fe-build` passed. |
| [ ] | 2 | `v0.11.9` | `v0.12.5` | Pending conflict review | Explain conflicts first, then wait for confirmation before upgrading. |
| [ ] | 3 | `v0.12.5` | `v0.12.10` | Pending conflict review | Explain conflicts first, then wait for confirmation before upgrading. |
| [ ] | 4 | `v0.12.10` | `v0.12.15` | Pending conflict review | Explain conflicts first, then wait for confirmation before upgrading. |
| [ ] | 5 | `v0.12.15` | `v0.13.2` | Pending conflict review | Explain conflicts first, then wait for confirmation before upgrading. |
| [ ] | 6 | `v0.13.2` | `v1.0.0-rc.8` | Pending conflict review | Explain conflicts first, then wait for confirmation before upgrading. |

## Verification Checklist

Run the relevant checks after each completed stage:

- Confirm the worktree contains only expected changes.
- Run backend tests or at least a backend build.
- Run the frontend production build.
- Verify relay, billing, subscription/topup, and OpenClaw/Hermes config paths did not regress.
- Record the completed upgrade commit and verification result in the plan table.
