# M047 S02 Graceful Listener Shutdown Evidence

Captured: 2026-06-15

## Scope

S02 covers GitHub issue #6 findings:

- #13 `ListenAndServe` fatal error calls `os.Exit(1)`, bypassing graceful shutdown.
- #32 direct `err != http.ErrServerClosed` comparison instead of `errors.Is`.

## Red Evidence

Command:

```bash
cd api && go test ./...
```

Expected red result after adding tests:

```text
build failed
undefined: reportHTTPServerError
undefined: serverErrorSignal
```

## Fix

- Added `serverErrorSignal`, a small `os.Signal` implementation whose `String()` is `server_error`.
- Added `reportHTTPServerError(logger, addr, listen, shutdownCh)`:
  - logs the listen address;
  - ignores wrapped `http.ErrServerClosed` with `errors.Is`;
  - logs fatal listener errors;
  - sends `serverErrorSignal` into the lifecycle signal channel instead of calling `os.Exit(1)`.
- Moved signal channel creation before listener startup and changed the listener goroutine to `go reportHTTPServerError(...)`.

## Green Evidence

Command:

```bash
cd api && gofmt -w main.go main_test.go && go test ./...
```

Result:

```text
285 passed in 9 packages
```

Static proof:

```text
gsd_exec 519aee78-cfa7-47d0-9fdf-aee5cddd1f83
PASS M047 S02 listener error routes to shutdown signal without goroutine os.Exit
```

## Requirement Outcome

- R035 validated for issue #6 findings #13 and #32.

## Residual Issue #6 Findings

Deferred to downstream M047 slices:

- #11 TEI retry/backoff/fast-fail.
- #14 warmup retry/backoff.
