---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T05: Final perf validation + concurrent load test

tools/verify_fd_v2_perf.sh: запускает T-P-1..T-P-5 против running fd v2. Проверяет: 1 input p95<50ms, 10 input p95<200ms, 32 input p95<1000ms, 100 sequential 0 errors, 4 concurrent × 8 input < 2s total. Также проверяет cache effectiveness: cache hit latency < 5ms, eviction работает. Результаты в benchmark-results/fd-v2-perf-validation-m041-s04.md. Спека: docs/fd-v2.md Section 5.4 T-P-1..T-P-5 + Section 6.3 F-4 (cache miss → hit).

## Inputs

- None specified.

## Expected Output

- `tools/verify_fd_v2_perf.sh`
- `benchmark-results/fd-v2-perf-validation-m041-s04.md`

## Verification

tools/verify_fd_v2_perf.sh exit 0, benchmark-results artifact содержит p50/p95/p99 для всех 5 test cases, все 4×8 concurrent requests complete < 2s.
