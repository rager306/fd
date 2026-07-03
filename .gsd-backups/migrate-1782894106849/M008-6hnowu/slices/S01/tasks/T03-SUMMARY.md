---
id: T03
parent: S01
milestone: M008-6hnowu
key_files:
  - .gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:04:38.957Z
blocker_discovered: false
---

# T03: Verified a Go ONNX path exists via `yalue/onnxruntime_go`, but native artifacts and dense-output equivalence are the real gates.

**Verified a Go ONNX path exists via `yalue/onnxruntime_go`, but native artifacts and dense-output equivalence are the real gates.**

## What Happened

Researched model-preserving ONNX Runtime Go path. `github.com/yalue/onnxruntime_go` is a mature Go wrapper around Microsoft ONNX Runtime with tagged versions, MIT license, current release metadata, and functions such as `SetSharedLibraryPath`, `InitializeEnvironment`, `GetInputOutputInfo`, `RegisterExecutionProviderLibrary`, and tensor/session APIs. It is directly relevant for a Go-native ONNX spike, but it does not eliminate native dependency management: the service still needs a compatible ONNX Runtime shared library, model/tokenizer artifacts, and provider libraries if non-default EPs are used. The safest path is not to adopt `go-bge-m3-embed` as a black box; instead, build a small fd-owned ONNX dense embedder adapter around a pinned BGE-M3 ONNX model/tokenizer and compare dense output to current TEI/Candle. Use S06 benchmark order: FP32 dense-only default CPU EP first, then threading/NUMA, then provider variants, then INT8 only with quality gates.

## Verification

Read yalue/onnxruntime_go GitHub/pkg.go.dev docs and connected findings to S06 ONNX provider research.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Fetched: https://github.com/yalue/onnxruntime_go` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Fetched: https://pkg.go.dev/github.com/yalue/onnxruntime_go` | -1 | unknown (coerced from string) | 0ms |
| 3 | `S06 research artifact: .gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

Go ONNX Runtime path still depends on native `libonnxruntime.so` packaging and tokenizer/model artifact provenance. BGE-M3 ONNX artifacts from community model cards must be pinned by hash and validated against current TEI output before adoption. Provider-specific EPs may be hard to access from Go bindings unless shared provider registration and native libraries are explicitly handled.

## Files Created/Modified

- `.gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md`
