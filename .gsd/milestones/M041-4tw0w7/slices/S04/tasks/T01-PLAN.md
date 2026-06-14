---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Baseline measurement fd perf захвачен до S04 T04 perf optimization

tools/measure_fd_baseline.sh: простой benchmark который шлёт 1/10/32 input запросы 100 раз каждый, меряет p50/p95/p99 latency, error rate. Запускается против текущего fd (после S01/S02/S03) для baseline numbers. Сохранить в benchmark-results/fd-v2-baseline-before-m041-s04.md. Это даст опорные цифры чтобы понять, нужен ли real optimization или достаточно validation fixes из S01. Спека target values: docs/fd-v2.md Section 5.4 T-P-1..T-P-5.

## Inputs

- None specified.

## Expected Output

- `tools/measure_fd_baseline.sh`
- `benchmark-results/fd-v2-baseline-before-m041-s04.md`

## Verification

Baseline artifact содержит: p50/p95/p99 для batch=1/10/32, error rate, throughput. Можно сравнить с target values (50/200/1000ms).
