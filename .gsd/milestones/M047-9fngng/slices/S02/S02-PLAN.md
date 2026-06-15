# S02: Graceful server error path

**Goal:** Refactor HTTP listener fatal error handling so non-ErrServerClosed failures are reported to main control flow, matched with errors.Is, and shut down through the same graceful path used by signals.
**Demo:** Fatal listener errors enter controlled shutdown handling instead of exiting inside the listener goroutine.

## Must-Haves

- Tests prove `http.ErrServerClosed` and wrapped ErrServerClosed are ignored using `errors.Is`.
- Tests prove a non-server-closed listener error is reported through a channel and does not call `os.Exit` inside the listener goroutine.
- `main` handles listener errors by triggering lifecycle graceful shutdown and resource cleanup.
- R035 is validated.

## Proof Level

- This slice proves: Focused unit tests around extracted listener helper plus `cd api && go test ./...`.

## Integration Closure

Signal-triggered shutdown remains unchanged; fatal listener error path reuses lifecycle shutdown semantics.

## Verification

- Listener fatal errors remain structured logs with error context and controlled shutdown state.

## Tasks

- [x] **T01: Pinned listener fatal-error routing with red tests.** `est:small`
  Extract or test a small listener helper contract. Add tests proving wrapped `http.ErrServerClosed` is ignored and arbitrary listener errors are forwarded for main control flow handling.
  - Files: `api/main.go`, `api/main_test.go`
  - Verify: cd api && go test ./... (expected red before implementation).

- [x] **T02: Replaced listener goroutine os.Exit with controlled shutdown signalling.** `est:medium`
  Add a listener-error channel/helper, replace direct goroutine `os.Exit(1)`, use `errors.Is(err, http.ErrServerClosed)`, and ensure fatal listener errors enter `lifecycle.AwaitSignalAndShutdown` via signal/channel path or equivalent shared drain path.
  - Files: `api/main.go`, `api/main_test.go`
  - Verify: cd api && go test ./...

- [x] **T03: Recorded S02 evidence and validated R035.** `est:small`
  Write S02 evidence artifact, validate R035, run full tests, and complete S02.
  - Files: `benchmark-results/m047-s02-graceful-listener-shutdown.md`, `.gsd/REQUIREMENTS.md`
  - Verify: cd api && go test ./... plus static check that listener goroutine no longer calls os.Exit.

## Files Likely Touched

- api/main.go
- api/main_test.go
- benchmark-results/m047-s02-graceful-listener-shutdown.md
- .gsd/REQUIREMENTS.md
