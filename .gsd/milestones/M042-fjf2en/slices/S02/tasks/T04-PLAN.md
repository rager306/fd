---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T04: Perf benchmark: async vs sync cold/warm path

tools/verify_fd_async_perf.sh: прогон cold path measurements × batch sizes × async on/off. Output benchmark-results/fd-v2-async-perf-m042.md с table (batch 1/10/32/64/128 cold, sync vs async) и conclusion (improvement factor, where it falls short of 1000ms target). Также test concurrent scenario: 4 parallel fd calls × batch=32 cold, sync vs async (per M041 T-P-5 spec).

## Inputs

- None specified.

## Expected Output

- `tools/verify_fd_async_perf.sh`
- `benchmark-results/fd-v2-async-perf-m042.md`

## Verification

tools/verify_fd_async_perf.sh exit 0. Artifact содержит: cold path table (batch × async on/off), concurrent test results, conclusion с comparison vs M041 baseline (25s → ≤10s for batch=128 cold).
