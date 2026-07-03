---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T06: Added executable lifecycle integration coverage for startup readiness, F-1 model_not_loaded, F-2 model_overloaded, and F-5 shutdown drain behavior.

tests/integration/fd_v2_lifecycle_test.go: воспроизвести F-1 (caller hit во время warmup → 503 model_not_loaded + Retry-After), F-2 (concurrent overload → 503 model_overloaded + Retry-After, после снижения load → 200), F-5 (SIGTERM → 503 shutting_down + drain). Также test: startup sequence — /live=200, /ready=503, /ready=200 после warmup, /health deep корректно меняется. Спека: docs/fd-v2.md Section 6.1 + 6.3 F-1/F-2/F-5.

## Inputs

- None specified.

## Expected Output

- `tests/integration/fd_v2_lifecycle_test.go`

## Verification

go test ./tests/integration/... -run TestFdV2Lifecycle -v: F-1, F-2, F-5, и startup sequence test все pass.
