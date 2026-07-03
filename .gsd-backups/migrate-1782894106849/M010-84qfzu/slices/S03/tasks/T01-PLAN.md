---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Prepare ONNX artifact workspace and provenance

Prepare `.gsd/runtime/onnx/m010-s03/` as the local ignored artifact workspace and capture source model provenance: model path, revision, source hashes, tokenizer/config hashes, available disk, and dependency plan. Verify no large artifacts are staged.

## Inputs

- `.gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md`
- `benchmark-results/fd-dense-comparator-m010-s02.txt`
- `tei-models/deepvk--USER-bge-m3/`

## Expected Output

- `.gsd/runtime/onnx/m010-s03/metadata.json or task summary provenance`

## Verification

Workspace exists under ignored `.gsd/runtime`; source hashes recorded; `git status --short` shows no large ONNX/model artifacts staged.

## Observability Impact

Creates the local provenance root for export/load evidence without touching production runtime.
