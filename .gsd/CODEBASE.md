# Codebase Map

Generated: 2026-05-28T12:59:41Z | Files: 141 | Described: 0/141
<!-- gsd:codebase-meta {"generatedAt":"2026-05-28T12:59:41Z","fingerprint":"0e6474fd681e17df6d3474bf02efac1e349ff354","fileCount":141,"truncated":false} -->

### (root)/
- `.gitignore`
- `.golangci.yml`
- `benchmark.py`
- `CHANGELOG.md`
- `docker-compose.override.yaml`
- `docker-compose.yaml`
- `Dockerfile.onnx`
- `README.md`

### .github/workflows/
- `.github/workflows/go-quality.yml`
- `.github/workflows/onnx-packaging.yml`

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
- `api/handlers/health_test.go`
- `api/handlers/health.go`

### benchmark-results/
- *(67 files: 64 .txt, 3 .md)*

### docs/
- `docs/same-host-embedding-service-contract.md`

### docs/onnx-artifacts/
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/OPERATIONS.md`
- `docs/onnx-artifacts/PROVISIONING.md`
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
- `tools/build_onnx_image.sh`
- `tools/compare_dense_embeddings.py`
- `tools/compare_onnx_dense_embeddings.py`
- `tools/compare_tokenizers.py`
- `tools/diagnose_onnx_sequence_length.py`
- `tools/evaluate_legal_model_quick_gate.py`
- `tools/evaluate_legal_retrieval.py`
- `tools/export_user_bge_m3_dense_onnx.py`
- `tools/profile_legal_divergence.py`
- `tools/provision_onnx_artifacts.py`
- `tools/run_m040_s02_docker_restart_proof.sh`
- `tools/verify_legal_model_quick_gate_artifact.py`
- `tools/verify_m040_s02_artifacts.py`
- `tools/verify_m040_s04_recommendation.py`
- `tools/verify_onnx_artifacts.py`
- `tools/verify_onnx_export_contract.py`
