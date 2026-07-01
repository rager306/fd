---
id: M021-4t2wpt
title: "ONNX 1024 Docker CI packaging contract"
status: complete
completed_at: 2026-05-20T10:27:22.318Z
key_decisions:
  - D019: default Docker/CI remains CGO-disabled TEI path; opt-in ONNX runtime requires explicit `onnx hf_tokenizers` tags and verified artifacts.
key_files:
  - tools/verify_onnx_artifacts.py
  - docs/onnx-artifacts/README.md
  - api/embed/onnx.go
  - api/embed/onnx_disabled.go
  - api/embed/onnx_types.go
  - api/embed/onnx_token_types.go
  - api/embed/onnx_tokenizer_default.go
  - api/embed/onnx_tokenizer_hf.go
  - api/embed/onnx_test.go
  - .gsd/DECISIONS.md
lessons_learned:
  - Default Docker build is the right guardrail for TEI-default production safety.
  - ONNX runtime and native tokenizer parity are separate build-tag concerns: `onnx` for runtime, `hf_tokenizers` for native tokenizer binding.
  - Artifact verification should happen before any tagged ONNX Docker/CI run, but default CI can remain artifact-free.
---

# M021-4t2wpt: ONNX 1024 Docker CI packaging contract

**M021 added ONNX artifact verification and fixed the build-tag boundary so default Docker stays TEI-safe while ONNX remains explicit opt-in.**

## What Happened

M021 created a safe artifact provisioning contract and validated Docker/CI boundaries for ONNX 1024. S01 added `tools/verify_onnx_artifacts.py` and `docs/onnx-artifacts/README.md`, proving local ignored ONNX and native tokenizer artifacts match tracked manifests and remain untracked. S2 then validated default Docker behavior and discovered a real regression: the default CGO-disabled image build imported `onnxruntime_go`. The fix moved the real ONNX embedder behind an explicit `onnx` build tag, added a default stub, scoped token types appropriately, and updated ONNX runtime documentation to use `onnx hf_tokenizers`. Final verification passed across artifact checks, default tests, lint, tagged tests, ONNX+native smoke tests, default Docker build, binary hygiene, and runtime cleanup. TEI remains production/default and ONNX remains opt-in experimental.

## Success Criteria Results

- Artifact contract: PASS.
- Default Docker boundary: PASS.
- Binary hygiene: PASS.
- Verification: PASS.
- Next gate explicit: PASS.

## Definition of Done Results

- Artifact verification contract exists: met.
- Default Docker build works without ONNX/native artifacts: met.
- ONNX runtime is opt-in behind build tags: met.
- No binaries tracked: met.
- Fresh verification passed: met.
- Next Docker/CI provisioning gate explicit: met.

## Requirement Outcomes

- ONNX artifact provisioning contract: validated.
- Default Docker/CI safety: validated and fixed.
- ONNX production readiness: remains blocked by dedicated artifact provisioning and packaged quality/performance gates.
- TEI default: preserved.

## Deviations

M021 uncovered a default Docker build regression and fixed it by adding an explicit `onnx` build tag boundary. This expanded S02 from validation-only into a small packaging-boundary code fix.

## Follow-ups

Create a dedicated ONNX Docker/CI artifact provisioning milestone: add an ONNX image target or CI job that provisions ONNX/native artifacts externally, verifies checksums with `tools/verify_onnx_artifacts.py`, builds with `onnx hf_tokenizers`, runs tagged tests, and reruns packaged quality/performance gates.
