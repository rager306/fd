---
id: M014-vjfs9f
title: "Tagged ONNX performance benchmark"
status: complete
completed_at: 2026-05-20T04:37:31.442Z
key_decisions:
  - D010: Continue ONNX as opt-in experimental path; do not switch production/default from TEI yet.
  - Benchmark.py restart behavior must be target-aware for non-Compose benchmark APIs.
  - Tokenizer env names should not be redacted as secret-like solely because they contain `TOKENIZER`.
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m014-tei-baseline.txt
  - benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt
  - benchmark-results/fd-benchmark-m014-comparison.txt
  - .gsd/DECISIONS.md
  - .gsd/milestones/M014-vjfs9f/M014-vjfs9f-VALIDATION.md
lessons_learned:
  - For benchmark harnesses, restart semantics are part of correctness; a benchmark target on a different port/process needs a target-specific restart mechanism.
  - ONNX cold-path speedups can be real while cache-dominated p95 behavior remains mixed; compare by scenario, not by one aggregate.
  - Over-broad secret redaction can hide non-secret tokenizer metadata; redaction needs to protect credentials without erasing reproducibility fields.
---

# M014-vjfs9f: Tagged ONNX performance benchmark

**M014 proved tagged ONNX is promising for cold/model-bound performance but not ready to replace TEI as default.**

## What Happened

M014 measured the fixed-probe-correct tagged ONNX path against the current TEI default runtime. S01 extended the benchmark harness with runtime/native/ONNX/ORT metadata. S02 captured a fresh TEI baseline. S03 captured a tagged ONNX benchmark after fixing restart semantics so Redis L2 checks restarted the actual ONNX server, not the default Compose API. S04 compared results and recorded decision D010. Evidence shows tagged ONNX is materially faster on cold/model-bound paths and slightly higher in peak throughput in this single local run, but cache-dominated metrics are mixed and operational blockers remain. TEI stays production/default.

## Success Criteria Results

- TEI and tagged ONNX artifacts produced: pass.
- Sanitized config and runtime metadata present: pass.
- Redis namespace/cache effects explicit: pass.
- Correctness gate referenced before performance interpretation: pass (`fd-go-onnx-hf-tokenizer-m013-s03.txt` PASS).
- Data-backed recommendation without production switch: pass.

## Definition of Done Results

- [x] Default TEI baseline benchmark artifact exists with sanitized config snapshot: `benchmark-results/fd-benchmark-m014-tei-baseline.txt`.
- [x] Tagged ONNX benchmark artifact exists with build tags/native/ONNX/ORT metadata: `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`.
- [x] Cold/warm/batch/cache/startup/memory signals captured or documented: artifacts plus S03 startup/RSS summary.
- [x] Benchmark comparison synthesis states faster/slower by scenario: `benchmark-results/fd-benchmark-m014-comparison.txt`.
- [x] TEI remains default; no production switch occurred.
- [x] Tests/lint/tagged tests/artifact hygiene/GitNexus checks passed.

## Requirement Outcomes

- Benchmark comparability: advanced/validated with snapshot v2/v3 artifacts and comparison report.
- ONNX performance evidence: validated locally for tagged HF-tokenizer path.
- Production safety: validated; no default switch.
- Artifact hygiene: validated; no tracked ONNX/native binaries and raw probe leaks 0.

## Deviations

S03 discovered and fixed a benchmark validity issue: the Redis L2 restart check originally restarted Docker Compose `api` regardless of benchmark target. `benchmark.py` now supports `BENCHMARK_API_RESTART_COMMAND`, and tagged ONNX uses snapshot v3. TEI baseline remains snapshot v2, with this caveat documented.

## Follow-ups

Recommended next milestone: reproducible tagged ONNX packaging and repeated tuning benchmark. Scope should include pinned native tokenizer artifact release instead of `latest`, Docker/CI artifact supply and checksum verification, repeated benchmark runs/concurrency tuning, startup/RSS tracking in artifacts, and broader Russian/legal retrieval quality evaluation before any production switch.
