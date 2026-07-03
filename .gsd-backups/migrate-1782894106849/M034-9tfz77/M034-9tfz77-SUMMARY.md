---
id: M034-9tfz77
title: "Hosted ONNX workflow input alignment"
status: complete
completed_at: 2026-05-21T08:00:36.108Z
key_decisions:
  - D032 — ONNX Runtime sha workflow input is optional override; manifest source_contract sha is default.
key_files:
  - .github/workflows/onnx-packaging.yml
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-workflow-input-alignment-m034-s02.txt
lessons_learned:
  - Workflow inputs should not duplicate tracked manifest checksums unless overriding them intentionally.
  - Documenting explicit external-action boundaries prevents workflow skeletons from being mistaken for hosted proof.
---

# M034-9tfz77: Hosted ONNX workflow input alignment

**Aligned the manual ONNX workflow input contract with manifest-derived runtime checksum provisioning.**

## What Happened

M034 aligned the manual hosted ONNX packaging workflow with the provisioning capabilities added in M033. The runtime sha input is now optional; when omitted, provisioning uses the tracked manifest's ONNX Runtime sha. The docs and outcome define the safe input contract and preserve the exact ONNX model source blocker. No push, dispatch, upload, or production promotion occurred.

## Success Criteria Results

- PASS — M033 runtime wheel provisioning usable from workflow input surface.
- PASS — Runtime checksums still verified.
- PASS — Exact ONNX model blocker explicit.
- PASS — No external state changes.
- PASS — Final verification passed; commit follows DB checkpoint.

## Definition of Done Results

- Done — Workflow input validation aligned with M033.
- Done — Workflow passes runtime sha override only when supplied.
- Done — Docs/outcome define safe input contract.
- Done — actionlint and full guardrails passed.
- Done — No external action or production switch occurred.

## Requirement Outcomes

No formal requirement IDs updated. The milestone advances hosted workflow readiness while preserving blockers.

## Deviations

None.

## Follow-ups

Next recommended milestone: exact ONNX model binary source. Either define exact-binary immutable hosting/mirroring inputs or design reproducible-export workflow. Any push/workflow dispatch/upload still requires explicit user approval.
