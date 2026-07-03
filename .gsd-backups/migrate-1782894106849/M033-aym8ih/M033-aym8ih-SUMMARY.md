---
id: M033-aym8ih
title: "ONNX Runtime wheel provisioning"
status: complete
completed_at: 2026-05-21T07:25:27.871Z
key_decisions:
  - D031 — ONNX Runtime wheel/zip provisioning extracts the configured member with size/sha verification and preserves direct-file fallback.
key_files:
  - tools/provision_onnx_artifacts.py
  - docs/onnx-artifacts/PROVISIONING.md
  - benchmark-results/fd-onnx-runtime-wheel-provisioning-m033-s02.txt
lessons_learned:
  - Source contracts must be executable by provisioning tooling, not only documented.
  - Zip member validation should reject explicit unsafe POSIX file types while tolerating normal entries without POSIX metadata.
---

# M033-aym8ih: ONNX Runtime wheel provisioning

**Made ONNX Runtime wheel sources safely provisionable without changing TEI defaults or ONNX rollout status.**

## What Happened

M033 made the previously selected PyPI ONNX Runtime wheel source candidate usable by provisioning tooling. The helper now reads ONNX Runtime member/size/sha metadata from the ONNX manifest source contract, extracts configured members from `.whl`/`.zip` sources safely, rejects unsafe/missing/oversized/wrong-checksum cases, and preserves direct shared-library fallback. The work is documented and verified while leaving exact ONNX model source, hosted proof, and production promotion blocked.

## Success Criteria Results

- PASS — Runtime wheel candidate actionable.
- PASS — Security hardening retained.
- PASS — Exact ONNX model binary blocker explicit.
- PASS — No external state changes.
- PASS — Final verification passed; commit follows DB checkpoint.

## Definition of Done Results

- Done — ONNX Runtime wheel extraction implemented.
- Done — Existing direct-file and native tokenizer behavior preserved.
- Done — Docs/outcome/decision recorded.
- Done — Positive and negative synthetic probes passed.
- Done — Guardrails passed.
- Done — No external actions or production switch occurred.

## Requirement Outcomes

No formal requirement IDs updated. The milestone advances operational source provisioning readiness.

## Deviations

Zip regular-file detection was corrected during final verification to reject symlinks/non-regular POSIX entries while tolerating normal zip entries with no POSIX file type. Verification was rerun after the change.

## Follow-ups

Next recommended milestone: exact ONNX model binary source. Choose either an exact-binary immutable hosting plan or a full reproducible-export workflow; both remain blocked from external action until explicit user confirmation.
