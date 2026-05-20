# Codebase Map

Generated: 2026-05-20T10:31:24Z | Files: 96 | Described: 0/96
<!-- gsd:codebase-meta {"generatedAt":"2026-05-20T10:31:24Z","fingerprint":"c481a9881434c6a42502c4f506c77606ffd1174c","fileCount":96,"truncated":false} -->

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
- `api/embed/onnx_disabled.go`
- `api/embed/onnx_manifest_test.go`
- `api/embed/onnx_manifest.go`
- `api/embed/onnx_test.go`
- `api/embed/onnx_token_types.go`
- `api/embed/onnx_tokenizer_default.go`
- `api/embed/onnx_tokenizer_hf.go`
- `api/embed/onnx_types.go`
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
- *(36 files: 36 .txt)*

### docs/onnx-artifacts/
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/README.md`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

### docs/superpowers/plans/
- `docs/superpowers/plans/2026-05-14-fd-tiered-cache-redesign.md`

### docs/superpowers/specs/
- `docs/superpowers/specs/2026-05-14-fd-embedding-service-v2-design.md`
- `docs/superpowers/specs/2026-05-14-fd-embedding-service-v2-legal-design.md`

### tests/
- `tests/44-FZ-2026-articles.jsonl`

### tests/integration/
- `tests/integration/api_test.go`

### tools/
- `tools/compare_dense_embeddings.py`
- `tools/compare_onnx_dense_embeddings.py`
- `tools/compare_tokenizers.py`
- `tools/diagnose_onnx_sequence_length.py`
- `tools/evaluate_legal_retrieval.py`
- `tools/export_user_bge_m3_dense_onnx.py`
- `tools/profile_legal_divergence.py`
- `tools/verify_onnx_artifacts.py`
