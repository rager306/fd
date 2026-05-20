---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Inspect artifact metadata requirements

Inspect M010 export metadata, current gitignore/runtime paths, README runtime docs, and any existing artifact conventions. Determine the smallest manifest shape needed for S02/S03 without committing the 1.43GB ONNX file.

## Inputs

- `.gsd/milestones/M010-84qfzu/M010-84qfzu-SUMMARY.md`
- `.gsd/milestones/M010-84qfzu/slices/S04/S04-RESEARCH.md`

## Expected Output

- `Task summary with manifest field list and storage decision`

## Verification

Required manifest fields listed; confirms no large artifact will be tracked.

## Observability Impact

Identifies which fields future startup/load errors must report.
