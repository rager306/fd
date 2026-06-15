# M047 S03 TEI Retry and Fast Fail Evidence

Captured: 2026-06-15

## Scope

S03 covers GitHub issue #6 finding:

- #11 TEI HTTP client has no retry / backoff / circuit breaker.

## Red Evidence

Command:

```bash
cd api && go test ./embed
```

Expected red result after adding tests:

```text
19 passed, 2 failed in 1 package
TestTEIClientEmbedRetriesRetriableStatus: status 503 returned without retry.
TestTEIClientEmbedFastFailsAfterRepeatedRetriableFailures: returned status 503 instead of TEI circuit open.
```

## Fix

- `TEIClient` now has a bounded retry policy:
  - default max attempts: 3;
  - retriable status codes: 502, 503, 504;
  - retriable network errors unless the request context is canceled or expired;
  - context-aware exponential backoff with small defaults.
- Non-retriable paths do not retry:
  - 4xx responses;
  - decode errors;
  - empty successful payloads;
  - context cancellation/deadline.
- Repeated retriable call failures open a short cooldown circuit.
- Circuit-open calls fail fast with `TEI circuit open ...` and do not issue another HTTP request.
- Successful calls reset consecutive failure and circuit state.

## Green Evidence

Commands:

```bash
cd api && gofmt -w embed/tei.go embed/tei_test.go && go test ./embed
cd api && go test ./...
```

Results:

```text
go test ./embed: 21 passed in 1 package
go test ./...: 288 passed in 9 packages
```

Static proof:

```text
gsd_exec 06c49705-07f9-4c63-add2-85eb6ef673c9
PASS M047 S03 TEI retry and circuit code invariants
```

## Requirement Outcome

- R033 validated for TEI retry/backoff and repeated-outage fast-fail behavior.

## Residual Issue #6 Findings

Deferred to S04:

- #14 warmup retry/backoff.
