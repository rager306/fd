---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Prepare benchmark command plan

Inspect benchmark.py options and existing M014 artifacts to build the exact ONNX 1024 benchmark command with isolated namespace and restart command.

## Inputs

- `benchmark.py`
- `benchmark-results/fd-benchmark-m014-tei-baseline.txt`
- `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`

## Expected Output

- `Task summary with benchmark command plan`

## Verification

Command plan identifies API_URL, output artifact, runtime label, namespace, and restart command.

## Observability Impact

Documents benchmark settings before running expensive performance work.
