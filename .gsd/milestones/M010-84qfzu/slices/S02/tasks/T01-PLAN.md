---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Define dense comparator contract and probes

Define the minimal comparator contract and probe set. Use non-sensitive Russian/legal-style fixed probes, assign stable labels, and document expected output fields: API URL, model, dimensions, probe label/length, finite values, L2 norm, vector hash, and cosine similarities. Ensure raw texts will not be printed in benchmark artifacts.

## Inputs

- `.gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md`
- `benchmark.py`

## Expected Output

- `Comparator contract captured in task summary and/or script constants`

## Verification

Review comparator contract/probes; confirm raw texts are not emitted by design.

## Observability Impact

Defines the sanitized fields that make future ONNX comparison debuggable without leaking probe text.
