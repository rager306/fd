---
id: T03
parent: S01
milestone: M039-aexhf5
key_files:
  - benchmark-results/fd-onnx-docker-smoke-m039-s01.txt
  - benchmark-results/fd-onnx-docker-smoke-rerun-m039-s01.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T11:21:48.741Z
blocker_discovered: false
---

# T03: Packaged ONNX Docker smoke passed twice with runtime library verification enabled.

**Packaged ONNX Docker smoke passed twice with runtime library verification enabled.**

## What Happened

Ran the packaged ONNX image on port 18000 with isolated Redis namespaces. The first health probe revealed `runtime_library_verified=false` because the sha env var was not set, so the container was restarted with `ONNX_RUNTIME_SHA256`. The accepted smoke and a user-requested rerun both passed with backend `onnx`, artifact/tokenizer/runtime-library verification true, provider `CPUExecutionProvider`, 1024-dimensional normalized embeddings, and `production_default=false`. Raw probe text was excluded from artifacts. Container was stopped and port 18000 is clean.

## Verification

Smoke checks and cleanup passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker run fd-api:onnx1024-m039 with ONNX_RUNTIME_SHA256 and namespace m039-docker-smoke` | 0 | ✅ pass — container started | 8100ms |
| 2 | `packaged smoke probe -> benchmark-results/fd-onnx-docker-smoke-m039-s01.txt` | 0 | ✅ pass — 1024 dims, norm 1.00000053, artifact/tokenizer/runtime verified | 9800ms |
| 3 | `docker restart/run with namespace m039-docker-smoke-rerun` | 0 | ✅ pass — rerun container started | 5800ms |
| 4 | `packaged smoke rerun -> benchmark-results/fd-onnx-docker-smoke-rerun-m039-s01.txt` | 0 | ✅ pass — 1024 dims, norm 1.00000049, artifact/tokenizer/runtime verified | 4400ms |
| 5 | `M039 packaged smoke artifact leak checks` | 0 | ✅ pass — missing_required=0, leak_markers=0, signed_url_like=0 | 5600ms |
| 6 | `docker rm -f fd-onnx-m039-smoke && port check` | 0 | ✅ pass — port_18000_clean | 5500ms |

## Deviations

Initial smoke was rerun because the first container start omitted `ONNX_RUNTIME_SHA256`, causing health metadata to report `runtime_library_verified=false`. The accepted smoke and explicit rerun both set `ONNX_RUNTIME_SHA256` and passed.

## Known Issues

Container runtime library verification requires `ONNX_RUNTIME_SHA256` to be explicitly set; Dockerfile embeds the runtime library path but not the sha env var.

## Files Created/Modified

- `benchmark-results/fd-onnx-docker-smoke-m039-s01.txt`
- `benchmark-results/fd-onnx-docker-smoke-rerun-m039-s01.txt`
