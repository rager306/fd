---
id: T02
parent: S01
milestone: M017-j10hmp
key_files: []
key_decisions:
  - Tagged ONNX 512 service is running as background process `6809eefc` in group `m017-onnx512`.
  - Runtime config: port 18000, `ONNX_MAX_SEQUENCE_LENGTH=512`, `EMBEDDING_CACHE_VERSION=m017-onnx-512-legal-quality`, native HF tokenizer build tag enabled.
duration: 
verification_result: passed
completed_at: 2026-05-20T07:25:56.468Z
blocker_discovered: false
---

# T02: Started and verified the tagged ONNX 512 service.

**Started and verified the tagged ONNX 512 service.**

## What Happened

Started the tagged Go ONNX service with HF native tokenizer, local ONNX artifact, and max sequence length 512. The service reached ready state on port 18000 and `/health` returned ok.

## Verification

`bg_shell wait_for_ready` reported ready and `curl -fsS http://localhost:18000/health` returned ok.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `bg_shell wait_for_ready id=6809eefc` | 0 | ✅ pass — process ready on port 18000 | 0ms |
| 2 | `curl -fsS http://localhost:18000/health` | 0 | ✅ pass — health ok | 0ms |

## Deviations

None.

## Known Issues

Gin reports debug mode warning; acceptable for local benchmark gate and not production deployment.

## Files Created/Modified

None.
