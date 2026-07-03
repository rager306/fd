---
id: T03
parent: S05
milestone: M008-6hnowu
key_files: []
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:10:55.167Z
blocker_discovered: false
---

# T03: C has the official ORT API and maximum control, but should be limited to narrow FFI only if a measured hotspot requires it.

**C has the official ORT API and maximum control, but should be limited to narrow FFI only if a measured hotspot requires it.**

## What Happened

Researched C option maturity and likely gain. ONNX Runtime's C API is official, cross-platform, and exposes core controls: environment creation, sessions, tensors, predefined execution providers, custom logging/allocators, thread pool size, graph optimization level, and model loading from file or array. `hiredis` and C HTTP libraries could support a very low-overhead service. However, a full C service would carry high operational and security cost: manual memory management, less ergonomic structured error handling, harder concurrency/cancellation, harder JSON/OpenAI-compatible API maintenance, harder test ergonomics, and more native packaging work. C is best reserved as a narrow FFI boundary or direct ORT integration reference if a specific hotspot is proven, not as a full rewrite. If native ORT C API is needed, Go or Rust wrappers should be preferred first, with C limited to a small audited adapter if bindings lack required provider/session controls.

## Verification

Read ONNX Runtime C API documentation and compared operational risk against fd's current Go API/cache surface.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Fetched: https://onnxruntime.ai/docs/get-started/with-c.html` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Search query: ONNX Runtime C API hiredis C HTTP server embedding service performance operational complexity` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

C/C++ can expose maximum ORT control but increases memory-safety and operational risk. The project currently has no evidence that C-specific overhead savings would beat ONNX provider/threading/cache improvements.

## Files Created/Modified

None.
