# S02: Health and metrics context — UAT

**Milestone:** M049-7dn2gp
**Written:** 2026-06-15T13:06:46.625Z

# S02: Health and metrics context — UAT

**Milestone:** M049-7dn2gp
**Written:** 2026-06-15

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: S02 validates backend health/metrics source and tests. Runtime `/health` and `/metrics` verification is explicitly deferred to S03 after aggregate container rebuild.

## Preconditions

- `benchmark-results/m049-s02-health-metrics-context.md` exists.
- Focused and full Go tests have passed.

## Smoke Test

Verify health diagnostics, metrics gauges, and evidence artifact completeness.

## Test Cases

### 1. Health diagnostic fields

1. Inspect `api/handlers/health.go`.
2. Expected: `last_error`, `dependencies`, `in_flight_capacity`, dependency probe plumbing, and options handler exist.

### 2. Metrics gauges

1. Inspect `api/observability/metrics.go`.
2. Expected: `fd_in_flight_requests`, `fd_in_flight_capacity`, `fd_cache_entries`, and `SetRuntimeObservers` exist.

### 3. Evidence artifact

1. Inspect `benchmark-results/m049-s02-health-metrics-context.md`.
2. Expected: artifact records health/metrics implementation and test proof.

## Requirements Proved By This UAT

- R041 advanced for health/metrics implementation.

## Not Proven By This UAT

- Live container `/health` and `/metrics` output. This is planned for S03 runtime UAT.

## Notes for Tester

Evidence IDs: `14dc034d-e147-49b6-9a35-0723f3553065`, `00a1dc71-1adb-469e-bdb2-c7af7b58e15b`, `f7e110b0-1e20-4efa-82ba-de00341b696a`.
