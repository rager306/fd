# S02: Graceful server error path — UAT

**Milestone:** M047-9fngng
**Written:** 2026-06-15T08:19:48.560Z

# S02: Graceful server error path — UAT

**Milestone:** M047-9fngng
**Written:** 2026-06-15

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: S02 changes backend process control flow around listener errors. The observable contract is code/test evidence; no browser surface is involved.

## Preconditions

- `benchmark-results/m047-s02-graceful-listener-shutdown.md` exists.

## Smoke Test

Verify listener fatal errors are routed through the lifecycle shutdown trigger and ErrServerClosed is matched with `errors.Is`.

## Test Cases

### 1. ErrServerClosed matching

1. Inspect `api/main.go`.
2. **Expected:** `reportHTTPServerError` uses `errors.Is(err, http.ErrServerClosed)` and contains no `os.Exit`.

### 2. Fatal listener error routing

1. Inspect `api/main.go`.
2. **Expected:** fatal listener errors send `serverErrorSignal{err: err}` into the signal channel and `main` starts `go reportHTTPServerError(logger, addr, srv.ListenAndServe, sigCh)`.

### 3. Evidence artifact complete

1. Inspect `benchmark-results/m047-s02-graceful-listener-shutdown.md`.
2. **Expected:** artifact covers #13/#32, red evidence, green evidence, and R035 validation.

## Edge Cases

- Wrapped `http.ErrServerClosed` is ignored.
- Non-server-closed listener error triggers `server_error` shutdown signal.
- Existing signal shutdown path remains the same lifecycle path.

## Failure Signals

- Listener goroutine calls `os.Exit(1)` again.
- Direct `err != http.ErrServerClosed` comparison returns.
- `go test ./...` fails.

## Requirements Proved By This UAT

- R035: fatal listener errors enter controlled graceful shutdown path and use `errors.Is`.

## Not Proven By This UAT

- TEI retry/fast-fail and warmup retry; these are downstream slices.

## Notes for Tester

UAT evidence: `969e79a5-18cb-426a-8b97-6bc13c4f079d`, `d61924dc-609c-41db-ac79-771c0577fccf`, `e3f5d9d0-e44d-4950-bb32-02ea8e775fc8`.
