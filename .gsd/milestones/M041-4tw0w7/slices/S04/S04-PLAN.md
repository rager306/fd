# S04: Performance baseline and LRU cache

**Goal:** Performance baseline (1<50ms, 10<200ms, 32<1000ms p95, 100 sequential zero errors, 4×8 concurrent < 2s) + in-memory LRU cache (10000 entries, 24h TTL) с X-Cache header и fd_cache_hits_total metric. Закрывает R-P0-6 и R-P1-4.
**Demo:** After this, cache/perf validation includes warm/cold baseline plus M043 gates; cache code must keep context propagation, gocyclo <=15 for production functions, and no new static-analysis suppressions without justification.

## Must-Haves

- T-P-1..T-P-5 (Section 5.4) все pass с p95 latency: 1<50ms, 10<200ms, 32<1000ms, 100 sequential zero errors, 4×8 concurrent total < 2s
- T-HDR-6/7 pass: повторный тот же input возвращает X-Cache: HIT с latency < 5ms; первый запрос MISS
- Cache eviction работает: > 10000 уникальных inputs в течение часа не вызывает OOM и recent entries остаются
- /metrics показывает fd_cache_hits_total{result=hit|miss} counter
- F-4 scenario: cache miss → hit переход работает и видим в логах/метриках
- benchmark-results/fd-v2-perf-validation-m041-s04.md final artifact
- golangci-lint pass

## Proof Level

- This slice proves: runtime + integration + benchmark

## Integration Closure

Cache слой между validation (S01) и model inference. Cache key = SHA256(input_text + | + dimensions). Cache lookup метрика HIT/MISS → X-Cache header (S03) и fd_cache_hits_total (S03 metrics).

## Verification

- Добавляет fd_cache_hits_total counter и расширяет fd_batch_size histogram (cache hit снижает effective batch size).

## Tasks

- [x] **T01: Baseline measurement fd perf захвачен до S04 T04 perf optimization** `est:2h`
  tools/measure_fd_baseline.sh: простой benchmark который шлёт 1/10/32 input запросы 100 раз каждый, меряет p50/p95/p99 latency, error rate. Запускается против текущего fd (после S01/S02/S03) для baseline numbers. Сохранить в benchmark-results/fd-v2-baseline-before-m041-s04.md. Это даст опорные цифры чтобы понять, нужен ли real optimization или достаточно validation fixes из S01. Спека target values: docs/fd-v2.md Section 5.4 T-P-1..T-P-5.
  - Files: `tools/measure_fd_baseline.sh`, `benchmark-results/fd-v2-baseline-before-m041-s04.md`
  - Verify: Baseline artifact содержит: p50/p95/p99 для batch=1/10/32, error rate, throughput. Можно сравнить с target values (50/200/1000ms).

- [x] **T02: Added goroutine-safe in-memory LRU vector cache with TTL, env configuration, SHA256 keys, and cache metrics hooks.** `est:3h`
  api/cache/lru.go: in-memory LRU cache на (string, int) → []float32. TTL 24h, size 10000, configurable через env FD_CACHE_SIZE, FD_CACHE_TTL_HOURS. Использовать hashicorp/golang-lru или свою реализацию с sync.RWMutex. Метрики: fd_cache_hits_total{result=hit|miss} counter, fd_cache_evictions_total counter. Cache key = SHA256(input_text + | + str(dimensions)).
  - Files: `api/cache/lru.go`, `api/cache/lru_test.go`
  - Verify: Unit tests: Get/Put корректны. Eviction на size limit. TTL expiration. Concurrent access safe (race detector). fd_cache_hits_total increments on hit.

- [x] **T03: Integrated cache HIT/MISS behavior into /v1/embeddings with X-Cache headers, LRU EmbeddingCache adapter methods, and metrics verification.** `est:3h`
  api/middleware/cache.go: gin middleware который проверяет cache перед вызовом model. На HIT — return cached embedding + X-Cache: HIT. На MISS — call model, store result в cache + return + X-Cache: MISS. Должен сидеть ПОСЛЕ validation (S01) и lifecycle gate (S02), ДО model call. Кэш key учитывает dimensions.
  - Files: `api/middleware/cache.go`, `api/middleware/cache_test.go`
  - Verify: Unit tests: первый запрос → X-Cache: MISS, повторный → X-Cache: HIT. Cache HIT latency < 5ms. fd_cache_hits_total{result=hit} increments.

- [x] **T04: Confirmed real cache-miss inference performance blocker: current fd + TEI CPU misses T-P latency targets when Redis namespace is isolated.** `est:8h (если нужен); 1h (если baseline уже ОК)`
  Если baseline из T01 НЕ достигает target (50/200/1000ms p95), выполнить targeted optimization. Возможные направления: (a) batch tensor packing — отправлять весь input array одним тензором в ONNX/TEI за раз, не per-item; (b) concurrent workers — если backend медленный, N goroutines обрабатывают chunks; (c) request coalescing — multiple requests с одинаковым input идут в один model call; (d) ONNX session warmup и graph optimization. Каждое изменение проверяется отдельным benchmark. Изменения оформляются как opt-in через env (FD_BATCH_TENSOR_PACKING=true, FD_CONCURRENT_WORKERS=4, etc).
  - Files: `api/embed/optimizations.go`, `api/embed/optimizations_test.go`, `api/embed/perf_test.go`
  - Verify: После optimization: повторный baseline показывает p95 latency в target. benchmark-results/fd-v2-perf-m041-s04.md содержит before/after numbers и какие optimization сработали.

- [ ] **T05: Final perf validation after backend remediation** `est:1`
  Run tools/verify_fd_v2_perf.sh against a current fd instance backed by real inference and an isolated Redis cache namespace. Passing cache-hot runs are insufficient. This task remains pending while the current TEI CPU backend misses T-P latency targets; complete only after backend/runtime remediation (for example ONNX/GPU/faster TEI runtime) or explicit requirement rescope.
  - Files: `tools/verify_fd_v2_perf.sh`, `benchmark-results/fd-v2-perf-validation-m041-s04.md`
  - Verify: FD_BASE_URL=http://localhost:8000 tools/verify_fd_v2_perf.sh exits 0 against current fd with isolated EMBEDDING_CACHE_VERSION and real inference. Artifact contains p50/p95/p99 for T-P cases, 100 sequential 0 errors, 4x8 concurrent <2s, and cache HIT <5ms.

## Files Likely Touched

- tools/measure_fd_baseline.sh
- benchmark-results/fd-v2-baseline-before-m041-s04.md
- api/cache/lru.go
- api/cache/lru_test.go
- api/middleware/cache.go
- api/middleware/cache_test.go
- api/embed/optimizations.go
- api/embed/optimizations_test.go
- api/embed/perf_test.go
- tools/verify_fd_v2_perf.sh
- benchmark-results/fd-v2-perf-validation-m041-s04.md
