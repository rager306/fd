---
id: T02
parent: S02
milestone: M008-6hnowu
key_files:
  - .gsd/milestones/M008-6hnowu/slices/S02/S02-RESEARCH.md
  - benchmark.py
  - benchmark-results/fd-environment-inxi-m008.txt
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:13:47.251Z
blocker_discovered: false
---

# T02: Defined the next benchmark matrix: baseline first, Redis evidence, MGET/pipeline A/B, ONNX FP32 dense, then provider/INT8/sidecar only behind gates.

**Defined the next benchmark matrix: baseline first, Redis evidence, MGET/pipeline A/B, ONNX FP32 dense, then provider/INT8/sidecar only behind gates.**

## What Happened

Defined the benchmark matrix and stop criteria for the next optimization milestone. The order is: comparable TEI+Go+Redis baseline; Redis long-lived cache and batch-hit baseline; Redis MGET/pipeline A/B; ONNX FP32 dense-only default CPU EP; ORT threading/provider matrix; optional language sidecar only if profiling justifies it. Quality-risking changes require Russian legal retrieval metrics such as Recall@k, nDCG@k, and MRR@k against a stable corpus. Benchmark artifacts must capture sanitized config snapshots including git/Docker/env/Redis/ONNX/model/tokenizer/hardware fields. Stop criteria prevent speculative rewrites, INT8, provider-stack changes, or model replacement without evidence.

## Verification

Saved S02 research artifact with integration seams, benchmark phases, metrics, config fields, and risk classes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Saved: .gsd/milestones/M008-6hnowu/slices/S02/S02-RESEARCH.md` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

Benchmark design is not yet implemented. Current `benchmark.py` still lacks the config snapshot, per-layer timing, Redis batch-hit isolation, and quality gate metrics described here.

## Files Created/Modified

- `.gsd/milestones/M008-6hnowu/slices/S02/S02-RESEARCH.md`
- `benchmark.py`
- `benchmark-results/fd-environment-inxi-m008.txt`
