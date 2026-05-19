# S04 — Research

**Date:** 2026-05-19

## Summary

M010 proves that a model-preserving FP32 dense-only ONNX path for the exact `deepvk/USER-bge-m3` model is locally feasible. S01 found no ready upstream FP32 ONNX artifact for this exact model and ranked local export from the current snapshot as the only safe primary path. S02 created a sanitized TEI/API dense comparator baseline. S03 exported a local FP32 dense-only ONNX candidate, loaded it with ONNX Runtime CPU EP, and compared its dense output against the TEI/API baseline.

The local ONNX artifact is `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`, size `1432482908` bytes, SHA256 `28538a17a99302e144149732d73fb273cd7c7a0468dc59167caa5a2d5ff2a3d4`. ONNX Runtime metadata: inputs `input_ids` and `attention_mask`, output `dense_vecs`, output shape `[batch_size, 1024]`, output type `tensor(float)`, provider `CPUExecutionProvider`. The S03 comparison artifact passed: live TEI hashes matched the S02 baseline and TEI-vs-ONNX cosine values were approximately `0.99999347` to `0.99999393` across fixed Russian/legal-style probes.

This is a feasibility proof, not a production switch. Production/default runtime remains TEI. The safe next step, if desired, is a separate milestone for a non-default ONNX adapter prototype and benchmark gate. Do not introduce INT8, provider variants, sparse/ColBERT, Rust/C rewrite, or model replacement from this milestone alone.

## Recommendation

Proceed, but only to a **non-default ONNX adapter/prototype milestone**, not a production runtime replacement.

Recommended next milestone scope:

1. Add a non-default ONNX embedder path behind explicit configuration, leaving TEI as default.
2. Define ONNX artifact storage/distribution outside git because the artifact is ~1.43GB.
3. Add repeatable export or artifact-download procedure with pinned dependencies, especially `transformers==4.51.3` for the current export path.
4. Benchmark ONNX CPU EP vs current TEI+Go+Redis under the existing M009 benchmark framework.
5. Run a larger Russian/legal quality gate before production use.
6. Only after dense FP32 passes quality/performance gates should provider/INT8/quantization work be planned.

Current decision: **continue research/prototype path, do not switch production runtime**.

## Implementation Landscape

### Key Files

- `.gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md` — candidate ranking and provenance requirements.
- `tools/compare_dense_embeddings.py` — S02 TEI/API baseline comparator.
- `benchmark-results/fd-dense-comparator-m010-s02.txt` — tracked TEI/API baseline artifact.
- `tools/export_user_bge_m3_dense_onnx.py` — local dense-only FP32 ONNX export tool.
- `.gsd/runtime/onnx/m010-s03/export-metadata.json` — ignored local export metadata, including ONNX hash and dependency versions.
- `tools/compare_onnx_dense_embeddings.py` — ONNX-vs-TEI comparison tool.
- `benchmark-results/fd-onnx-fp32-m010-s03.txt` — tracked ONNX comparison artifact.

### Build Order

For any follow-up implementation:

1. Keep TEI as default.
2. Decide artifact storage and retrieval first; git cannot carry the 1.43GB ONNX artifact.
3. Add an ONNX adapter behind explicit env/config such as `EMBEDDING_BACKEND=onnx`, not enabled by default.
4. Reuse S02/S03 comparator probes in tests, then expand to larger Russian/legal corpus evaluation.
5. Run performance benchmarks after integration, comparing cold/hot/cache behavior against M009/M010 baselines.
6. Only then consider INT8/provider experiments.

### Verification Approach

Minimum follow-up gates before production consideration:

- Export reproducibility: same source model revision and hashes, pinned dependencies, ONNX hash recorded.
- API equivalence: `/v1/embeddings` dimensions, normalization, and response shape match TEI default.
- Dense similarity: ONNX vs TEI cosine threshold on fixed probes and expanded Russian/legal corpus.
- Performance: benchmark latency/throughput under current Docker/CPU environment, including cold/hot cache behavior.
- Operations: artifact download/checksum verification, startup failure messages, health/metadata endpoint, and fallback behavior.
- Security: no raw input text in logs/artifacts; no external Redis exposure regression; no model artifact committed to git.

## Don't Hand-Roll

| Problem | Existing Solution | Why Use It |
|---------|------------------|------------|
| Dense ONNX execution | ONNX Runtime CPUExecutionProvider | Baseline provider and simplest CPU proof before provider variants. |
| Export | PyTorch ONNX export with explicit wrapper, or Optimum if wrapper support is improved | Keeps pooling/normalization explicit and model-preserving. |
| Comparator | Existing `tools/compare_dense_embeddings.py` and `tools/compare_onnx_dense_embeddings.py` | Avoids subjective equivalence claims and raw text leakage. |
| Benchmarking | Existing `benchmark.py` M009 snapshot framework | Captures sanitized config/runtime metadata for comparability. |

## Constraints

- Production/default runtime was not changed and should not be changed by this spike.
- The ONNX artifact is large (~1.43GB) and must stay out of git.
- Current export path requires `transformers==4.51.3`; unpinned `transformers 5.8.1` failed during legacy torch.onnx trace.
- The proof is dense-only. Sparse and ColBERT remain out of scope.
- The proof uses fixed short probes. It is not a full Russian/legal retrieval quality benchmark.
- Provider/INT8 claims still require separate evidence.

## Common Pitfalls

- **Treating load success as production readiness** — S03 proves feasibility, not service integration or throughput.
- **Committing artifacts** — ONNX files are multi-GB-scale runtime artifacts and must remain ignored.
- **Dropping dependency pins** — current export is sensitive to Transformers version.
- **Skipping quality corpus** — fixed probes catch shape/norm/regression issues but do not prove retrieval quality.
- **Defaulting to ONNX too early** — an adapter must be opt-in until quality/performance/ops gates pass.

## Open Risks

- Runtime performance may be worse or better than TEI; no throughput benchmark has been run yet.
- Memory footprint at service startup is unknown for a Go-integrated or sidecar ONNX runtime path.
- Artifact distribution/checksum verification must be designed before production deploy.
- Export may need maintenance as Torch/Transformers ONNX APIs evolve away from legacy tracing.

## Sources

- S01 candidate ranking and provenance requirements (source: `.gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md`).
- TEI/API dense baseline evidence (source: `benchmark-results/fd-dense-comparator-m010-s02.txt`).
- Successful local ONNX export/load metadata (source: `.gsd/runtime/onnx/m010-s03/export-metadata.json`).
- ONNX-vs-TEI comparison evidence (source: `benchmark-results/fd-onnx-fp32-m010-s03.txt`).
