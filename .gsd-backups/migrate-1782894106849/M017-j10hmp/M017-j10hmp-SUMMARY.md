---
id: M017-j10hmp
title: "ONNX 512 legal quality gate"
status: complete
completed_at: 2026-05-20T07:32:18.356Z
key_decisions:
  - D015: 512-token ONNX is necessary but insufficient; next remediation must handle >512-token legal fragments before production promotion.
key_files:
  - benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt
  - benchmark-results/fd-onnx-512-outcome-m017-s02.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - A gate can successfully complete with a measured FAIL when the objective is risk retirement rather than promotion.
  - 512-token ONNX restores ranking parity but not strict vector equivalence for long legal text.
  - Future ONNX legal quality work must treat chunking/long-sequence policy as correctness, not optimization.
---

# M017-j10hmp: ONNX 512 legal quality gate

**M017 proved tagged Go ONNX 512 greatly improves ranking parity but still fails strict legal vector equivalence, requiring chunking or longer-sequence remediation next.**

## What Happened

M017 validated the actual tagged Go ONNX runtime at 512 tokens on the Russian legal corpus. S01 started the tagged Go ONNX service with native HF tokenizer, `ONNX_MAX_SEQUENCE_LENGTH=512`, and isolated namespace `m017-onnx-512-legal-quality`, then ran the full legal retrieval evaluator. The result was a measured strict FAIL: minimum cross-backend cosine was 0.98982302 below the 0.999 threshold. However, ranking parity was excellent: top-1 agreement 1.0, mean overlap@5 0.997701, and ONNX recall ratio 1.0. S02 interpreted this outcome and recorded D015: 512-token ONNX is necessary but insufficient; the next remediation must handle legal fragments above 512 tokens through deterministic chunking or longer-sequence handling. TEI remains the production/default runtime and ONNX remains experimental.

## Success Criteria Results

- Tagged Go ONNX 512 gate: PASS as executed measurement.
- Metrics compared: PASS.
- Raw text hygiene: PASS.
- TEI default preserved: PASS.
- Next remediation explicit: PASS.

## Definition of Done Results

- Tagged Go ONNX 512 legal gate was run: met.
- Isolated cache namespace was used: met (`m017-onnx-512-legal-quality`).
- Raw legal text was excluded from artifacts: met by hygiene checks.
- TEI remains production/default: met.
- ONNX remains experimental: met.
- Runtime cleanup verified: met.
- Fresh tests/lint/tagged checks passed: met.

## Requirement Outcomes

- ONNX 512 quality evidence: produced and validated as measured strict FAIL.
- ONNX production readiness: remains blocked.
- Long-text remediation requirement: refined to chunking or longer sequence for >512-token legal fragments.
- Cache isolation requirement: validated in artifact configuration.

## Deviations

The evaluator returned exit code 2 because the measured quality verdict was FAIL. This is not a milestone failure; the purpose of M017 was to measure whether 512 passes strict legal quality, and the answer is no.

## Follow-ups

Create the next GSD milestone for deterministic chunking or longer-sequence handling. Recommended path: 512-token ONNX baseline plus chunking for fragments above 512 tokens, include chunking version in cache namespace, rerun full legal retrieval gate, then benchmark performance only after quality passes.
