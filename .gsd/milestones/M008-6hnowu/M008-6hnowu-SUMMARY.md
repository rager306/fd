---
id: M008-6hnowu
title: "Embedding runtime optimization research"
status: complete
completed_at: 2026-05-19T17:16:28.070Z
key_decisions:
  - Optimize the current Russian/legal-capable model first; model replacement requires Russian legal retrieval benchmarks.
  - Treat Redis L2 as a long-lived reusable embedding cache with model/version-aware invalidation.
  - Benchmark artifacts must include sanitized effective config snapshots.
  - Do not rewrite fd in Rust or C now; keep Go and profile before sidecar/FFI decisions.
  - ONNX FP32 dense-only default CPU EP is the first native runtime candidate; provider tuning and INT8 are later gated experiments.
key_files:
  - .gsd/milestones/M008-6hnowu/M008-6hnowu-ROADMAP.md
  - .gsd/milestones/M008-6hnowu/M008-6hnowu-VALIDATION.md
  - .gsd/milestones/M008-6hnowu/slices/S02/S02-RESEARCH.md
  - .gsd/milestones/M008-6hnowu/slices/S03/S03-RESEARCH.md
  - .gsd/milestones/M008-6hnowu/slices/S04/S04-RESEARCH.md
  - .gsd/milestones/M008-6hnowu/slices/S05/S05-RESEARCH.md
  - .gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md
  - benchmark-results/fd-environment-inxi-m008.txt
  - .gsd/DECISIONS.md
  - .gsd/REQUIREMENTS.md
lessons_learned:
  - Performance optimization research needs benchmark comparability before runtime experiments.
  - Redis cache key correctness becomes more important as retention increases.
  - Language/runtime rewrite decisions should follow pprof/per-layer evidence, not assumptions.
  - BGE-M3 multi-output capability is useful future context but current fd API/cache should remain dense-only.
---

# M008-6hnowu: Embedding runtime optimization research

**M008 produced a measured optimization roadmap: benchmark/config and long-lived Redis cache foundation first, ONNX/provider/language experiments later behind quality and performance gates.**

## What Happened

M008 researched embedding runtime optimization paths without making speculative runtime changes. It verified Go/ONNX options, defined the Russian legal quality gate, designed a long-lived Redis vector-cache strategy, evaluated Go/Rust/C rewrite tradeoffs, researched ONNX CPU provider/threading/INT8 options, captured hardware/environment context, and synthesized everything into an implementation-ready recommendation. The final path is deliberately conservative: build benchmark/config/cache evidence first, then test Redis MGET/pipeline, then ONNX FP32 dense-only, with provider tuning, INT8, Rust sidecar, and C FFI behind explicit gates.

## Success Criteria Results

- ✅ Go embedding alternatives verified from current sources.
- ✅ Integration and benchmark risks documented.
- ✅ Redis/cache throughput optimization opportunities researched and benchmark-scoped.
- ✅ Go vs C vs Rust rewrite/targeted-native options researched and benchmark-scoped.
- ✅ ONNX Runtime CPU acceleration and quantization options researched and benchmark-scoped.
- ✅ Measured next-step plan produced without unverified runtime migration, model replacement, provider-stack change, or language rewrite.
- ✅ GSD artifacts complete; local commit follows after checkpoint.

## Definition of Done Results

- ✅ All slices complete in GSD status.
- ✅ Research artifacts saved for S02/S03/S04/S05/S06 and summaries for completed slices/tasks.
- ✅ Fresh Go tests passed: 49 tests in 4 packages.
- ✅ Pinned GolangCI-Lint passed: 0 issues.
- ✅ `docker compose config` passed.
- ✅ Artifact existence check passed for required M008 files.
- ✅ GitNexus detect_changes passed with low risk and no changed symbols/processes.
- ⏳ Local commit will be created after DB checkpoint.

## Requirement Outcomes

- R001 advanced: Russian/legal embedding quality gate defined and carried into final recommendation.
- R002 advanced: Redis L2 long-lived reusable cache architecture defined.
- R003 advanced: env-configurable cache/runtime settings specified.
- R004 advanced: sanitized benchmark config snapshot fields specified.

No requirements were invalidated.

## Deviations

Milestone scope expanded during execution to include Redis long-lived cache architecture, env/config snapshot requirements, language rewrite research, ONNX CPU provider/INT8 research, and hardware baseline capture. The expansion was incorporated into additional slices S04/S05/S06 and final S02/S03 synthesis.

## Follow-ups

Create next implementation milestone for measured cache and benchmark foundation: sanitized benchmark config snapshots, env-configured model-aware Redis retention, Redis persistence hardening, Redis batch-hit benchmarks, then MGET/pipeline A/B if measured. Do not push without explicit user confirmation.
