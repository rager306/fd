---
id: T02
parent: S01
milestone: M002-d2au23
key_files:
  - .gitignore
  - docker-compose.yaml
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T07:20:55.215Z
blocker_discovered: false
---

# T02: Removed Compose obsolete-version warning and ignored local runtime noise.

**Removed Compose obsolete-version warning and ignored local runtime noise.**

## What Happened

Removed the obsolete top-level Compose `version` field. Added narrow .gitignore rules for local runtime artifacts: .bg-shell, transient GSD state/audit/event/exec/graphs/journal/notifications/runtime/state-manifest files, and SQLite WAL/SHM files. Durable `.gsd/gsd.db`, `.gsd/milestones/**`, and `.gsd/quick/**` remain unignored and trackable.

## Verification

`docker compose config` completed without the obsolete-version warning; `git status --short` no longer shows runtime noise.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker compose config >/tmp/fd-compose-clean.txt 2>/tmp/fd-compose-clean.err && if grep -q 'obsolete' /tmp/fd-compose-clean.err; then echo 'obsolete warning still present'; cat /tmp/fd-compose-clean.err; exit 1; fi && echo 'compose config clean' && git status --short` | 0 | ✅ pass: compose config clean and runtime noise hidden | 0ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `.gitignore`
- `docker-compose.yaml`
