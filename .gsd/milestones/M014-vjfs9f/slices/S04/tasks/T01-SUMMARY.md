---
id: T01
parent: S04
milestone: M014-vjfs9f
key_files:
  - benchmark-results/fd-benchmark-m014-comparison.txt
  - benchmark-results/fd-benchmark-m014-tei-baseline.txt
  - benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt
key_decisions:
  - Tagged ONNX is materially faster on uncached/cold model-bound paths in the local benchmark, but not uniformly faster on every cache-dominated metric.
  - Peak throughput was higher for ONNX in this single run, but concurrency shape differs and needs repeated/tuned runs before capacity claims.
  - Benchmark evidence supports continuing ONNX as opt-in experimental work, not a production default switch.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:34:09.941Z
blocker_discovered: false
---

# T01: Created the M014 TEI vs tagged ONNX comparison artifact with deltas and recommendation.

**Created the M014 TEI vs tagged ONNX comparison artifact with deltas and recommendation.**

## What Happened

Parsed the TEI and tagged ONNX benchmark artifacts and wrote a comparison artifact. ONNX shows strong cold-path improvement: best cold latency 59.0ms → 10.2ms, warm mean 2.25ms → 1.63ms, and max throughput ~750 → ~891 req/s. Cache-dominated metrics are mixed: ONNX wins Redis L2 batch p95, TEI wins batch L1 p95 and chunk reuse p95. The artifact documents correctness gate, runtime modes, caveats, and recommendation.

## Verification

Comparison artifact exists and includes scenario deltas, interpretation, caveats, and recommendation.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `test ! -e benchmark-results/fd-benchmark-m014-comparison.txt` | 0 | ✅ pass — new artifact path verified before write | 0ms |
| 2 | `uv run --python 3.13 python comparison generator` | 0 | ✅ pass — comparison artifact written | 0ms |
| 3 | `read benchmark-results/fd-benchmark-m014-comparison.txt` | 0 | ✅ pass — deltas/caveats/recommendation present | 0ms |

## Deviations

None.

## Known Issues

Comparison uses one TEI run and one ONNX run on local KVM/QEMU host; it is not a repeated statistical benchmark or production-host result.

## Files Created/Modified

- `benchmark-results/fd-benchmark-m014-comparison.txt`
- `benchmark-results/fd-benchmark-m014-tei-baseline.txt`
- `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`
