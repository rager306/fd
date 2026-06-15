---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Pinned S01 issue #6 contract failures with red tests.

Use issue #6 and current code to pin S01 scope. Add tests showing safe fallback expectations for `getEnvInt` overflow/negative values and registry expectations for canonical error codes.

## Inputs

- `documents/issue-6-current-m047.md`
- `api/main.go`
- `api/handlers/errors.go`

## Expected Output

- `api/main_env_test.go`
- `api/handlers/errors_test.go`

## Verification

cd api && go test ./... (expected red before fixes), then focused failing tests are recorded.

## Observability Impact

Red tests document the failure mode instead of relying on prose.
