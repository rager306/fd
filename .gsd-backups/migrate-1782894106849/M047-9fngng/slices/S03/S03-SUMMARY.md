---
id: S03
parent: M047-9fngng
milestone: M047-9fngng
provides:
  - R033 validated.
  - Issue #6 finding #11 closed.
requires:
  - slice: S01
    provides: Issue input artifact and contract cleanup baseline.
affects:
  []
key_files:
  - api/embed/tei.go
  - api/embed/tei_test.go
  - benchmark-results/m047-s03-tei-retry-fast-fail.md
  - .gsd/REQUIREMENTS.md
key_decisions:
  - Use lightweight process-local circuit state inside TEIClient rather than introducing an external dependency or global breaker.
patterns_established:
  - Dependency retry policies should classify retriable errors explicitly and expose fast-fail behavior during repeated outages.
observability_surfaces:
  - TEI retry exhaustion and circuit-open errors are distinguishable in returned errors; S03 artifact records proof.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-15T08:28:06.003Z
blocker_discovered: false
---

# S03: TEI retry and fast fail

**TEI calls now retry bounded transient failures and fast-fail during repeated outages.**

## What Happened

S03 addressed issue #6 #11. Red tests first proved 503 responses were not retried and repeated failures did not open a circuit. The implementation added bounded retry policy to `TEIClient`: three attempts, context-aware exponential backoff, retry classification for network errors and 502/503/504 statuses, no retry for context cancellation/deadline, 4xx, decode errors, or malformed successful payloads, and a short cooldown circuit after repeated retriable call failures. Successful calls reset failure/circuit state. R033 was validated with focused tests, full tests, static proof, artifact evidence, and UAT.

## Verification

Red evidence: `go test ./embed` failed on retry and circuit tests. Green evidence: `cd api && gofmt -w embed/tei.go embed/tei_test.go && go test ./embed` passed with 21 tests; `cd api && go test ./...` passed with 288 tests. Static proof `06c49705-07f9-4c63-add2-85eb6ef673c9` passed. Artifact completeness `5dcac7b4-b032-42f6-9ed2-7adfc19219cd` passed. UAT PASS was saved with evidence `b7f0d8fc-6db3-4e92-833a-ccf4ac74796c`, `d606b15c-39eb-46f3-a785-e7867d2a4a3f`, and `8ed72f5f-2c26-4a22-bbe6-5d9b134f47c8`.

## Requirements Advanced

None.

## Requirements Validated

- R033 — S03 httptest tests and static proof validate bounded TEI retry and circuit fast-fail behavior.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Used deterministic context-aware exponential backoff instead of jitter to keep behavior bounded and testable in-process.

## Known Limitations

Circuit state is process-local and intentionally lightweight; distributed circuit coordination is out of scope for same-host fd.

## Follow-ups

Proceed to S04 warmup retry and milestone closure.

## Files Created/Modified

- `api/embed/tei.go` — Added retry policy and circuit state to TEIClient.
- `api/embed/tei_test.go` — Added retry, no-retry, and circuit tests.
- `benchmark-results/m047-s03-tei-retry-fast-fail.md` — S03 evidence artifact.
- `.gsd/REQUIREMENTS.md` — R033 validated.
