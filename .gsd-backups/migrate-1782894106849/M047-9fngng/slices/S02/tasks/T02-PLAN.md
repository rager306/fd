---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Replaced listener goroutine os.Exit with controlled shutdown signalling.

Add a listener-error channel/helper, replace direct goroutine `os.Exit(1)`, use `errors.Is(err, http.ErrServerClosed)`, and ensure fatal listener errors enter `lifecycle.AwaitSignalAndShutdown` via signal/channel path or equivalent shared drain path.

## Inputs

- `api/main.go`
- `api/lifecycle/shutdown.go`

## Expected Output

- `api/main.go`
- `api/main_test.go`

## Verification

cd api && go test ./...

## Observability Impact

Fatal listener errors are logged once with controlled shutdown context.
