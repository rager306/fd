---
id: S01
parent: M040-pbp9z1
milestone: M040-pbp9z1
provides:
  - Same-host local HTTP embedding service contract: endpoints, env/runtime requirements, health metadata, timeout/retry guidance, cache namespace guidance, and no-silent-fallback rules.
  - Programmatically visible TEI/default runtime metadata in `/health` with safe fields.
  - Documented and regression-covered model-field compatibility semantics for `/v1/embeddings`.
requires:
  []
affects:
  - S02
  - S03
  - S04
key_files:
  - docs/same-host-embedding-service-contract.md
  - README.md
  - api/handlers/health.go
  - api/handlers/health_test.go
  - api/main.go
  - api/handlers/embeddings_integration_test.go
  - .gsd/milestones/M040-pbp9z1/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M040-pbp9z1/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M040-pbp9z1/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M040-pbp9z1/slices/S01/tasks/T04-SUMMARY.md
key_decisions:
  - TEI/default `/health` now exposes safe runtime metadata fields for same-host clients.
  - `/v1/embeddings` request `model` remains compatibility metadata; response model and `/health.runtime.model` are authoritative.
  - `/health` must not be treated as a live inference probe; full readiness requires a smoke embedding request.
  - Cache namespace guidance is part of the contract to avoid TEI/ONNX cache-contamination assumptions.
patterns_established:
  - Canonical docs contract linked from README rather than duplicated.
  - Safe health metadata exposes operational identity without paths, tokens, signed URLs, raw text, or secrets.
  - Compatibility behavior is locked down with explicit tests and contract language when implementation hardening is not the safest choice.
observability_surfaces:
  - `/health.runtime.backend` identifies TEI/default vs ONNX.
  - `/health.runtime.model` reports the configured model.
  - `/health.runtime.dimensions` reports embedding dimensionality.
  - `/health.runtime.production_default` identifies default runtime posture.
  - `/health.runtime.cache_namespace` supports cache isolation diagnostics.
drill_down_paths:
  - .gsd/milestones/M040-pbp9z1/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M040-pbp9z1/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M040-pbp9z1/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M040-pbp9z1/slices/S01/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-22T05:35:34.077Z
blocker_discovered: false
---

# S01: Same-host service contract

**Established and verified the canonical same-host HTTP embedding service contract, including safe runtime health metadata, model-field semantics, cache guidance, and non-secret public docs.**

## What Happened

S01 produced `docs/same-host-embedding-service-contract.md` as the canonical contract for neighboring same-host services. The contract covers `/health`, `/v1/embeddings`, `/embeddings/batch`, dimensions, request/response shapes, status/error behavior, timeout and retry guidance, runtime and environment expectations, cache namespace guidance, no-silent-fallback rules, readiness limitations, and non-goals. README now links to the contract without duplicating it. The health path was updated so TEI/default runtime metadata is programmatically visible with safe fields only: backend, configured model, dimensions, production_default, and cache_namespace. ONNX health metadata remains safe and non-leaky. The `/v1/embeddings` request `model` field was resolved as compatibility metadata rather than a selector; the response model and `/health.runtime.model` are authoritative, with regression coverage and contract documentation. S01 stayed scoped to contract/readiness semantics and did not perform S02 benchmark proof or S04 runtime recommendation work.

## Verification

Fresh closeout verification passed through `gsd_exec`: `cd api && go test ./... -short` passed for fd-api, fd-api/cache, fd-api/embed, and fd-api/handlers (exec e9d77c24-186c-4c34-bbd9-ef007f4a8ad5). The required broad leak-audit command over docs, README, benchmark-results, and S01 artifacts completed (exec 5bfcc7f2-6462-44c3-aff6-864f0fa2b0e9); remaining matches are policy/history references to audit terms, not secret material. A focused high-risk check over current public docs/artifacts passed with no signed URL token parameters, private key blocks, or prohibited raw sample text (exec b3526174-14f4-49e1-9f11-9b692ef3c499). The earlier failed root-level `go test ./... -short` gate was caused by running outside the Go module; the slice plan's module-scoped command is `cd api && go test ./... -short`, which now passes.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

- None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

The closeout verification gate initially ran `go test ./... -short` from the repository root, which is outside the Go module and fails with `directory prefix . does not contain main module`. The slice plan's intended verification is `cd api && go test ./... -short`; that module-scoped verification passed. T03 also chose documentation plus regression coverage over handler hardening because the safest behavior was to keep request `model` as compatibility metadata.

## Known Limitations

`/health` is operational metadata and Redis/API startup visibility, not a live inference probe; clients still need a smoke embedding request for full runtime readiness. S01 does not prove Docker restart/cache behavior, alternative model quality, or the final runtime recommendation.

## Follow-ups

S02 should consume the contract to measure packaged Docker restart and Redis L2 cache behavior. S03 should use the contract scope boundary while performing the bounded legal model quick gate. S04 should use S01's contract and S02/S03 evidence to produce the final runtime recommendation and operating contract.

## Files Created/Modified

- `docs/same-host-embedding-service-contract.md` — Canonical same-host embedding service contract covering endpoints, runtime readiness, cache, model semantics, and non-goals.
- `README.md` — Added a link to the canonical same-host embedding service contract and kept examples non-sensitive.
- `api/handlers/health.go` — Added safe runtime metadata support to health responses.
- `api/handlers/health_test.go` — Covered TEI/default and ONNX health metadata behavior.
- `api/main.go` — Wired embedding runtime health metadata into service startup configuration.
- `api/handlers/embeddings_integration_test.go` — Added regression coverage for `/v1/embeddings` model-field compatibility behavior.
