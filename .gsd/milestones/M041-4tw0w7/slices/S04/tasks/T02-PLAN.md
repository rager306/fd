---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: LRU cache implementation

api/cache/lru.go: in-memory LRU cache на (string, int) → []float32. TTL 24h, size 10000, configurable через env FD_CACHE_SIZE, FD_CACHE_TTL_HOURS. Использовать hashicorp/golang-lru или свою реализацию с sync.RWMutex. Метрики: fd_cache_hits_total{result=hit|miss} counter, fd_cache_evictions_total counter. Cache key = SHA256(input_text + | + str(dimensions)).

## Inputs

- None specified.

## Expected Output

- `api/cache/lru.go`
- `api/cache/lru_test.go`

## Verification

Unit tests: Get/Put корректны. Eviction на size limit. TTL expiration. Concurrent access safe (race detector). fd_cache_hits_total increments on hit.
