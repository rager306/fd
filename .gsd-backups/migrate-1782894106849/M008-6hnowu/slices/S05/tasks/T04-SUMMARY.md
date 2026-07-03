---
id: T04
parent: S05
milestone: M008-6hnowu
key_files:
  - .gsd/milestones/M008-6hnowu/slices/S05/S05-RESEARCH.md
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:11:27.743Z
blocker_discovered: false
---

# T04: Recommended no rewrite now: keep Go, profile first, use Rust sidecar only if native inference evidence justifies it, and avoid full C service.

**Recommended no rewrite now: keep Go, profile first, use Rust sidecar only if native inference evidence justifies it, and avoid full C service.**

## What Happened

Recommended language/runtime rewrite strategy. Keep fd in Go for API/cache orchestration because current code is small, tested, and most likely bottlenecks are TEI/model inference, Redis round trips, cache retention/namespace, or ONNX provider/threading—not Go itself. Add pprof/per-layer timing before any rewrite: handler validation, cache lookup, Redis round trip, model call, marshal/unmarshal, response serialization. If cached-hit p95 is Redis-bound, implement MGET/pipeline first. If cold p95 is model-bound, test ONNX FP32/threading/provider path. If Go ONNX wrapper or Go allocations become a proven bottleneck after moving inference in-process, build a Rust sidecar behind the same OpenAI-compatible API for A/B comparison. C should be used only as a narrow audited FFI/reference when wrappers lack required ORT controls.

## Verification

Saved S05 research artifact and aligned recommendation with measured-bottleneck gating.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Saved: .gsd/milestones/M008-6hnowu/slices/S05/S05-RESEARCH.md` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Evidence sources include ort docs, ONNX Runtime C API docs, and current fd bottleneck map` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

This recommendation depends on future profiling evidence. Current benchmark lacks pprof/per-layer timing, so rewrite decisions remain gated.

## Files Created/Modified

- `.gsd/milestones/M008-6hnowu/slices/S05/S05-RESEARCH.md`
