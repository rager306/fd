---
id: S03
parent: M010-84qfzu
milestone: M010-84qfzu
provides:
  - Feasible local FP32 dense-only ONNX candidate.
  - Export script and comparison script.
  - Evidence for S04 recommendation and future adapter design.
requires:
  []
affects:
  - S04
key_files:
  - tools/export_user_bge_m3_dense_onnx.py
  - tools/compare_onnx_dense_embeddings.py
  - benchmark-results/fd-onnx-fp32-m010-s03.txt
key_decisions:
  - Model-preserving FP32 dense-only ONNX path is feasible locally.
  - `transformers==4.51.3` must be pinned for the current export path; latest `transformers 5.8.1` failed during torch.onnx trace.
  - ONNX output is `dense_vecs` with shape `[batch_size,1024]`, CPUExecutionProvider, CLS pooling, and L2 normalization.
patterns_established:
  - Pin export dependencies when tracing transformer models; unpinned latest packages can break export.
  - Compare ONNX outputs to TEI only after verifying TEI live hashes match the S02 baseline.
  - Keep large ONNX artifacts under ignored runtime storage and track only metadata/results.
observability_surfaces:
  - `.gsd/runtime/onnx/m010-s03/source-provenance.json` local source hashes.
  - `.gsd/runtime/onnx/m010-s03/export-metadata.json` local export/load metadata.
  - `benchmark-results/fd-onnx-fp32-m010-s03.txt` tracked TEI-vs-ONNX comparison evidence.
drill_down_paths:
  - .gsd/milestones/M010-84qfzu/slices/S03/tasks/T01-SUMMARY.md
  - .gsd/milestones/M010-84qfzu/slices/S03/tasks/T02-SUMMARY.md
  - .gsd/milestones/M010-84qfzu/slices/S03/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T18:47:05.940Z
blocker_discovered: false
---

# S03: ONNX FP32 CPU feasibility

**S03 proved local FP32 dense-only ONNX feasibility for the exact USER-bge-m3 model, with CPU EP load and TEI cosine agreement around 0.999993.**

## What Happened

S03 prepared an ignored local ONNX artifact workspace, captured source model provenance, implemented a dense-only FP32 export tool, and compared the exported ONNX candidate against the S02 TEI/API baseline. The first unpinned dependency run failed with `transformers 5.8.1` during legacy torch.onnx tracing, but the exact same export succeeded with `transformers==4.51.3`. The successful ONNX artifact is `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`, size `1432482908` bytes, SHA256 `28538a17a99302e144149732d73fb273cd7c7a0468dc59167caa5a2d5ff2a3d4`. ONNX Runtime CPU EP loaded it with inputs `input_ids`/`attention_mask` and output `dense_vecs` shape `[batch_size,1024]`. The comparison artifact passed: all live TEI hashes matched S02, ONNX vectors were finite and normalized, and TEI-vs-ONNX cosine values were all above `0.999` with observed values around `0.999993`. Production runtime remains unchanged.

## Verification

S03 verification passed: export script compiled; unpinned export failure recorded; pinned export succeeded; ONNX Runtime CPU EP dummy load succeeded; ONNX comparator command exited 0 and produced a PASS artifact with no raw probe text leakage.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

- Future production ONNX integration needs an artifact distribution/storage requirement because git cannot track the 1.43GB ONNX artifact.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

The first unpinned export attempt failed under `transformers 5.8.1`; the slice adapted by pinning `transformers==4.51.3`, which is consistent with BGE-M3 ONNX reference tooling. No production runtime default was changed.

## Known Limitations

The proof covers fixed short Russian/legal probes only, not a full retrieval-quality benchmark or production throughput benchmark. The local ONNX artifact is ignored and not distributable from git. Export currently uses legacy TorchScript-based ONNX export with a deprecation warning. No ONNX adapter is wired into the Go API.

## Follow-ups

S04 should recommend proceeding to an optional non-default ONNX adapter/prototype only if the project wants runtime integration. Before production consideration, add larger Russian/legal corpus evaluation, benchmark latency/throughput against TEI, and decide artifact storage/distribution strategy.

## Files Created/Modified

- `tools/export_user_bge_m3_dense_onnx.py` — Local ONNX export tool for dense-only USER-bge-m3 FP32 candidate.
- `tools/compare_onnx_dense_embeddings.py` — ONNX-vs-TEI comparison tool using fixed S02 probes and baseline hash verification.
- `benchmark-results/fd-onnx-fp32-m010-s03.txt` — Tracked comparison evidence for S03.
- `.gsd/runtime/onnx/m010-s03/export-metadata.json` — Ignored local ONNX export metadata and 1.43GB ONNX artifact.
