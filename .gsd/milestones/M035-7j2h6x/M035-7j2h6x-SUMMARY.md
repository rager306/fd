---
id: M035-7j2h6x
title: "Exact ONNX binary hosting contract"
status: complete
completed_at: 2026-05-21T09:27:33.833Z
key_decisions:
  - D033 — Represent exact ONNX binary source blocker as planned_not_uploaded hosting contract while keeping source_status blocked and no source_url.
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - A planned artifact object key must be explicitly marked not uploaded and not usable as a source until real hosting exists.
  - GSD milestone completion/commit/reindex should remain post-slice actions, not required task work inside the final slice.
---

# M035-7j2h6x: Exact ONNX binary hosting contract

**Defined the planned exact ONNX binary hosting contract while preserving the source blocker and external-action boundary.**

## What Happened

M035 converted the remaining exact ONNX model binary source blocker into a precise hosting contract. The manifest now includes a planned_not_uploaded hosting contract with exact size/sha, recommended object key and release filename, allowed/forbidden source forms, and pre-dispatch requirements while preserving source_status=blocked and no source_url. Docs and outcome artifacts mirror the contract, D033 records the decision, and final guardrails passed. No external actions occurred.

## Success Criteria Results

- PASS — Exact ONNX model binary blocker became an actionable hosting contract.
- PASS — No fake URL or unavailable source was marked immutable_selected.
- PASS — Workflow dispatch remains blocked until actual source exists and user approves.
- PASS — TEI remains production/default and ONNX remains opt-in experimental.

## Definition of Done Results

- Done — Exact ONNX binary hosting contract documented.
- Done — Manifest keeps source blocked and no source_url.
- Done — Outcome and D033 recorded.
- Done — No upload, push, workflow dispatch, signed URL, or production switch.
- Done — Final guardrails passed.

## Requirement Outcomes

No formal requirement IDs updated. The source blocker is now actionable but not resolved.

## Deviations

S02 T03 closure/commit work was recorded as an ordering correction because GSD slice completion must precede milestone completion and commit/reindex. The actual completion/checkpoint/commit/reindex is performed after milestone completion.

## Follow-ups

Next gate can either perform an explicitly approved external upload/mirroring step for the exact binary, or design a reproducible-export workflow. Any push, upload, or GitHub workflow dispatch requires explicit user approval.
