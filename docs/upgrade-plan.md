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

## Production Preflight Gate

The user owns the actual production deployment. Before asking the user to deploy
an upgraded build, perform extra read-only production checks and report the
result clearly. Do not restart, recreate, or deploy production services unless
the user explicitly asks for that operation.

Required checks before the final "ready to deploy" handoff:

1. Inspect the production Compose merge result, not only individual Compose
   files:
   ```bash
   sudo docker compose \
     --env-file /root/new-api/.env.prod \
     -f /root/new-api/docker-compose.yml \
     -f /root/new-api/docker-compose.prod.yml \
     config
   ```
2. Check whether production override files change important runtime settings
   from the base Compose file, especially:
   - `REDIS_CONN_STRING`
   - `SQL_DSN`
   - published ports
   - volume mounts
   - image names and tags
   - service commands
3. Compare live Redis requirements with the API container's effective Redis
   connection string:
   ```bash
   sudo docker inspect redis --format '{{json .Args}}'
   sudo docker inspect nebula-api --format '{{range .Config.Env}}{{println .}}{{end}}' | grep REDIS
   ```
   If Redis uses `--requirepass`, `REDIS_CONN_STRING` must include the password.
4. Check live container state and recent fatal logs:
   ```bash
   sudo docker ps
   sudo docker logs --tail 100 nebula-api
   ```
5. Probe public health endpoints without mutating production state:
   ```bash
   curl -fsS https://ai.nebulatrip.com/api/status
   curl -fsS https://image-api.nebulatrip.com/api/status
   ```

If any check shows a mismatch, stop before the deployment handoff and explain the
risk. The May 24, 2026 outage was caused by this class of issue: the base Compose
file had the correct password-bearing Redis URL, but `docker-compose.prod.yml`
overrode it with the old `redis://redis` value while Redis required
`--requirepass 123456`.

### Stage 5 Production Notes

For the `v0.12.15 -> v0.13.2` deployment handoff, check these items before the
user deploys:

- Database compatibility is expected to be additive: this stage introduces
  fields such as `payment_provider` on payment records and user audit fields
  such as creation and last-login timestamps. Existing rows should remain
  readable, and old payment rows without `payment_provider` are intentionally
  guarded by legacy order-prefix compatibility logic.
- Token key length is relaxed for old-token migration compatibility. This
  should not invalidate existing tokens, but token login and token list pages
  must be checked after deployment.
- `gpt-5.5` support is mainly billing-ratio recognition, not automatic channel
  replacement. Production channels that already manually include `gpt-5.5`
  should remain, but model ratio, completion ratio, and model pricing should be
  reviewed after deployment.
- Brave Search channel, `/v1/search`, and Firecrawl-compatible `/v2/search`
  paths were removed from `main-plus` because they are no longer used.
- If deployment fails before the app becomes healthy, roll back the image/code
  version first. The additive database columns can normally remain in place for
  the old version, but do not manually drop production columns during incident
  recovery.

## Plan

| Done | Step | From | To | Status | Notes |
|---|---:|---|---|---|---|
| [x] | 0 | `v0.11.8` | `v0.11.8` | Baseline confirmed | `main-plus` merge-base is `v0.11.8`. |
| [x] | 1 | `v0.11.8` | `v0.11.9` | Completed | Merged `v0.11.9`; removed obsolete fingerprint/Claude OAuth UI paths; `go test ./...` and `mise run fe-build` passed. |
| [x] | 2 | `v0.11.9` | `v0.12.5` | Completed | Merged `v0.12.5`; removed obsolete OpenClaw/Hermes token config route and UI; conflict handling followed `C2-01` through `C2-05`; `go test ./...` and `mise run fe-build` passed. |
| [x] | 3 | `v0.12.5` | `v0.12.10` | Completed | Merged `v0.12.10`; conflicts handled per `C3-01` through `C3-05`; kept payment-provider isolation while adding Stripe async webhook handling; `go test ./...` and `mise run fe-build` passed. |
| [x] | 4 | `v0.12.10` | `v0.12.15` | Completed | Merged `v0.12.15`; conflicts handled per `C4-01` through `C4-08`; kept payment-provider isolation, direct WeChat/Alipay, and Claude empty-text handling while adding Waffo Pancake/payment log/Codex/Gemini/passkey updates; `go test ./...` and `mise run fe-build` passed. |
| [x] | 5 | `v0.12.15` | `v0.13.2` | Completed | Merged `v0.13.2`; conflicts handled per `C5-01` through `C5-07`; kept payment-provider isolation and direct WeChat/Alipay compatibility, removed Brave Search routes/channel, added tiered billing/model sync updates; `go test ./...` and `mise run fe-build` passed. |
| [ ] | 6 | `v0.13.2` | `v1.0.0-rc.8` | Pending conflict review | Explain conflicts first, then wait for confirmation before upgrading. |

## Verification Checklist

Run the relevant checks after each completed stage:

- Confirm the worktree contains only expected changes.
- Run backend tests or at least a backend build.
- Run the frontend production build.
- Verify relay, billing, subscription/topup, and OpenClaw/Hermes config paths did not regress.
- Run the production preflight gate before telling the user the upgrade is ready for production deployment.
- Record the completed upgrade commit and verification result in the plan table.
