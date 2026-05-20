---
id: S01
parent: M022-i079tk
milestone: M022-i079tk
provides:
  - A local ONNX packaging proof for S02 CI boundary decisions.
requires:
  []
affects:
  - S02 CI artifact provisioning boundary
key_files:
  - Dockerfile.onnx
  - tools/build_onnx_image.sh
  - docs/onnx-artifacts/README.md
key_decisions:
  - Use generated `.gsd/runtime/docker/onnx1024-context` staging context.
  - Keep default `api/Dockerfile` as the TEI/default image.
  - Build ONNX image with `onnx hf_tokenizers` tags only.
patterns_established:
  - Generate minimal Docker contexts under `.gsd/runtime/` for local ignored artifacts.
  - Run checksum verifier before any ONNX packaging build.
  - Keep ONNX Docker packaging separate from default TEI image.
observability_surfaces:
  - Artifact verifier JSON output, build script echo of built image/context dir, Docker build logs, health and embedding smoke proof.
drill_down_paths:
  - .gsd/milestones/M022-i079tk/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M022-i079tk/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M022-i079tk/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T10:43:10.332Z
blocker_discovered: false
---

# S01: Dedicated ONNX Docker packaging proof

**S01 proved a dedicated opt-in ONNX Docker image can build and serve embeddings locally.**

## What Happened

S01 implemented and proved the dedicated ONNX Docker packaging path. The build script verifies local artifacts, stages only required files, builds a CGO-enabled binary with `onnx hf_tokenizers`, and packages ONNX model, tokenizer JSON, manifests, and ONNX Runtime shared library into a dedicated runtime image. The built container served `/health` and a smoke embedding with 1024 dimensions. Default Docker build remained passing and binaries remained untracked.

## Verification

ONNX build/run proof, default Docker build, verifier, binary hygiene, and cleanup passed.

## Requirements Advanced

- onnx-1024-docker-packaging — Implemented and proved local opt-in Docker packaging for ONNX 1024.

## Requirements Validated

- m022-local-onnx-image — `fd-api:onnx1024-m022` built and served `/health` plus a 1024-dim embedding smoke response.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None. The ONNX image proof passed; no blocker was needed.

## Known Limitations

Build context is large (1.569GB). CI needs an artifact store/cache before this can become a normal required job.

## Follow-ups

S02 should define the CI boundary: either add a safe contract-only CI check or document that full ONNX image CI requires external artifact provisioning/cache before it can be automated truthfully.

## Files Created/Modified

- `Dockerfile.onnx` — Dedicated opt-in ONNX Docker image definition.
- `tools/build_onnx_image.sh` — Build script that verifies artifacts and stages a minimal Docker context under `.gsd/runtime/`.
- `docs/onnx-artifacts/README.md` — Documented local ONNX Docker image proof command and staging behavior.
