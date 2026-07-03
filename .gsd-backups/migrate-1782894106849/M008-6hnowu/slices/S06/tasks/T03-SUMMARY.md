---
id: T03
parent: S06
milestone: M008-6hnowu
key_files: []
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:02:00.740Z
blocker_discovered: false
---

# T03: INT8 BGE-M3 ONNX is feasible as an experiment, but only after FP32 ONNX baseline and Russian legal quality gates pass.

**INT8 BGE-M3 ONNX is feasible as an experiment, but only after FP32 ONNX baseline and Russian legal quality gates pass.**

## What Happened

Researched INT8 quantization feasibility for BGE-M3 ONNX. ONNX Runtime supports 8-bit quantization through QOperator or QDQ representations and provides dynamic, static, and QAT quantization APIs. ORT docs note dynamic quantization computes activation quantization parameters during inference, often preserving accuracy but adding overhead, while static quantization uses calibration data and may be faster but risks more accuracy loss. ORT docs also recommend transformer-specific preprocessing/optimization before quantization, with shape inference and graph optimization as separate stages for easier debugging. For BGE-M3 specifically, `gpahal/bge-m3-onnx-int8` exists and claims BGE-M3 ONNX Runtime INT8 with dense, sparse, and ColBERT outputs plus O2 optimization; however, that does not prove suitability for Russian legal retrieval. Any INT8 spike must compare FP32 ONNX dense vectors and INT8 dense vectors against current TEI output, evaluate retrieval quality on a Russian legal corpus, and record model artifact hashes, quantization method, calibration corpus/version, ORT version/provider, and output similarity. INT8 should be a later experiment after FP32 ONNX baseline works.

## Verification

Read ONNX Runtime quantization docs and BGE-M3 INT8 ONNX model card evidence.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Fetched: https://onnxruntime.ai/docs/performance/model-optimizations/quantization.html` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Search query: ONNX Runtime quantization INT8 transformer BGE-M3 ONNX dynamic quantization static calibration accuracy dense sparse ColBERT` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Read source evidence for https://huggingface.co/gpahal/bge-m3-onnx-int8` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

INT8 is not automatically safe for fd because retrieval quality can regress. The found `gpahal/bge-m3-onnx-int8` has minimal popularity signal and must be treated as a reproducible candidate only if export scripts/model files are inspected and output equivalence is tested. Static quantization needs representative calibration data; for this project that means Russian legal text, not generic English samples.

## Files Created/Modified

None.
