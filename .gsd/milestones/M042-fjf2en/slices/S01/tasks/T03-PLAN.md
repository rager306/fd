---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Сформулировать hypothesis tree + verdict document

documents/te-perf-root-cause-m042.md: введение (TEI cold path 6s, queue_time 2.7s — почему), evidence section (из T01 + T02), hypothesis tree с минимум 3 гипотезами: (H1) single backend thread — TEI Rust default 1 worker despite max_concurrent_requests; (H2) lock contention в batcher scheduler; (H3) q_time metrics measures not real concurrency wait but internal scheduling. Каждая hypothesis с testable prediction, evidence supporting/refuting, verdict. Final verdict + recommended action для S02 (async) и S03 (ONNX).

## Inputs

- None specified.

## Expected Output

- `documents/te-perf-root-cause-m042.md`

## Verification

Документ ≥2KB, содержит: snapshot, hypothesis tree (≥3 с testable predictions), verdict, recommended action. Cross-references M019 ONNX как comparison baseline. Файл может быть прочитан S02 executor'ом перед implementation.
