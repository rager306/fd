---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Error envelope + 17 error codes catalog с table-driven tests

api/handlers/errors.go: type ErrorResponse {Error ErrorDetail}; type ErrorDetail {Code, Type, Param, Message}. Конструктор NewError(code, httpStatus) возвращает (httpStatus, body). Все 16 кодов из Section 3 каталога: input_required, input_too_long, batch_too_large, dimensions_invalid, dimensions_required, invalid_json, unauthorized, not_found, payload_too_large, rate_limit_exceeded, internal_error, model_not_loaded, model_overloaded, shutting_down, request_timeout, dimensions_mismatch. Каждый код маппится на (type, httpStatus) одной таблицей. Константы для каждого code value (goconst compliance). Helper WriteError(c *gin.Context, code string, param string, details ...any).

## Inputs

- None specified.

## Expected Output

- `api/handlers/errors.go`
- `api/handlers/errors_test.go`

## Verification

go test ./api/handlers/... -run TestErrorEnvelope покрывает все 16 кодов: проверяет HTTP status, code, type, message format. Все тесты pass. golangci-lint pass.
