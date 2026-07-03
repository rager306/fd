# S01: Quality tooling baseline

**Goal:** Inspect current Go testing patterns, tool availability, and lint baseline before adding dependencies/config.
**Demo:** After this, current Go test/lint/tooling baseline is known before dependency/config edits.

## Must-Haves

- Existing Go test command is identified and run.
- Current lint/staticcheck availability is known.
- Existing test style and candidate Testify migration files are identified.
- Git status and unpushed commits are understood.

## Proof Level

- This slice proves: command evidence

## Integration Closure

Establishes a reliable baseline for Testify and lint changes.

## Verification

- Records current quality gate state before changes.

## Tasks

- [x] **T01: Inspect current Go tests and git state** `est:small`
  Inspect Go module, current tests, and git state to identify Testify migration candidates and pending commits.
  - Files: `api/go.mod`, `api/cache/tiered_cache_test.go`, `api/handlers/embeddings_integration_test.go`
  - Verify: Findings recorded.

- [x] **T02: Run baseline Go tests** `est:small`
  Run existing Go tests before changes to establish baseline.
  - Files: `api`
  - Verify: `cd api && go test ./... -short` passes.

- [x] **T03: Inspect lint tool availability** `est:small`
  Check availability/version of golangci-lint and staticcheck, and run current fallback lint if available.
  - Files: `api`
  - Verify: Tool availability and baseline command result recorded.

## Files Likely Touched

- api/go.mod
- api/cache/tiered_cache_test.go
- api/handlers/embeddings_integration_test.go
- api
