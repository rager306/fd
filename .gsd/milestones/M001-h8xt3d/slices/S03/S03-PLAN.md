# S03: Local cache semantics

**Goal:** Make LocalCache behavior match its configuration and cache expectations.
**Demo:** L1 cache enforces documented size/overwrite semantics with tests.

## Must-Haves

- LocalCache.Set refreshes existing key value and TTL.
- maxSize behavior is explicit and tested.
- Existing TTL/delete behavior remains passing.

## Proof Level

- This slice proves: unit tests plus full Go short suite

## Integration Closure

Tiered cache continues using LocalCache without API changes.

## Verification

- Reduces memory surprise from unbounded local entries.

## Tasks

- [x] **T01: Assess LocalCache blast radius** `est:small`
  Run impact analysis for LocalCache.Set and related size bookkeeping before editing.
  - Files: `api/cache/local.go`, `api/cache/local_test.go`, `api/cache/tiered.go`
  - Verify: No code changes; document findings.

- [x] **T02: Implement LocalCache semantics** `est:medium`
  Change LocalCache.Set to overwrite value/TTL for existing keys, enforce maxSize for new keys, and add tests for overwrite and capacity behavior.
  - Files: `api/cache/local.go`, `api/cache/local_test.go`
  - Verify: cd api && go test ./cache -run 'TestLocalCache' -count=1

- [x] **T03: Verify LocalCache slice** `est:small`
  Run full short suite and commit S03 changes if passing.
  - Verify: cd api && go test ./... -short

## Files Likely Touched

- api/cache/local.go
- api/cache/local_test.go
- api/cache/tiered.go
