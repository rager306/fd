---
id: T03
parent: S04
milestone: M009-zjrq6j
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m009-s04.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T18:12:34.449Z
blocker_discovered: false
---

# T03: Verified S04: benchmark now isolates L1, Redis L2, cached batch, and repeated chunk reuse with Redis deltas; evidence does not justify MGET/pipeline yet.

**Verified S04: benchmark now isolates L1, Redis L2, cached batch, and repeated chunk reuse with Redis deltas; evidence does not justify MGET/pipeline yet.**

## What Happened

Verified S04 benchmark evidence. Go tests passed with 60 tests across 4 packages, pinned lint passed with 0 issues, and Docker compose config passed. The full benchmark generated `benchmark-results/fd-benchmark-m009-s04.txt`. The artifact includes Redis INFO deltas for L1-hot repeated single requests, Redis L2 after API restart, batch L1 hot, batch Redis L2 after API restart, and repeated chunk reuse. Parser checks confirmed the new sections, Redis config snapshot, and batch L2 delta. Evidence: L1 repeated requests produced zero Redis hits/misses; Redis L2 single request after API restart produced one Redis hit; 16-item batch L2 after API restart produced 16 Redis hits, zero misses, 18 Redis commands, and p95 about 5.61ms; repeated chunk reuse had first cold round about 1622ms and warm reuse p95 about 4.22ms.

## Verification

Fresh verification passed for tests, lint, compose config, full benchmark, artifact parser, and GitNexus detect_changes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./... -short` | 0 | ✅ pass: 60 tests in 4 packages | 17100ms |
| 2 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 17100ms |
| 3 | `docker compose config >/tmp/fd-compose-config-m009-s04.txt` | 0 | ✅ pass | 17000ms |
| 4 | `docker compose up -d --force-recreate redis api; health; Redis CONFIG check` | 0 | ✅ pass: Redis restored to default maxmemory 2147483648, allkeys-lru, save 300 1, appendonly no | 35400ms |
| 5 | `uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-m009-s04.txt` | 0 | ✅ pass: benchmark completed and artifact written | 29800ms |
| 6 | `artifact parser for S04 sections and Redis deltas` | 0 | ✅ pass: m009 s04 artifact parser passed | 7000ms |
| 7 | `gitnexus_detect_changes(scope: all, repo: fd)` | 0 | ✅ pass: medium risk limited to benchmark flows | 0ms |

## Deviations

The first S04 full artifact was discarded and rerun after Redis was force-recreated with default Compose settings. This corrected an effective-config mismatch left over from S03's explicit Redis test environment.

## Known Issues

S04 evidence does not show a strong Redis round-trip bottleneck for batch cache hits. Batch L2 after API restart had 16 Redis hits and p95 about 5.61ms for 16 items, while cold batch was seconds-scale and model-bound. This suggests S05 MGET/pipeline should be skipped or deferred unless the user wants speculative micro-optimization.

## Files Created/Modified

- `benchmark.py`
- `benchmark-results/fd-benchmark-m009-s04.txt`
