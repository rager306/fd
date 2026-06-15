# M050 S01 Existing Test Actuality Audit

Date: 2026-06-15
Milestone: M050-rfqm1p
Slice: S01

## Goal

Verify the existing test and verification surface against the current fd service contract before adding new e2e or mutation layers.

## Inventory

Source evidence: `.gsd/exec/c71d1245-b800-414e-9131-7614998c0e21.stdout`.

- Go test files under `api/`: 44
  - api root/internal: 6
  - buildinfo: 1
  - cache: 8
  - embed: 4
  - handlers: 11
  - lifecycle: 3
  - middleware: 8
  - observability: 2
  - openapi: 1
- Explicit integration-named Go tests inside `api/`:
  - `api/fd_v2_cache_integration_test.go`
  - `api/fd_v2_lifecycle_integration_test.go`
  - `api/fd_v2_observability_integration_test.go`
  - `api/handlers/embeddings_integration_test.go`
- Property-based tests:
  - `api/embed/codec_rapid_test.go`
- Root black-box integration test:
  - `tests/integration/api_test.go`
- Verification scripts:
  - `tools/verify_fd_v2_contract.py`
  - `tools/verify_legal_model_quick_gate_artifact.py`
  - `tools/verify_m040_s02_artifacts.py`
  - `tools/verify_m040_s04_recommendation.py`
  - `tools/verify_onnx_artifacts.py`
  - `tools/verify_onnx_export_contract.py`
- CI regular gate:
  - `.github/workflows/go-quality.yml` runs `cd api && go test ./... -short`, golangci-lint, and govulncheck.

## Fresh Baseline Commands

| Command | Result | Notes |
|---|---:|---|
| `cd api && go test ./...` | PASS | 295 passed in 10 packages. |
| `cd api && go test ./... -short` | PASS | 295 passed in 10 packages. |
| `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | PASS | 0 issues. |
| `cd api && go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | PASS | 0 reachable vulnerabilities. |
| `go test ./tests/integration` | FAIL before fix | Root has no Go module, so the command could not discover a main module. |
| `cd tests/integration && go test -v .` | FAIL during first fix attempt | Standalone `go 1.25` module attempted unavailable `go1.25` toolchain; corrected to `go 1.25.0`. |
| `cd tests/integration && go test -v .` | FAIL during second fix attempt | Protected endpoint tests used inherited `FD_API_KEY` that did not match running container; service correctly returned 401. |
| `cd tests/integration && go test -v .` | PASS after fix | 2 public/fail-closed tests passed; protected happy-path tests skip unless `FD_INTEGRATION_API_KEY` is explicitly supplied. |

## Actuality Findings

### Current

- `api/**/*_test.go` is current against the in-process API contracts and passes both normal and `-short` modes.
- CI quality gate is current for `api/`: test, lint, and govulncheck commands pass.
- `api/embed/codec_rapid_test.go` is property-based coverage, not mutation testing.
- ONNX artifact verification scripts remain historical/future-runtime artifact gates; TEI-only runtime is the current product path.

### Fixed

- `tests/integration/api_test.go` was stale as a runnable test layer:
  - It lived outside any Go module, so root `go test ./tests/integration` did not run.
  - It assumed protected `/v1/embeddings` requests could be sent without auth.
  - It used `FD_API_KEY` implicitly, which can be an unrelated shell value and cause false failures against the running container.
- Fixes applied:
  - Added `tests/integration/go.mod` with `go 1.25.0`.
  - Added `FD_BASE_URL` support with default `http://localhost:8000`.
  - Added explicit `TestEmbeddingsEndpoint_RequiresAuth` expecting `401` without a bearer token.
  - Protected happy-path checks now require `FD_INTEGRATION_API_KEY`; without it they skip rather than using an accidental secret or stale shell value.
  - Existing optional request-model expectation was preserved: request `model` is compatibility metadata; runtime response model is authoritative.

### Deferred to S02

- A full black-box Docker Compose suite with authenticated happy-path embeddings, dimensions, cache hit, cache invalidation, metrics, and dependency diagnostics is intentionally deferred to S02.
- S01 only made the existing root integration test honest and runnable; it did not expand it into the new e2e suite.

## Current Commands for Existing Tests

```bash
cd api && go test ./...
cd api && go test ./... -short
cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...
cd api && go run golang.org/x/vuln/cmd/govulncheck@latest ./...
cd tests/integration && go test -v .
```

For protected endpoint checks in `tests/integration`:

```bash
cd tests/integration && FD_INTEGRATION_API_KEY=<matching-running-service-key> go test -v .
```

Do not print or commit the key. If the key is not supplied, the protected happy-path integration checks skip; the public `/health` and fail-closed unauthenticated `/v1/embeddings` checks still run.

## Verdict

S01 found one stale existing test layer (`tests/integration`) and corrected it to match the current service posture. The primary `api` test/lint/vulnerability gates were already current and remain passing after the fix.
