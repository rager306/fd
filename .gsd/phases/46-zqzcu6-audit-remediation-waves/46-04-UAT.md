# S04: Exposure posture policy — UAT

**Milestone:** M046-zqzcu6
**Written:** 2026-06-14T19:08:50.671Z

# S04: Exposure posture policy — UAT

**Milestone:** M046-zqzcu6
**Written:** 2026-06-14

## UAT Type

- UAT mode: runtime-executable
- Why this mode is sufficient: S04 changes backend HTTP exposure behavior. Runtime HTTP checks against the rebuilt local API prove the default no-key posture, while unit/static tests prove trusted-proxy and limiter internals.

## Preconditions

- API container has been rebuilt with S04 changes using `docker compose up -d --build api`.
- Local compose does not provide `FD_API_KEY`.

## Smoke Test

Request public probe endpoints and protected endpoints without credentials.

## Test Cases

### 1. Public probes remain unauthenticated

1. GET `/live`, `/ready`, `/health`, and `/v1/healthcheck`.
2. **Expected:** all return 200 without `Authorization`.

### 2. Protected inference fails closed without API key

1. POST a valid body to `/v1/embeddings` without `Authorization` while `FD_API_KEY` is unset.
2. **Expected:** HTTP 401 with error code `unauthorized`, type `authentication_error`, and param `authorization`.

### 3. Metrics is protected

1. GET `/metrics` without `Authorization` while `FD_API_KEY` is unset.
2. **Expected:** HTTP 401 with error code `unauthorized`.

### 4. OpenAPI remains public

1. GET `/openapi.json` without `Authorization`.
2. **Expected:** OpenAPI JSON is returned.

## Edge Cases

### Forwarded header spoofing

1. Unit test sends `X-Forwarded-For` with a different `RemoteAddr` after router trusted-proxy configuration.
2. **Expected:** `ClientIP()` uses direct remote address.

### Rate limiter state bound

1. Unit test fills limiter state to `maxRateLimitKeys`, advances time beyond the window, then inserts a new key.
2. **Expected:** expired buckets are pruned and the fresh key is retained without exceeding the cap.

## Failure Signals

- `/v1/embeddings` returns 200 without `FD_API_KEY`.
- `/metrics` returns 200 without `FD_API_KEY`.
- Probe endpoints require auth.
- `ClientIP()` trusts `X-Forwarded-For` by default.
- Limiter `items` grows beyond the configured cap.

## Requirements Proved By This UAT

- R030 — exposure posture is fail-closed for protected endpoints while public probes stay available.

## Not Proven By This UAT

- S05 LocalCache correctness.
- S06 residual P1 #6 and P2/P3 closure matrix.
- Authorized runtime inference with a configured key; this is covered by unit tests to avoid introducing secrets into runtime artifacts.

## Notes for Tester

Runtime evidence: `3b80d6c3-11aa-41ff-bf6a-e72531677268`, `237f9c00-e72a-4a19-8fb3-7f907075c417`, `85b9e147-57dc-4104-8e62-3bcbeea259a3`, `9484221c-844c-4aa6-9569-e59dfec157c7`.
Static proof: `15cb6196-085b-4882-8e4b-18d17008ee4d`.
Completeness proof: `beddf2e8-b429-4e99-99a5-cccd72fbcebe`.
