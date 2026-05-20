---
id: S02
parent: M021-4t2wpt
milestone: M021-4t2wpt
provides:
  - Validated default Docker boundary and opt-in ONNX build-tag contract.
requires:
  []
affects:
  - Future dedicated ONNX Docker/CI artifact provisioning milestone
key_files:
  - api/embed/onnx.go
  - api/embed/onnx_disabled.go
  - api/embed/onnx_types.go
  - api/embed/onnx_token_types.go
  - api/embed/onnx_tokenizer_default.go
  - api/embed/onnx_tokenizer_hf.go
  - api/embed/onnx_test.go
  - .gsd/DECISIONS.md
key_decisions:
  - D019: default Docker/CI remains CGO-disabled TEI path; opt-in ONNX runtime requires explicit `onnx hf_tokenizers` tags and verified artifacts.
  - Do not change production/default runtime; ONNX remains experimental.
patterns_established:
  - Real ONNX backend builds require `onnx`; parity-correct native tokenizer requires `hf_tokenizers`; validated ONNX-native runtime uses both.
  - Default Docker builds are guardrails for production/default TEI safety.
observability_surfaces:
  - D019 decision, verifier output, Docker build evidence, and build-tag-separated failure mode for ONNX disabled builds.
drill_down_paths:
  - .gsd/milestones/M021-4t2wpt/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M021-4t2wpt/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M021-4t2wpt/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T10:26:34.724Z
blocker_discovered: false
---

# S02: Docker CI boundary validation

**S02 fixed and validated the Docker/CI boundary so default builds no longer require ONNX runtime artifacts.**

## What Happened

S02 validated and fixed the Docker/CI boundary. The initial default Docker build failed because ONNX runtime code was included in the CGO-disabled default build. The fix made ONNX runtime implementation explicit via the `onnx` build tag and left default builds with a safe stub. Fresh verification then passed: default Go tests, lint, native tokenizer tests, ONNX+native smoke tests, default Docker build, artifact verifier, binary hygiene, runtime cleanup, and GitNexus scope.

## Verification

Fresh verification passed across artifact verifier, default and tagged Go tests, lint, default Docker build, binary hygiene, runtime cleanup, and GitNexus scope.

## Requirements Advanced

- onnx-1024-packaging-contract — Protected default Docker/CI from ONNX/native artifact dependencies and documented opt-in ONNX build requirements.

## Requirements Validated

- m021-default-docker-boundary — `docker build -f api/Dockerfile -t fd-api:m021-default api` passes after ONNX build-tag split.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Default Docker validation uncovered and fixed an ONNX build-tag boundary regression: default CGO-disabled builds were importing `onnxruntime_go`. The real ONNX embedder is now behind `onnx` build tag with a default stub.

## Known Limitations

Dedicated ONNX Docker image, CI artifact provisioning, and packaged quality/performance reruns are not implemented in M021.

## Follow-ups

Next milestone should implement a dedicated ONNX image/CI artifact provisioning path using the verifier, explicit `onnx hf_tokenizers` tags, and external artifact supply. Packaged environment must rerun quality and performance gates before any promotion.

## Files Created/Modified

- `api/embed/onnx.go` — Real ONNX embedder moved behind `onnx` build tag.
- `api/embed/onnx_disabled.go` — Default ONNX stub for non-ONNX builds.
- `api/embed/onnx_types.go` — Shared ONNX options for default and ONNX builds.
- `api/embed/onnx_token_types.go` — Token encoding/interface types scoped to ONNX/HF-tokenizer builds.
- `api/embed/onnx_tokenizer_default.go` — Tokenizer build tags updated for ONNX backend.
- `api/embed/onnx_tokenizer_hf.go` — HF tokenizer ONNX bridge now requires ONNX and native tokenizer tags.
- `api/embed/onnx_test.go` — ONNX embedder tests now require ONNX build tag.
- `.gsd/DECISIONS.md` — Decision register updated with D019.
