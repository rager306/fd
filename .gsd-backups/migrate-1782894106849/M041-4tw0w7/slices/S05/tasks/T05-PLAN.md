---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T05: Added ETag and Cache-Control middleware for `/v1/embeddings` and `/info` with If-None-Match 304 support.

api/middleware/cache_headers.go: на /v1/embeddings и /info responses вычислять ETag = SHA256(response body) и выставлять Cache-Control: public, max-age=86400. Поддержка If-None-Match: если request header matches ETag → 304 Not Modified без body.

## Inputs

- None specified.

## Expected Output

- `api/middleware/cache_headers.go`
- `api/middleware/cache_headers_test.go`

## Verification

Unit tests: первый request → ETag: <hash>, Cache-Control: public, max-age=86400. Повторный с If-None-Match: <hash> → 304.
