---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Pinned listener fatal-error routing with red tests.

Extract or test a small listener helper contract. Add tests proving wrapped `http.ErrServerClosed` is ignored and arbitrary listener errors are forwarded for main control flow handling.

## Inputs

- `api/main.go`
- `documents/issue-6-current-m047.md`

## Expected Output

- `api/main_test.go`

## Verification

cd api && go test ./... (expected red before implementation).

## Observability Impact

Tests pin shutdown behavior before changing entrypoint control flow.
