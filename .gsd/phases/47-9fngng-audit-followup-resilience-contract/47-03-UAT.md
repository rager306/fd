# S03: TEI retry and fast fail — UAT

**Milestone:** M047-9fngng
**Written:** 2026-06-15T08:28:06.004Z

# S03: TEI retry and fast fail — UAT

**Milestone:** M047-9fngng
**Written:** 2026-06-15

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: S03 changes backend TEI dependency behavior. The observable contract is covered by httptest/fake server unit tests, full Go tests, static proof, and code artifact checks; no browser surface is involved.

## Preconditions

- `benchmark-results/m047-s03-tei-retry-fast-fail.md` exists.

## Smoke Test

Verify TEI retry and circuit invariants.

## Test Cases

### 1. Bounded retry policy

1. Inspect `api/embed/tei.go`.
2. **Expected:** `TEIClient` includes retry attempts and retriable 502/503/504 classification.

### 2. Fast-fail circuit

1. Inspect `api/embed/tei.go`.
2. **Expected:** repeated retriable failures are recorded, circuit state opens, and circuit-open calls return a `TEI circuit open` error.

### 3. Evidence artifact complete

1. Inspect `benchmark-results/m047-s03-tei-retry-fast-fail.md`.
2. **Expected:** artifact covers #11, red evidence, green evidence, and R033 validation.

## Edge Cases

- 400 responses do not retry.
- Decode errors do not retry.
- Context cancellation/deadline does not retry.
- Successful calls reset circuit state.

## Failure Signals

- `httpClient.Do(req)` appears without retry orchestration.
- `TEI circuit open` path disappears.
- `go test ./embed` or `go test ./...` fails.

## Requirements Proved By This UAT

- R033: TEI transient failures have bounded retry/backoff and repeated outage fast-fail behavior.

## Not Proven By This UAT

- Warmup retry; this is S04.

## Notes for Tester

UAT evidence: `b7f0d8fc-6db3-4e92-833a-ccf4ac74796c`, `d606b15c-39eb-46f3-a785-e7867d2a4a3f`, `8ed72f5f-2c26-4a22-bbe6-5d9b134f47c8`.
