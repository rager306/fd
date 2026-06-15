# S03: TEI retry and fast fail

**Goal:** Add bounded retry/backoff and repeated-outage fast-fail behavior to the TEI HTTP client without changing successful embedding semantics.
**Demo:** TEI transient dependency failures retry within a bounded policy and repeated outages fail quickly with a clear error.

## Must-Haves

- Tests prove network errors and 502/503/504 retry with bounded attempts.
- Tests prove 4xx and malformed successful responses are not retried unnecessarily.
- Tests prove repeated retriable failures open a short-circuit state and subsequent calls fail quickly until cooldown.
- R033 is validated.

## Proof Level

- This slice proves: httptest/fake-transport tests with tiny injected sleeps plus `cd api && go test ./...`.

## Integration Closure

Existing TEI success, request shape, response ordering, and batch behavior remain green.

## Verification

- Retry and fast-fail errors include context sufficient to diagnose TEI dependency outages without secrets.

## Tasks

- [x] **T01: Pinned TEI retry and fast-fail gaps with red tests.** `est:medium`
  Add tests in `api/embed/tei_test.go` that expect bounded retries for transient network/5xx failures, no retry for non-retriable responses, and fast-fail behavior after repeated retriable failures.
  - Files: `api/embed/tei_test.go`
  - Verify: cd api && go test ./embed (expected red before implementation).

- [x] **T02: Added bounded TEI retry and repeated-outage fast-fail behavior.** `est:medium`
  Add a small retry policy to `TEIClient`: bounded attempts, jitter/backoff injection or tiny defaults, classification for network errors and 502/503/504, and a lightweight circuit breaker that short-circuits repeated retriable failures for a cooldown. Preserve existing constructor compatibility.
  - Files: `api/embed/tei.go`, `api/embed/tei_test.go`
  - Verify: cd api && go test ./embed && cd api && go test ./...

- [x] **T03: Recorded S03 evidence and validated R033.** `est:small`
  Write S03 evidence artifact, validate R033, run full tests, and complete S03.
  - Files: `benchmark-results/m047-s03-tei-retry-fast-fail.md`, `.gsd/REQUIREMENTS.md`
  - Verify: cd api && go test ./embed && cd api && go test ./... plus static proof.

## Files Likely Touched

- api/embed/tei_test.go
- api/embed/tei.go
- benchmark-results/m047-s03-tei-retry-fast-fail.md
- .gsd/REQUIREMENTS.md
