# Codebase Map

Generated: 2026-05-19T18:22:26Z | Files: 47 | Described: 0/47
<!-- gsd:codebase-meta {"generatedAt":"2026-05-19T18:22:26Z","fingerprint":"5c809e87249560d30dea765d0c5f3c50446d58a1","fileCount":47,"truncated":false} -->

### (root)/
- `.gitignore`
- `.golangci.yml`
- `benchmark.py`
- `CHANGELOG.md`
- `docker-compose.override.yaml`
- `docker-compose.yaml`
- `README.md`

### .github/workflows/
- `.github/workflows/go-quality.yml`

### api/
- `api/.dockerignore`
- `api/Dockerfile`
- `api/Dockerfile.tests`
- `api/go.mod`
- `api/go.sum`
- `api/main_test.go`
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
- `api/handlers/constants.go`
- `api/handlers/embeddings_integration_test.go`
- `api/handlers/embeddings.go`
- `api/handlers/health.go`

### benchmark-results/
- `benchmark-results/fd-benchmark-baseline-py313.txt`
- `benchmark-results/fd-benchmark-m004-final.txt`
- `benchmark-results/fd-benchmark-m004-s01.txt`
- `benchmark-results/fd-benchmark-m004-s03.txt`
- `benchmark-results/fd-benchmark-m009-s01.txt`
- `benchmark-results/fd-benchmark-m009-s02.txt`
- `benchmark-results/fd-benchmark-m009-s03.txt`
- `benchmark-results/fd-benchmark-m009-s04.txt`
- `benchmark-results/fd-environment-inxi-m008.txt`
- `benchmark-results/fd-runtime-stats-logs.txt`

### docs/superpowers/plans/
- `docs/superpowers/plans/2026-05-14-fd-tiered-cache-redesign.md`

### docs/superpowers/specs/
- `docs/superpowers/specs/2026-05-14-fd-embedding-service-v2-design.md`
- `docs/superpowers/specs/2026-05-14-fd-embedding-service-v2-legal-design.md`

### tests/integration/
- `tests/integration/api_test.go`
