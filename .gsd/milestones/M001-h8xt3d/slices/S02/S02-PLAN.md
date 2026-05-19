# S02: API validation and handler tests

**Goal:** Make API validation consistent and improve handler test fidelity.
**Demo:** Invalid batch dimensions or encoding formats return HTTP 400; production handler paths are tested directly.

## Must-Haves

- /embeddings/batch rejects invalid dimensions.
- /embeddings/batch rejects invalid encoding_format.
- Handler tests exercise production handler code where practical.
- Existing valid request behavior remains covered.

## Proof Level

- This slice proves: handler tests plus full Go short suite

## Integration Closure

HTTP contract remains compatible for valid requests and stricter for invalid requests.

## Verification

- Bad input is rejected with clear 400 errors rather than silently coerced defaults.

## Tasks

- [x] **T01: Assess handler blast radius** `est:small`
  Run impact analysis for handler constructor/signature and batch validation changes, documenting direct call sites and tests affected.
  - Files: `api/handlers/embeddings.go`, `api/handlers/batch.go`, `api/handlers/*_test.go`, `api/main.go`
  - Verify: No code changes; document GitNexus/LSP/text-search findings.

- [x] **T02: Implement handler validation and tests** `est:medium`
  Introduce minimal handler dependency interfaces so tests can instantiate production handlers with mocks; add strict batch validation for dimensions and encoding_format.
  - Files: `api/handlers/embeddings.go`, `api/handlers/batch.go`, `api/handlers/embeddings_integration_test.go`
  - Verify: cd api && go test ./handlers -count=1

- [x] **T03: Verify API validation slice** `est:small`
  Run full Go short suite and commit S02 changes if passing.
  - Verify: cd api && go test ./... -short

## Files Likely Touched

- api/handlers/embeddings.go
- api/handlers/batch.go
- api/handlers/*_test.go
- api/main.go
- api/handlers/embeddings_integration_test.go
