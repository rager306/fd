---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Research NUMA and threading controls

Research NUMA/threading tuning for ONNX Runtime CPU inference and Linux VPS deployment: intra/inter op threads, affinity, OMP/MKL/ORT env vars, numactl, CPU pinning, and how benchmark config snapshots should record them.

## Inputs

- `.gsd/REQUIREMENTS.md`
- `.gsd/DECISIONS.md`

## Expected Output

- `S06 T02 summary`

## Verification

Threading/NUMA knobs and benchmark controls recorded.

## Observability Impact

Defines fair benchmark controls for EPYC/NUMA and thread tuning.
