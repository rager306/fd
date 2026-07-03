---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Write artifact provisioning contract

Write artifact provisioning contract documenting required artifact sources, destination paths, checksum policy, cache layout, and current blockers.

## Inputs

- `docs/onnx-artifacts/README.md`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`

## Expected Output

- `docs/onnx-artifacts/PROVISIONING.md`

## Verification

Contract exists and states blockers without raw/secrets/binaries.

## Observability Impact

Makes missing external artifact source explicit and discoverable.
