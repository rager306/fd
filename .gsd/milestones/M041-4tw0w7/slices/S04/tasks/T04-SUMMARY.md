---
id: T04
parent: S04
milestone: M041-4tw0w7
key_files:
  - benchmark-results/fd-v2-perf-validation-m041-s04.md
  - benchmark-results/m041-s04-t04-bottleneck-measurement.txt
  - benchmark-results/m041-s04-t05-cache-namespace.txt
  - benchmark-results/m041-s04-t05-runtime-current-fd-isolated-cache.txt
  - benchmark-results/m041-s04-t05-runtime-current-fd.txt
key_decisions:
  - Real inference validation must use an isolated Redis cache namespace; cache-hot latency is not accepted as proof for T-P targets.
  - Current TEI CPU backend is the blocker for real cache-miss latency targets; fd wrapper overhead is not the bottleneck.
duration: 
verification_result: passed
completed_at: 2026-06-14T07:33:07.600Z
blocker_discovered: true
---

# T04: Confirmed real cache-miss inference performance blocker: current fd + TEI CPU misses T-P latency targets when Redis namespace is isolated.

**Confirmed real cache-miss inference performance blocker: current fd + TEI CPU misses T-P latency targets when Redis namespace is isolated.**

## What Happened

Re-ran S04 performance validation against a rebuilt current `fd_api` connected to real `fd_tei` and `fd_redis`. First pass on the default `v2` namespace passed, but that was cache-warm from prior runs. Following the project cache-isolation gotcha, recreated `fd_api` with a fresh `EMBEDDING_CACHE_VERSION=m041-s04-t05-real-*` namespace and reran `tools/verify_fd_v2_perf.sh`. The isolated-cache run failed real inference targets: batch=1 p95 425ms vs 50ms, batch=10 p95 about 1993ms vs 200ms, batch=32 p95 about 5671ms vs 1000ms, and 4x8 concurrent took about 6.987s vs 2s. Cache effectiveness itself passed (`X-Cache: HIT`, about 2ms). A direct TEI comparison showed fd latency is approximately the TEI backend latency (fd batch 10 p50 about 2947ms; TEI direct batch 10 p50 about 3276ms), so the bottleneck is real TEI CPU inference rather than fd cache/header/middleware code. User selected the honest blocker path rather than treating targets as cache-hot or starting ONNX provisioning in this slice.

## Verification

Evidence generated on current fd + real TEI/Redis: `benchmark-results/m041-s04-t05-cache-namespace.txt` records the isolated namespace; `benchmark-results/m041-s04-t05-runtime-current-fd-isolated-cache.txt` records verifier exit 1; `benchmark-results/fd-v2-perf-validation-m041-s04.md` records FAIL with p50/p95/p99; `benchmark-results/m041-s04-t04-bottleneck-measurement.txt` records fd vs TEI direct latency comparison. This verifies the blocker, not a passing final performance result.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd && EMBEDDING_CACHE_VERSION=m041-s04-t05-real-* docker compose up -d --force-recreate api` | 0 | ✅ pass: current fd recreated with isolated cache namespace | 300000ms |
| 2 | `cd /root/fd && FD_BASE_URL=http://localhost:8000 tools/verify_fd_v2_perf.sh` | 1 | ✅ expected fail confirming blocker: isolated-cache real inference misses T-P latency targets while cache HIT passes | 600000ms |
| 3 | `python fd-vs-tei direct bottleneck measurement script` | 0 | ✅ pass: fd latency is comparable to direct TEI latency, identifying backend/runtime bottleneck | 240000ms |

## Deviations

The prior T04 no-op was invalidated because it relied on cache-warm latency evidence. The correct real-inference measurement requires an isolated cache namespace and shows targets fail on TEI CPU.

## Known Issues

T05 and S04 remain incomplete. Final perf validation cannot pass under the current TEI CPU backend unless requirements are rescoped to cache-hot latency or backend/runtime changes are made (for example ONNX/GPU/faster TEI runtime).

## Files Created/Modified

- `benchmark-results/fd-v2-perf-validation-m041-s04.md`
- `benchmark-results/m041-s04-t04-bottleneck-measurement.txt`
- `benchmark-results/m041-s04-t05-cache-namespace.txt`
- `benchmark-results/m041-s04-t05-runtime-current-fd-isolated-cache.txt`
- `benchmark-results/m041-s04-t05-runtime-current-fd.txt`
