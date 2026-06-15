---
id: T01
parent: S02
milestone: M049-7dn2gp
key_files:
  - api/handlers/health_test.go
  - api/observability/metrics_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T13:00:16.086Z
blocker_discovered: false
---

# T01: Pinned health/metrics diagnostic gaps with red tests.

**Pinned health/metrics diagnostic gaps with red tests.**

## What Happened

Added tests that expect `NewHealthHandlerWithOptions`, dependency probes, last_error, in_flight_capacity, dependency health fields, and metrics runtime/cache gauges. The refined red run fails because these APIs and response fields are not implemented yet, matching issue #8 AN-B/AN-C.

## Verification

`cd api && go test ./handlers ./observability` failed as expected with missing `NewHealthHandlerWithOptions`, `HealthOptions`, `DependencyChecks`, `DependencyProbeFunc`, `DependencyStatus`, `DeepHealthResponse.LastError`, and `Metrics.SetRuntimeObservers`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./handlers ./observability` | 1 | ✅ expected red | 8800ms |

## Deviations

Fixed an initial red-test harness mistake by using existing `state.TrackRequest()` signature before recording the red result.

## Known Issues

S02 remains red until health options and metrics observers are implemented.

## Files Created/Modified

- `api/handlers/health_test.go`
- `api/observability/metrics_test.go`
