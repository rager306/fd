# S01: TEI queue_time root cause analysis

**Goal:** RCA TEI queue_time bottleneck: почему 2.7s queue wait при max_concurrent_requests=512. Документ обосновывает mitigation strategy для S02/S03.
**Demo:** After this, документ docs/te-perf-root-cause-m042.md объясняет почему TEI queue_time=2.7s несмотря на max_concurrent_requests=512, и обосновывает выбор mitigation strategy (async vs ONNX vs both) с evidence.

## Must-Haves

- documents/te-perf-root-cause-m042.md существует, ≥2KB. Содержит: TEI cold telemetry snapshot, hypothesis tree (≥3 гипотез с testable predictions), evidence collected, verdict + recommended action. Cross-references M019. Может служить input для S02 implementation.

## Proof Level

- This slice proves: artifact

## Integration Closure

Document consumed by S02 (async implementation) to decide concurrency level and rate-limiting. Also consumed by S03 (ONNX) to decide if ONNX is necessary vs async alone.

## Verification

- No code change; document only.

## Tasks

- [x] **T01: Captured direct TEI telemetry snapshot showing batch-size-sensitive queue_time growth for batch 1, 8, and 32.** `est:1h`
  Извлечь из benchmark-results/fd-v2-baseline-before-m041-s04.md (post-S01 section) + live TEI /info + docker logs fd_tei за последний час: total_time, tokenization_time, queue_time, inference_time per call. Построить таблицу cold latency by batch size (1, 8, 32) и correlation: queue_time vs batch_tokens. Capture concurrency limits: max_concurrent_requests=512, max_batch_requests=4, max_client_batch_size=32. Это evidence section в документе.
  - Files: `benchmark-results/fd-v2-baseline-before-m041-s04.md`, `docs/fd-v2-compatibility-report.md`
  - Verify: Snapshot таблица с минимум 3 batch sizes, ≥10 data points каждый. Готова для вставки в RCA документ.

- [ ] **T02: Live profile TEI: варьировать concurrency, batch size, timing** `est:1h`
  Active measurement: (1) одновременно 4 curl в parallel с batch=32, измерить max queue_time vs sequential. (2) Одновременно 16 curl с batch=1, измерить queue_time degradation. (3) Sleep 30s, один curl batch=32 — sanity. (4) Перезапустить fd_tei (down/up), один curl batch=32 — измерить true cold start. Цель: понять TEI behavior при разной нагрузке.
  - Files: `tools/profile_tei_concurrency.sh`
  - Verify: Профиль собран, ≥4 сценария, queue_time pattern задокументирован.

- [ ] **T03: Сформулировать hypothesis tree + verdict document** `est:1h`
  documents/te-perf-root-cause-m042.md: введение (TEI cold path 6s, queue_time 2.7s — почему), evidence section (из T01 + T02), hypothesis tree с минимум 3 гипотезами: (H1) single backend thread — TEI Rust default 1 worker despite max_concurrent_requests; (H2) lock contention в batcher scheduler; (H3) q_time metrics measures not real concurrency wait but internal scheduling. Каждая hypothesis с testable prediction, evidence supporting/refuting, verdict. Final verdict + recommended action для S02 (async) и S03 (ONNX).
  - Files: `documents/te-perf-root-cause-m042.md`
  - Verify: Документ ≥2KB, содержит: snapshot, hypothesis tree (≥3 с testable predictions), verdict, recommended action. Cross-references M019 ONNX как comparison baseline. Файл может быть прочитан S02 executor'ом перед implementation.

## Files Likely Touched

- benchmark-results/fd-v2-baseline-before-m041-s04.md
- docs/fd-v2-compatibility-report.md
- tools/profile_tei_concurrency.sh
- documents/te-perf-root-cause-m042.md
