---
id: T02
parent: S01
milestone: M022-i079tk
key_files:
  - Dockerfile.onnx
  - tools/build_onnx_image.sh
  - docs/onnx-artifacts/README.md
key_decisions:
  - Dedicated ONNX packaging path uses `Dockerfile.onnx` plus `tools/build_onnx_image.sh`.
  - The script creates a temporary context under `.gsd/runtime/docker/onnx1024-context`, runs the artifact verifier first, copies only required source/manifests/artifacts, and builds with `onnx hf_tokenizers` tags.
  - Default image remains `api/Dockerfile`; ONNX image is explicit opt-in.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:35:43.921Z
blocker_discovered: false
---

# T02: Implemented the dedicated ONNX Docker packaging script and Dockerfile.

**Implemented the dedicated ONNX Docker packaging script and Dockerfile.**

## What Happened

Implemented the dedicated ONNX Docker packaging path. `Dockerfile.onnx` builds a CGO-enabled ONNX binary with `onnx hf_tokenizers` tags and packages the verified ONNX model, tokenizer JSON, native tokenizer static library, manifests, and ONNX Runtime shared library into a Debian runtime image. `tools/build_onnx_image.sh` verifies artifacts, generates a minimal staging context under `.gsd/runtime/`, and builds the image. The README now documents the local image proof command.

## Verification

Shell syntax for the build script passed and the existing default image smoke command works.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `chmod +x tools/build_onnx_image.sh && bash -n tools/build_onnx_image.sh` | 0 | ✅ pass — script syntax valid | 0ms |
| 2 | `docker run --rm --entrypoint /bin/sh fd-api:m021-default -c 'echo default-image-smoke-ok'` | 0 | ✅ pass — default-image-smoke-ok | 0ms |

## Deviations

None.

## Known Issues

The ONNX image build has not yet been run in T02; T03 will provide build/run proof or concrete blocker.

## Files Created/Modified

- `Dockerfile.onnx`
- `tools/build_onnx_image.sh`
- `docs/onnx-artifacts/README.md`
