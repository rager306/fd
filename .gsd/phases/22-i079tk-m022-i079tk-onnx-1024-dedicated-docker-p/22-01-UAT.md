# S01: Dedicated ONNX Docker packaging proof — UAT

**Milestone:** M022-i079tk
**Written:** 2026-05-20T10:43:10.332Z

# S01 UAT — Dedicated ONNX Docker packaging proof

## Checks

- [x] Packaging strategy chosen.
- [x] `Dockerfile.onnx` exists and is opt-in.
- [x] `tools/build_onnx_image.sh` verifies artifacts before staging/building.
- [x] ONNX image builds with `onnx hf_tokenizers` tags.
- [x] ONNX container starts on port 18000.
- [x] `/health` returns ok.
- [x] Embedding smoke returns 1024 dimensions for `deepvk/USER-bge-m3`.
- [x] Default Docker image still builds.
- [x] No ONNX/native/runtime binaries are tracked.
- [x] Port 18000 is clean after cleanup.

## Result

Pass. Local ONNX Docker packaging proof is established.

