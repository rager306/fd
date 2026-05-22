# S01: Same-host service contract

**Goal:** Define and verify the same-host local HTTP embedding service contract so neighboring services know how to call fd, identify the configured runtime, interpret health/readiness metadata, and avoid silent fallback or cache-contamination assumptions.
**Demo:** After this, local services have a concrete contract for calling fd and interpreting runtime readiness.

## Must-Haves

- Canonical same-host service contract document exists and covers endpoints, request/response shapes, status/error behavior, runtime/env expectations, health metadata, timeout/retry guidance, cache namespace guidance, and no-silent-fallback rules.
- README links to the canonical contract without duplicating the full operating contract.
- TEI/default health metadata is made programmatically visible if the contract requires runtime identification from `/health`, while ONNX safe metadata remains non-leaky.
- `/v1/embeddings` model-field behavior is either hardened with tests or explicitly documented as compatibility metadata with response/health model authoritative.
- Verification includes Go tests when code changes and a leak check over docs/artifacts.

## Proof Level

- This slice proves: Doc plus implementation verification. Static contract review plus Go tests for any code changes and a leak check proving no secrets, raw legal corpus text, or signed URLs were added to docs/artifacts.

## Integration Closure

S01 produces the service contract consumed by S02 restart/cache proof and S04 final runtime recommendation. It must not perform benchmark/runtime recommendation work itself.

## Verification

- Improves agent/operator visibility by making health metadata and readiness semantics explicit for same-host clients; avoids overclaiming TEI inference readiness unless implementation adds a probe.

## Tasks

- [x] **T01: Write same-host service contract document** `est:1h`
  Create `docs/same-host-embedding-service-contract.md` as the canonical local HTTP consumer contract. Cover `/health`, `/v1/embeddings`, `/embeddings/batch`, dimensions, encoding caveats, status/error behavior, runtime/env expectations, timeout/retry guidance, cache namespace guidance, no-silent-fallback rules, readiness limitations, and non-goals. Keep claims grounded in current code and S01 research.
  - Files: `docs/same-host-embedding-service-contract.md`, `.gsd/milestones/M040-pbp9z1/slices/S01/S01-RESEARCH.md`
  - Verify: Document inspection: required sections are present and do not overclaim `/health` as full TEI inference readiness unless code changes add that behavior.

- [x] **T02: Expose TEI runtime metadata in health** `est:1.5h`
  Update the runtime health path so TEI/default reports safe runtime metadata from `/health` comparable to ONNX basics: backend, configured model, dimensions, production/default flag, and cache namespace. Preserve ONNX safe metadata behavior and avoid exposing paths, signed URLs, tokens, raw input text, or secrets. Add/update health tests for TEI metadata.
  - Files: `api/main.go`, `api/handlers/health.go`, `api/handlers/health_test.go`
  - Verify: cd api && go test ./... -short

- [x] **T03: Resolve embeddings model-field contract** `est:1h`
  Choose the smallest safe behavior for `/v1/embeddings` request `model`: either harden the handler to reject non-empty model values that do not match the configured model with a 400 and tests, or document that request `model` is compatibility metadata and response model plus `/health.runtime.model` are authoritative. Prefer implementation hardening only if local tests and existing API expectations support it without broad breakage.
  - Files: `api/handlers/embeddings.go`, `api/handlers/embeddings_integration_test.go`, `api/handlers/health_test.go`, `docs/same-host-embedding-service-contract.md`
  - Verify: cd api && go test ./... -short

- [x] **T04: Link contract and verify S01 safety** `est:45m`
  Link the new contract from README. Run Go tests if code changed, run a docs/artifact leak check, and record verification evidence in the task summary. Ensure S01 remains scoped to contract/readiness metadata and does not perform S02 benchmark or S04 runtime recommendation work.
  - Files: `README.md`, `docs/same-host-embedding-service-contract.md`
  - Verify: cd api && go test ./... -short
rg -n "signed|token=|X-Amz|BEGIN|PRIVATE|юридическая справка" docs README.md benchmark-results .gsd/milestones/M040-pbp9z1/slices/S01

## Files Likely Touched

- docs/same-host-embedding-service-contract.md
- .gsd/milestones/M040-pbp9z1/slices/S01/S01-RESEARCH.md
- api/main.go
- api/handlers/health.go
- api/handlers/health_test.go
- api/handlers/embeddings.go
- api/handlers/embeddings_integration_test.go
- README.md
