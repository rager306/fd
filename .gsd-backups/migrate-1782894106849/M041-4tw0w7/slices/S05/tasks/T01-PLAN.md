---
estimated_steps: 1
estimated_files: 5
skills_used: []
---

# T01: Added OpenAI-compatible `user` and `priority` request fields, priority validation, and tests for base64/user/priority behavior.

api/handlers/embeddings.go: расширить request struct: EncodingFormat *string (valid: float|base64), User *string, Priority *string (valid: low|normal|high). Validation в S01 middleware: невалидный encoding_format → 400. Base64 encoding для response: при encoding_format=base64, кодировать []float32 в base64-encoded float32 LE array.

## Inputs

- None specified.

## Expected Output

- `api/handlers/embeddings.go`
- `api/handlers/batch.go`
- `api/embed/base64.go`
- `api/middleware/validation.go`
- `api/handlers/embeddings_test.go`

## Verification

Unit tests: T-H-5 (encoding_format=base64 → base64 string в response), T-H-6 (priority=high принимается), user field принимается. Невалидный encoding_format → 400.
