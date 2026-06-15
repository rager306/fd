# M049 S02 Health and Metrics Context Evidence

Captured: 2026-06-15

## Scope

S02 implements GitHub issue #8 AN-B and AN-C:

- `/health` exposes lifecycle last error and dependency context.
- `/health` exposes in-flight capacity alongside current in-flight requests.
- `/metrics` exposes cheap runtime/capacity/cache occupancy gauges.

## Implemented

### Health

- Added `last_error {code,message,at}` when lifecycle state records an error.
- Added `dependencies.tei {reachable,latency_ms,error?}`.
- Added `dependencies.redis {reachable,latency_ms,namespace,error?}`.
- Added `in_flight_capacity`.
- Added `HealthOptions`, `DependencyChecks`, and `DependencyProbe` abstractions so health callers can inject lightweight probes.
- `main.go` wires TEI `/health` and Redis `PING` probes with a short timeout.

### Metrics

- Added `fd_in_flight_requests` gauge.
- Added `fd_in_flight_capacity` gauge.
- Added `fd_cache_entries{tier="l1"}` gauge.
- Metrics update these gauges at scrape time through `SetRuntimeObservers`.

## Safety Notes

- TEI dependency check uses TEI `/health`; it does not perform embedding inference.
- Redis dependency check uses `PING`; it does not scan keys.
- Cache occupancy is L1-only because Redis namespace counts would require `SCAN` on every metrics scrape. This keeps metrics cheap and appropriate for the user's solo deployment.
- Readiness semantics are unchanged: lifecycle readiness remains based on warmup/shutdown/error state.

## Verification

Red test command:

```bash
cd api && go test ./handlers ./observability
```

Expected red result before implementation:

```text
missing NewHealthHandlerWithOptions, HealthOptions, DependencyChecks, DependencyStatus, DeepHealthResponse.LastError, and Metrics.SetRuntimeObservers
```

Green commands:

```bash
cd api && go test ./handlers ./observability
cd api && go test ./...
```

Results:

```text
go test ./handlers ./observability: 92 passed in 2 packages
go test ./...: 295 passed in 10 packages
```

Static proof:

```text
gsd_exec 7b258850-f9dc-4b91-a7bf-706f4872eff5
PASS M049 S02 health and metrics invariants
```

## Requirement Outcome

- R041 advanced: health and metrics context are implemented and tested.
- Runtime proof for rebuilt-container `/health` and `/metrics` output is deferred to M049 S03.
