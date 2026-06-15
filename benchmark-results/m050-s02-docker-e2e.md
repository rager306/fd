# M050 S02 Docker E2E Suite

Date: 2026-06-15
Milestone: M050-rfqm1p
Slice: S02

## Goal

Create and verify a maintained black-box Docker Compose e2e suite for the current fd service contract.

## E2E Contract

The suite lives in `tests/integration` and runs as an independent Go module.

Checks implemented in `tests/integration/api_test.go`:

- Public runtime diagnostics:
  - `GET /live` returns 200.
  - `GET /ready` returns 200.
  - `GET /health` returns status `ok` and exposes `runtime`, `dependencies`, and `in_flight_capacity`.
- Auth fail-closed:
  - `POST /v1/embeddings` without bearer token returns 401.
  - `POST /v1/cache/flush` without bearer token returns 401.
  - `POST /v1/cache/delete` without bearer token returns 401.
- Authenticated diagnostics:
  - `GET /metrics` with bearer token returns Prometheus text containing `fd_in_flight_requests`, `fd_in_flight_capacity`, and `fd_cache_entries{tier="l1"}`.
- Authenticated embeddings:
  - Batch `/v1/embeddings` with `dimensions=512` returns object `list`, two embedding objects, `dimensions=512`, and vector length 512.
  - Invalid JSON returns 400 with auth.
  - Empty input returns 400 with auth.
  - Missing request `model` remains accepted, because request model is compatibility metadata.
- Authenticated cache behavior:
  - `MISS -> HIT` for repeated input.
  - `POST /v1/cache/flush` makes the same input `MISS` again.
  - `POST /v1/cache/delete` deletes one input/dimension and makes the same input `MISS` again.

## Secret Handling

The e2e suite uses `FD_INTEGRATION_API_KEY` for client auth and never prints the value.

For local proof, the API container is recreated with a temporary Compose override that sets `FD_API_KEY` from an in-shell generated value, and the same value is passed to `FD_INTEGRATION_API_KEY`. The generated value is not printed or written to tracked files.

## Commands

No-key compile/public/fail-closed mode:

```bash
cd tests/integration && go test -v .
```

Authenticated runtime proof mode:

```bash
# Generate a local key without printing it, recreate api with a temporary compose override,
# then run: cd tests/integration && FD_INTEGRATION_API_KEY=<same key> go test -v .
```

## Results

### No-key mode

Command:

```bash
cd tests/integration && go test -v .
```

Result: PASS after S02 metrics fix.

- 5 passed in 1 package.
- Authenticated checks skipped when `FD_INTEGRATION_API_KEY` was not set.
- This mode proves public diagnostics, auth fail-closed, and test-code compilation without depending on a secret.

### Authenticated Docker Compose mode

Command shape:

```bash
# key generated in-shell and not printed
FD_E2E_API_KEY=<generated> docker compose -f docker-compose.yaml -f docker-compose.override.yaml -f /tmp/fd-m050-e2e-override.yaml up -d --force-recreate api
cd tests/integration && FD_INTEGRATION_API_KEY=<same generated key> go test -json .
```

Result: PASS.

```text
fd_api healthy after recreate
fd_redis healthy
fd_tei healthy
SUMMARY pass=9 fail=0 skip=0
PASS TestPublicRuntimeDiagnostics
PASS TestAuthenticatedMetricsDiagnostics
PASS TestProtectedEndpointsRequireAuth/embeddings
PASS TestProtectedEndpointsRequireAuth/cache_flush
PASS TestProtectedEndpointsRequireAuth/cache_delete
PASS TestProtectedEndpointsRequireAuth
PASS TestAuthenticatedEmbeddingDimensionsAndBatch
PASS TestAuthenticatedInvalidRequestsUseCurrentValidation
PASS TestAuthenticatedEmbeddingCacheInvalidation
```

### Correction made during S02

The initial no-key run expected `/metrics` to be public and failed with 401. The current service protects `/metrics`, so the test was corrected: public diagnostics cover `/live`, `/ready`, and `/health`; metrics diagnostics now run in the authenticated test path.

## Verdict

S02 established a maintained current-service Docker Compose e2e suite. The suite verifies public probes, health context, auth fail-closed, authenticated metrics, embeddings dimensions and batch behavior, validation errors, missing-model compatibility, and cache HIT/flush/delete invalidation against real `fd_api`, `fd_redis`, and `fd_tei` containers.
