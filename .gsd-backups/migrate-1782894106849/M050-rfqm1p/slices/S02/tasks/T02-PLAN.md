---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Расширен `tests/integration` до auth-aware Docker e2e suite.

Расширить `tests/integration/api_test.go` текущими Docker e2e checks. Использовать helper functions for JSON requests, auth, cache status extraction and metrics text assertions. Не логировать секреты.

## Inputs

- `T01 contract`

## Expected Output

- `tests/integration/api_test.go`

## Verification

`cd tests/integration && go test -v .` passes with skips when no `FD_INTEGRATION_API_KEY`; authenticated run planned in T03.

## Observability Impact

E2E failures будут указывать, какой runtime contract сломан.
