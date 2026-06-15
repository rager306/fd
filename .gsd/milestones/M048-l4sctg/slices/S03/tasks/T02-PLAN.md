---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Fixed validation message clarity and OpenAPI helper fail-loud behavior.

Guard empty `json.UnmarshalTypeError.Field` in validation and change `openapi.m()` non-string key behavior from silent continue to panic/fail-loud. Keep existing valid OpenAPI spec tests green.

## Inputs

- `api/middleware/validation.go`
- `api/openapi/spec.go`

## Expected Output

- `api/middleware/validation.go`
- `api/openapi/spec.go`

## Verification

cd api && go test ./middleware ./openapi && cd api && go test ./...

## Observability Impact

Malformed client payloads and generator mistakes become easier to diagnose.
