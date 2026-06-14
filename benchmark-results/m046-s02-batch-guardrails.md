# M046 S02 Batch Endpoint Guardrails

Captured: 2026-06-14

## Scope

Fix issue #3 P0 #2 and P0 #3: batch endpoints were mounted outside the full request guardrail posture.

## Changes

- Added shared `middleware.LimitRequestBody()` for endpoints with non-`/v1/embeddings` JSON shapes.
- Mounted both batch routes with:
  - `middleware.LimitRequestBody()`
  - `middleware.UserRateLimitFromEnv()`
  - `middleware.LifecycleGateWithCapacity(...)`
- Added batch-specific input validation before cache/TEI work:
  - legacy `/embeddings/batch`: non-empty inputs, max 128 inputs, max 2048 chars per input.
  - `/v1/batch`: existing max groups and max inner inputs plus max 2048 chars per input.
- Added MaxBytesError handling in batch handlers so body cap failures return `payload_too_large` instead of generic invalid JSON.

## Verification

### Red test before implementation

`cd api && go test ./handlers -run 'Test(V1BatchHandlerRejectsTooLongInputBeforeEmbedder|CreateBatchEmbeddingsRejectsTooLongInputBeforeEmbedder)'`

Result before implementation:

```text
handlers [build failed]
handlers/embeddings_integration_test.go:357:73: undefined: maxBatchInputChars
handlers/v1batch_test.go:87:77: undefined: maxBatchInputChars
```

### Targeted tests after implementation

`cd api && go test ./middleware ./handlers`

Result:

```text
ok fd-api/middleware
ok fd-api/handlers
```

### Full Go suite

`cd api && go test ./...`

Result after implementation and lint cleanup:

```text
ok fd-api
ok fd-api/buildinfo
ok fd-api/cache
ok fd-api/embed
ok fd-api/handlers
ok fd-api/lifecycle
ok fd-api/middleware
ok fd-api/observability
ok fd-api/openapi
```

### Lint

`cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...`

Result:

```text
0 issues.
```

### Vulnerability scan

`cd api && go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...`

Result:

```text
No vulnerabilities found.
Your code is affected by 0 vulnerabilities.
```

### Static route guardrail check

Evidence: `gsd_exec:070f8b99-679d-45e2-a90e-597c350f6837`

```text
PASS batch routes include body, rate-limit, lifecycle guardrails
```

## Findings Closed

- P0 #2: `/embeddings/batch` no longer mounts as a bare handler.
- P0 #3: `/v1/batch` no longer has lifecycle-only route guardrails.

## Still Open For S03

- P1 #4 and P1 #5: batch handlers still perform per-input backend calls on cache misses. Input guardrails now make that bounded, but S03 should reshape backend work to collect misses and issue chunked TEI calls.
