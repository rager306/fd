---
id: T01
parent: S01
milestone: M023-myx9u4
key_files: []
key_decisions:
  - Use existing healthy TEI baseline at `http://localhost:8000`.
  - Use packaged ONNX image `fd-api:onnx1024-m022-final` on `http://localhost:18000`.
  - Use isolated ONNX cache namespace `m023-onnx-docker-legal`.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:53:33.126Z
blocker_discovered: false
---

# T01: Prepared healthy TEI and packaged ONNX endpoints for the M023 legal gate.

**Prepared healthy TEI and packaged ONNX endpoints for the M023 legal gate.**

## What Happened

Prepared the packaged legal gate environment. The TEI baseline API on port 8000 is healthy. The packaged ONNX Docker image exists, started on port 18000 with isolated cache namespace `m023-onnx-docker-legal`, returned healthy status, and produced a 1024-dimensional non-legal smoke embedding for `deepvk/USER-bge-m3`.

## Verification

TEI and packaged ONNX endpoints are healthy; ONNX embedding smoke returned 1024 dimensions.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `curl -fsS http://localhost:8000/health` | 0 | ✅ pass — tei_api_8000_health=pass | 0ms |
| 2 | `docker image inspect fd-api:onnx1024-m022-final` | 0 | ✅ pass — onnx_image_present=fd-api:onnx1024-m022-final | 0ms |
| 3 | `docker run --rm --name fd-onnx-m023-legal --network host ... EMBEDDING_CACHE_VERSION=m023-onnx-docker-legal fd-api:onnx1024-m022-final` | 0 | ✅ pass — ready on port 18000 | 6000ms |
| 4 | `curl -fsS http://localhost:18000/health` | 0 | ✅ pass — onnx_docker_18000_health=pass | 0ms |
| 5 | `curl -fsS http://localhost:18000/v1/embeddings ...` | 0 | ✅ pass — embedding_dims=1024 model=deepvk/USER-bge-m3 | 0ms |

## Deviations

None.

## Known Issues

Existing TEI/Redis stack was already running; M023 did not restart it.

## Files Created/Modified

None.
