---
id: T01
parent: S01
milestone: M040-pbp9z1
key_files:
  - docs/same-host-embedding-service-contract.md (created, 16 KB)
  - README.md (updated: added contract link at line 59)
key_decisions:
  - TEI /health omits runtime block by design — absence of runtime block = TEI backend active. ONNX opt-in adds full runtime metadata block with artifact/tokenizer/runtime SHA256 verification fields.
  - Request model field in /v1/embeddings is compatibility-only — service does not validate against configured MODEL_ID; response model is authoritative and documented as such.
  - encoding_format=float in batch returns stringified JSON array literals in a string array — NOT a nested JSON array — documented as a known shape quirk.
  - /health is not a live inference probe for either backend — full readiness requires a smoke embedding request; documented for clients and future agents.
duration: 
verification_result: passed
completed_at: 2026-05-22T04:32:41.759Z
blocker_discovered: false
---

# T01: Created docs/same-host-embedding-service-contract.md covering endpoints, runtime/env, health semantics, timeout/retry, cache, no-fallback rules, and non-goals

**Created docs/same-host-embedding-service-contract.md covering endpoints, runtime/env, health semantics, timeout/retry, cache, no-fallback rules, and non-goals**

## What Happened

T01 created `docs/same-host-embedding-service-contract.md` (16 KB, 11 sections) as the canonical local HTTP consumer contract for fd. The document covers all required areas from the task plan: `/health` (with explicit no-overclaim language — TEI liveness is API startup + Redis ping, not live inference probe), `/v1/embeddings` (OpenAI-compatible, response model authoritative, request model field is compatibility-only), `/embeddings/batch` (FalkorDB internal endpoint, base64/float encoding quirks documented), runtime identification (TEI default vs ONNX opt-in via EMBEDDING_BACKEND), health metadata semantics (runtime block present for ONNX, absent for TEI), timeout/retry guidance (30 s single, 120 s batch, transport/503 retry only), cache namespace isolation (EMBEDDING_CACHE_VERSION + model/revision/tokenizer/chunking components), no-silent-fallback rules (startup-only backend selection, no per-request switch), status/error codes, encoding caveats (base64 binary float32 little-endian, float is stringified JSON array, NOT nested JSON), and explicit non-goals. README.md was updated to link the new contract. Go tests (`go test ./... -short`) and lint (`golangci-lint`) both pass with zero issues. No secrets or raw benchmark/legal text leaked — the only probe string (`"юридическая справка"`) is a standard API example already established in README.md.

## Verification

Document inspection confirmed 11 required sections present. `/health` overclaim prevention verified: explicit "does NOT prove" section states TEI liveness is API startup + Redis ping only, not live inference probe. Go tests and lint both pass with zero issues. README link added and confirmed at line 59. No secrets leaked (only accepted example string present).

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go test ./... -short` | 0 | ✅ pass | 12000ms |
| 2 | `cd /root/fd/api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass | 30000ms |
| 3 | `grep -c '^## [0-9]' docs/same-host-embedding-service-contract.md` | 0 | ✅ pass — 11 sections found | 0ms |
| 4 | `grep -c 'does NOT prove\|not a live\|does not perform' docs/same-host-embedding-service-contract.md` | 0 | ✅ pass — /health no-overclaim language present | 0ms |
| 5 | `grep -c 'EMBEDDING_CACHE_VERSION\|cache_namespace' docs/same-host-embedding-service-contract.md` | 0 | ✅ pass — cache namespace guidance present | 0ms |
| 6 | `grep 'same-host-embedding-service-contract' README.md` | 0 | ✅ pass — README links contract at line 59 | 0ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `docs/same-host-embedding-service-contract.md (created, 16 KB)`
- `README.md (updated: added contract link at line 59)`
