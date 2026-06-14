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

- [ ] **S02: Lifecycle warmup readiness and graceful shutdown** `risk:medium` `depends:[]`
  > After this: After this, fd pre-warms model при старте, /live отвечает 200 сразу, /ready отвечает 200 только после warmup, deep /health отвечает 503 status=down до warmup и 200 status=ok после, SIGTERM приводит к 503 shutting_down для новых запросов и корректному drain in-flight за ≤30s.

- [ ] **S03: Observability surface endpoints headers and deep health** `risk:low` `depends:[S01,S02]`
  > After this: After this, /version возвращает semver+model+build_hash+uptime, /info возвращает список моделей с dims/limits/device/loaded/warmup, /metrics возвращает Prometheus text format с requests_total, request_duration_seconds histogram, batch_size, cache_hits (после S04), errors_total, model_loaded gauge, /v1/healthcheck работает как alias, и каждый response несёт Server, X-Request-Id, X-Model-Id, X-Dimensions, Connection: keep-alive, Retry-After на 429/503.

- [ ] **S04: Performance baseline and LRU cache** `risk:medium` `depends:[S02,S03]`
  > After this: After this, fd достигает perf baseline на warm model (1<50ms, 10<200ms, 32<1000ms p95, 100 sequential zero errors, 4×8 concurrent < 2s), in-memory LRU cache 10000 entries/24h TTL сокращает повторные inference до < 5ms, и cache status репортится через X-Cache header и fd_cache_hits_total metric.

- [ ] **S05: OpenAI v2 compat features OpenAPI schema and P2 enhancements** `risk:low` `depends:[S01]`
  > After this: After this, /v1/embeddings принимает encoding_format=base64, user, priority. FD_API_KEY env включает bearer auth с 401 unauthorized. CORS headers на responses. /openapi.json возвращает валидный OpenAPI 3.1 spec, /docs рендерит Swagger UI. /v1/batch принимает batches:[[..]] и возвращает batches:[[..]]. Rate limiting (если включен) даёт 429 с X-RateLimit-*. ETag на /v1/embeddings responses. /v1/traces возвращает последние N requests.

## Boundary Map

Not provided.
