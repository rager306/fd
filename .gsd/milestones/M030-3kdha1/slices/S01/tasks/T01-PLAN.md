---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Harden Go manifest artifact path policy

Implement Go ONNX manifest path policy and safe path display: restrict `artifact.local_path` to approved repo-relative roots; adjust errors to avoid absolute paths where possible; add tests for allowed and rejected paths.

## Inputs

- `.gsd/milestones/M028-y63tog/slices/S01/S01-RESEARCH.md`
- `api/embed/onnx_manifest.go`
- `api/embed/onnx_manifest_test.go`

## Expected Output

- `api/embed/onnx_manifest.go`
- `api/embed/onnx_manifest_test.go`

## Verification

Targeted manifest tests pass.

## Observability Impact

Startup errors become safer while preserving artifact labels.
