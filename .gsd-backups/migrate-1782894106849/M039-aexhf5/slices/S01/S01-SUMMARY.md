---
id: S01
parent: M039-aexhf5
milestone: M039-aexhf5
provides:
  - Fresh image `fd-api:onnx1024-m039` with id `sha256:40b80c47491d27402b0213a76e86c46332968dabba7ff1c55b75555ee6ca79dc`.
  - Packaged smoke proof and rerun proof for downstream legal/performance gates.
requires:
  []
affects:
  - S02
key_files:
  - benchmark-results/fd-onnx-docker-smoke-m039-s01.txt
  - benchmark-results/fd-onnx-docker-smoke-rerun-m039-s01.txt
key_decisions: []
patterns_established:
  - Packaged ONNX runtime health verification must include `ONNX_RUNTIME_SHA256` in the container environment.
observability_surfaces:
  - Docker image id, health runtime metadata, probe hash, embedding dimension/norm, cache namespace, artifact/tokenizer/runtime verification flags.
drill_down_paths:
  - .gsd/milestones/M039-aexhf5/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M039-aexhf5/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M039-aexhf5/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M039-aexhf5/slices/S01/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T11:22:26.553Z
blocker_discovered: false
---

# S01: Packaged ONNX smoke proof

**Built and smoke-tested packaged ONNX Docker runtime with verified artifact, tokenizer, and runtime library.**

## What Happened

S01 verified local packaging prerequisites, built a fresh dedicated ONNX image `fd-api:onnx1024-m039`, and smoke-tested the packaged Go ONNX API. A first probe exposed that runtime-library verification is only reported when `ONNX_RUNTIME_SHA256` is supplied, so the smoke was rerun with the sha env var. Both accepted smoke runs passed with verified artifact, tokenizer, runtime library, CPU provider, 1024-dimensional normalized embeddings, isolated namespaces, and production_default=false. The container was stopped and port 18000 was clean.

## Verification

S01 verification passed: artifact prerequisites, Docker build, packaged smoke, packaged smoke rerun, leak checks, cleanup, port check, and GitNexus detect all passed.

## Requirements Advanced

- implicit target-runtime validation requirement — Added packaged runtime smoke evidence through the Dockerized Go ONNX API.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

The first packaged health check showed `runtime_library_verified=false` because `ONNX_RUNTIME_SHA256` was not set. The accepted smoke and explicit rerun restarted the container with the sha env var and passed.

## Known Limitations

Packaged legal/performance gates remain for S02. Hosted workflow proof and production promotion remain out of scope.

## Follow-ups

S02 should reuse image `fd-api:onnx1024-m039` and always pass `ONNX_RUNTIME_SHA256` when starting the packaged runtime.

## Files Created/Modified

- `benchmark-results/fd-onnx-docker-smoke-m039-s01.txt` — Packaged smoke result with runtime verification enabled.
- `benchmark-results/fd-onnx-docker-smoke-rerun-m039-s01.txt` — User-requested packaged smoke rerun result.
