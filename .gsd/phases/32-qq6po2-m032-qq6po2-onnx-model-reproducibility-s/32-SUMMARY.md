---
id: M032-qq6po2
title: "ONNX model reproducibility source proof"
status: complete
completed_at: 2026-05-21T07:03:19.679Z
key_decisions:
  - D030 — Use bounded local existing-artifact verifier; do not treat it as regenerated export reproducibility proof.
key_files:
  - tools/verify_onnx_export_contract.py
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - benchmark-results/fd-onnx-export-contract-verifier-m032-s02.txt
  - .gsd/milestones/M032-qq6po2/slices/S01/S01-RESEARCH.md
lessons_learned:
  - Reproducibility language must distinguish existing-artifact contract verification from regenerated-export proof.
  - Machine-readable `claim_scope` helps prevent future overclaims.
---

# M032-qq6po2: ONNX model reproducibility source proof

**Added and documented a bounded local ONNX export contract verifier without promoting ONNX.**

## What Happened

M032 converted the remaining exact ONNX model blocker into an executable local proof surface. The new verifier checks the existing ignored ONNX artifact against tracked manifest, source provenance, export metadata, and toolchain pins. Negative probes show it catches artifact checksum, model revision, and package drift. Docs, manifest metadata, outcome, and D030 now make the proof boundary explicit: this is not a fresh regenerated export proof, and hosted packaging still requires exact-binary immutable hosting or a separate reproducible-export workflow.

## Success Criteria Results

- PASS — Existing ONNX export contract verifier available.
- PASS — Claim bounded and machine-readable.
- PASS — TEI default/ONNX opt-in unchanged.
- PASS — No binaries committed.
- PASS — Local verification passed; commit follows DB checkpoint.

## Definition of Done Results

- Done — Local verifier added and verified.
- Done — Proof boundary documented.
- Done — Outcome and decision recorded.
- Done — Guardrails passed.
- Done — TEI default and ONNX opt-in status unchanged.
- Done — No push/workflow dispatch/upload occurred.

## Requirement Outcomes

No formal requirement IDs updated. The milestone advances operational/source readiness while preserving production blockers.

## Deviations

Verifier initially expected Python inside export metadata packages; corrected to compare the top-level Python runtime string separately. No plan-invalidating deviations.

## Follow-ups

Choose one next gate: (1) exact-binary hosting/mirroring of the current ONNX artifact to a non-secret immutable source, requiring external destination and explicit user approval before upload/push; or (2) reproducible-export workflow that regenerates the ONNX binary from pinned source/toolchain and reruns legal/performance/package gates.
