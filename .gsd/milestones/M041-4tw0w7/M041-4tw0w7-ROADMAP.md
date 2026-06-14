# M041-4tw0w7: fd v2 service hardening and observability

**Vision:** Реализовать fd v2: production-ready локальный HTTP embedding service с предсказуемыми ошибками, deep health/observability surface, корректным lifecycle, и performance baseline. v1 caller (только POST /v1/embeddings) продолжает работать. Все 45 acceptance test cases из /root/fd-v2.md Section 5 проходят.

## Success Criteria

- Все 12 probe bugs из /root/fd-v2.md Section 1.4 (B1..B12) устранены, кроме B1/B2/B3/B5/B6 которые уже валидны (нужно сохранить поведение).
- Все 6 P0 endpoints реализованы и возвращают 200: /version, /info, /metrics, /v1/healthcheck, /live, /ready.
- Все P0 response headers (Server, X-Request-Id, X-Model-Id, X-Dimensions, X-Cache, Retry-After, Connection) присутствуют и покрыты тестами.
- Все 16 error codes из Section 3 catalog реализованы с правильным HTTP статусом, code, type.
- Performance baseline (T-P-1..T-P-5) держится на warm model.
- 10 behavior scenarios (Section 6.3 F-1..F-10) воспроизводятся тестами.
- Backward compatibility: v1 caller (только POST /v1/embeddings с OpenAI shape) работает.
- /openapi.json валидный OpenAPI 3.1 spec, /docs рендерит Swagger UI.
- Final artifact: `benchmark-results/fd-v2-validation-m041.md` со всеми 45 test results и perf numbers.

## Slices

- [x] **S01: Validation and OpenAI style error envelope** `risk:low` `depends:[]`
  > After this: After this, /v1/embeddings возвращает OpenAI style error envelope с machine readable code/type, oversized batch и oversized input дают 413 (не 500), malformed JSON даёт clean 400 invalid_json (не сырую Gin ошибку), и все 16 кодов из Section 3 catalog работают.

- [x] **S02: Lifecycle warmup readiness and graceful shutdown** `risk:medium` `depends:[]`
  > After this: After this, fd pre-warms model at startup, /live is cheap, /ready transitions 503 to 200 after warmup, shutdown gates new requests with 503+Retry-After, and the slice passes M043 gates: go test ./..., golangci-lint 18 linters, no reachable govulncheck findings.

- [x] **S03: Observability surface endpoints headers and deep health** `risk:low` `depends:[S01,S02]`
  > After this: After this, /version, /info, /metrics, /v1/healthcheck, deep /health, /warmup, and response headers are implemented and tested; new exported observability APIs have godoc and pass M043 lint/test/govulncheck gates.

- [x] **S04: Performance baseline and LRU cache** `risk:medium` `depends:[S02,S03]`
  > After this: After this, cache/perf validation includes warm/cold baseline plus M043 gates; cache code must keep context propagation, gocyclo <=15 for production functions, and no new static-analysis suppressions without justification.

- [ ] **S05: OpenAI v2 compat features OpenAPI schema and P2 enhancements** `risk:low` `depends:[S01]`
  > After this: After this, remaining OpenAI v2 compat work is complete without reimplementing encoding_format already covered by S01; user/priority/auth/CORS/OpenAPI/rate-limit/traces work passes M043 gates and documents any new API godoc.

## Boundary Map

Not provided.
