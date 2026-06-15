---
id: S02
parent: M047-9fngng
milestone: M047-9fngng
provides:
  - R035 validated.
  - Issue #6 findings #13 and #32 closed.
requires:
  - slice: S01
    provides: Contract cleanup baseline and issue input artifact.
affects:
  []
key_files:
  - api/main.go
  - api/main_test.go
  - benchmark-results/m047-s02-graceful-listener-shutdown.md
  - .gsd/REQUIREMENTS.md
key_decisions:
  - Use a synthetic `server_error` signal to preserve the existing lifecycle shutdown orchestration.
patterns_established:
  - Fatal process-control errors should be routed into lifecycle state transitions instead of exiting from background goroutines.
observability_surfaces:
  - Structured listener error log remains; S02 artifact records red/green/static/UAT proof.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-15T08:19:48.560Z
blocker_discovered: false
---

# S02: Graceful server error path

**Fatal HTTP listener errors now trigger controlled lifecycle shutdown instead of exiting from the listener goroutine.**

## What Happened

S02 addressed issue #6 #13 and #32. Red tests first established the desired contract for a listener helper. The implementation added `serverErrorSignal` and `reportHTTPServerError`, which uses `errors.Is(err, http.ErrServerClosed)` for wrapped server-closed errors, logs fatal listener errors, and sends a synthetic `server_error` signal into the same channel consumed by `lifecycle.AwaitSignalAndShutdown`. The direct `os.Exit(1)` from the listener goroutine was removed. R035 was validated with full tests, static proof, artifact evidence, and UAT.

## Verification

Red evidence: `go test ./...` failed with undefined helper/type before implementation. Green evidence: `cd api && gofmt -w main.go main_test.go && go test ./...` passed with 285 tests. Static proof `519aee78-cfa7-47d0-9fdf-aee5cddd1f83` passed. Artifact completeness `99df63d5-02ac-4c7f-a076-b97f4b3b1da5` passed. UAT PASS was saved with evidence `969e79a5-18cb-426a-8b97-6bc13c4f079d`, `d61924dc-609c-41db-ac79-771c0577fccf`, and `e3f5d9d0-e44d-4950-bb32-02ea8e775fc8`.

## Requirements Advanced

None.

## Requirements Validated

- R035 — S02 tests and static proof validate controlled listener fatal error shutdown and `errors.Is` matching.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Used a synthetic `os.Signal` to reuse existing lifecycle shutdown path rather than introducing a parallel shutdown API.

## Known Limitations

S03/S04 still need TEI retry/fast-fail and warmup retry.

## Follow-ups

Proceed to S03 TEI retry and fast-fail behavior.

## Files Created/Modified

- `api/main.go` — Added listener helper and routed fatal errors into lifecycle signal channel.
- `api/main_test.go` — Added listener helper red/green tests.
- `benchmark-results/m047-s02-graceful-listener-shutdown.md` — S02 evidence artifact.
- `.gsd/REQUIREMENTS.md` — R035 validated.
