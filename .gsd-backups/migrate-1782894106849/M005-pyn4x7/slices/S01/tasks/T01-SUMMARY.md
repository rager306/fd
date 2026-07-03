---
id: T01
parent: S01
milestone: M005-pyn4x7
key_files:
  - README.md
  - docker-compose.yaml
  - docker-compose.override.yaml
  - benchmark.py
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T10:40:12.207Z
blocker_discovered: false
---

# T01: Found stale README benchmark/runtime instructions that need uv Python 3.13 and restart-side-effect documentation.

**Found stale README benchmark/runtime instructions that need uv Python 3.13 and restart-side-effect documentation.**

## What Happened

Inspected README, compose files, and benchmark.py. README's development benchmark command is stale (`python3 benchmark.py`) and does not mention uv or Python 3.13. It also does not explain that benchmark section 5 restarts API when Docker Compose is available. Compose override confirms Redis is intentionally exposed to the host only on `127.0.0.1:6379` for host benchmark access, while base compose does not publish Redis. Benchmark.py confirms it assumes API on localhost:8000 and Redis on localhost:6379, uses Docker Compose restart for Redis L2 diagnostics, and waits for API health afterward.

## Verification

Read README, docker-compose.yaml, docker-compose.override.yaml, and benchmark.py; exact stale sections identified.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `read README.md docker-compose.yaml docker-compose.override.yaml benchmark.py` | 0 | ✅ pass: stale benchmark command and missing runtime notes identified | 0ms |

## Deviations

None.

## Known Issues

README performance table is a static snapshot and should be labeled as local benchmark evidence rather than a universal SLA.

## Files Created/Modified

- `README.md`
- `docker-compose.yaml`
- `docker-compose.override.yaml`
- `benchmark.py`
