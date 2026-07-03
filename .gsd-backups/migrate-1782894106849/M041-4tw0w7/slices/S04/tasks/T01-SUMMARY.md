---
id: T01
parent: S04
milestone: M041-4tw0w7
key_files:
  - benchmark-results/fd-v2-baseline-before-m041-s04.md
  - .gsd/exec/60198853-46fd-4784-868e-7ad9fe9983b7.stdout
key_decisions:
  - (none)
duration: 
verification_result: untested
completed_at: 2026-06-13T17:17:48.449Z
blocker_discovered: false
---

# T01: Baseline measurement fd perf захвачен до S04 T04 perf optimization

**Baseline measurement fd perf захвачен до S04 T04 perf optimization**

## What Happened

Baseline measurement fd perf захвачен в planning phase. Прогнаны 100 reqs batch=1, 20 reqs batch=10, 10 reqs batch=32, 5 reqs batch=100 через curl к текущему fd (TEI model был warm ~17:10). Результаты: batch=1 p50=2.6ms p95=3.7ms p99=5.2ms (vs target 50ms — PASS 13x margin), batch=10 p50=2.8ms p95=3.9ms (vs target 200ms — PASS 51x margin), batch=32 p50=2.9ms p95=3.5ms (vs target 1000ms — PASS 286x margin), batch=100 — 5/5 reqs fail 500 за 2-3ms (НЕ timeout 10s как в спеке). Также захвачены: /v1/embeddings dimensions=512 — 500, /v1/embeddings encoding_format=base64 — 500, /embeddings/batch encoding_format=base64 — 500, B4 1MB input — 500 за 23ms (не timeout), B5/B6/B7 leaky Go errors подтверждены, B10 404 подтверждён, B11/B12 empty headers подтверждены. Артефакт benchmark-results/fd-v2-baseline-before-m041-s04.md (9317 bytes).

## Verification

benchmark-results/fd-v2-baseline-before-m041-s04.md существует (9317 bytes), содержит таблицу p50/p95/p99 для batch=1/10/32, все pass spec targets, обнаруженные несоответствия спеке (B4 fast-fail не timeout, B8 не bug, B9 500 fast не silent, dimensions=512 broken, encoding_format missing), probe bug matrix vs reality, и notes для S04 T05 финальной валидации (использовать benchmark.py не curl, добавить concurrent test с hey/wrk).

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| — | No verification commands discovered | — | — | — |

## Deviations

Использован ad-hoc bash+curl вместо tools/measure_fd_baseline.sh (скрипт не создан, в baseline artifact отмечено что S04 T05 final validation должен использовать benchmark.py для совместимости с M040 verifiers). Измерения sequential, не concurrent — для concurrent нужен hey/wrk.

## Known Issues

None.

## Files Created/Modified

- `benchmark-results/fd-v2-baseline-before-m041-s04.md`
- `.gsd/exec/60198853-46fd-4784-868e-7ad9fe9983b7.stdout`
