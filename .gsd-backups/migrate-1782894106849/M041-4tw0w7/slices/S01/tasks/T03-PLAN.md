---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Validation middleware (повторное complete после replan)

api/middleware/validation.go: gin middleware который читает body через http.MaxBytesReader(w, body, 10*1024*1024) для size limit, парсит JSON, валидирует input array (non-empty, len<=32, все strings, каждый string <=2048 chars), валидирует dimensions (если указан то 512 или 1024). На любой failure — abort с правильным error envelope из T02 через WriteError. Validate BEFORE model call. Использовать struct tags binding:"required,max=32,dive,max=2048" где возможно.

## Inputs

- None specified.

## Expected Output

- `api/middleware/validation.go`
- `api/middleware/validation_test.go`

## Verification

Unit tests: T-E-1 ({} missing input → 400 input_required), T-E-2 (input:[] → 400 input_required), T-E-3 (dimensions:99999 → 400 dimensions_invalid), T-E-4 (input:[123] → 400 invalid_request_error НЕ 'json: cannot unmarshal'), T-E-5 (malformed JSON → 400 invalid_json), T-E-6 (100 inputs → 413 batch_too_large), T-E-7 (10000 char string → 413 input_too_long), T-E-14 (50MB body → 413 payload_too_large). Все 8 test cases pass.
