# S01: Compose and git hygiene — UAT

**Milestone:** M002-d2au23
**Written:** 2026-05-19T07:21:51.817Z

# UAT: S01 Compose and git hygiene

## Verification performed

- `docker compose config >/tmp/fd-compose-clean.txt 2>/tmp/fd-compose-clean.err` followed by grep for `obsolete` — passed.
- `git check-ignore` confirmed `.bg-shell/`, `.gsd/runtime/`, `.gsd/exec/`, `.gsd/journal/`, `.gsd/audit/`, `.gsd/graphs/`, `.gsd/gsd.db-shm`, and `.gsd/gsd.db-wal` are ignored.
- `git check-ignore` produced no matches for `.gsd/gsd.db`, `.gsd/milestones/...`, and `.gsd/quick/...`, confirming durable artifacts remain trackable.

