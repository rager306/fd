# Codebase Map

Generated: 2026-06-14T12:46:46Z | Files: 318 | Described: 0/318
<!-- gsd:codebase-meta {"generatedAt":"2026-06-14T12:46:46Z","fingerprint":"d686ce4930d77bfc344f3a83c17241ee269860c7","fileCount":318,"truncated":false} -->

### (root)/
- `.gitignore`
- `.golangci.yml`
- `benchmark.py`
- `CHANGELOG.md`
- `docker-compose.override.yaml`
- `docker-compose.yaml`
- `README.md`

### .artifacts/browser/2026-06-14T12-23-36-082Z-session/
- `.artifacts/browser/2026-06-14T12-23-36-082Z-session/m045-s03-browser-health-timeline.json`

### .github/workflows/
- `.github/workflows/go-quality.yml`

### api/
- `api/.dockerignore`
- `api/Dockerfile`
- `api/Dockerfile.tests`
- `api/fd_v2_cache_integration_test.go`
- `api/fd_v2_lifecycle_integration_test.go`
- `api/fd_v2_observability_integration_test.go`
- `api/go.mod`
- `api/go.sum`
- `api/main_env_test.go`
- `api/main_test.go`
- `api/main.go`

### api/buildinfo/
- `api/buildinfo/info_test.go`
- `api/buildinfo/info.go`

### api/cache/
- `api/cache/local_test.go`
- `api/cache/local.go`
- `api/cache/lru_rapid_test.go`
- `api/cache/lru_test.go`
- `api/cache/lru.go`
- `api/cache/redis_bench_test.go`
- `api/cache/redis_binary_test.go`
- `api/cache/redis_test.go`
- `api/cache/redis.go`
- `api/cache/tiered_cache_test.go`
- `api/cache/tiered_constructor_test.go`
- `api/cache/tiered_test.go`
- `api/cache/tiered.go`

### api/embed/
- `api/embed/codec_rapid_test.go`
- `api/embed/codec.go`
- `api/embed/dimensions_test.go`
- `api/embed/tei_test.go`
- `api/embed/tei.go`
- `api/embed/types_test.go`
- `api/embed/types.go`

### api/handlers/
- *(23 files: 23 .go)*

### api/lifecycle/
- `api/lifecycle/default_state_test.go`
- `api/lifecycle/shutdown_test.go`
- `api/lifecycle/shutdown.go`
- `api/lifecycle/state_test.go`
- `api/lifecycle/state.go`
- `api/lifecycle/warmup_test.go`
- `api/lifecycle/warmup.go`

### api/middleware/
- `api/middleware/auth_test.go`
- `api/middleware/auth.go`
- `api/middleware/cache_headers_test.go`
- `api/middleware/cache_headers.go`
- `api/middleware/cors_test.go`
- `api/middleware/cors.go`
- `api/middleware/env_constructors_test.go`
- `api/middleware/headers_test.go`
- `api/middleware/headers.go`
- `api/middleware/lifecycle_test.go`
- `api/middleware/lifecycle.go`
- `api/middleware/ratelimit_test.go`
- `api/middleware/ratelimit.go`
- `api/middleware/validation_test.go`
- `api/middleware/validation.go`

### api/observability/
- `api/observability/metrics_test.go`
- `api/observability/metrics.go`
- `api/observability/traces_test.go`
- `api/observability/traces.go`

### api/openapi/
- `api/openapi/spec_test.go`
- `api/openapi/spec.go`

### benchmark-results/
- *(186 files: 177 .txt, 9 .md)*

### docs/
- `docs/fd-v2.md`
- `docs/same-host-embedding-service-contract.md`
- `docs/static-analysis-phase1-report-m043.md`
- `docs/static-analysis-phase2-report-m043.md`
- `docs/static-analysis-recommendation.md`

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

### documents/
- `documents/onnx-deactivation-inventory-m042.md`
- `documents/te-perf-root-cause-m042.md`
- `documents/te-perf-snapshot-m042-s01.md`
- `documents/tei-startup-mitigation-m045.md`
- `documents/tei-startup-recon-m045.md`

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
- `tools/profile_tei_concurrency.sh`
- `tools/provision_onnx_artifacts.py`
- `tools/run_m040_s02_docker_restart_proof.sh`
- `tools/verify_fd_v2_contract.py`
- `tools/verify_fd_v2_perf.sh`
- `tools/verify_legal_model_quick_gate_artifact.py`
- `tools/verify_m040_s02_artifacts.py`
- `tools/verify_m040_s04_recommendation.py`
- `tools/verify_onnx_artifacts.py`
- `tools/verify_onnx_export_contract.py`
