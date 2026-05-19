---
id: T01
parent: S04
milestone: M003-xx4yc3
key_files:
  - benchmark.py
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T08:21:10.497Z
blocker_discovered: false
---

# T01: Benchmark Python dependencies are available through a temp PYTHONPATH target.

**Benchmark Python dependencies are available through a temp PYTHONPATH target.**

## What Happened

Prepared the benchmark Python environment. System Python is 3.12.3. The global environment had `requests` but lacked `redis`. Creating a venv failed because ensurepip is unavailable. Installed `requests` and `redis` into `/tmp/fd-benchmark-deps` using `pip --target` and verified imports via `PYTHONPATH`. Versions: requests 2.34.2, redis 7.4.0.

## Verification

`PYTHONPATH=/tmp/fd-benchmark-deps python3` imports requests and redis successfully.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 imports requests/redis` | 1 | ❌ initial fail: redis module missing | 0ms |
| 2 | `python3 -m venv /tmp/fd-benchmark-venv` | 1 | ⚠️ unavailable: ensurepip/python3-venv missing | 6000ms |
| 3 | `python3 -m pip install --target /tmp/fd-benchmark-deps requests redis && PYTHONPATH=/tmp/fd-benchmark-deps python3 -c imports` | 0 | ✅ pass: requests 2.34.2, redis 7.4.0 | 26500ms |

## Deviations

`python3 -m venv` failed because `ensurepip`/python3-venv is unavailable. Used `pip --target /tmp/fd-benchmark-deps` and `PYTHONPATH` instead to avoid modifying the repo or system packages beyond a temp dependency target.

## Known Issues

pip emitted a root-user warning; dependencies were installed to `/tmp/fd-benchmark-deps`, not into the project tree.

## Files Created/Modified

- `benchmark.py`
