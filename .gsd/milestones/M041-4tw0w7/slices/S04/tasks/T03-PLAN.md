---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Cache middleware integration

api/middleware/cache.go: gin middleware который проверяет cache перед вызовом model. На HIT — return cached embedding + X-Cache: HIT. На MISS — call model, store result в cache + return + X-Cache: MISS. Должен сидеть ПОСЛЕ validation (S01) и lifecycle gate (S02), ДО model call. Кэш key учитывает dimensions.

## Inputs

- None specified.

## Expected Output

- `api/middleware/cache.go`
- `api/middleware/cache_test.go`

## Verification

Unit tests: первый запрос → X-Cache: MISS, повторный → X-Cache: HIT. Cache HIT latency < 5ms. fd_cache_hits_total{result=hit} increments.
