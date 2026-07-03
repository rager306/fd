---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Record 512 outcome decision

Record a GSD decision that 512-token ONNX is necessary but insufficient for strict legal equivalence, so the next implementation gate must add chunking or longer sequence handling.

## Inputs

- `benchmark-results/fd-onnx-512-outcome-m017-s02.txt`

## Expected Output

- `.gsd/DECISIONS.md`

## Verification

Decision saved through GSD.

## Observability Impact

Prevents future agents from treating 512 ranking parity as production readiness.
