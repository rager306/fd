---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Record exact binary hosting decision

Record GSD decision for exact ONNX binary hosting contract and update outcome if needed to reference decision.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt`

## Expected Output

- `.gsd/DECISIONS.md`
- `benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt`

## Verification

Decision and outcome checks pass.

## Observability Impact

Durable decision keeps future agents from treating the planned key as an available source.
