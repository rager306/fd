---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Persist source contract in docs and manifests

Update ONNX artifact manifests and provisioning docs with pinned source candidates, ONNX model blocker, and runtime/source checksum details without changing TEI/default runtime behavior.

## Inputs

- `.gsd/milestones/M031-gn517a/slices/S01/S01-RESEARCH.md`

## Expected Output

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/PROVISIONING.md`

## Verification

JSON parses; docs contain source statuses and no production promotion language.

## Observability Impact

Makes source status visible in durable docs/manifests.
