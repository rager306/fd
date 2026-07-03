---
id: M009-zjrq6j
title: "Measured cache and benchmark foundation"
status: complete
completed_at: 2026-05-19T18:18:09.584Z
key_decisions:
  - Use allowlisted sanitized benchmark snapshots, not full environment dumps.
  - Preserve default Redis `v2` keys and 24h TTL unless explicit env settings are provided.
  - Use hashed model/revision/tokenizer/chunk namespace fields in Redis keys and clear values in sanitized benchmark snapshots.
  - Use RDB-first Redis persistence by default; keep AOF disabled for rebuildable cache data.
  - Skip MGET/pipeline in M009 because S04 measured batch L2 p95 was already low-ms and cold path remains model-bound.
key_files:
  - benchmark.py
  - api/cache/redis.go
  - api/cache/redis_test.go
  - api/main.go
  - api/main_test.go
  - docker-compose.yaml
  - docker-compose.override.yaml
  - README.md
  - benchmark-results/fd-benchmark-m009-s01.txt
  - benchmark-results/fd-benchmark-m009-s02.txt
  - benchmark-results/fd-benchmark-m009-s03.txt
  - benchmark-results/fd-benchmark-m009-s04.txt
lessons_learned:
  - Benchmark effective server config must be checked, not inferred from host env; stale containers can keep old Redis command settings.
  - Strict config validation can expose silent defaulting bugs that previously passed unnoticed.
  - Redis INFO deltas are useful to distinguish L1-hot paths from Redis L2 hits.
  - Cold batch latency is model-bound in current local runs; cached batch path is already low milliseconds at tested sizes.
---

# M009-zjrq6j: Measured cache and benchmark foundation

**M009 made benchmark results comparable and Redis cache behavior tunable, persistent, observable, and measured; MGET/pipeline was skipped based on evidence.**

## What Happened

M009 implemented the measured cache and benchmark foundation. It added sanitized benchmark config snapshots, env-configurable Redis cache namespace and retention, explicit RDB-first Redis persistence configuration, Redis CONFIG/INFO benchmark visibility, and batch-hit/repeated-chunk benchmark sections. It verified live TTL and no-expire behavior, Redis key namespace hashing, Redis restart persistence, benchmark snapshot parsing, and cached batch performance. Based on S04 evidence, it skipped MGET/pipeline implementation rather than adding speculative complexity.

## Success Criteria Results

- ✅ Benchmark config snapshots generated and sanitized: S01.
- ✅ Redis cache namespace and retention env-configurable with safe defaults/tests: S02.
- ✅ Redis persistence and memory policy documented/reflected in benchmark artifacts: S03.
- ✅ Batch cache-hit benchmark evidence exists before MGET/pipeline optimization: S04.
- ✅ No ONNX, INT8, provider, Rust, C, or model replacement work introduced.
- ✅ MGET/pipeline avoided because evidence did not justify it.

## Definition of Done Results

- ✅ All executed slices complete; S05 skipped with measured rationale.
- ✅ Go tests passed: 60 tests in 4 packages.
- ✅ Pinned GolangCI-Lint passed: 0 issues.
- ✅ Docker compose config passed.
- ✅ Benchmark artifacts S01-S04 exist and parse.
- ✅ Redis TTL/no-expire, namespace, persistence, and batch-hit evidence were verified live.
- ✅ GitNexus detect_changes reported no uncommitted changes before milestone completion.

## Requirement Outcomes

- R002 advanced: Redis L2 supports long-lived TTL/no-expire behavior and RDB-first persistence for research/chunk reuse.
- R003 advanced: cache/runtime/Redis settings are configurable via env/Compose with safe validation.
- R004 advanced: benchmark artifacts record sanitized effective config, Redis INFO/CONFIG, and section-level Redis deltas.
- R001 preserved: no model replacement, ONNX runtime change, INT8, sparse/ColBERT, Rust, or C change was introduced.

## Deviations

S05 was skipped rather than implemented because S04 benchmark evidence did not justify MGET/pipeline complexity. During S02 verification, strict Redis config validation surfaced and fixed a pre-existing `getEnvInt` defaulting bug. During S04 verification, an inherited Redis test config mismatch was detected and corrected by force-recreating Redis/API with default settings before recording the final artifact.

## Follow-ups

Next milestone should move to either (a) ONNX FP32 dense-only spike now that benchmark comparability exists, or (b) larger-scale batch/cache experiments if production chunk batch sizes are materially larger than S04's local synthetic workload. MGET/pipeline should remain deferred unless larger evidence shows Redis round-trip pressure.
