# M046 S04 Exposure Posture

Captured: 2026-06-14

## Scope

Fix issue #3 exposure/auth posture findings:

- P0 #1: empty `FD_API_KEY` disabled authentication by default.
- P1 #7: `/metrics` was auth-exempt.
- P1 #8: rate-limit client identity was spoofable via forwarded headers and limiter state could grow unbounded.

## Changes

- Protected endpoints now fail closed when `FD_API_KEY` is missing.
- Public auth carve-outs remain limited to `/live`, `/ready`, `/health`, `/v1/healthcheck`, `/docs`, `/docs/*`, and `/openapi.json`.
- `/metrics` is no longer public.
- Gin trusted proxies are set to `nil` by default so `ClientIP()` does not trust spoofable `X-Forwarded-For` headers.
- Rate limiter key state is bounded by `maxRateLimitKeys` with expired-bucket pruning and oldest-bucket eviction.
- README documents the fail-closed auth posture without committing any secret value.

## Red Evidence

```bash
cd api && go test ./middleware -run 'TestAPIKeyAuth|TestRateLimiterPrunesExpiredBuckets'
```

Initial result:

```text
middleware/ratelimit_test.go: undefined: maxRateLimitKeys
```

```bash
cd api && go test . -run 'TestRouterDoesNotTrustForwardedForByDefault'
```

Initial result:

```text
./main_test.go: undefined: configureTrustedProxies
```

## Verification

### Targeted tests

```bash
cd api && go test ./middleware -run 'TestAPIKeyAuth|TestRateLimiterPrunesExpiredBuckets'
cd api && go test . -run 'TestRouterDoesNotTrustForwardedForByDefault'
```

Result:

```text
ok fd-api/middleware
ok fd-api
```

### Full Go suite

```bash
cd api && go test ./...
```

Result:

```text
Go test: 279 passed in 9 packages
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

### Static posture proof

Evidence: `gsd_exec:15cb6196-085b-4882-8e4b-18d17008ee4d`

```text
PASS S04 current static posture proof
```

### Runtime UAT after API rebuild

Evidence:

- `gsd_uat_exec:3b80d6c3-11aa-41ff-bf6a-e72531677268` — public probes remain unauthenticated.
- `gsd_uat_exec:237f9c00-e72a-4a19-8fb3-7f907075c417` — /v1/embeddings returns 401 without `FD_API_KEY`.
- `gsd_uat_exec:85b9e147-57dc-4104-8e62-3bcbeea259a3` — /metrics returns 401 without `FD_API_KEY`.
- `gsd_uat_exec:9484221c-844c-4aa6-9569-e59dfec157c7` — `/openapi.json` remains public.

## Findings Closed

- P0 #1: protected endpoints no longer default-open when `FD_API_KEY` is missing.
- P1 #7: `/metrics` is no longer auth-exempt.
- P1 #8: default `ClientIP()` behavior no longer trusts forwarded headers, and limiter state is bounded/pruned.

## Still Open

- S05: `LocalCache` size-counter/concurrency/lifecycle correctness.
- S06: residual P1 #6 and P2/P3 triage/closure matrix.
