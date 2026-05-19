# Codebase Map

Generated: 2026-05-19T10:38:38Z | Files: 37 | Described: 0/37
<!-- gsd:codebase-meta {"generatedAt":"2026-05-19T10:38:38Z","fingerprint":"9955347fd754ae1541afe0f6e45e9c8a7d6eb278","fileCount":37,"truncated":false} -->

### (root)/
- `.gitignore`
- `benchmark.py`
- `CHANGELOG.md`
- `docker-compose.override.yaml`
- `docker-compose.yaml`
- `README.md`

### api/
- `api/.dockerignore`
- `api/Dockerfile`
- `api/Dockerfile.tests`
- `api/go.mod`
- `api/main.go`

### api/cache/
- `api/cache/local_test.go`
- `api/cache/local.go`
- `api/cache/redis_bench_test.go`
- `api/cache/redis_binary_test.go`
- `api/cache/redis_test.go`
- `api/cache/redis.go`
- `api/cache/tiered_cache_test.go`
- `api/cache/tiered_test.go`
- `api/cache/tiered.go`

### api/embed/
- `api/embed/dimensions_test.go`
- `api/embed/tei.go`
- `api/embed/types_test.go`
- `api/embed/types.go`

### api/handlers/
- `api/handlers/batch.go`
- `api/handlers/embeddings_integration_test.go`
- `api/handlers/embeddings.go`
- `api/handlers/health.go`

### benchmark-results/
- `benchmark-results/fd-benchmark-baseline-py313.txt`
- `benchmark-results/fd-benchmark-m004-final.txt`
- `benchmark-results/fd-benchmark-m004-s01.txt`
- `benchmark-results/fd-benchmark-m004-s03.txt`
- `benchmark-results/fd-runtime-stats-logs.txt`

### docs/superpowers/plans/
- `docs/superpowers/plans/2026-05-14-fd-tiered-cache-redesign.md`

### docs/superpowers/specs/
- `docs/superpowers/specs/2026-05-14-fd-embedding-service-v2-design.md`
- `docs/superpowers/specs/2026-05-14-fd-embedding-service-v2-legal-design.md`

### tests/integration/
- `tests/integration/api_test.go`
