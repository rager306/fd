---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Inspect local model artifact storage

Inspect local TEI/Docker model storage and project config to determine what model artifacts are already present locally and where large ONNX artifacts would live without being committed.

## Inputs

- `.gsd/milestones/M009-zjrq6j/M009-zjrq6j-SUMMARY.md`

## Expected Output

- `.gsd/milestones/M010-84qfzu/slices/S01/tasks/T01-SUMMARY.md`

## Verification

Local Docker volumes/config inspected; no large artifacts staged.

## Observability Impact

Records local artifact directories and gitignore safety for large files.
