---
id: M037-d23oz4
title: "Target runtime validation contract"
status: complete
completed_at: 2026-05-21T10:23:35.669Z
key_decisions:
  - D035 — Python helper evidence is setup/provenance only; target-runtime acceptance is required for Go and any future Rust runtime.
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-target-runtime-validation-contract-m037-s01.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - Python scripts can drive evidence collection only when pointed at actual runtime endpoints; Python-only metadata/provisioning checks do not prove deployed runtime behavior.
  - Future runtime implementations need independent evidence because tokenizer, linking, provider, normalization, and packaging failure modes are runtime-specific.
---

# M037-d23oz4: Target runtime validation contract

**Defined the target-runtime validation boundary for ONNX artifacts across Go and future Rust runtimes.**

## What Happened

M037 codified the target-runtime validation boundary the user raised. The manifest and docs now say Python helper checks are setup/provenance evidence only. Current ONNX acceptance requires Go fd API and packaged Go Docker runtime gates, including tokenizer parity, embedding API behavior, health metadata, legal quality, performance, package validation, and Redis namespace isolation. Any future Rust backend must independently pass equivalent gates. D035 records this collaborative decision. No new artifact was generated and no external actions occurred.

## Success Criteria Results

- PASS — Python checks are no longer framed as production runtime acceptance.
- PASS — Go target-runtime gates are explicit.
- PASS — Future Rust gates are explicit.
- PASS — TEI remains production/default and ONNX remains opt-in experimental.

## Definition of Done Results

- Done — Target-runtime boundary documented.
- Done — Manifest records Go and future Rust validation gates.
- Done — Docs/outcome require Go API/legal/performance/package validation before acceptance.
- Done — Final guardrails passed.
- Done — No external action or production switch occurred.

## Requirement Outcomes

New implicit requirement surfaced: ONNX artifact acceptance must validate the intended target runtime implementation, not only Python helpers.

## Deviations

None.

## Follow-ups

Next runtime-focused gate can implement or run a consolidated Go target-runtime acceptance harness for the existing artifact or a regenerated candidate. If Rust is introduced, create a separate Rust target-runtime acceptance milestone.
