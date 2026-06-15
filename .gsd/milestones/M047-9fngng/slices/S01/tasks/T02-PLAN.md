---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Fixed safe env integer parsing and removed un-emitted error codes from the public registry.

Replace hand-rolled `getEnvInt` parsing with safe `strconv.Atoi` behavior that falls back on invalid, overflow, and negative values. Resolve the error-code registry contract by removing un-emitted reserved codes or adding explicit reserved metadata with test coverage; prefer removal unless current docs/tests require them.

## Inputs

- `api/main.go`
- `api/handlers/errors.go`
- `api/handlers/errors_test.go`
- `api/main_env_test.go`

## Expected Output

- `api/main.go`
- `api/handlers/errors.go`
- `api/handlers/errors_test.go`
- `api/main_env_test.go`

## Verification

cd api && go test ./...

## Observability Impact

Safe env parsing reduces silent misconfiguration risk.
