---
id: T02
parent: S01
milestone: M041-4tw0w7
key_files:
  - api/handlers/errors.go
  - api/handlers/errors_test.go
key_decisions:
  - encoding_format_invalid добавлен как 17-й code для support R-P1-5 в /v1/embeddings.
  - Unknown code в WriteError fails closed as internal_error — никогда не отдаём non-canonical envelope.
  - Использую goconst-friendly константы для каждого code value (goconst в golangci config).
  - Helper WriteErrorWithRetryAfter отдельный от WriteError чтобы 429/503 paths явно задавали retry, и легко тестировались.
duration: 
verification_result: untested
completed_at: 2026-06-13T18:02:23.724Z
blocker_discovered: false
---

# T02: Error envelope + 17 error codes catalog с table-driven tests

**Error envelope + 17 error codes catalog с table-driven tests**

## What Happened

Создан api/handlers/errors.go (5018 bytes) с ErrorResponse/ErrorDetail types, 17 констант для canonical error codes (16 из spec + dimensions_mismatch + encoding_format_invalid), errorCodeRegistry map (code → type+httpStatus), WriteError helper с unknown-code-fails-closed, WriteErrorWithRetryAfter для 429/503. Создан api/handlers/errors_test.go (5670 bytes) с 21 test case: TestAllErrorCodesRegistered, TestErrorEnvelopeShape (17 sub-tests на каждый code, проверяет HTTP status, code, type, param, message), TestWriteErrorUnknownCodeFailsClosedAsInternal, TestWriteErrorWithRetryAfter, TestWriteErrorAbortsContext, TestHTTPStatusForUnknownReturns500. Conventions: gin.SetMode(TestMode), gin.CreateTestContext, httptest, table-driven subtests, goconst-friendly константы для каждого code value.

## Verification

go test ./api/handlers/... -run "TestAllErrorCodesRegistered|TestErrorEnvelopeShape|TestWriteError" -v: 21/21 PASS. HTTPStatusFor unknown returns 500. AllErrorCodes matches registry. WriteError aborts context. WriteErrorWithRetryAfter sets Retry-After header.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| — | No verification commands discovered | — | — | — |

## Deviations

Добавлен 17-й error code encoding_format_invalid (не в spec Section 3) — нужен потому что encoding_format переносится в /v1/embeddings в S01 T04 (по user request), и нужен machine-readable code для невалидного значения. Также добавлен dimensions_mismatch (был в spec) который отсутствовал в первоначальном списке 16. Итого 17 codes, spec Section 3 имеет 16 + 1 наш = 17 в registry.

## Known Issues

None.

## Files Created/Modified

- `api/handlers/errors.go`
- `api/handlers/errors_test.go`
