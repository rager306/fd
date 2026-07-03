---
id: T02
parent: S01
milestone: M019-opzh2g
key_files: []
key_decisions:
  - Benchmark service PID is tracked in `.gsd/runtime/m019-onnx1024.pid`.
  - Service is running on port 18000 with ONNX 1024 config via the M019 runtime helper.
duration: 
verification_result: passed
completed_at: 2026-05-20T08:12:20.887Z
blocker_discovered: false
---

# T02: Verified the ONNX 1024 benchmark service is running and healthy.

**Verified the ONNX 1024 benchmark service is running and healthy.**

## What Happened

Verified the ONNX 1024 benchmark service is healthy on port 18000 and that the PID file points to a running `go run -tags hf_tokenizers .` process.

## Verification

`curl -fsS http://localhost:18000/health` returned ok and PID check found the running process.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `curl -fsS http://localhost:18000/health && ps -p $(cat .gsd/runtime/m019-onnx1024.pid)` | 0 | ✅ pass — health ok and process running | 0ms |

## Deviations

Service was started during T01 while smoke-testing the restart helper; T02 verified the already-running service rather than starting a second copy.

## Known Issues

None.

## Files Created/Modified

None.
