---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Pinned TEI retry and fast-fail gaps with red tests.

Add tests in `api/embed/tei_test.go` that expect bounded retries for transient network/5xx failures, no retry for non-retriable responses, and fast-fail behavior after repeated retriable failures.

## Inputs

- `api/embed/tei.go`
- `documents/issue-6-current-m047.md`

## Expected Output

- `api/embed/tei_test.go`

## Verification

cd api && go test ./embed (expected red before implementation).

## Observability Impact

Tests document retry classification and fast-fail contract.
