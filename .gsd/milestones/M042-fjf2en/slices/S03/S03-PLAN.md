# S03: ONNX conditional fallback and speed measurement

**Goal:** Opt-in ONNX mode (FD_BACKEND=onnx). fd binary rebuild с -tags onnx. Cold path ≤500ms per M019. Default off (TEI per R001). Legal quality gate deferred.
**Demo:** After this, FD_BACKEND=onnx (requires rebuilding fd binary с -tags onnx) переключает fd на Go ONNX runtime. Per M019: cold path batch=32 ≤500ms, warm path ≤10ms. Default off (TEI остаётся production). Legal quality gate deferred per M015/M016.

## Must-Haves

- ONNX binary builds clean: go build -tags onnx exit 0. FD_BACKEND=onnx: cold path batch=32 ≤500ms, warm path batch=1 ≤10ms (per M019 baseline). FD_BACKEND=tei (default): cold path не regressed vs M041 baseline. All M041 acceptance tests (45 test cases) pass в FD_BACKEND=onnx. All M042 S02 async tests pass в FD_BACKEND=onnx. documents/onnx-mode-m042.md: legal quality gate deferred (reference M015/M016), production rollout requires separate milestone. benchmark-results/fd-v2-onnx-perf-m042.md с cold/warm numbers vs TEI baseline.

## Proof Level

- This slice proves: runtime + integration + benchmark

## Integration Closure

Reuses existing api/embed/onnx.go (from M008/M019). Adds env-gated selection in main.go. Async pipeline from S02 works transparently over both backends (handler doesn't care about Embedder impl).

## Verification

- No new metrics — backend-specific metrics come from ONNX runtime itself. Add FD_BACKEND info to /info endpoint.

## Tasks

- [ ] **T01: Audit existing ONNX implementation + build matrix** `est:1h`
  Прочитать api/embed/onnx.go (уже есть с M008/M019), api/embed/onnx_manifest.go, api/embed/onnx_tokenizer_*.go. Проверить что impl поддерживает batched input (Embed(ctx, texts []string) ([][]float32, error) — это contract). Build matrix: (1) go build . (default, no onnx) — работает, (2) go build -tags onnx . — компилируется, runtime требует libonnxruntime.so. Создать Makefile target `make build-onnx`.
  - Files: `Makefile (new)`, `api/embed/onnx.go (read-only audit)`, `api/embed/onnx_manifest.go (read-only audit)`
  - Verify: go build -tags onnx -o /tmp/fd-api-onnx ./api exit 0. Audit doc: ONNXEmbedder implements Embedder interface (batch input supported), handles tokenizer, has artifact manifest validation.

- [ ] **T02: FD_BACKEND env selection в main.go** `est:2h`
  api/main.go: parse FD_BACKEND env ("tei" | "onnx", default "tei"). Если "onnx" — instantize ONNXEmbedder вместо TEIClient, передать в handlers.NewEmbeddingsHandler. Должно работать БЕЗ перекомпиляции если бинарь уже собран с -tags onnx. Если FD_BACKEND=onnx но бинарь без onnx tag — fail fast с clear error message. Также env FD_ONNX_ARTIFACT_MANIFEST, FD_ONNX_RUNTIME_LIBRARY, FD_ONNX_TOKENIZER_PATH уже читаются (M008). Создать /info endpoint data update: добавить backend: "tei"|"onnx" поле.
  - Files: `api/main.go (FD_BACKEND env)`, `api/handlers/observability.go (M041 S03 /info endpoint)`
  - Verify: FD_BACKEND=onnx → ONNXEmbedder used. FD_BACKEND=tei (default) → TEIClient used. /info endpoint показывает backend field. Unit test: выбор backend based on env.

- [ ] **T03: ONNX cold/warm perf benchmark + async combo** `est:2h`
  tools/verify_fd_onnx_perf.sh: build fd-onnx binary (go build -tags onnx), start with FD_BACKEND=onnx + ONNX_ARTIFACT_MANIFEST=docs/onnx-artifacts/user-bge-m3-dense-fp32.json + ONNX_RUNTIME_LIBRARY=libonnxruntime.so + ONNX_TOKENIZER_PATH=..., measure cold/warm path × batch sizes 1/8/32/64/128, with FD_ASYNC_CHUNKS=true AND false. Output benchmark-results/fd-v2-onnx-perf-m042.md с comparison table (TEI vs ONNX × async on/off).
  - Files: `tools/verify_fd_onnx_perf.sh`, `benchmark-results/fd-v2-onnx-perf-m042.md`
  - Verify: ONNX mode exit 0. Artifact содержит: cold start latency (first request), warm p95 by batch, async on/off comparison. ONNX cold ≤500ms, warm ≤10ms. TEI baseline (from S02) shown for comparison.

- [ ] **T04: Regression suite: M041 + M042 S02 acceptance в ONNX mode** `est:1h`
  tools/verify_fd_v2_contract.py запустить с FD_BACKEND=onnx + FD_ASYNC_CHUNKS=true. Все 45 M041 test cases + все M042 S02 async tests pass в ONNX mode. Особенно: encoding_format=base64 (response uses ONNX []float32), dimensions=512 (Matryoshka truncation), validation envelope (413, 400 etc), error propagation (ONNX errors → OpenAI envelope). Любой failure — bug в ONNX impl или в adapter.
  - Files: `tools/verify_fd_v2_contract.py (run with FD_BACKEND=onnx)`, `tests/integration/fd_v2_onnx_test.go`
  - Verify: go test ./tests/integration/... -run TestFdV2ONNXMode: все M041 + M042 S02 acceptance pass в ONNX mode. 0 regressions.

- [ ] **T05: ONNX mode docs + legal quality gate reference** `est:1h`
  documents/onnx-mode-m042.md: (1) что такое ONNX mode (env flags, build command), (2) perf numbers (из S03 T03 benchmark), (3) legal quality gate deferred — explicit reference to M015/M016 findings (128-token truncation causes divergence), (4) production rollout checklist (legal quality gate close-out required as separate milestone, contact for opting in, monitoring recommendations). Обновить docs/fd-v2.md Section 5.4 с consolidated "after M042" perf table (TEI sync/async, ONNX sync/async) и known limitations.
  - Files: `documents/onnx-mode-m042.md`, `docs/fd-v2.md (Section 5.4 update)`
  - Verify: Документ существует, cross-references M015/M016, M019, M041. Production rollout checklist explicit. docs/fd-v2.md Section 5.4 updated.

## Files Likely Touched

- Makefile (new)
- api/embed/onnx.go (read-only audit)
- api/embed/onnx_manifest.go (read-only audit)
- api/main.go (FD_BACKEND env)
- api/handlers/observability.go (M041 S03 /info endpoint)
- tools/verify_fd_onnx_perf.sh
- benchmark-results/fd-v2-onnx-perf-m042.md
- tools/verify_fd_v2_contract.py (run with FD_BACKEND=onnx)
- tests/integration/fd_v2_onnx_test.go
- documents/onnx-mode-m042.md
- docs/fd-v2.md (Section 5.4 update)
