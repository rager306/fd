---
id: T02
parent: S01
milestone: M039-aexhf5
key_files:
  - Dockerfile.onnx
  - tools/build_onnx_image.sh
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T11:05:00.032Z
blocker_discovered: false
---

# T02: Built fresh packaged ONNX image `fd-api:onnx1024-m039`.

**Built fresh packaged ONNX image `fd-api:onnx1024-m039`.**

## What Happened

Built a fresh dedicated ONNX Docker image from current local artifacts using `tools/build_onnx_image.sh`. Image tag is `fd-api:onnx1024-m039`, image id `sha256:40b80c47491d27402b0213a76e86c46332968dabba7ff1c55b75555ee6ca79dc`, size `1377538338` bytes. Build context was generated under `.gsd/runtime/docker/onnx1024-context`.

## Verification

Docker image build and inspect passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `IMAGE_TAG='fd-api:onnx1024-m039' tools/build_onnx_image.sh` | 0 | ✅ pass — image built successfully | 110900ms |
| 2 | `docker image inspect fd-api:onnx1024-m039 --format '{{.Id}} {{.Size}} {{.Created}}'` | 0 | ✅ pass — image id sha256:40b80c47491d..., size 1377538338 | 8900ms |

## Deviations

None.

## Known Issues

Image is local only; no registry push or hosted workflow proof.

## Files Created/Modified

- `Dockerfile.onnx`
- `tools/build_onnx_image.sh`
