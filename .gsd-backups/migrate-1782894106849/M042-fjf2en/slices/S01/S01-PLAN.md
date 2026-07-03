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

- [x] **T02: Documented the attempted TEI concurrency profile and the stronger finding: TEI restart/recreate can spend ~48 minutes in backend startup before becoming ready.** `est:1h`
  Use evidence already captured from the failed profile run, docker health/logs, process state, and T01 direct TEI timings to document TEI startup/restart fragility and concurrency observations. Do not perform additional TEI restarts unless explicitly required; restore/leave service state clearly documented. Write `benchmark-results/te-concurrency-profile-m042-s01.md` with scenarios attempted, successful T01/T02 signals, restart timeout evidence, and limitations.
  - Files: `benchmark-results/te-concurrency-profile-m042-s01.md`, `benchmark-results/te-concurrency-profile-m042-s01-run.txt`
  - Verify: Artifact exists, includes >=4 scenarios/attempts (sequential batch32, parallel batch32 attempt, parallel batch1 attempt, idle batch32 attempt, restart timeout), records TEI health/startup evidence, and does not claim missing metrics as pass.

- [x] **T03: Wrote the M042 TEI RCA, concluding that TEI queue/startup behavior is the root target and ONNX should be deferred from the current milestone.** `est:1h`
  Write `documents/te-perf-root-cause-m042.md` with introduction, evidence from T01/T02, hypothesis tree (>=3 hypotheses with testable predictions), verdict, and recommended action. The recommendation must explicitly defer ONNX runtime implementation from M042 and focus on TEI stabilization/mitigation. Cross-reference M019/M040 as historical ONNX research only, not current implementation scope.
  - Files: `documents/te-perf-root-cause-m042.md`
  - Verify: Document is >=2KB and contains snapshot, hypothesis tree with >=3 hypotheses and testable predictions, verdict, recommended TEI-first action, and M019/M040 cross-references. R020 can be validated from this artifact.

## Files Likely Touched

- benchmark-results/fd-v2-baseline-before-m041-s04.md
- docs/fd-v2-compatibility-report.md
- benchmark-results/te-concurrency-profile-m042-s01.md
- benchmark-results/te-concurrency-profile-m042-s01-run.txt
- documents/te-perf-root-cause-m042.md
