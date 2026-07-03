---
id: T01
parent: S05
milestone: M003-xx4yc3
key_files:
  - docker-compose.override.yaml
  - benchmark-results/fd-benchmark-baseline-py313.txt
  - benchmark-results/fd-runtime-stats-logs.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T08:25:02.341Z
blocker_discovered: false
---

# T01: Classified runtime findings: one Redis exposure fix, one local stale-container cleanup, no remaining correctness blockers.

**Classified runtime findings: one Redis exposure fix, one local stale-container cleanup, no remaining correctness blockers.**

## What Happened

Classified runtime findings. Fixed defect: Redis was exposed on all host interfaces via docker-compose.override.yaml and logs showed real external attack attempts; this was fixed by binding to 127.0.0.1. Fixed operational blocker: stale exited fd_api container caused startup name conflict; it was removed locally. No additional correctness defects appeared in health, endpoint, cache, or benchmark gates. Non-blocking notes: Redis host overcommit warning should be handled on deployment hosts; TEI lacks ONNX artifacts and falls back to Candle CPU; API logs one INFO line per successful embedding, noisy during throughput tests. Future optimizations should focus on cache observability/log sampling, dependency-aware readiness, TEI batching constraints, and benchmark-driven transport/batch tuning.

## Verification

Classification is based on S01-S04 evidence: Compose startup/logs, live API smoke tests, cache validation, benchmark output, and docker stats/logs.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Review S01-S04 summaries and benchmark/log artifacts` | 0 | ✅ pass: findings classified | 0ms |

## Deviations

None.

## Known Issues

Redis host memory overcommit warning remains a host-level tuning issue. TEI ONNX fallback is expected for this model unless ONNX artifacts are produced. API INFO logs are noisy under throughput but not a correctness defect.

## Files Created/Modified

- `docker-compose.override.yaml`
- `benchmark-results/fd-benchmark-baseline-py313.txt`
- `benchmark-results/fd-runtime-stats-logs.txt`
