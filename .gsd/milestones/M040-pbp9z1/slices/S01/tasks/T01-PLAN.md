---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Write same-host service contract document

Create `docs/same-host-embedding-service-contract.md` as the canonical local HTTP consumer contract. Cover `/health`, `/v1/embeddings`, `/embeddings/batch`, dimensions, encoding caveats, status/error behavior, runtime/env expectations, timeout/retry guidance, cache namespace guidance, no-silent-fallback rules, readiness limitations, and non-goals. Keep claims grounded in current code and S01 research.

## Inputs

- `.gsd/milestones/M040-pbp9z1/M040-pbp9z1-CONTEXT.md`
- `.gsd/milestones/M040-pbp9z1/slices/S01/S01-RESEARCH.md`
- `README.md`
- `api/embed/types.go`
- `api/handlers/health.go`
- `api/handlers/embeddings.go`
- `api/handlers/batch.go`

## Expected Output

- `docs/same-host-embedding-service-contract.md`

## Verification

Document inspection: required sections are present and do not overclaim `/health` as full TEI inference readiness unless code changes add that behavior.

## Observability Impact

Defines the health/readiness and runtime-identification contract local clients and future agents will use.
