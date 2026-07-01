# S01 Same-host service contract — Research

**Date:** 2026-05-21

## Summary

S01 should produce a concrete same-host HTTP contract for neighboring local services and close two contract-readiness gaps found in the current implementation. The endpoint surface is already small and stable: `GET /health`, `POST /v1/embeddings`, and `POST /embeddings/batch` are wired in `api/main.go`; request/response shapes are defined in `api/embed/types.go`; validation and handler timeouts live in `api/handlers/embeddings.go` and `api/handlers/batch.go`. The README documents the basics, but there is no dedicated same-host service contract document yet.

The current code already enforces the most important no-silent-runtime-fallback rule: exactly one backend is selected at startup by `EMBEDDING_BACKEND`, invalid ONNX config exits non-zero, and the request path never falls back from ONNX to TEI or from TEI to ONNX. However, `/health` only exposes runtime metadata for ONNX; TEI/default health omits backend/model/cache namespace, so local clients cannot reliably inspect the runtime they are using. Also, `/v1/embeddings` accepts a `model` request field but does not validate it against the configured model; the response model is authoritative, but this should either be documented explicitly or hardened to reject mismatches.

Active S01 requirements are R003, R005, and R009. The highest-risk research finding for these is not endpoint implementation complexity; it is avoiding an overclaim in the contract. Today `/health` is a liveness/config surface after API startup: Redis is pinged before serving, ONNX model/artifact preflight completes before serving, but TEI reachability is not actively probed in `main.go` before `ListenAndServe`. The contract must state exactly what `/health` proves, or the implementation must add a shallow TEI readiness check if planner chooses a stronger readiness definition.

## Recommendation

Build S01 as a small docs-plus-health-contract slice, not a benchmark or runtime recommendation slice. Create a dedicated service contract doc, link it from README, and make the health metadata consistent enough that neighboring services can programmatically identify backend/model/cache namespace before sending embeddings. Keep runtime selection explicit and fail-fast; do not introduce request-level fallback.

Recommended implementation approach:

1. Add `docs/same-host-embedding-service-contract.md` as the canonical consumer contract. It should define endpoint shapes, status codes, dimensions, batch encoding caveats, env/runtime requirements, health metadata, retry/timeout guidance, cache namespace expectations, and no-silent-fallback rules.
2. Update `api/main.go` / `api/handlers/health.go` so `main` exposes safe runtime metadata for TEI as well as ONNX. TEI/default should report at least `backend: "tei"`, `model`, `dimensions: 1024`, `production_default: true`, and `cache_namespace`; ONNX keeps existing artifact/provider/tokenizer/runtime fields.
3. Decide whether to harden `/v1/embeddings` model handling now. If changed, reject non-empty `req.Model` that does not equal configured `modelID` with `400`, and add tests. If not changed, the contract must say request `model` is compatibility metadata only and clients must trust the response model plus `/health.runtime.model`.
4. Do not change `/embeddings/batch` response shape in this slice. It currently returns `embeddings: []string`; `encoding_format=base64` returns binary float32 vectors as base64, while `encoding_format=float` returns JSON-encoded vector strings, not nested arrays. Document that exactly to avoid breaking FalkorDB/local consumers.

## Implementation Landscape

### Key Files

- `api/main.go` — owns env loading, backend selection, fail-fast ONNX preflight, Redis ping, handler wiring, and `runtimeConfig.Health(modelID, cacheNamespace)`. Main gap: `Health` currently returns `nil` for TEI/default, so `/health` in the running service omits runtime metadata for the production/default runtime.
- `api/handlers/health.go` — defines the safe public health metadata shape. Existing `RuntimeHealth` already has most fields required by M040 (`backend`, `model`, `dimensions`, sequence length, verification booleans, provider, `cache_namespace`) and intentionally omits filesystem paths/secrets. It may need comments or small field additions only if planner wants to distinguish liveness vs dependency readiness.
- `api/handlers/health_test.go` — tests default health shape and ONNX safe metadata. Update/add tests for TEI runtime metadata through `NewHealthHandler`; keep `HealthHandler` default-shape compatibility test if that helper remains intentionally minimal.
- `api/handlers/embeddings.go` — OpenAI-compatible endpoint. Accepts `input` as string or array via custom unmarshal, supports only `dimensions` 1024 or 512, uses a 30s request context, returns float arrays. Potential hardening seam: validate `req.Model` against `h.modelID` if the contract should fail fast on model mismatch.
- `api/handlers/batch.go` — internal batch endpoint. Supports `inputs`, `dimensions` 1024/512, `encoding_format` base64/default or float, 120s timeout. It loops through inputs with cache-aside loading and returns `[]string` embeddings.
- `api/embed/types.go` — request/response structs for both endpoints. Important contract detail: `/v1/embeddings` response includes per-item `dimensions`; batch response has top-level `dimensions` and `count` only.
- `api/cache/redis.go` — env-driven Redis namespace and retention contract. `EMBEDDING_CACHE_VERSION`, model/revision/tokenizer/chunking namespace fields, `REDIS_CACHE_TTL`, and `REDIS_CACHE_NO_EXPIRE` are correctness/operability fields that should appear in the service contract. Invalid TTL/no-expire combinations fail startup.
- `api/cache/tiered.go` — cache behavior: L1 local, Redis L2, singleflight, no raw text logs. Contract should state cache is transparent and does not change vector semantics, but cache namespace contamination across runtimes/models is a correctness risk.
- `README.md` — currently has basic API/config docs. Link the new contract here rather than expanding README into the full operating contract.
- `docs/onnx-artifacts/OPERATIONS.md` — existing ONNX-specific preflight/health/rollback contract. S01 should reference and align with it, not duplicate all ONNX artifact details.
- `benchmark.py` — downstream S02 consumes `BENCHMARK_API_RESTART_COMMAND` and records sanitized env/runtime metadata. S01 only needs to document the restart/cache contract boundary that S02 will prove.
- `benchmark-results/fd-onnx-go-target-runtime-acceptance-m038-s02.txt` — prior actual Go ONNX endpoint evidence: `/health` reported ONNX metadata, legal parity passed, Redis namespaces were isolated, Redis L2 restart was skipped.
- `benchmark-results/fd-onnx-docker-target-runtime-acceptance-m039-s02.txt` — prior packaged Docker evidence: packaged smoke/legal/performance passed; `ONNX_RUNTIME_SHA256` was required for `runtime_library_verified=true`; Redis L2 restart subcheck was skipped and is left for S02.

### Build Order

1. **Contract doc first** — write the consumer-facing service contract from existing code evidence. This clarifies what the implementation must prove and prevents overclaiming `/health` semantics.
2. **Health metadata parity second** — add TEI/default runtime metadata to the runtime health path if the contract requires clients to identify runtime readiness from `/health`. This is the main S01 code seam and unlocks a truthful same-host readiness contract.
3. **Optional model-mismatch hardening third** — decide whether to reject mismatched `model` in `/v1/embeddings`. If implemented, it is a small handler/test change; if not, document the response model as authoritative.
4. **README link and tests last** — link the canonical doc from README and update handler/health tests. Avoid touching benchmark behavior in S01; S02 owns restart/cache proof.

### Verification Approach

- Static/doc checks:
  - New doc includes required sections for endpoints, env/runtime requirements, health metadata, timeout/retry guidance, cache namespace guidance, no-silent-fallback rule, and non-goals.
  - README links to the new contract.
- Go tests after code changes:
  - `cd api && go test ./... -short`
  - If lint is in scope: `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...`
- Local smoke if a service is running:
  - `curl -s http://localhost:8000/health` should show `status: ok` and runtime metadata for the configured backend.
  - `curl -s -X POST http://localhost:8000/v1/embeddings -H 'Content-Type: application/json' -d '{"model":"deepvk/USER-bge-m3","input":"юридическая справка","dimensions":512}'` should return `data[0].dimensions == 512` and response `model` matching health.
  - Batch smoke: `POST /embeddings/batch` with `encoding_format=base64` returns `count`, `dimensions`, and string embeddings.
- Safety/leak check:
  - `rg -n "signed|token=|X-Amz|BEGIN|PRIVATE|юридическая справка" docs README.md benchmark-results .gsd/milestones/M040-pbp9z1/slices/S01` after writing artifacts; only intentionally tiny sample strings in docs should remain, and no secrets/signed URLs should appear.
- Graph/process check for executor: because S01 likely edits symbols (`embeddingRuntimeConfig.Health`, handlers, tests), run GitNexus impact before editing modified symbols and `gitnexus_detect_changes()` before closure per project policy.

## Constraints

- TEI remains production/default unless a later decision changes it. The S01 contract must not imply ONNX promotion.
- ONNX remains opt-in with `EMBEDDING_BACKEND=onnx` plus `onnx hf_tokenizers` build tags and artifact/runtime/tokenizer env. Missing/invalid ONNX config is a startup failure, not fallback.
- `ONNX_RUNTIME_SHA256` is optional in code but required by the M040 evidence envelope when claiming `runtime_library_verified=true` in health for packaged ONNX.
- Redis namespace isolation is correctness-critical for TEI-vs-ONNX and candidate comparisons. Use distinct `EMBEDDING_CACHE_VERSION` or model/tokenizer namespace fields; do not rely on shared Redis keys when comparing runtime outputs.
- `/health` must not expose manifest paths, runtime library paths, tokenizer paths, raw input text, signed URLs, or secrets. Existing tests already assert path fields are omitted for ONNX health.
- Current `/health` is not a deep TEI inference check. If the contract says health means full inference readiness, implementation must add a TEI probe; otherwise document health as API/config/cache readiness plus backend metadata and require a smoke embedding for full end-to-end readiness.
- Existing batch `encoding_format=float` returns JSON-encoded vector strings inside the `embeddings` string array. Changing this is a breaking API change and should not be bundled into S01 unless explicitly planned.

## Common Pitfalls

- **Overclaiming readiness** — `status: ok` currently does not prove TEI inference is healthy. Either add the probe or write the contract with this limitation.
- **Mixed-vector cache contamination** — TEI and ONNX can appear equivalent if they reuse Redis entries. Always isolate `EMBEDDING_CACHE_VERSION` or namespace fields for comparisons and record the namespace in artifacts.
- **Silent model mismatch** — `/v1/embeddings` currently ignores request `model` for routing. If the contract promises model mismatch failures, implement validation; if not, document that the service configuration and response model are authoritative.
- **Fallback language ambiguity** — TEI's internal Candle/CPU fallback warning from earlier milestones is not the same as fd request-level fallback. The S01 no-silent-fallback rule should focus on fd not changing backend/model/tokenizer per request within one service run.
- **Health path leaks** — ONNX health should continue exposing safe metadata only, not artifact filesystem paths or signed artifact locations.

## Open Risks

- A stronger readiness definition may require code beyond metadata: TEI startup/dependency probing or a `/ready` endpoint. That would be a small but real API semantics change and should be planned deliberately.
- Model request validation could break ad-hoc clients/tests that send placeholder model names. If implemented, update tests and document that clients may omit `model` or must send `deepvk/USER-bge-m3`.

## Skills Discovered

| Technology | Skill | Status |
|------------|-------|--------|
| HTTP/API contract design | `api-design` (installed in prompt) | Relevant; use contract-first endpoint/status/error-shape thinking. |
| Health/observability contract | `observability` (installed in prompt) | Relevant; health metadata must be safe, explicit, and useful to the next agent/operator. |
| Go Gin HTTP API | `bobmatnyc/claude-mpm-skills@golang-http-frameworks` and `henriqueatila/golang-gin-best-practices@golang-gin-api` found | Not installed; S01 is simple enough to follow local Gin patterns. |
| Redis cache | Several generic/Azure/Laravel Redis skills found | Not installed; not directly useful for this Go same-host contract. |
| ONNX Runtime Go | none found | No install. |
| Hugging Face TEI | generic Hugging Face skills found | Not installed; S01 relies on local TEI/ONNX evidence rather than new HF integration work. |
