---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Write remediation assessment

Write a remediation assessment artifact comparing longer ONNX max sequence length, explicit chunking, and longer-sequence export/runtime options using S01/S02 evidence.

## Inputs

- `benchmark-results/fd-legal-divergence-profile-m016-s01.txt`
- `benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt`
- `benchmark-results/fd-legal-retrieval-m015-s03.txt`

## Expected Output

- `benchmark-results/fd-onnx-remediation-plan-m016-s03.txt`

## Verification

Artifact exists, includes S01/S02 metrics, and contains no raw legal text.

## Observability Impact

Creates a durable decision artifact for future agents implementing ONNX remediation.
