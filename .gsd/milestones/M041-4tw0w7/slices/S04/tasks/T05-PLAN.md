---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T05: Final perf validation after backend remediation

Run tools/verify_fd_v2_perf.sh against a current fd instance backed by real inference and an isolated Redis cache namespace. Passing cache-hot runs are insufficient. This task remains pending while the current TEI CPU backend misses T-P latency targets; complete only after backend/runtime remediation (for example ONNX/GPU/faster TEI runtime) or explicit requirement rescope.

## Inputs

- `tools/verify_fd_v2_perf.sh`
- `docs/fd-v2.md`
- `benchmark-results/m041-s04-t04-bottleneck-measurement.txt`

## Expected Output

- `benchmark-results/fd-v2-perf-validation-m041-s04.md`

## Verification

FD_BASE_URL=http://localhost:8000 tools/verify_fd_v2_perf.sh exits 0 against current fd with isolated EMBEDDING_CACHE_VERSION and real inference. Artifact contains p50/p95/p99 for T-P cases, 100 sequential 0 errors, 4x8 concurrent <2s, and cache HIT <5ms.
