# Codebase Map

Generated: 2026-05-20T04:08:30Z | Files: 70 | Described: 0/70
<!-- gsd:codebase-meta {"generatedAt":"2026-05-20T04:08:30Z","fingerprint":"af578b903aa8e28ff87c7bb4bf74f84951c2a4f1","fileCount":70,"truncated":false} -->

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
- `api/embed/hf_tokenizer_native_test.go`
- `api/embed/hf_tokenizer_native.go`
- `api/embed/onnx_manifest_test.go`
- `api/embed/onnx_manifest.go`
- `api/embed/onnx_test.go`
- `api/embed/onnx_tokenizer_default.go`
- `api/embed/onnx_tokenizer_hf.go`
- `api/embed/onnx.go`
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
- `benchmark-results/fd-dense-comparator-m010-s02.txt`
- `benchmark-results/fd-environment-inxi-m008.txt`
- `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt`
- `benchmark-results/fd-go-onnx-m011-s03.txt`
- `benchmark-results/fd-onnx-fp32-m010-s03.txt`
- `benchmark-results/fd-runtime-stats-logs.txt`
- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`
- `benchmark-results/fd-tokenizer-go-current-m012-s02.txt`
- `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt`
- `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt`
- `benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt`

### docs/onnx-artifacts/
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

### docs/superpowers/plans/
- `docs/superpowers/plans/2026-05-14-fd-tiered-cache-redesign.md`

### docs/superpowers/specs/
- `docs/superpowers/specs/2026-05-14-fd-embedding-service-v2-design.md`
- `docs/superpowers/specs/2026-05-14-fd-embedding-service-v2-legal-design.md`

### tests/integration/
- `tests/integration/api_test.go`

### tools/
- `tools/compare_dense_embeddings.py`
- `tools/compare_onnx_dense_embeddings.py`
- `tools/compare_tokenizers.py`
- `tools/export_user_bge_m3_dense_onnx.py`
