# S06: Research ONNX CPU acceleration and quantization

**Goal:** Verify ONNX Runtime CPU optimization paths for current-model inference: ZenDNN/oneDNN/OpenVINO providers, NUMA/threading, INT8 quantization, and BGE-M3 dense/sparse/ColBERT outputs.
**Demo:** After this, ONNX CPU-level optimization options are verified and benchmark-scoped for the current model.

## Must-Haves

- Verify whether ZenDNN is available as an ONNX Runtime EP/build option and what AMD EPYC support looks like.
- Compare default CPU EP, oneDNN/DNNL, OpenVINO, and ZenDNN where applicable.
- Define NUMA/threading/env knobs and benchmark controls.
- Verify BGE-M3 ONNX dense/sparse/ColBERT output claims and implications for fd API/cache.
- Define INT8 quantization feasibility, correctness risks, and Russian/legal quality gate.

## Proof Level

- This slice proves: source research plus benchmark design

## Integration Closure

Separates provider/build/runtime tuning from model replacement and records what can be tested on the target VPS.

## Verification

- Defines config snapshot fields and runtime metrics needed to compare CPU EP, thread, NUMA, and quantization variants.

## Tasks

- [x] **T01: Verify ONNX Runtime CPU provider options** `est:medium`
  Research ONNX Runtime CPU execution provider options relevant to AMD EPYC/VPS: default CPU EP, oneDNN/DNNL, OpenVINO, ZenDNN if available, build/package requirements, and whether prebuilt artifacts include them.
  - Verify: Current docs/source evidence recorded with practical availability classification.

- [x] **T02: Research NUMA and threading controls** `est:small`
  Research NUMA/threading tuning for ONNX Runtime CPU inference and Linux VPS deployment: intra/inter op threads, affinity, OMP/MKL/ORT env vars, numactl, CPU pinning, and how benchmark config snapshots should record them.
  - Verify: Threading/NUMA knobs and benchmark controls recorded.

- [x] **T03: Research INT8 quantization feasibility** `est:medium`
  Research INT8 quantization feasibility for BGE-M3 ONNX: dynamic/static quantization, calibration data needs, dense/sparse/ColBERT output correctness risks, Russian/legal quality benchmark implications, and rollback criteria.
  - Verify: Quantization path includes quality/correctness gate and model-output compatibility risks.

- [x] **T04: Verify BGE-M3 ONNX multi-output implications** `est:medium`
  Verify BGE-M3 ONNX output claims: dense, sparse, and ColBERT outputs, how they differ from fd's current dense embedding API/cache, and whether fd should ignore, expose, or separately store non-dense outputs.
  - Verify: Output-shape implications and compatibility decisions recorded.

- [x] **T05: Recommend ONNX CPU optimization benchmark path** `est:small`
  Produce ONNX CPU optimization recommendation with provider/quantization/NUMA benchmark matrix, required env/config snapshot fields, success metrics, and stop criteria.
  - Files: `.gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md`
  - Verify: Research artifact includes matrix, ranked options, benchmark config fields, and exclusions.

## Files Likely Touched

- .gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md
