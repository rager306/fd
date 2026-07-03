# S03: Opt in ONNX dense backend prototype

**Goal:** Prototype a real opt-in ONNX dense backend using the S01 manifest and S02 validation seam, while preserving TEI as default and recording a blocker if Go ONNX/tokenizer dependencies are not viable.
**Demo:** After this, an explicit ONNX backend either returns dense embeddings locally or is blocked with concrete binding/runtime evidence.

## Must-Haves

- ONNX backend is never default.
- ONNX artifact checksum is validated before use.
- Go dependency/runtime requirements are documented and tested.
- Dense output is 1024-dimensional and normalized if backend works.
- Comparator checks against TEI baseline pass or blocker is documented.
- No INT8/provider/sparse/ColBERT work is introduced.

## Proof Level

- This slice proves: Executable opt-in backend path in Go or structured blocker evidence, plus comparator output against TEI baseline.

## Integration Closure

Feeds S04 benchmark/recommendation with either a working opt-in ONNX backend or concrete dependency/runtime blocker evidence.

## Verification

- ONNX backend must report artifact_id, provider/library path, dimensions, and load/validation errors without raw input text.

## Tasks

- [x] **T01: Check ONNX Go dependency feasibility** `est:small`
  Run impact analysis for the runtime wiring symbols and verify Go dependency feasibility for `github.com/yalue/onnxruntime_go` and `github.com/sugarme/tokenizer`. Check whether a usable `libonnxruntime.so` path exists and document required env vars such as `ONNXRUNTIME_SHARED_LIBRARY_PATH` or future `ONNX_RUNTIME_LIBRARY`.
  - Files: `api/go.mod`, `api/main.go`, `api/embed/`
  - Verify: GitNexus impact recorded; dependency docs and local shared-library path checked; no runtime code changed.

- [x] **T02: Implement ONNX dense embedder or blocker** `est:large`
  Add an ONNX dense embedder package implementation behind the existing `handlers.Embedder` interface. Use `sugarme/tokenizer` to load local tokenizer JSON and `yalue/onnxruntime_go` to load/run `dense_vecs`. The implementation should validate manifest first, initialize ONNX Runtime with an explicit shared-library path, tokenize each input, run CPU EP, return normalized 1024-dimensional `[]float32`, and close native resources. If cgo/runtime blocks implementation, record structured blocker instead of stubbing fake embeddings.
  - Files: `api/embed/onnx.go`, `api/embed/onnx_test.go`, `api/go.mod`, `api/go.sum`
  - Verify: Focused ONNX embedder tests pass or blocker artifact explains exact dependency/runtime failure.

- [x] **T03: Wire opt in ONNX backend** `est:medium`
  Wire the opt-in ONNX backend into `api/main.go`: default TEI path unchanged; when `EMBEDDING_BACKEND=onnx`, validate manifest, require shared library/tokenizer paths, construct ONNX embedder, and route handlers to it. Add tests for default TEI behavior and ONNX config error paths. Do not expose ONNX unless explicitly configured.
  - Files: `api/main.go`, `api/main_test.go`
  - Verify: `cd api && go test ./... -short` passes; default env still selects TEI.

- [x] **T04: Run opt in ONNX API comparison** `est:medium`
  Run the opt-in ONNX backend against the local artifact if dependencies are viable, compare `/v1/embeddings` output against M010 TEI baseline using existing comparator flow or a small API smoke script, and save tracked evidence. If backend cannot run, save a blocker artifact with exact failure. Verify no raw texts are logged.
  - Files: `benchmark-results/fd-go-onnx-m011-s03.txt`
  - Verify: ONNX API comparison artifact exists with PASS/FAIL/BLOCKED and no raw probe text leakage.

## Files Likely Touched

- api/go.mod
- api/main.go
- api/embed/
- api/embed/onnx.go
- api/embed/onnx_test.go
- api/go.sum
- api/main_test.go
- benchmark-results/fd-go-onnx-m011-s03.txt
