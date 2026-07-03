---
id: M031-gn517a
title: "ONNX artifact source contract"
status: complete
completed_at: 2026-05-21T06:42:15.990Z
key_decisions:
  - D029 — Select pinned checksum-matched supporting artifact candidates while keeping exported ONNX model binary blocked.
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
  - docs/onnx-artifacts/PROVISIONING.md
  - benchmark-results/fd-onnx-source-contract-m031-s02.txt
  - .gsd/milestones/M031-gn517a/slices/S01/S01-RESEARCH.md
lessons_learned:
  - Separate `immutable_selected` source candidates from hosted workflow rollout proof.
  - Do not treat upstream HF model files as equivalent to the exported ONNX model binary.
---

# M031-gn517a: ONNX artifact source contract

**Defined and persisted the ONNX artifact source contract: three pinned supporting candidates and one explicit ONNX model binary blocker.**

## What Happened

M031 turned the post-M030 artifact-source blocker into a concrete source contract. S01 inventoried all required artifacts and verified candidate sources for native tokenizer, tokenizer JSON, and ONNX Runtime. S02 persisted the source contract in tracked manifests/docs, recorded an outcome artifact and D029, and ran the final guardrails. The exact exported ONNX binary remains the remaining artifact-source blocker, so hosted packaging proof and production promotion remain blocked.

## Success Criteria Results

- PASS — Source contract truthful and no overclaiming.
- PASS — Security policies preserved.
- PASS — TEI default and ONNX opt-in status unchanged.
- PASS — No external state changes.
- PASS — Verification passed; local commit follows GSD DB checkpoint.

## Definition of Done Results

- Done — Artifact source statuses are explicit.
- Done — Supporting artifacts have pinned checksum-matched candidates.
- Done — ONNX model binary remains blocked and not overclaimed.
- Done — Docs/manifests/outcome/decision updated.
- Done — Verification passed.
- Done — No push/workflow dispatch/production switch occurred.

## Requirement Outcomes

No formal requirement IDs were updated. The milestone advances operational rollout readiness but intentionally does not validate ONNX production readiness.

## Deviations

Initial verification-command mistakes were corrected and rerun: provisioning dry-run does not accept `--allow-missing`, and binary hygiene must use tracked files rather than ignored runtime caches.

## Follow-ups

Next recommended milestone: create an immutable source for the exact ONNX model binary, either by mirroring/uploading the current binary to a non-secret immutable source or by designing a reproducible-export workflow that regenerates and revalidates the artifact. After that, run hosted workflow proof only after explicit push approval.
