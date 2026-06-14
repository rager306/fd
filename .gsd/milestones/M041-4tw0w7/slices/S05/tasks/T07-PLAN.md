---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T07: OpenAPI 3.1 schema generation и Swagger UI

api/openapi/spec.go: генерирует OpenAPI 3.1 spec программно (на основе реальных routes и типов). Включает все endpoints: /health, /live, /ready, /warmup, /version, /info, /metrics, /v1/embeddings, /v1/batch, /v1/healthcheck, /v1/traces, /openapi.json, /docs. Все request/response schemas, headers, error envelope. GET /openapi.json → 200 application/json. GET /docs → 200 text/html с Swagger UI (swagger-ui-dist или CDN). Валидация через openapi-spec-validator.

## Inputs

- None specified.

## Expected Output

- `api/openapi/spec.go`
- `api/openapi/spec_test.go`
- `api/handlers/openapi.go`
- `api/handlers/docs.go`

## Verification

Unit tests: /openapi.json возвращает валидный OpenAPI 3.1 (validate через openapi-spec-validator). /docs возвращает HTML с swagger-ui в body. Integration test: curl /openapi.json | openapi-spec-validator — exit 0.
