---
id: M018-vq2ttb
title: "ONNX 1024 legal quality gate"
status: complete
completed_at: 2026-05-20T07:42:43.128Z
key_decisions:
  - D016: tagged Go ONNX with max sequence length 1024 passes the selected Russian/legal quality gate, but remains experimental until performance, memory, packaging, CI, and operational gates pass.
key_files:
  - benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt
  - benchmark-results/fd-onnx-1024-outcome-m018-s02.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - The strict legal quality blocker can be retired for the selected corpus by using 1024 max sequence length.
  - The next risk moved from quality to runtime operationalization and performance.
  - Longer sequence length can be validated before implementing chunking, avoiding unnecessary complexity when it passes the target gate.
---

# M018-vq2ttb: ONNX 1024 legal quality gate

**M018 proved tagged Go ONNX 1024 passes the selected Russian/legal quality gate while keeping ONNX experimental.**

## What Happened

M018 measured whether tagged Go ONNX with max sequence length 1024 can pass the Russian legal quality gate. S01 started the tagged Go ONNX service with native HF tokenizer, `ONNX_MAX_SEQUENCE_LENGTH=1024`, and isolated namespace `m018-onnx-1024-legal-quality`, then ran the full legal retrieval evaluator. The result passed strict thresholds: minimum cross-backend cosine 0.99989883, top-1 agreement 1.0, mean overlap@5 0.997701, and ONNX recall ratio 1.0. S02 interpreted this outcome and recorded D016: ONNX 1024 passes selected legal quality but remains experimental until performance, memory, packaging, CI, and operational gates pass. TEI remains the production/default runtime.

## Success Criteria Results

- Tagged Go ONNX 1024 gate: PASS.
- Metrics compared: PASS.
- Raw text hygiene: PASS.
- TEI default preserved: PASS.
- Next gate explicit: PASS.

## Definition of Done Results

- Tagged Go ONNX 1024 legal gate was run: met.
- Isolated cache namespace was used: met (`m018-onnx-1024-legal-quality`).
- Raw legal text was excluded from artifacts: met by hygiene checks.
- TEI remains production/default: met.
- ONNX remains experimental: met.
- Runtime cleanup verified: met.
- Fresh tests/lint/tagged checks passed: met.

## Requirement Outcomes

- ONNX 1024 legal quality evidence: produced and validated as PASS.
- ONNX production readiness: remains blocked by performance, memory, packaging, CI, and operational gates.
- Long-text remediation: sequence length 1024 is sufficient for the selected corpus gate; chunking remains future policy for unbounded documents.
- Cache isolation requirement: validated in artifact configuration.

## Deviations

None.

## Follow-ups

Create the next GSD milestone for ONNX 1024 performance, memory, artifact contract, Docker/CI packaging, and operational diagnostics. Do not promote ONNX until those gates pass. Chunking remains future policy for unbounded legal documents beyond 1024 tokens.
