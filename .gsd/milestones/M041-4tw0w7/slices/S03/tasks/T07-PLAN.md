---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T07: Added executable observability integration coverage for health/version/info/metrics/warmup endpoints and response headers.

tests/integration/fd_v2_observability_test.go: автоматизировать Section 5.1 T-H-7..T-H-10, Section 5.3 T-HDR-1..T-HDR-10 (кроме T-HDR-6/7 которые зависят от cache в S04), Section 5.5 T-E-1..T-E-3 (endpoints existence). Спека: docs/fd-v2.md Section 5.1, 5.3, 5.5.

## Inputs

- None specified.

## Expected Output

- `tests/integration/fd_v2_observability_test.go`

## Verification

go test ./tests/integration/... -run TestFdV2Observability -v: все 22 test cases pass.
