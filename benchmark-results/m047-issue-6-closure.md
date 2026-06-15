# M047 Issue #6 Closure Matrix

Captured: 2026-06-15

Issue: https://github.com/rager306/fd/issues/6
Input artifact: `documents/issue-6-current-m047.md`

## Verification Summary

Final gates:

```text
go test ./...: 290 passed in 9 packages
golangci-lint: 0 issues
govulncheck: 0 reachable vulnerabilities
```

## Closure Matrix

| Issue #6 finding | Priority | Status | Evidence |
|---|---:|---|---|
| #11 TEI HTTP client has no retry / backoff / circuit breaker | P2 | Fixed | S03 added bounded retry for network and 502/503/504 failures, no-retry classification for non-retriable paths, and repeated-outage fast-fail circuit behavior. Evidence: `benchmark-results/m047-s03-tei-retry-fast-fail.md`; R033 validated. |
| #14 Warmup failure permanently degrades readiness; no auto-retry | P2 | Fixed | S04 added bounded warmup retry policy, per-attempt error recording/logging, and readiness recovery after later success. Evidence: `benchmark-results/m047-s04-warmup-retry-closure.md`; R034 validated. |
| #13 `ListenAndServe` fatal error calls `os.Exit(1)`, bypassing graceful shutdown | P2 | Fixed | S02 routes fatal listener errors through `serverErrorSignal` into the lifecycle shutdown path instead of exiting inside the listener goroutine. Evidence: `benchmark-results/m047-s02-graceful-listener-shutdown.md`; R035 validated. |
| #32 direct sentinel comparison instead of `errors.Is` | P3 | Fixed | S02 `reportHTTPServerError` uses `errors.Is(err, http.ErrServerClosed)` and tests wrapped ErrServerClosed behavior. Evidence: `benchmark-results/m047-s02-graceful-listener-shutdown.md`; R035 validated. |
| #25 error codes registered but never emitted | P3 | Fixed | S01 removed `dimensions_required`, `dimensions_mismatch`, and `request_timeout` from the public error registry and added `TestAllErrorCodesHaveNonTestEmitters`. Evidence: `benchmark-results/m047-s01-contract-cleanup.md`; R036 validated. |
| #15 `getEnvInt` overflow silently disables in-flight capacity gate | P2 | Fixed | S01 replaced manual digit accumulation with `strconv.Atoi` and fallback for invalid, overflow, and negative values. Evidence: `benchmark-results/m047-s01-contract-cleanup.md`; R036 validated. |

## Requirement Outcomes

- R033 validated: TEI retry/backoff and repeated-outage fast-fail.
- R034 validated: warmup retry and readiness recovery.
- R035 validated: graceful listener fatal-error path and `errors.Is`.
- R036 validated: safe env parsing and emitted error registry policy.

## Notes

- Browser/runtime UI verification is not applicable; issue #6 is backend reliability and API contract work.
- No GitHub issue comments or closure actions were performed by this milestone.
