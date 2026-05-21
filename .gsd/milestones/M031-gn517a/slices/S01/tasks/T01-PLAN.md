---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T01: Inventory artifact source requirements

Inventory current tracked manifests, local provenance, provisioning docs, and workflow requirements into an artifact source matrix with checksums/sizes and current blockers.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/PROVISIONING.md`
- `.gsd/runtime/onnx/m010-s03/source-provenance.json`
- `.gsd/runtime/onnx/m010-s03/export-metadata.json`

## Expected Output

- `.gsd/milestones/M031-gn517a/slices/S01/S01-RESEARCH.md`

## Verification

Inventory artifact contains all four required artifacts and exact local checksums.

## Observability Impact

Creates source inventory for future workflow proof.
