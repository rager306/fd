---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T05: Live integration verification: 12 probe bugs против running fd container, все pass или в correct envelope

tests/integration/fd_v2_validation_test.go: автоматизировать все 15 error path test cases (T-E-1..T-E-15) из docs/fd-v2.md Section 5.2. Также regression test для backward compat: v1 caller POST /v1/embeddings {input:[hello]} → 200, response object/data/embedding/usage/model fields присутствуют. Также test для encoding_format=base64 (T-H-5) и dimensions=512 fix.

## Inputs

- None specified.

## Expected Output

- `tests/integration/fd_v2_validation_test.go`
- `tests/integration/fd_v1_backward_compat_test.go`

## Verification

go test ./tests/integration/... -run TestFdV2Validation -v: все 15 test cases pass. Backward compat test pass. encoding_format и dimensions=512 tests pass.
