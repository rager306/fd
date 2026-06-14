# Requirements

This file is the explicit capability and coverage contract for the project.

## Active

### R010 — Все входящие /v1/embeddings запросы должны пройти предварительную валидацию (input length, batch size, dimensions, JSON shape) и при ошибке возвращать OpenAI-style error envelope с машинно-читаемым code/type и корректным HTTP-статусом (400/413/503), а не сырое сообщение Gin/Go. Реализует R-P0-1, R-P0-2, R-P0-18, R-P0-19.
- Class: core-capability
- Status: active
- Description: Все входящие /v1/embeddings запросы должны пройти предварительную валидацию (input length, batch size, dimensions, JSON shape) и при ошибке возвращать OpenAI-style error envelope с машинно-читаемым code/type и корректным HTTP-статусом (400/413/503), а не сырое сообщение Gin/Go. Реализует R-P0-1, R-P0-2, R-P0-18, R-P0-19.
- Why it matters: Текущая реализация отдает сырые ошибки Go-парсера и 500 на oversized batch; caller не может отличить caller-bug (400/413) от server-bug (500) и retry-логика ломается.
- Source: /root/fd-v2.md Section 2.4 + Section 3 error catalog
- Primary owning slice: M041-4tw0w7/S01
- Validation: 45 test cases Section 5 (T-E-1..T-E-15): все 400/413/405/500 ошибки возвращают правильный code/type, batch_too_large и input_too_long НЕ возвращают 500, валидация происходит ДО model inference.

### R012 — На warm service с предварительно прогретым cache для измеряемых payload: 1 input p95 < 50ms, 10 inputs p95 < 200ms, 32 inputs (max batch) p95 < 1000ms, 100 sequential cache-hot requests без ошибок, 4 concurrent callers × 8 cache-hot inputs < 2s total. Реализует R-P0-6 после D045 rescope.
- Class: quality-attribute
- Status: active
- Description: На warm service с предварительно прогретым cache для измеряемых payload: 1 input p95 < 50ms, 10 inputs p95 < 200ms, 32 inputs (max batch) p95 < 1000ms, 100 sequential cache-hot requests без ошибок, 4 concurrent callers × 8 cache-hot inputs < 2s total. Реализует R-P0-6 после D045 rescope.
- Why it matters: Daily-archive pipeline требует ≥7 papers/min end-to-end; B8 (10 inputs timeout 10s) и B9 (100 inputs 500) — blocking performance баги.
- Source: /root/fd-v2.md Section 1.4 B4/B8/B9 + Section 2.1 R-P0-6 + Section 5.4 T-P-1..T-P-5
- Primary owning slice: M041-4tw0w7/S04
- Supporting slices: M041-4tw0w7/S02
- Validation: D045 фиксирует cache-hot трактовку T-P-1..T-P-5. `tools/verify_fd_v2_perf.sh` prewarm-ит measured payload через real inference и затем требует `X-Cache: HIT` для latency cases. Evidence: `benchmark-results/fd-v2-perf-validation-m041-s04.md` PASS (batch=1 p95 2.236ms, batch=10 p95 3.468ms, batch=32 p95 7.595ms, sequential/concurrent/cache HIT pass) плюс non-blocking cache-miss diagnostics.
- Notes: Real cache-miss TEI CPU latency is intentionally diagnostic only for M041 S04; backend remediation was explicitly descoped by the user.

### R021 — fd handler отправляет chunked TEI calls в ПАРАЛЛЕЛЬ (bounded concurrency 4, matches TEI max_batch_requests=4) вместо sequential. Cold path for batch=128 должен упасть с 25s до ≤10s; batch=32 cold с 6s до ≤4s. Env FD_ASYNC_CHUNKS=true включает async mode (default off для backward compat). Каждый chunk error агрегируется, partial response не отдаётся.
- Class: quality-attribute
- Status: active
- Description: fd handler отправляет chunked TEI calls в ПАРАЛЛЕЛЬ (bounded concurrency 4, matches TEI max_batch_requests=4) вместо sequential. Cold path for batch=128 должен упасть с 25s до ≤10s; batch=32 cold с 6s до ≤4s. Env FD_ASYNC_CHUNKS=true включает async mode (default off для backward compat). Каждый chunk error агрегируется, partial response не отдаётся.
- Why it matters: TEI queue_time=2.7s создаёт sequential bottleneck: каждый chunk of 32 ждёт ~6s. Async pipeline в fd позволит параллельно слать несколько chunks — TEI может обрабатывать max_batch_requests=4 sub-batches параллельно. Reduction с 25s до 10s = 60% improvement для batch=128, без изменения TEI config.
- Source: M041-4tw0w7 S04 perf measurement; M042 CONTEXT decision on bounded concurrency
- Primary owning slice: M042-fjf2en/S02
- Validation: tools/verify_fd_async_perf.sh: FD_ASYNC_CHUNKS=true vs false perf comparison. Cold path batch=128 ≤10s (was 25s sequential). Cold path batch=32 ≤4s (was 6s sequential). Cache hit path не regressed (≤5ms per request). Benchmark artifact в benchmark-results/fd-v2-async-perf-m042.md.

### R026 — Upgrade fd `/openapi.json` and `/docs` contract from OpenAPI 3.1.0 to OAS 3.2.0, including verifier and validation evidence.
- Class: integration
- Status: active
- Description: Upgrade fd `/openapi.json` and `/docs` contract from OpenAPI 3.1.0 to OAS 3.2.0, including verifier and validation evidence.
- Why it matters: OAS 3.2.0 is the newer official OpenAPI specification version; adopting it keeps fd's schema surface current while preserving M041's already-validated 3.1.0 baseline as historical evidence.
- Source: User follow-up after reviewing https://spec.openapis.org/oas/v3.2.0.html#openapi-specification
- Validation: `GET /openapi.json` returns an OAS 3.2.0 document; docs render it; the final contract verifier asserts `openapi == "3.2.0"`; external schema validation or compatibility checks pass; mandatory Go gates (`go test ./...`, golangci-lint v2.12.2, govulncheck) pass.
- Notes: Implement as a new follow-up milestone/slice, not by editing M041 closure claims.

### R027 — Current fd product/runtime scope is TEI-first: ONNX runtime branch must be disabled or removed from active build, CI, docs, and runtime selection paths; ONNX may remain only as future research history/artifacts.
- Class: constraint
- Status: active
- Description: Current fd product/runtime scope is TEI-first: ONNX runtime branch must be disabled or removed from active build, CI, docs, and runtime selection paths; ONNX may remain only as future research history/artifacts.
- Why it matters: ONNX has not passed operational readiness and adds build/runtime/artifact complexity. The project should advance the working TEI path without ONNX code paths, binary/dependency noise, or confusing operator choices.
- Source: User directive during M042 TEI perf investigation
- Primary owning slice: M042-fjf2en/S02
- Validation: Default build, Docker image, docs, and runtime config expose only TEI as current backend; ONNX build/runtime selectors are absent or fail closed as explicitly research-only; `go test ./...`, golangci-lint v2.12.2, and govulncheck pass.
- Notes: This supersedes M042/S03 ONNX implementation scope; R022 is deferred.

## Validated

### R001 — Embedding runtime optimizations must preserve Russian-language and legal-domain retrieval/embedding quality for the current model; any model replacement requires benchmark evidence on a Russian legal corpus.
- Class: quality-attribute
- Status: validated
- Description: Embedding runtime optimizations must preserve Russian-language and legal-domain retrieval/embedding quality for the current model; any model replacement requires benchmark evidence on a Russian legal corpus.
- Why it matters: Latency gains are not useful if Russian legal-domain semantic quality regresses.
- Source: user clarification during M008 optimization research
- Primary owning slice: M040-pbp9z1/S02
- Supporting slices: M040-pbp9z1/S03, M040-pbp9z1/S04
- Validation: M040-pbp9z1 validation covered Russian/legal quality preservation: S02 legal retrieval guard passed for packaged ONNX with deepvk/USER-bge-m3, and S03/S04 reject unproven alternative model replacement fail-closed without legal/runtime evidence.
- Notes: Validated by M040 closeout. deepvk/USER-bge-m3 remains the model baseline; any future runtime or candidate replacement still requires fresh Russian/legal-domain evidence.

### R002 — Research/chunking workflows must use a sufficiently long-lived embedding cache so repeated chunk processing can reuse vectors and reduce model load.
- Class: quality-attribute
- Status: validated
- Description: Research/chunking workflows must use a sufficiently long-lived embedding cache so repeated chunk processing can reuse vectors and reduce model load.
- Why it matters: During research, chunks and vectors may be reused several times; short cache retention increases model load and slows experimentation.
- Source: user clarification during M008 Redis optimization research
- Primary owning slice: M040-pbp9z1/S02
- Supporting slices: M040-pbp9z1/S04
- Validation: M040-pbp9z1 S02 proved Redis L2 reuse across an API-only packaged Docker restart using isolated cache namespace m040-s02-onnx-restart; S04 requires cache namespace isolation in the operating contract.
- Notes: Validated for the same-host service readiness boundary. Future larger-corpus retention policy tuning remains an operational sizing concern, not an unvalidated core capability.

### R003 — Performance/cache/runtime tuning parameters should be configurable through environment variables with safe defaults and validation.
- Class: operability
- Status: validated
- Description: Performance/cache/runtime tuning parameters should be configurable through environment variables with safe defaults and validation.
- Why it matters: Research and VPS deployment need fine tuning without rebuilding code or editing source files.
- Source: user clarification during M008 Redis/cache architecture research
- Primary owning slice: M040-pbp9z1/S01
- Supporting slices: M040-pbp9z1/S02, M040-pbp9z1/S04
- Validation: M040-pbp9z1 S01 documented runtime/environment expectations and cache namespace guidance in the same-host service contract; S02 artifacts recorded sanitized runtime/cache configuration; S04 documents explicit ONNX opt-in, artifact/tokenizer/runtime preflight, optional ONNX_RUNTIME_SHA256, and EMBEDDING_CACHE_VERSION.
- Notes: Validated for current same-host TEI/default and explicit ONNX opt-in operation. Additional tuning knobs can be added later as new requirements if needed.

### R004 — Benchmark artifacts must record the effective environment/configuration parameters used for the run so results remain comparable across tuning experiments.
- Class: operability
- Status: validated
- Description: Benchmark artifacts must record the effective environment/configuration parameters used for the run so results remain comparable across tuning experiments.
- Why it matters: Performance results are not comparable if env tuning differs invisibly between runs.
- Source: user clarification during M008 benchmark/config research
- Primary owning slice: M040-pbp9z1/S02
- Supporting slices: M040-pbp9z1/S04
- Validation: M040-pbp9z1 S02 benchmark/preflight artifacts recorded sanitized effective configuration, ONNX backend, model, dimensions, namespace, restart evidence, and legal/audit PASS; S04 verifier checked evidence artifacts and redaction boundaries.
- Notes: Validated for M040 benchmark/recommendation artifacts. Continue excluding raw legal/probe text, secrets, and signed URLs from artifacts.

### R005 — fd must provide a same-host local HTTP embedding service contract for neighboring services, centered on `/v1/embeddings`, batch embeddings, and `/health`.
- Class: core-capability
- Status: validated
- Description: fd must provide a same-host local HTTP embedding service contract for neighboring services, centered on `/v1/embeddings`, batch embeddings, and `/health`.
- Why it matters: Neighboring services need a clear, stable local integration surface rather than an open-ended runtime experiment.
- Source: user
- Primary owning slice: M040-pbp9z1/S01
- Supporting slices: M040-pbp9z1/S04
- Validation: M040-pbp9z1 S01 created docs/same-host-embedding-service-contract.md covering /health, /v1/embeddings, /embeddings/batch, dimensions, request/response shapes, status/error behavior, timeout/retry guidance, runtime/env expectations, cache guidance, and non-goals; S04 links final recommendation from the contract.
- Notes: Validated by the same-host service contract. The scope remains local HTTP service integration, not embedded/library API or hosted CI proof.

### R006 — The TEI-vs-ONNX runtime recommendation must be based on an evidence envelope covering legal quality, same-host performance, restart/cache behavior, health/preflight clarity, and operational simplicity.
- Class: quality-attribute
- Status: validated
- Description: The TEI-vs-ONNX runtime recommendation must be based on an evidence envelope covering legal quality, same-host performance, restart/cache behavior, health/preflight clarity, and operational simplicity.
- Why it matters: The project goal is the best local embedding service for quality and speed, not ONNX experimentation for its own sake.
- Source: user
- Primary owning slice: M040-pbp9z1/S04
- Supporting slices: M040-pbp9z1/S01, M040-pbp9z1/S02, M040-pbp9z1/S03
- Validation: M040 S04 final artifact `benchmark-results/fd-runtime-recommendation-m040-s04.md` passed `tools/verify_m040_s04_recommendation.py` against S02/S03 evidence inputs in gsd_exec c52073f9-7ea0-4b13-9efa-99d54193c6f0, proving the final TEI-vs-ONNX evidence envelope and recommendation caveats are machine-checkable.
- Notes: Validated by S04 closeout verification; recommendation keeps TEI default until explicit ONNX switch and defers alternative model replacement fail-closed.

### R007 — M040 must not treat hosted GitHub Actions proof, remote workflow dispatch, push, upload, or artifact mirroring as required readiness gates.
- Class: constraint
- Status: validated
- Description: M040 must not treat hosted GitHub Actions proof, remote workflow dispatch, push, upload, or artifact mirroring as required readiness gates.
- Why it matters: The target deployment is same-host local service readiness, so hosted CI proof is outside the relevant acceptance boundary.
- Source: user
- Primary owning slice: M040-pbp9z1/S04
- Supporting slices: none
- Validation: M040 S04 verifier and final artifact passed in gsd_exec c52073f9-7ea0-4b13-9efa-99d54193c6f0, including hosted/remote CI readiness-gate rejection semantics and final artifact language that keeps hosted CI proof out of the same-host readiness gate.
- Notes: Validated by S04 closeout verification; readiness depends on same-host contract, preflight, cache namespace isolation, and smoke `POST /v1/embeddings`, not hosted CI.

### R008 — Alternative embedding model checks must be bounded to 1-2 plausible candidates and must use legal-domain evidence before any model can challenge `deepvk/USER-bge-m3`.
- Class: constraint
- Status: validated
- Description: Alternative embedding model checks must be bounded to 1-2 plausible candidates and must use legal-domain evidence before any model can challenge `deepvk/USER-bge-m3`.
- Why it matters: The user wants excellent legal-domain quality and speed without open-ended model experimentation.
- Source: user
- Primary owning slice: M040-pbp9z1/S03
- Supporting slices: M040-pbp9z1/S04
- Validation: M040-pbp9z1/S03 produced `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`, validated by `tools/verify_legal_model_quick_gate_artifact.py --max-candidates 2` plus closeout schema checks. The artifact caps candidates to BAAI/bge-m3 and intfloat/multilingual-e5-large, records sanitized legal-corpus hashes/counts, rejects cross-model cosine parity as a replacement criterion, and defers candidate replacement fail-closed because baseline `/health` lacks runtime metadata.
- Notes: S03 validated the bounded alternative-model gate contract. Candidate replacement remains deferred until a baseline and candidate endpoints expose contract-required runtime metadata and live legal retrieval metrics can run.

### R009 — The local embedding service must avoid silent per-request fallback between TEI and ONNX runtimes or between different tokenizers/models within one service run.
- Class: operability
- Status: validated
- Description: The local embedding service must avoid silent per-request fallback between TEI and ONNX runtimes or between different tokenizers/models within one service run.
- Why it matters: Neighboring services must know which embedding contract is serving requests to avoid mixed-vector correctness issues.
- Source: inferred
- Primary owning slice: M040-pbp9z1/S01
- Supporting slices: M040-pbp9z1/S04
- Validation: M040-pbp9z1 S01 contract and health metadata establish runtime identity and no-silent-fallback rules; /v1/embeddings request model is compatibility metadata, not a selector. S04 final stance and verifier require no request-level fallback and a smoke embedding readiness check beyond /health.
- Notes: Validated for current service semantics: runtime fallback or model switching must be an operator-level restart/reconfiguration path, not hidden per-request behavior.

### R011 — Сервис должен поддерживать корректный lifecycle: pre-warm model при старте (1 dummy inference), отдавать /live (cheap) и /ready (200 только после warmup), маппить model-not-loaded/overloaded на 503+Retry-After, и завершаться по SIGTERM за ≤30s с in-flight drain. Реализует R-P0-3, R-P0-4, R-P0-5.
- Class: operability
- Status: validated
- Description: Сервис должен поддерживать корректный lifecycle: pre-warm model при старте (1 dummy inference), отдавать /live (cheap) и /ready (200 только после warmup), маппить model-not-loaded/overloaded на 503+Retry-After, и завершаться по SIGTERM за ≤30s с in-flight drain. Реализует R-P0-3, R-P0-4, R-P0-5.
- Why it matters: Сейчас при cold start caller получает 500 silent timeout; SIGTERM рвёт in-flight запросы; нет k8s probe surface.
- Source: /root/fd-v2.md Section 2.1 P0 lifecycle + Section 6.1 startup sequence + Section 6.3 F-1/F-5
- Primary owning slice: M041-4tw0w7/S02
- Validation: Validated by M041-4tw0w7/S02. Evidence: `benchmark-results/m041-s02-t06-lifecycle-integration.txt` covers startup readiness, F-1 model_not_loaded, F-2 model_overloaded via FD_MAX_IN_FLIGHT, and F-5 shutdown drain; `benchmark-results/m041-s02-t06-go-test.txt`, `m041-s02-t06-lint.txt`, and `m041-s02-t06-govulncheck.txt` pass the mandatory Go gates.
- Notes: S02 implements lifecycle pre-warm, /live, /ready, /v1/embeddings lifecycle gate, graceful SIGTERM/SIGINT drain, and default-off capacity overload control via FD_MAX_IN_FLIGHT.

### R013 — Должны быть доступны: GET /version (semver+model+build_hash+uptime), GET /info или /v1/models (список моделей с dims/limits/device/loaded/warmup), GET /metrics (Prometheus text: requests_total, request_duration_seconds histogram, batch_size histogram, cache_hits_total, errors_total, model_loaded gauge), GET /v1/healthcheck (alias). Реализует R-P0-7, R-P0-8, R-P0-9, R-P0-10.
- Class: failure-visibility
- Status: validated
- Description: Должны быть доступны: GET /version (semver+model+build_hash+uptime), GET /info или /v1/models (список моделей с dims/limits/device/loaded/warmup), GET /metrics (Prometheus text: requests_total, request_duration_seconds histogram, batch_size histogram, cache_hits_total, errors_total, model_loaded gauge), GET /v1/healthcheck (alias). Реализует R-P0-7, R-P0-8, R-P0-9, R-P0-10.
- Why it matters: Сейчас нет machine-readable way узнать версию/модель/limits/метрики — caller и оператор работают вслепую.
- Source: /root/fd-v2.md Section 1.3 missing endpoints + Section 2.2 P0 observability + Section 4 OpenAPI spec
- Primary owning slice: M041-4tw0w7/S03
- Validation: Validated by M041-4tw0w7/S03. Evidence: `benchmark-results/m041-s03-t07-observability-integration.txt` covers `/version`, `/info`, `/metrics`, and `/v1/healthcheck`; `m041-s03-t07-go-test.txt`, `m041-s03-t07-lint.txt`, and `m041-s03-t07-govulncheck.txt` pass mandatory Go gates.
- Notes: Build/version metadata, model info, Prometheus metrics, and healthcheck alias are implemented under the executable `api/` module layout.

### R014 — Каждый response должен нести: Server: fd/<version>, X-Request-Id (echo caller-passed или generated UUIDv4), X-Model-Id (на /v1/embeddings), X-Dimensions (на /v1/embeddings), X-Cache: HIT|MISS (если cache включен), Retry-After (на 429/503), Connection: keep-alive. Реализует R-P0-11..R-P0-17.
- Class: operability
- Status: validated
- Description: Каждый response должен нести: Server: fd/<version>, X-Request-Id (echo caller-passed или generated UUIDv4), X-Model-Id (на /v1/embeddings), X-Dimensions (на /v1/embeddings), X-Cache: HIT|MISS (если cache включен), Retry-After (на 429/503), Connection: keep-alive. Реализует R-P0-11..R-P0-17.
- Why it matters: B11/B12 показывают пустые headers — нет request correlation, нет cache observability, нет server identity.
- Source: /root/fd-v2.md Section 1.4 B11/B12 + Section 2.3 P0 headers + Section 5.3 T-HDR-1..T-HDR-10
- Primary owning slice: M041-4tw0w7/S03
- Validation: S03 validated Server, X-Request-Id, X-Model-Id, X-Dimensions, Retry-After, and Connection headers. S04 validated `X-Cache: MISS` on first request and `X-Cache: HIT` on repeated input via `benchmark-results/fd-v2-perf-validation-m041-s04.md` and `api/fd_v2_cache_integration_test.go`; final verifier requires non-HIT count 0 for cache-hot cases.
- Notes: Fully validated after S04 added X-Cache integration. Real cache-miss latency remains diagnostic per D045, but header presence/semantics are validated.

### R015 — GET /health — deep check (model_loaded, warmup_done, device, last_inference_at, in_flight_requests, status=ok|degraded|down, 503 если degraded/down). GET /warmup — status/progress. POST /warmup — trigger on-demand warmup. Реализует R-P1-1, R-P1-2, R-P1-3.
- Class: failure-visibility
- Status: validated
- Description: GET /health — deep check (model_loaded, warmup_done, device, last_inference_at, in_flight_requests, status=ok|degraded|down, 503 если degraded/down). GET /warmup — status/progress. POST /warmup — trigger on-demand warmup. Реализует R-P1-1, R-P1-2, R-P1-3.
- Why it matters: Текущий /health — shallow (только timestamp); не различает "процесс жив" vs "model loaded vs not"; нет способа дождаться или форсировать warmup.
- Source: /root/fd-v2.md Section 1.1 /health shallow + Section 2.5 P1 health checks + Section 5.1 T-H-7
- Primary owning slice: M041-4tw0w7/S03
- Validation: Validated by M041-4tw0w7/S03. Evidence: `benchmark-results/m041-s03-t03-deep-health.txt` covers deep `/health` status fields and `last_inference_at`; `m041-s03-t06-warmup.txt` covers GET/POST `/warmup`; `m041-s03-t07-observability-integration.txt` covers integration behavior.
- Notes: Deep health reports ok/degraded/down, model_loaded, warmup_done, device, last_inference_at, and in_flight_requests. Warmup endpoints expose status/progress and trigger background pre-warm.

### R016 — In-memory LRU cache на (input_text, dimensions) → embedding, size 10000, TTL 24h, настраивается через env. Cache HIT skip model inference, отдаёт < 5ms. Метрики: fd_cache_hits_total{result=hit|miss}. Реализует R-P1-4.
- Class: differentiator
- Status: validated
- Description: In-memory LRU cache на (input_text, dimensions) → embedding, size 10000, TTL 24h, настраивается через env. Cache HIT skip model inference, отдаёт < 5ms. Метрики: fd_cache_hits_total{result=hit|miss}. Реализует R-P1-4.
- Why it matters: Daily-archive обрабатывает 12k+ papers с overlapping chunks — без cache один и тот же текст гоняется через model повторно, что замедляет pipeline и нагружает GPU.
- Source: /root/fd-v2.md Section 2.6 P1 R-P1-4 + Section 6.3 F-4
- Primary owning slice: M041-4tw0w7/S04
- Validation: S04 implemented LRU cache with size/TTL env config, copy-on-read/write, eviction metrics, and EmbeddingCache adapter methods. Evidence: `api/cache/lru_test.go`, `api/fd_v2_cache_integration_test.go`, `benchmark-results/m041-s04-t03-cache-integration.txt`, and final `benchmark-results/fd-v2-perf-validation-m041-s04.md` showing repeated input `X-Cache: HIT` in 1.870ms and cache-hot latency targets passing.
- Notes: Validated for fd-controlled cache-hot behavior. Backend cache-miss TEI CPU latency is outside R016 after D045.

### R017 — Расширить /v1/embeddings request schema: encoding_format: float|base64 (~30% bandwidth savings), user field (для abuse tracking и per-user rate limits), priority: low|normal|high. Реализует R-P1-5, R-P1-6, R-P1-7.
- Class: differentiator
- Status: validated
- Description: Расширить /v1/embeddings request schema: encoding_format: float|base64 (~30% bandwidth savings), user field (для abuse tracking и per-user rate limits), priority: low|normal|high. Реализует R-P1-5, R-P1-6, R-P1-7.
- Why it matters: Daily-archive scripts используют urllib и не используют headers, но более толстые callers (например, web UI) выиграют от base64 и user-based rate limits.
- Source: /root/fd-v2.md Section 2.6 P1 R-P1-5/6/7 + Section 4 OpenAPI spec
- Primary owning slice: M041-4tw0w7/S05
- Validation: `benchmark-results/fd-v2-validation-m041.md` 45/45 PASS: T019 validates `encoding_format=base64`; T020 validates `priority=high`; T021 validates `user`; T033 validates invalid `encoding_format` returns `encoding_format_invalid`; T034 validates invalid priority returns `priority_invalid`. Unit evidence in `benchmark-results/m041-s05-t01-go-test.txt`.
- Notes: The `user` field is accepted and used by optional per-user rate limiting; it is not persisted.

### R018 — Если env FD_API_KEY задан, все endpoints кроме /live, /metrics, /docs требуют Authorization: Bearer <key>, иначе 401 unauthorized. CORS headers (Access-Control-Allow-Origin/Methods/Headers) для web clients. Реализует R-P1-8, R-P1-9.
- Class: compliance/security
- Status: validated
- Description: Если env FD_API_KEY задан, все endpoints кроме /live, /metrics, /docs требуют Authorization: Bearer <key>, иначе 401 unauthorized. CORS headers (Access-Control-Allow-Origin/Methods/Headers) для web clients. Реализует R-P1-8, R-P1-9.
- Why it matters: Локальный сервис, но если exposed через reverse proxy — нужна хоть какая-то auth; CORS нужен для будущих web UI.
- Source: /root/fd-v2.md Section 2.6 P1 R-P1-8/9
- Primary owning slice: M041-4tw0w7/S05
- Validation: Unit evidence in `benchmark-results/m041-s05-t02-go-test.txt` validates FD_API_KEY bearer auth: missing/wrong token -> 401 unauthorized, correct token -> 200, public endpoints skipped. `benchmark-results/fd-v2-validation-m041.md` validates CORS preflight (T027) and docs/openapi public surfaces.
- Notes: Auth remains opt-in via FD_API_KEY; CORS defaults to `*` unless FD_CORS_ORIGINS is set.

### R019 — GET /openapi.json (OpenAPI 3.1 spec), GET /docs (Swagger UI), POST /v1/batch (batches:[[..],[..]] → batches:[[..],[..]]), rate limiting (per-IP 100 req/min, per-user 1000 req/min, headers X-RateLimit-*), ETag+Cache-Control на responses, /v1/traces (recent N requests с latency/status), optional SSE streaming. Реализует R-P2-1..R-P2-6.
- Class: quality-attribute
- Status: validated
- Description: GET /openapi.json (OpenAPI 3.1 spec), GET /docs (Swagger UI), POST /v1/batch (batches:[[..],[..]] → batches:[[..],[..]]), rate limiting (per-IP 100 req/min, per-user 1000 req/min, headers X-RateLimit-*), ETag+Cache-Control на responses, /v1/traces (recent N requests с latency/status), optional SSE streaming. Реализует R-P2-1..R-P2-6.
- Why it matters: Caller и операторы хотят self-describing API (OpenAPI), явный batch endpoint, rate limiting для защиты, traces для debugging.
- Source: /root/fd-v2.md Section 2.7 P2 nice-to-have + Section 4 OpenAPI
- Primary owning slice: M041-4tw0w7/S05
- Validation: `benchmark-results/fd-v2-validation-m041.md` 45/45 PASS validates `/openapi.json` (T010/T043), `/docs` (T011), `/v1/batch` (T038-T040), ETag/Cache-Control and 304 (T024-T026), `/v1/traces` (T012/T042), and cache-hot performance checks (T044-T045). T03/T05/T06/T07 task evidence covers rate limiting, cache validators, traces, OpenAPI validator, and docs unit tests.
- Notes: SSE streaming remains optional and was not implemented in M041 S05; core R-P2 surfaces were validated.

### R020 — Письменный root cause analysis объясняет почему TEI queue_time=2.7s несмотря на max_concurrent_requests=512. Документ содержит hypothesis, evidence (TEI logs, /info metrics, профилирование), и рекомендации. Возможные выводы: (a) TEI single backend thread — fixed by source change (out of fd scope), (b) ONNX fallback, (c) async pipeline.
- Class: quality-attribute
- Status: validated
- Description: Письменный root cause analysis объясняет почему TEI queue_time=2.7s несмотря на max_concurrent_requests=512. Документ содержит hypothesis, evidence (TEI logs, /info metrics, профилирование), и рекомендации. Возможные выводы: (a) TEI single backend thread — fixed by source change (out of fd scope), (b) ONNX fallback, (c) async pipeline.
- Why it matters: M041 perf measurements показали TEI cold path 6s per chunk — bottleneck для fd v2 latency target. Без RCA дальнейшие fixes (async/ONNX) могут не решить root problem. Документ нужен чтобы: (1) объяснить stakeholders почему current perf так себе, (2) обосновать выбор mitigation strategy, (3) capture знание для future M0xx milestones.
- Source: M041-4tw0w7 perf measurement (2026-06-13 18:59) + M042 CONTEXT
- Primary owning slice: M042-fjf2en/S01
- Validation: `documents/te-perf-root-cause-m042.md` explains TEI queue/startup behavior with T01/T02 evidence: direct TEI batch=32 queue p50 ~2434.795ms despite `max_concurrent_requests=512`, and restart/recreate spent ~48 minutes from `Starting model backend` to `Ready` after delayed missing-ONNX ORT fallback. The RCA includes hypothesis tree, verdict, and TEI-first recommendation.
- Notes: RCA rejects ONNX as M042 mitigation and defers ONNX implementation per D047/R022.

### R023 — Расширить .golangci.yml Tier 1 linters: gosec, bodyclose, prealloc, errorlint, revive. Каждый новый линтер в warn mode в первом проходе, потом fail mode после cleanup. Fix найденных issues в существующем fd коде (M041 новый код, M042 S02 async) с явными justification comments. Интегрировать в CI go-quality.yml. Acceptance: golangci-lint run exit 0 на всём fd repo, нет issues от Tier 1.
- Class: quality-attribute
- Status: validated
- Description: Расширить .golangci.yml Tier 1 linters: gosec, bodyclose, prealloc, errorlint, revive. Каждый новый линтер в warn mode в первом проходе, потом fail mode после cleanup. Fix найденных issues в существующем fd коде (M041 новый код, M042 S02 async) с явными justification comments. Интегрировать в CI go-quality.yml. Acceptance: golangci-lint run exit 0 на всём fd repo, нет issues от Tier 1.
- Why it matters: fd имеет консервативный baseline из 7 linters (errcheck, govet, ineffassign, staticcheck, unused, goconst, misspell). Без gosec, bodyclose, prealloc, errorlint, revive: (a) security issues не ловятся (gosec G107/G110), (b) HTTP body leaks могут проскочить (bodyclose), (c) slice allocations regressions возможны (prealloc), (d) error wrapping regressions возможны (errorlint), (e) code documentation quality не enforced (revive). 2026 Go community consensus: golangci-lint + staticcheck + targeted security/optim linters — стандарт де-факто.
- Source: https://github.com/dpolivaev/static-analysis Go section + 2026 community consensus (Reddit r/golang, analysis-tools.dev); M041-4tw0w7 baseline (7 linters)
- Primary owning slice: M043-dpr0cq/S01
- Validation: M043 S01: Tier 1 linters enabled and fixed; final lint 0 issues. Evidence: docs/static-analysis-phase1-report-m043.md, benchmark-results/m043-tier1-baseline.txt.

### R024 — Добавить Tier 2 linters: gocyclo (cyclomatic complexity, custom threshold для fd), gocritic (selective enabled-tags: diagnostic, performance, style), durationcheck (time.Duration conversions), unparam (unused parameters), contextcheck (context propagation), nilnil (nil error returns). Каждый в warn mode первым проходом, потом fail mode. Fix issues в существующем коде (особенно gocyclo в CreateEmbedding handler в M041 S01, M042 S02 async orchestrator).
- Class: quality-attribute
- Status: validated
- Description: Добавить Tier 2 linters: gocyclo (cyclomatic complexity, custom threshold для fd), gocritic (selective enabled-tags: diagnostic, performance, style), durationcheck (time.Duration conversions), unparam (unused parameters), contextcheck (context propagation), nilnil (nil error returns). Каждый в warn mode первым проходом, потом fail mode. Fix issues в существующем коде (особенно gocyclo в CreateEmbedding handler в M041 S01, M042 S02 async orchestrator).
- Why it matters: Phase 1 закрывает critical gaps (security, body leaks, allocations, error wrapping, docs). Phase 2 покрывает medium-value checks: complexity (CreateEmbedding уже ~150 LOC с nested loops после M041 S04 chunking), quality (unused params, context), style (naming, structure). 2026 best practice: gocyclo с threshold=15-20 для service code, gocritic selective to avoid noise.
- Source: https://github.com/dpolivaev/static-analysis Go section + 2026 community consensus; fd complexity analysis (M041 S04 chunked handler ~150 LOC)
- Primary owning slice: M043-dpr0cq/S02
- Validation: M043 S02: Tier 2 linters enabled; 17 baseline issues fixed; final lint 0 issues. Evidence: docs/static-analysis-phase2-report-m043.md, benchmark-results/m043-s02-final-lint.txt.

### R025 — Standalone govulncheck (golang.org/x/vuln) integrated в CI go-quality.yml. Команда: go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./... Scan dependencies (api/go.mod) + stdlib usage в fd. Включить как required CI step (fail on known vulnerabilities). Также: docs/static-analysis-recommendation.md обновить с финальной M043 phased plan: что реализовано (Phase 1, 2, 3), что deferred (govulncheck always-on после CI integration, future tiers), Phase 3 opt-in linters (gofumpt, structslop, etc).
- Class: quality-attribute
- Status: validated
- Description: Standalone govulncheck (golang.org/x/vuln) integrated в CI go-quality.yml. Команда: go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./... Scan dependencies (api/go.mod) + stdlib usage в fd. Включить как required CI step (fail on known vulnerabilities). Также: docs/static-analysis-recommendation.md обновить с финальной M043 phased plan: что реализовано (Phase 1, 2, 3), что deferred (govulncheck always-on после CI integration, future tiers), Phase 3 opt-in linters (gofumpt, structslop, etc).
- Why it matters: govulncheck — официальный Go vuln scanner (golang.org/x/vuln). Catches known vulnerabilities в stdlib и зависимостях. golangci-lint НЕ покрывает эту функциональность (отдельный tool). 2026 best practice: govulncheck как required CI step. Phase 3 закрывает оставшиеся gaps из analysis (govulncheck CI integration + documentation finalization).
- Source: https://github.com/dpolivaev/static-analysis Go section (govulncheck); 2026 Go security best practice; M041 fd depends on gin, redis/go-redis, onnxruntime, etc.
- Primary owning slice: M043-dpr0cq/S03
- Validation: M043 S03: govulncheck CI step added and local govulncheck exits 0 with 0 reachable vulnerabilities; docs finalized. Evidence: benchmark-results/m043-s03-govulncheck-final.txt, docs/static-analysis-recommendation.md.

## Deferred

### R022 — Opt-in ONNX mode (FD_BACKEND=onnx, requires onnx build tag) — fd serves embeddings из Go ONNX runtime вместо TEI HTTP. Per M019: cold latency 8.3ms, warm 1.19ms, throughput 858 req/s. Опционально через env, default off (TEI остаётся production per R001/M015). fd binary должен билдиться с onnx tag без regression: все M041 acceptance criteria должны pass в обоих режимах.
- Class: quality-attribute
- Status: deferred
- Description: Opt-in ONNX mode (FD_BACKEND=onnx, requires onnx build tag) — fd serves embeddings из Go ONNX runtime вместо TEI HTTP. Per M019: cold latency 8.3ms, warm 1.19ms, throughput 858 req/s. Опционально через env, default off (TEI остаётся production per R001/M015). fd binary должен билдиться с onnx tag без regression: все M041 acceptance criteria должны pass в обоих режимах.
- Why it matters: M019 measurements показывают ONNX Go runtime на 100-700x быстрее TEI на warm path. Даже с ограничениями legal-quality gate (M015/M016), opt-in mode даёт операторам speed-first option для workloads где legal quality менее критична (например, прототипирование, non-legal embeddings). Production default остаётся TEI.
- Source: M019 ONNX 1024 perf benchmark; M015/M016 legal quality gate (stays as blocking concern for production); M042 CONTEXT decision on opt-in
- Primary owning slice: M042-fjf2en/S03
- Validation: Deferred by user decision during M042: ONNX will not be implemented as current opt-in runtime. Prior M019 evidence remains research-only, not production readiness proof.
- Notes: Replace M042/S03 ONNX implementation work with TEI-first stabilization and/or ONNX branch deactivation. Future ONNX work requires a separate research milestone and explicit packaging/readiness gates.

## Out of Scope

## Traceability

| ID | Class | Status | Primary owner | Supporting | Proof |
|---|---|---|---|---|---|
| R001 | quality-attribute | validated | M040-pbp9z1/S02 | M040-pbp9z1/S03, M040-pbp9z1/S04 | M040-pbp9z1 validation covered Russian/legal quality preservation: S02 legal retrieval guard passed for packaged ONNX with deepvk/USER-bge-m3, and S03/S04 reject unproven alternative model replacement fail-closed without legal/runtime evidence. |
| R002 | quality-attribute | validated | M040-pbp9z1/S02 | M040-pbp9z1/S04 | M040-pbp9z1 S02 proved Redis L2 reuse across an API-only packaged Docker restart using isolated cache namespace m040-s02-onnx-restart; S04 requires cache namespace isolation in the operating contract. |
| R003 | operability | validated | M040-pbp9z1/S01 | M040-pbp9z1/S02, M040-pbp9z1/S04 | M040-pbp9z1 S01 documented runtime/environment expectations and cache namespace guidance in the same-host service contract; S02 artifacts recorded sanitized runtime/cache configuration; S04 documents explicit ONNX opt-in, artifact/tokenizer/runtime preflight, optional ONNX_RUNTIME_SHA256, and EMBEDDING_CACHE_VERSION. |
| R004 | operability | validated | M040-pbp9z1/S02 | M040-pbp9z1/S04 | M040-pbp9z1 S02 benchmark/preflight artifacts recorded sanitized effective configuration, ONNX backend, model, dimensions, namespace, restart evidence, and legal/audit PASS; S04 verifier checked evidence artifacts and redaction boundaries. |
| R005 | core-capability | validated | M040-pbp9z1/S01 | M040-pbp9z1/S04 | M040-pbp9z1 S01 created docs/same-host-embedding-service-contract.md covering /health, /v1/embeddings, /embeddings/batch, dimensions, request/response shapes, status/error behavior, timeout/retry guidance, runtime/env expectations, cache guidance, and non-goals; S04 links final recommendation from the contract. |
| R006 | quality-attribute | validated | M040-pbp9z1/S04 | M040-pbp9z1/S01, M040-pbp9z1/S02, M040-pbp9z1/S03 | M040 S04 final artifact `benchmark-results/fd-runtime-recommendation-m040-s04.md` passed `tools/verify_m040_s04_recommendation.py` against S02/S03 evidence inputs in gsd_exec c52073f9-7ea0-4b13-9efa-99d54193c6f0, proving the final TEI-vs-ONNX evidence envelope and recommendation caveats are machine-checkable. |
| R007 | constraint | validated | M040-pbp9z1/S04 | none | M040 S04 verifier and final artifact passed in gsd_exec c52073f9-7ea0-4b13-9efa-99d54193c6f0, including hosted/remote CI readiness-gate rejection semantics and final artifact language that keeps hosted CI proof out of the same-host readiness gate. |
| R008 | constraint | validated | M040-pbp9z1/S03 | M040-pbp9z1/S04 | M040-pbp9z1/S03 produced `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`, validated by `tools/verify_legal_model_quick_gate_artifact.py --max-candidates 2` plus closeout schema checks. The artifact caps candidates to BAAI/bge-m3 and intfloat/multilingual-e5-large, records sanitized legal-corpus hashes/counts, rejects cross-model cosine parity as a replacement criterion, and defers candidate replacement fail-closed because baseline `/health` lacks runtime metadata. |
| R009 | operability | validated | M040-pbp9z1/S01 | M040-pbp9z1/S04 | M040-pbp9z1 S01 contract and health metadata establish runtime identity and no-silent-fallback rules; /v1/embeddings request model is compatibility metadata, not a selector. S04 final stance and verifier require no request-level fallback and a smoke embedding readiness check beyond /health. |
| R010 | core-capability | active | M041-4tw0w7/S01 | none | 45 test cases Section 5 (T-E-1..T-E-15): все 400/413/405/500 ошибки возвращают правильный code/type, batch_too_large и input_too_long НЕ возвращают 500, валидация происходит ДО model inference. |
| R011 | operability | validated | M041-4tw0w7/S02 | none | Validated by M041-4tw0w7/S02. Evidence: `benchmark-results/m041-s02-t06-lifecycle-integration.txt` covers startup readiness, F-1 model_not_loaded, F-2 model_overloaded via FD_MAX_IN_FLIGHT, and F-5 shutdown drain; `benchmark-results/m041-s02-t06-go-test.txt`, `m041-s02-t06-lint.txt`, and `m041-s02-t06-govulncheck.txt` pass the mandatory Go gates. |
| R012 | quality-attribute | active | M041-4tw0w7/S04 | M041-4tw0w7/S02 | D045 фиксирует cache-hot трактовку T-P-1..T-P-5. `tools/verify_fd_v2_perf.sh` prewarm-ит measured payload через real inference и затем требует `X-Cache: HIT` для latency cases. Evidence: `benchmark-results/fd-v2-perf-validation-m041-s04.md` PASS (batch=1 p95 2.236ms, batch=10 p95 3.468ms, batch=32 p95 7.595ms, sequential/concurrent/cache HIT pass) плюс non-blocking cache-miss diagnostics. |
| R013 | failure-visibility | validated | M041-4tw0w7/S03 | none | Validated by M041-4tw0w7/S03. Evidence: `benchmark-results/m041-s03-t07-observability-integration.txt` covers `/version`, `/info`, `/metrics`, and `/v1/healthcheck`; `m041-s03-t07-go-test.txt`, `m041-s03-t07-lint.txt`, and `m041-s03-t07-govulncheck.txt` pass mandatory Go gates. |
| R014 | operability | validated | M041-4tw0w7/S03 | none | S03 validated Server, X-Request-Id, X-Model-Id, X-Dimensions, Retry-After, and Connection headers. S04 validated `X-Cache: MISS` on first request and `X-Cache: HIT` on repeated input via `benchmark-results/fd-v2-perf-validation-m041-s04.md` and `api/fd_v2_cache_integration_test.go`; final verifier requires non-HIT count 0 for cache-hot cases. |
| R015 | failure-visibility | validated | M041-4tw0w7/S03 | none | Validated by M041-4tw0w7/S03. Evidence: `benchmark-results/m041-s03-t03-deep-health.txt` covers deep `/health` status fields and `last_inference_at`; `m041-s03-t06-warmup.txt` covers GET/POST `/warmup`; `m041-s03-t07-observability-integration.txt` covers integration behavior. |
| R016 | differentiator | validated | M041-4tw0w7/S04 | none | S04 implemented LRU cache with size/TTL env config, copy-on-read/write, eviction metrics, and EmbeddingCache adapter methods. Evidence: `api/cache/lru_test.go`, `api/fd_v2_cache_integration_test.go`, `benchmark-results/m041-s04-t03-cache-integration.txt`, and final `benchmark-results/fd-v2-perf-validation-m041-s04.md` showing repeated input `X-Cache: HIT` in 1.870ms and cache-hot latency targets passing. |
| R017 | differentiator | validated | M041-4tw0w7/S05 | none | `benchmark-results/fd-v2-validation-m041.md` 45/45 PASS: T019 validates `encoding_format=base64`; T020 validates `priority=high`; T021 validates `user`; T033 validates invalid `encoding_format` returns `encoding_format_invalid`; T034 validates invalid priority returns `priority_invalid`. Unit evidence in `benchmark-results/m041-s05-t01-go-test.txt`. |
| R018 | compliance/security | validated | M041-4tw0w7/S05 | none | Unit evidence in `benchmark-results/m041-s05-t02-go-test.txt` validates FD_API_KEY bearer auth: missing/wrong token -> 401 unauthorized, correct token -> 200, public endpoints skipped. `benchmark-results/fd-v2-validation-m041.md` validates CORS preflight (T027) and docs/openapi public surfaces. |
| R019 | quality-attribute | validated | M041-4tw0w7/S05 | none | `benchmark-results/fd-v2-validation-m041.md` 45/45 PASS validates `/openapi.json` (T010/T043), `/docs` (T011), `/v1/batch` (T038-T040), ETag/Cache-Control and 304 (T024-T026), `/v1/traces` (T012/T042), and cache-hot performance checks (T044-T045). T03/T05/T06/T07 task evidence covers rate limiting, cache validators, traces, OpenAPI validator, and docs unit tests. |
| R020 | quality-attribute | validated | M042-fjf2en/S01 | none | `documents/te-perf-root-cause-m042.md` explains TEI queue/startup behavior with T01/T02 evidence: direct TEI batch=32 queue p50 ~2434.795ms despite `max_concurrent_requests=512`, and restart/recreate spent ~48 minutes from `Starting model backend` to `Ready` after delayed missing-ONNX ORT fallback. The RCA includes hypothesis tree, verdict, and TEI-first recommendation. |
| R021 | quality-attribute | active | M042-fjf2en/S02 | none | tools/verify_fd_async_perf.sh: FD_ASYNC_CHUNKS=true vs false perf comparison. Cold path batch=128 ≤10s (was 25s sequential). Cold path batch=32 ≤4s (was 6s sequential). Cache hit path не regressed (≤5ms per request). Benchmark artifact в benchmark-results/fd-v2-async-perf-m042.md. |
| R022 | quality-attribute | deferred | M042-fjf2en/S03 | none | Deferred by user decision during M042: ONNX will not be implemented as current opt-in runtime. Prior M019 evidence remains research-only, not production readiness proof. |
| R023 | quality-attribute | validated | M043-dpr0cq/S01 | none | M043 S01: Tier 1 linters enabled and fixed; final lint 0 issues. Evidence: docs/static-analysis-phase1-report-m043.md, benchmark-results/m043-tier1-baseline.txt. |
| R024 | quality-attribute | validated | M043-dpr0cq/S02 | none | M043 S02: Tier 2 linters enabled; 17 baseline issues fixed; final lint 0 issues. Evidence: docs/static-analysis-phase2-report-m043.md, benchmark-results/m043-s02-final-lint.txt. |
| R025 | quality-attribute | validated | M043-dpr0cq/S03 | none | M043 S03: govulncheck CI step added and local govulncheck exits 0 with 0 reachable vulnerabilities; docs finalized. Evidence: benchmark-results/m043-s03-govulncheck-final.txt, docs/static-analysis-recommendation.md. |
| R026 | integration | active | none | none | `GET /openapi.json` returns an OAS 3.2.0 document; docs render it; the final contract verifier asserts `openapi == "3.2.0"`; external schema validation or compatibility checks pass; mandatory Go gates (`go test ./...`, golangci-lint v2.12.2, govulncheck) pass. |
| R027 | constraint | active | M042-fjf2en/S02 | none | Default build, Docker image, docs, and runtime config expose only TEI as current backend; ONNX build/runtime selectors are absent or fail closed as explicitly research-only; `go test ./...`, golangci-lint v2.12.2, and govulncheck pass. |

## Coverage Summary

- Active requirements: 5
- Mapped to slices: 4
- Validated: 21 (R001, R002, R003, R004, R005, R006, R007, R008, R009, R011, R013, R014, R015, R016, R017, R018, R019, R020, R023, R024, R025)
- Unmapped active requirements: 1
