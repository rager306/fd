---
id: T01
parent: S03
milestone: M008-6hnowu
key_files:
  - .gsd/milestones/M008-6hnowu/slices/S03/S03-RESEARCH.md
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:14:41.503Z
blocker_discovered: false
---

# T01: Recommended the next path: benchmark/config snapshots plus model-aware long-lived Redis cache first; ONNX/provider/language changes later behind evidence gates.

**Recommended the next path: benchmark/config snapshots plus model-aware long-lived Redis cache first; ONNX/provider/language changes later behind evidence gates.**

## What Happened

Synthesized all M008 research into a ranked recommendation. The next milestone should not jump to ONNX, INT8, provider stacks, Rust, C, or model replacement. It should first implement the measurement/cache foundation: sanitized benchmark config snapshots, env-configured model-aware Redis retention, Redis persistence/deployment hardening, Redis batch-hit benchmarks, and only then MGET/pipeline A/B if the baseline shows round-trip pressure. ONNX FP32 dense-only follows later once quality gates and config snapshots exist. Provider tuning, INT8, and Rust sidecar remain explicitly gated. Full C service and model replacement are rejected for now.

## Verification

Saved S03 final recommendation artifact with ranked options, next milestone proposal, required verification, and non-goals.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Saved: .gsd/milestones/M008-6hnowu/slices/S03/S03-RESEARCH.md` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Synthesized S01/S02/S04/S05/S06 outputs` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

This is a research recommendation. No code has been changed for benchmark snapshots, env cache retention, or MGET/pipeline yet.

## Files Created/Modified

- `.gsd/milestones/M008-6hnowu/slices/S03/S03-RESEARCH.md`
