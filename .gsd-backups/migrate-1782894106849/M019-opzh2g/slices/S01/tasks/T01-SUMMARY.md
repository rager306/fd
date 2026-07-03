---
id: T01
parent: S01
milestone: M019-opzh2g
key_files:
  - .gsd/runtime/restart-fd-api-onnx-m019-1024.sh
key_decisions:
  - Benchmark target is `http://localhost:18000` with runtime label `tagged-onnx-hf-1024`.
  - Use `EMBEDDING_CACHE_VERSION=m019-onnx-1024-benchmark`.
  - Use `.gsd/runtime/restart-fd-api-onnx-m019-1024.sh restart` as `BENCHMARK_API_RESTART_COMMAND` so Redis L2 restart checks target the actual ONNX service.
duration: 
verification_result: passed
completed_at: 2026-05-20T08:12:03.164Z
blocker_discovered: false
---

# T01: Prepared the ONNX 1024 benchmark command plan and restart helper.

**Prepared the ONNX 1024 benchmark command plan and restart helper.**

## What Happened

Prepared the benchmark command plan and runtime restart helper. The helper starts/restarts the tagged Go ONNX 1024 API on port 18000 with native HF tokenizer, local ONNX artifact, and isolated cache namespace. It was smoke-tested with `/health` returning ok.

## Verification

Restart helper was created, made executable, started the ONNX 1024 service, and `/health` returned ok.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `chmod +x .gsd/runtime/restart-fd-api-onnx-m019-1024.sh && .gsd/runtime/restart-fd-api-onnx-m019-1024.sh start && curl -fsS http://localhost:18000/health` | 0 | ✅ pass — ONNX 1024 restart helper starts healthy service | 0ms |

## Deviations

Created a temporary `.gsd/runtime/` restart helper for benchmark restart semantics. It is runtime-only and not intended for commit.

## Known Issues

The helper backgrounds the benchmark server because benchmark.py invokes restart commands as subprocesses. It is confined to `.gsd/runtime/` and will be cleaned up after benchmarking.

## Files Created/Modified

- `.gsd/runtime/restart-fd-api-onnx-m019-1024.sh`
