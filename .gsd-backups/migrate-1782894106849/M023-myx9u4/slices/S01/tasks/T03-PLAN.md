---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Validate artifact hygiene and cleanup

Validate legal artifact hygiene and clean packaged runtime: check no raw corpus lines leak, stop ONNX container, confirm port 18000 clean, and record pass/fail outcome.

## Inputs

- `tests/44-FZ-2026-articles.jsonl`

## Expected Output

- `Task summary with hygiene and cleanup evidence`

## Verification

Raw text leak check passes; no background processes; port 18000 clean.

## Observability Impact

Ensures sensitive legal text remains out of artifact and runtime cleanup is explicit.
