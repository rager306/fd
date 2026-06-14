---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T05: Added response headers middleware for request IDs, server/version, connection, and embedding model/dimensions headers.

api/middleware/headers.go: gin middleware. Server: fd/<version> (из buildinfo). X-Request-Id: echo caller-passed X-Request-Id header (любой case), иначе generate UUIDv4. На error pathе — обязательно сохранить X-Request-Id (recovery middleware из S01 должен его прочитать). Connection: keep-alive (default для HTTP/1.1, но explicit). X-Model-Id и X-Dimensions — выставляются на /v1/embeddings responses (model=<deepvk/USER-bge-m3>, dimensions=actual). Retry-After на 429/503 responses.

## Inputs

- None specified.

## Expected Output

- `api/middleware/headers.go`
- `api/middleware/headers_test.go`

## Verification

Unit tests: T-HDR-1 (Server: fd/2.0.0), T-HDR-2/3 (X-Request-Id echo vs generated), T-HDR-4/5 (X-Model-Id + X-Dimensions на /v1/embeddings), T-HDR-8 (Retry-After на 503), T-HDR-9 (Connection: keep-alive).
