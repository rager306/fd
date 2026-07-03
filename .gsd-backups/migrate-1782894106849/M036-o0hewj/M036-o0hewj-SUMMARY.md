---
id: M036-o0hewj
title: "ONNX reproducible export workflow contract"
status: complete
completed_at: 2026-05-21T09:41:45.967Z
key_decisions:
  - D034 — Represent the no-upload source resolution path as a planned reproducible-export workflow contract, not proof.
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-reproducible-export-contract-m036-s01.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - A regenerated-export plan only becomes source proof after actual regeneration plus downstream quality/performance/package evidence.
  - M032's verifier must stay framed as existing-artifact contract verification unless it is extended to execute a fresh export.
---

# M036-o0hewj: ONNX reproducible export workflow contract

**Defined the planned no-upload reproducible-export workflow contract without claiming proof.**

## What Happened

M036 defined the no-upload reproducible-export alternative for resolving the ONNX model source blocker. The manifest now records a planned_not_proven workflow contract with pinned model/source/toolchain inputs, expected artifact contract, acceptance gates, success/failure interpretations, and forbidden claims. Docs and outcome artifacts mirror the boundary, D034 records the decision, and final guardrails passed. No export was regenerated, no external actions occurred, and ONNX remains opt-in experimental.

## Success Criteria Results

- PASS — No-upload reproducible-export alternative is actionable and truthful.
- PASS — No byte-for-byte regenerated proof is claimed.
- PASS — Exact-binary hosting remains separate and blocked until approval.
- PASS — TEI remains production/default and ONNX remains opt-in experimental.

## Definition of Done Results

- Done — Reproducible-export path documented as planned, not proof.
- Done — Manifest records planned_not_proven contract with pinned inputs and gates.
- Done — Outcome/docs make clear no export was regenerated and no hosted workflow ran.
- Done — Final guardrails passed.
- Done — No external action or production switch occurred.

## Requirement Outcomes

No formal requirement IDs updated. The no-upload path is now actionable but not executed.

## Deviations

None.

## Follow-ups

Next gate can implement a local regenerated-export dry-run/workflow in an ignored workspace, but it must rerun fixed-probe, legal quality, performance, packaging, and hosted proof gates before resolving the blocker. Any push/upload/workflow dispatch still requires explicit user approval.
