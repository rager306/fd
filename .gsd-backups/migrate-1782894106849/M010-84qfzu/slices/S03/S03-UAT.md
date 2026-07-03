# S03: ONNX FP32 CPU feasibility — UAT

**Milestone:** M010-84qfzu
**Written:** 2026-05-19T18:47:05.940Z

# S03 UAT — ONNX FP32 CPU feasibility

## Checks

- [x] Ignored local workspace exists under `.gsd/runtime/onnx/m010-s03/`.
- [x] Exact local source model hashes/provenance recorded.
- [x] First unpinned export failure captured and explained.
- [x] Pinned `transformers==4.51.3` export succeeded.
- [x] ONNX Runtime CPU EP loaded the exported model.
- [x] Output metadata records `dense_vecs`, shape `[batch_size,1024]`, type `tensor(float)`.
- [x] S03 comparison artifact exists at `benchmark-results/fd-onnx-fp32-m010-s03.txt`.
- [x] TEI live hashes matched S02 baseline hashes.
- [x] TEI-vs-ONNX cosine values exceeded `0.999` for all fixed probes.
- [x] Raw probe texts were excluded from tracked artifacts.
- [x] Production API/runtime defaults were not changed.

## UAT Result

Pass. The model-preserving FP32 dense-only ONNX path is locally feasible. This is not yet a production integration decision; it is an evidence-backed input for S04.

