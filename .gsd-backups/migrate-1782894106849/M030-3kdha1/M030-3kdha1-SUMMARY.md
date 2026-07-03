---
id: M030-3kdha1
title: "ONNX artifact path security remediation"
status: complete
completed_at: 2026-05-21T05:36:47.552Z
key_decisions:
  - D028: M028 LOW findings are remediated by M030 for default tooling/startup behavior; immutable sources and hosted proof remain blockers.
key_files:
  - api/embed/onnx_manifest.go
  - api/embed/onnx_manifest_test.go
  - api/main_test.go
  - tools/provision_onnx_artifacts.py
  - tools/verify_onnx_artifacts.py
  - tools/build_onnx_image.sh
  - docs/onnx-artifacts/PROVISIONING.md
  - benchmark-results/fd-onnx-path-security-remediation-m030-s02.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - Security root policies should be reflected in test fixtures so tests exercise the real approved artifact layout.
  - Approved artifact roots are now a shared contract between Go startup validation and Python provisioning/verifier tooling.
  - Safer default diagnostics can reduce host path disclosure without removing artifact labels or repo-relative context.
---

# M030-3kdha1: ONNX artifact path security remediation

**M030 remediated ONNX artifact path root-policy and default path-disclosure findings.**

## What Happened

M030 remediated the remaining M028 low-severity ONNX artifact path findings. Go manifest validation now rejects absolute, traversal, and unapproved-root artifact paths while preserving existing `.gsd/runtime/onnx/...` layout. Python provisioning and verifier enforce the same approved-root policy and default to safer path display. The Docker packaging script no longer prints full configured tokenizer/runtime-library paths on missing-file errors. Documentation and outcome artifacts record approved roots, remediated findings, evidence, and remaining rollout blockers. TEI remains default and ONNX remains opt-in experimental.

## Success Criteria Results

- M028 LOW-3 remediation: PASS.
- M028 LOW-4 remediation: PASS.
- Guardrails: PASS.
- Production safety: PASS (no switch).

## Definition of Done Results

- Manifest artifact path approved-root policy: met.
- Safer default diagnostics: met.
- Tests/probes: met.
- Docs/outcome/decision: met.
- Production/default switch: not performed.

## Requirement Outcomes

- ONNX path security: validated.
- M028 security remediation: medium and low findings are now addressed for default tooling/startup behavior.
- Immutable artifact source selection: still pending.
- Hosted workflow proof: still pending.
- ONNX production rollout: still blocked.

## Deviations

Initial default Go test run failed due to a test fixture still using arbitrary temp artifact paths. The fixture was corrected to use the approved `.gsd/runtime/onnx/...` layout and all rerun checks passed.

## Follow-ups

Next recommended gate: immutable artifact source selection/pinning for ONNX model, native tokenizer, tokenizer JSON, and ONNX Runtime, followed by hosted workflow proof only after explicit push approval and real sources.
