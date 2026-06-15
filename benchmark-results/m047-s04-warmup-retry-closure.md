# M047 S04 Warmup Retry and Closure Evidence

Captured: 2026-06-15

## Scope

S04 covers GitHub issue #6 finding:

- #14 Warmup failure permanently degrades readiness; no auto-retry.

S04 also runs final milestone gates for M047.

## Red Evidence

Command:

```bash
cd api && go test ./...
```

Expected red result after adding tests:

```text
build failed
undefined: startModelWarmupWithPolicy
undefined: warmupRetryPolicy
```

## Fix

- Added `warmupRetryPolicy`.
- `startModelWarmup` now delegates to `startModelWarmupWithPolicy` with default max attempts of 3.
- Warmup attempts use bounded exponential backoff.
- Each failed attempt records `state.SetLastError(err)` and logs attempt/max/latency.
- A later successful attempt calls `state.MarkWarmupDone()`, which clears prior lifecycle errors.
- Terminal failure logs final error and attempt count.
- Tests use zero-delay injectable policy for deterministic retry proof.

## Green Evidence

Command:

```bash
cd api && gofmt -w main.go main_test.go && go test ./...
```

Result:

```text
290 passed in 9 packages
```

Static proof:

```text
gsd_exec 7ee9815e-9837-40f9-8430-8ef343422cdf
PASS M047 S04 warmup retry code invariants
```

## Final Gates

```bash
cd api && go test ./...
cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...
cd api && go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...
```

Results:

```text
go test ./...: 290 passed in 9 packages
golangci-lint: 0 issues
govulncheck: 0 reachable vulnerabilities
```

## Requirement Outcome

- R034 validated for bounded warmup retry and readiness recovery after later success.

## Residual Issue #6 Findings

None in M047 scope. Full closure matrix is in `benchmark-results/m047-issue-6-closure.md`.
