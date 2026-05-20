---
id: T01
parent: S03
milestone: M011-33b7wf
key_files:
  - api/go.mod
  - api/main.go
  - api/embed/
key_decisions:
  - Use `github.com/yalue/onnxruntime_go@v1.30.1` as the first Go ONNX Runtime binding candidate.
  - Use `github.com/sugarme/tokenizer@v0.3.0` and `pretrained.FromFile` for tokenizer JSON loading.
  - Require explicit ONNX Runtime shared library path because `onnxruntime_go` requires cgo and a compatible shared library before initialization.
duration: 
verification_result: passed
completed_at: 2026-05-19T19:09:26.905Z
blocker_discovered: false
---

# T01: Confirmed S03 ONNX Go dependencies are plausible but require explicit shared-library path and careful tokenizer transitive deps.

**Confirmed S03 ONNX Go dependencies are plausible but require explicit shared-library path and careful tokenizer transitive deps.**

## What Happened

Checked S03 dependency feasibility before runtime edits. GitNexus impact for `Function:api/main.go:main` and `Function:api/embed/tei.go:NewTEIClient` is LOW with no affected processes. Local Go environment supports cgo (`CGO_ENABLED=1`, `GOOS=linux`, `GOARCH=amd64`, `CC=gcc`, `/usr/bin/gcc`). Go module versions are available for `github.com/yalue/onnxruntime_go` through `v1.30.1` and `github.com/sugarme/tokenizer` through `v0.3.0`. Local ONNX Runtime shared-library candidates exist in uv caches, including `libonnxruntime.so.1.26.0` and `libonnxruntime.so.1.24.4`. A temporary module compile probe confirmed `onnxruntime_go` APIs (`SetSharedLibraryPath`, `InitializeEnvironment`, `NewAdvancedSession`, `NewTensor`, `NewEmptyTensor`) and `pretrained.FromFile` compile together after adding tokenizer BPE/unigram transitive module entries.

## Verification

GitNexus impacts were LOW. `go list -m -versions` found required modules. Local cgo/gcc is available. Temp-module compile probe passed after adding tokenizer BPE/unigram transitive modules.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gitnexus_impact({target:'Function:api/main.go:main', direction:'upstream', repo:'fd'})` | 0 | ✅ pass — LOW risk | 0ms |
| 2 | `gitnexus_impact({target:'Function:api/embed/tei.go:NewTEIClient', direction:'upstream', repo:'fd'})` | 0 | ✅ pass — LOW risk | 0ms |
| 3 | `cd api && go list -m -versions github.com/yalue/onnxruntime_go github.com/sugarme/tokenizer && go env CGO_ENABLED GOOS GOARCH CC && command -v gcc` | 0 | ✅ pass — modules available; cgo/gcc available | 0ms |
| 4 | `find /root/.cache -path '*onnxruntime/capi/libonnxruntime.so.*'` | 0 | ✅ pass — local uv cache contains ONNX Runtime shared libraries | 0ms |
| 5 | `temporary Go module compile probe importing yalue/onnxruntime_go and sugarme/tokenizer/pretrained` | 0 | ✅ pass after adding tokenizer model/bpe and model/unigram module entries | 0ms |

## Deviations

Dependency probe found `github.com/sugarme/tokenizer@v0.3.0` needs extra module entries for BPE/unigram transitive packages in a fresh module; adding `github.com/sugarme/tokenizer/model/bpe@v0.3.0` and `github.com/sugarme/tokenizer/model/unigram@v0.3.0` resolved the compile probe.

## Known Issues

Available local ONNX Runtime shared libraries came from uv Python caches and include versions `1.24.4`, `1.26.0`, and `1.20.1`; yalue README notes shared library/header version compatibility matters. A stable project-local or externally managed runtime library path is still needed before production.

## Files Created/Modified

- `api/go.mod`
- `api/main.go`
- `api/embed/`
