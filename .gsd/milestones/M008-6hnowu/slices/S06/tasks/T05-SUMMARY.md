---
id: T05
parent: S06
milestone: M008-6hnowu
key_files:
  - .gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md
  - benchmark-results/fd-environment-inxi-m008.txt
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:03:17.522Z
blocker_discovered: false
---

# T05: Produced the ONNX CPU optimization benchmark plan and exclusions.

**Produced the ONNX CPU optimization benchmark plan and exclusions.**

## What Happened

Produced the S06 research artifact synthesizing ONNX provider, threading/NUMA, INT8 quantization, and BGE-M3 multi-output findings into a benchmark plan. The recommendation is to start with current TEI/Candle baseline, then BGE-M3 ONNX FP32 dense-only on default CPU EP, then threading/NUMA tuning, then oneDNN/OpenVINO/ZenDNN only if practical, and INT8 only after FP32 ONNX plus Russian legal quality gates. The artifact lists required config snapshot fields and stop criteria.

## Verification

Saved S06-RESEARCH.md with provider matrix, benchmark order, config fields, and stop criteria.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Artifact: .gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

ZenDNN remains unverified as a current supported easy-to-use provider; it should not be benchmarked until a current package/build path is found. INT8 should not be attempted before FP32 ONNX dense output is validated.

## Files Created/Modified

- `.gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md`
- `benchmark-results/fd-environment-inxi-m008.txt`
