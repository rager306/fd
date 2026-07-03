---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Added bounded warmup retry with deterministic tests.

Add a small warmup retry policy and update `startModelWarmup` to retry failed prewarm attempts with bounded backoff. Clear readiness error on success via existing lifecycle state behavior and record terminal errors after max attempts.

## Inputs

- `api/main.go`
- `api/main_test.go`

## Expected Output

- `api/main.go`
- `api/main_test.go`

## Verification

cd api && go test ./...

## Observability Impact

Warmup attempt logs include attempt/max and elapsed duration.
