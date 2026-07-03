---
id: T03
parent: S01
milestone: M022-i079tk
key_files:
  - Dockerfile.onnx
  - tools/build_onnx_image.sh
  - docs/onnx-artifacts/README.md
key_decisions:
  - The dedicated ONNX image packages the model artifact and ONNX Runtime shared library inside the image, while the native tokenizer static library is build-time only.
  - Runtime smoke uses isolated `EMBEDDING_CACHE_VERSION=m022-onnx-docker-smoke`.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:42:48.363Z
blocker_discovered: false
---

# T03: Proved the dedicated ONNX Docker image builds and serves embeddings locally.

**Proved the dedicated ONNX Docker image builds and serves embeddings locally.**

## What Happened

Ran the ONNX Docker packaging proof. The build script verified artifacts, generated a staging context, built `fd-api:onnx1024-m022`, and the container started successfully on port 18000. `/health` returned ok and a non-legal smoke embedding request returned 1024 dimensions with model `deepvk/USER-bge-m3`. The default Docker build also passed, binary hygiene stayed clean, and the smoke container was stopped with port 18000 clean.

## Verification

ONNX image build, container health, embedding smoke, default Docker build, artifact verifier, binary hygiene, and cleanup all passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `IMAGE_TAG=fd-api:onnx1024-m022 tools/build_onnx_image.sh` | 0 | ✅ pass — Successfully tagged fd-api:onnx1024-m022 | 140300ms |
| 2 | `docker run --rm --name fd-onnx-m022 --network host ... fd-api:onnx1024-m022` | 0 | ✅ pass — process ready on port 18000 | 6000ms |
| 3 | `curl -fsS http://localhost:18000/health` | 0 | ✅ pass — {"status":"ok"} | 0ms |
| 4 | `curl -fsS http://localhost:18000/v1/embeddings ...` | 0 | ✅ pass — embedding_dims=1024 model=deepvk/USER-bge-m3 | 0ms |
| 5 | `docker build -f api/Dockerfile -t fd-api:m022-default api` | 0 | ✅ pass — Successfully tagged fd-api:m022-default | 4500ms |
| 6 | `python3 -m py_compile tools/verify_onnx_artifacts.py && tools/verify_onnx_artifacts.py` | 0 | ✅ pass — m022_artifact_verifier=pass | 4500ms |
| 7 | `git ls-files binary hygiene; bg_shell list; lsof port 18000` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0; no background processes; port_18000_clean | 0ms |

## Deviations

None. The ONNX image build and container smoke both passed.

## Known Issues

Docker build context is large (`1.569GB`) because it includes the FP32 ONNX artifact and ONNX Runtime shared library. This is acceptable for local proof but future CI should provision/cache artifacts explicitly.

## Files Created/Modified

- `Dockerfile.onnx`
- `tools/build_onnx_image.sh`
- `docs/onnx-artifacts/README.md`
