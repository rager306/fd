---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Captured direct TEI telemetry snapshot showing batch-size-sensitive queue_time growth for batch 1, 8, and 32.

Извлечь из benchmark-results/fd-v2-baseline-before-m041-s04.md (post-S01 section) + live TEI /info + docker logs fd_tei за последний час: total_time, tokenization_time, queue_time, inference_time per call. Построить таблицу cold latency by batch size (1, 8, 32) и correlation: queue_time vs batch_tokens. Capture concurrency limits: max_concurrent_requests=512, max_batch_requests=4, max_client_batch_size=32. Это evidence section в документе.

## Inputs

- None specified.

## Expected Output

- `documents/te-perf-snapshot-m042-s01.md (intermediate, вливается в основной RCA)`

## Verification

Snapshot таблица с минимум 3 batch sizes, ≥10 data points каждый. Готова для вставки в RCA документ.
