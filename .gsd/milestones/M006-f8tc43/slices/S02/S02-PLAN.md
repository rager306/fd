# S02: Add Testify testing helpers

**Goal:** Add Testify and migrate a small representative set of tests to assert/require style without broad churn.
**Demo:** After this, Testify is a project dependency and representative tests use it.

## Must-Haves

- `github.com/stretchr/testify` is added to `api/go.mod`/`api/go.sum`.
- At least one cache test and one handler test use Testify.
- Test semantics are preserved.
- `cd api && go test ./... -short` passes.

## Proof Level

- This slice proves: Go tests

## Integration Closure

Establishes Testify as project test helper while preserving existing coverage.

## Verification

- Test failures in migrated tests become clearer and less boilerplate-heavy.

## Tasks

- [x] **T01: Add Testify dependency** `est:small`
  Run GitNexus impact analysis for representative test functions before editing and add Testify dependency to api module.
  - Files: `api/go.mod`, `api/go.sum`
  - Verify: `go get github.com/stretchr/testify` succeeds and impact analysis recorded.

- [x] **T02: Migrate cache test assertions** `est:small`
  Migrate representative cache test assertions to Testify require/assert while preserving behavior.
  - Files: `api/cache/tiered_cache_test.go`
  - Verify: `cd api && go test ./cache -short` passes.

- [x] **T03: Migrate handler test assertions** `est:small`
  Migrate representative handler test assertions to Testify require/assert while preserving response behavior.
  - Files: `api/handlers/embeddings_integration_test.go`
  - Verify: `cd api && go test ./... -short` passes.

## Files Likely Touched

- api/go.mod
- api/go.sum
- api/cache/tiered_cache_test.go
- api/handlers/embeddings_integration_test.go
