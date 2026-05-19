---
id: T03
parent: S01
milestone: M009-zjrq6j
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m009-s01.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T17:30:15.683Z
blocker_discovered: false
---

# T03: Verified S01 end-to-end: benchmark artifact now contains a sanitized effective config snapshot and no raw benchmark texts.

**Verified S01 end-to-end: benchmark artifact now contains a sanitized effective config snapshot and no raw benchmark texts.**

## What Happened

Ran full S01 verification. `docker compose config` succeeded. Go short tests passed with 49 tests across 4 packages. The full benchmark ran via uv/Python 3.13 and wrote `benchmark-results/fd-benchmark-m009-s01.txt`. The artifact includes the new sanitized config snapshot with model/API/dimensions, git metadata, Docker compose hash/images, allowlisted environment, Redis INFO summary, environment baseline hash, and redaction policy. The parser check confirmed the snapshot exists, required fields are present, raw Russian benchmark texts are not printed, and no secret-like env keys are included. The benchmark also now uses corrected label character counts (`short 17`, `medium 73`, `long 422`, `very_long 693`).

## Verification

Fresh verification passed after final edits: compile/snapshot parser, docker compose config, Go tests, full uv benchmark, and artifact parser/redaction check.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests --with redis python -m py_compile benchmark.py && uv run --python 3.13 --with requests --with redis python - <<'PY' ...` | 0 | ✅ pass: compile and snapshot parser passed | 4900ms |
| 2 | `docker compose config >/tmp/fd-compose-config-m009-s01.txt` | 0 | ✅ pass | 4900ms |
| 3 | `cd api && go test ./... -short` | 0 | ✅ pass: 49 passed in 4 packages | 4800ms |
| 4 | `uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-m009-s01.txt` | 0 | ✅ pass: benchmark completed and artifact written | 23100ms |
| 5 | `uv run --python 3.13 --with requests --with redis python - <<'PY' ...` | 0 | ✅ pass: artifact snapshot parser passed | 7300ms |

## Deviations

The first artifact parser check was overly strict because it rejected the word `PASSWORD` inside the printed redaction policy. The final parser correctly verifies no secret-like environment keys or secret values are emitted while allowing the redaction policy to document its patterns.

## Known Issues

The benchmark snapshot reports `git.dirty: true` during local verification because the current slice changes are intentionally uncommitted. After commit, future runs should show clean state when no local changes exist.

## Files Created/Modified

- `benchmark.py`
- `benchmark-results/fd-benchmark-m009-s01.txt`
