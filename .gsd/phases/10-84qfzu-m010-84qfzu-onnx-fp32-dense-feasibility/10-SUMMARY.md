---
id: M010-84qfzu
title: "ONNX FP32 dense feasibility spike"
status: complete
completed_at: 2026-05-19T18:52:03.307Z
key_decisions:
  - D006: ONNX FP32 dense path is locally feasible, but TEI remains production/default until adapter, performance, artifact, and quality gates pass.
  - Pin `transformers==4.51.3` for current export path; latest `transformers 5.8.1` failed.
  - Proceed only to opt-in prototype, not production switch.
key_files:
  - tools/compare_dense_embeddings.py
  - tools/export_user_bge_m3_dense_onnx.py
  - tools/compare_onnx_dense_embeddings.py
  - benchmark-results/fd-dense-comparator-m010-s02.txt
  - benchmark-results/fd-onnx-fp32-m010-s03.txt
  - .gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md
  - .gsd/milestones/M010-84qfzu/slices/S04/S04-RESEARCH.md
  - .gsd/milestones/M010-84qfzu/M010-84qfzu-VALIDATION.md
lessons_learned:
  - Comparator-before-export prevented a load-only success from being mistaken for semantic feasibility.
  - Unpinned ML export dependencies are risky; dependency version belongs in artifact metadata and future reproducibility gates.
  - Large ONNX artifacts require an artifact distribution design before production integration.
---

# M010-84qfzu: ONNX FP32 dense feasibility spike

**M010 proved local FP32 dense-only ONNX feasibility for the exact USER-bge-m3 model and produced a gated recommendation to prototype, not switch production runtime.**

## What Happened

M010 verified a model-preserving FP32 dense-only ONNX path for the current `deepvk/USER-bge-m3` model without changing production runtime. S01 established the provenance path and ruled out community BAAI/INT8 artifacts as replacements. S02 built a sanitized TEI/API dense comparator baseline. S03 exported a local dense-only ONNX model from the exact local snapshot, loaded it with ONNX Runtime CPU EP, and compared its output to the TEI baseline with cosine values around `0.999993` on fixed Russian/legal-style probes. S04 synthesized the evidence into a bounded recommendation: continue only to a gated opt-in ONNX adapter/prototype; do not switch production/default runtime from TEI yet.

## Success Criteria Results

- Model-preserving ONNX path: met, local FP32 dense export succeeded from exact model snapshot.
- Dense comparator baseline: met, S02 artifact saved.
- ONNX provenance/hashes/output metadata: met, S03 metadata and benchmark artifact saved.
- No production runtime/model/provider/language change: met.
- Evidence-based recommendation: met, S04 and D006 saved.

## Definition of Done Results

- All slices complete: S01-S04 complete.
- Verification passed: Go tests, lint, Compose config, Python py_compile/artifact checks, raw-text leakage checks, and GitNexus detect changes.
- Production runtime unchanged: no Go runtime behavior or Compose runtime default switched to ONNX.
- Evidence persisted: S01/S04 research, S02/S03 benchmark artifacts, D006, validation artifact.

## Requirement Outcomes

No existing requirements were invalidated. A future production ONNX adapter requirement should cover artifact distribution/checksum handling, opt-in runtime configuration, quality corpus, and performance gates.

## Deviations

The first ONNX export attempt with unpinned/latest `transformers 5.8.1` failed during legacy torch.onnx tracing. The successful path pins `transformers==4.51.3`. This deviation became part of the recommendation and future gate list.

## Follow-ups

Recommended next milestone: gated non-default ONNX adapter/prototype. Required first decisions: artifact distribution/checksum storage for the 1.43GB ONNX file, explicit opt-in backend config, performance benchmark against TEI+Redis baseline, larger Russian/legal quality corpus, and startup/health/failure observability.
