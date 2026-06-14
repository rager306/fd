# M046 S03 Batch Backend Chunking

Captured: 2026-06-14

## Scope

Fix issue #3 P1 #4 and P1 #5: batch endpoints performed one TEI/backend call per input on cache misses.

## Red Evidence

Command:

```bash
cd api && go test ./handlers -run 'Test(V1BatchUsesSingleEmbedCallPerInnerBatch|CreateBatchEmbeddingsUsesSingleEmbedCallForMisses)'
```

Initial result:

```text
--- FAIL: TestCreateBatchEmbeddingsUsesSingleEmbedCallForMisses
    embedder calls = 4, want 1
--- FAIL: TestV1BatchUsesSingleEmbedCallPerInnerBatch
    embedder calls = 8, want 2
```

## Changes

- Added `api/handlers/batch_backend.go` with package-local miss chunking helpers.
- `/embeddings/batch` now:
  - checks cache with `GetIfPresent`,
  - calls TEI once per bounded miss chunk of 32,
  - backfills cache with `Set`,
  - preserves input order and legacy base64-by-default response encoding.
- `/v1/batch` now:
  - checks cache with `GetIfPresent`,
  - calls TEI once per bounded inner batch of up to 32,
  - backfills cache with `Set`,
  - preserves nested response order.
- Wrong vector counts from the backend fail closed as internal errors.

## Verification

### Focused tests

```bash
cd api && go test ./handlers -run 'Test(V1BatchUsesSingleEmbedCallPerInnerBatch|CreateBatchEmbeddingsUsesSingleEmbedCallForMisses)'
```

Result:

```text
ok fd-api/handlers
```

The tests also repeat the same request and assert the second request adds no embedder calls, proving cache hits skip TEI.

### Handler package

```bash
cd api && go test ./handlers
```

Result:

```text
ok fd-api/handlers
```

### Full Go suite

```bash
cd api && go test ./...
```

Result:

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

```bash
cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...
```

Result:

```text
0 issues.
```

### Vulnerability scan

```bash
cd api && go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...
```

Result:

```text
No vulnerabilities found.
Your code is affected by 0 vulnerabilities.
```

### Static chunking proof

Evidence: `gsd_exec:6591611c-d4d4-4485-b17e-ac2be3aa5d6d`

```text
PASS S03 batch handlers use miss chunking and no per-input GetOrLoad
```

## Findings Closed

- P1 #4: `/embeddings/batch` no longer makes one TEI call per input on cache misses.
- P1 #5: `/v1/batch` no longer makes one TEI call per input on cache misses.

## Still Open

- S04: default-open auth/exposure posture, including metrics exposure and rate-limit trust boundary.
- S05: `LocalCache` size-counter/concurrency lifecycle correctness.
- S06: residual P2/P3 recheck and closure matrix.
