# S01 — Research

**Date:** 2026-05-19

## Summary

S01 found no ready FP32 ONNX artifact for the exact production model `deepvk/USER-bge-m3`. Both local TEI storage and the upstream Hugging Face file tree contain safetensors/tokenizer/config/SentenceTransformer artifacts, but no ONNX model. The local model revision is `0cc6cfe48e260fb0474c753087a69369e88709ae`, with `model.safetensors` SHA256 `e6aa9c8e51a60ff383186a2f28f658305ba4ad23d2fa24296607885458ef2733`, `tokenizer.json` SHA256 `068d9f7ed9dd190a00a567e5f7750fdc591b93bc623072ac8050a303c25f5937`, and `config.json` SHA256 `f3552b70cacff0f14829896d9021372a7667676f900111d1e68664e021ab3f7f`.

The safest next path is a model-preserving FP32 dense-only export from this exact `deepvk/USER-bge-m3` snapshot, not adoption of a community BAAI artifact. Community ONNX repositories are useful for implementation clues: `aapot/bge-m3-onnx` shows Optimum export with explicit `dense_vecs`, `sparse_vecs`, and `colbert_vecs`; `hotchpotch/vespa-onnx-BAAI-bge-m3-only-dense` shows dense-only serving conventions; `yuniko-software/bge-m3-onnx` shows external-data and cross-language ONNX Runtime patterns. They are not model-preserving artifacts for fd.

No production runtime change is recommended or made by this slice. S02 should build the dense comparator first; S03 should then attempt export/load using ignored local artifact storage and record provenance/hashes/output metadata. INT8, OpenVINO, oneDNN, ZenDNN, sparse/ColBERT integration, and model replacement remain out of scope.

## Recommendation

Proceed to S02 with a TEI/API dense-output comparator before attempting ONNX runtime adoption. The comparator should establish fixed Russian/legal probes, expected dimensionality (`1024`), normalization checks, cosine similarity thresholds, and artifact metadata output. Only after that should S03 attempt FP32 dense-only export/load from the local `deepvk/USER-bge-m3` snapshot.

Preferred candidate ranking:

| Rank | Candidate | Use | Reason | Decision |
|---:|---|---|---|---|
| 1 | Local export from exact `deepvk/USER-bge-m3` snapshot | Primary path | Preserves current model/revision and M010 boundary; can use local TEI cache hashes | Proceed after S02 comparator |
| 2 | `aapot/bge-m3-onnx` | Export/reference pattern | Shows Optimum custom config, output names, external data, CLS pooling + L2 normalization | Reference only; base is `BAAI/bge-m3` |
| 3 | `yuniko-software/bge-m3-onnx` | ONNX Runtime/cross-language reference | Demonstrates full BGE-M3 ONNX pipeline and external-data handling | Reference only; base is `BAAI/bge-m3` |
| 4 | `hotchpotch/vespa-onnx-BAAI-bge-m3-only-dense` | Dense-only serving convention | Shows dense-only CLS pooling + normalize for Vespa | Reference only; FP16/INT8 and BAAI, not FP32 USER |
| 5 | `skatzR/USER-BGE-M3-ONNX-INT8` | Future quantization reference | Model-family relevant but INT8 only | Out of scope until FP32 comparator/export proves baseline |

## Implementation Landscape

### Key Files

- `.gsd/milestones/M010-84qfzu/slices/S01/tasks/T01-SUMMARY.md` — local TEI/HF cache inspection, model revision, file hashes, storage constraints.
- `.gsd/milestones/M010-84qfzu/slices/S01/tasks/T02-SUMMARY.md` — source-backed ONNX candidate ranking and dense-output compatibility risks.
- `benchmark.py` — existing benchmark/config snapshot pattern to reuse later for ONNX provenance snapshots; no S01 changes required.
- `api/embed/tei.go` — current TEI baseline boundary for comparator calls; no S01 changes required.
- `.gitignore` — already excludes GSD runtime noise and `tei-models/`; S03 should use `.gsd/runtime/onnx/` or add a dedicated ignored artifact path before large ONNX output.

### Build Order

1. Build S02 comparator against current TEI/API baseline, not ONNX first. This avoids mistaking an export/load success for semantic equivalence.
2. Define fixed Russian/legal probes and expected dense output properties: 1024 dimensions, finite float values, L2-normalized vector, stable cache-independent response, and cosine similarity comparison support.
3. In S03, attempt FP32 dense export from exact local snapshot. Record export command, package versions, opset, dtype, optimization level, artifact hashes, input/output names, output shapes, and comparison metrics.
4. Only if S03 proves export/load and comparator alignment should S04 recommend adapter implementation work.

### Verification Approach

- For S02: run a comparator script against `http://localhost:8000/v1/embeddings` or current API equivalent and verify fixed Russian/legal probes produce 1024-dimensional finite normalized vectors. Save no raw sensitive texts in benchmark artifacts unless the probes are explicitly non-sensitive fixtures.
- For S03: if ONNX export is attempted, verify with `onnxruntime` CPU EP load/run and compare ONNX dense output to the S02 TEI/API baseline by cosine similarity and shape checks.
- For every ONNX attempt: record artifact hashes and output metadata in a machine-readable sidecar or benchmark section.
- Final slice verification: `S01-RESEARCH.md` exists and explicitly states no production runtime change.

## Don't Hand-Roll

| Problem | Existing Solution | Why Use It |
|---------|------------------|------------|
| ONNX export | Hugging Face Optimum ONNX exporter | Supports Sentence Transformers/XLM-Roberta-family exports and records export options. |
| Runtime load/profiling | ONNX Runtime CPU Execution Provider | Baseline provider before oneDNN/OpenVINO/ZenDNN experiments. |
| Tokenization | Existing `tokenizer.json`/HF tokenizer from exact model snapshot | Prevents accidental tokenizer drift and preserves cache namespace semantics. |
| Dense pooling reference | Model card direct-transformers CLS pooling + L2 normalization; `aapot` implementation | Keeps output comparable to current SentenceTransformer/TEI behavior. |

## Constraints

- The production/default runtime must remain TEI/API; S01/S02/S03 may create local artifacts only.
- The model must remain `deepvk/USER-bge-m3`; BAAI BGE-M3 artifacts are not replacements.
- ONNX artifacts are large and must remain untracked.
- Sparse and multi-vector encoding for `deepvk/USER-bge-m3` were not thoroughly evaluated by the model authors; dense-only comes first.
- INT8/provider experiments require a later FP32 baseline and Russian/legal quality gate.
- Current host baseline is KVM/QEMU on AMD EPYC with no swap; provider/quantization claims must be benchmarked on this environment.

## Common Pitfalls

- **Exporting token hidden states instead of sentence embeddings** — direct `AutoModel` ONNX may output `last_hidden_state`; comparator must verify pooling and normalization.
- **Using BAAI artifacts as if they were USER-bge-m3** — community artifacts may load and run but violate the model-preservation requirement.
- **Untracked artifact leakage** — ONNX external-data files can be multi-GB; use ignored local storage and verify `git status` before committing.
- **Normalization mismatch** — some references emit normalized dense vectors and others require post-processing; comparator must check L2 norm and cosine similarity.
- **Over-reading INT8 benchmarks** — INT8 `USER-BGE-M3-ONNX-INT8` may be useful later, but M010 is FP32 dense-only.

## Open Risks

- Optimum's generic SentenceTransformer export may not include the exact pooling/normalization needed for byte-for-byte or high-cosine alignment with TEI.
- Export may need a custom wrapper similar to `aapot/bge-m3-onnx` to expose only dense vectors.
- ONNX external-data handling may complicate artifact hashing and local storage.
- S03 may fail due to memory, package compatibility, or unsupported graph export details on this host; that is an acceptable spike outcome if recorded with evidence.

## Sources

- Target model identity, 1024-dim dense output, Russian-focused training, direct transformers CLS pooling and normalization, and sparse/multi-vec limitations (source: [deepvk/USER-bge-m3](https://huggingface.co/deepvk/USER-bge-m3)).
- Upstream file tree showing safetensors/tokenizer/config but no ONNX artifact (source: [deepvk/USER-bge-m3 tree](https://huggingface.co/deepvk/USER-bge-m3/tree/main)).
- Optimum ONNX export support for Sentence Transformers and XLM-Roberta-family architectures (source: [Optimum ONNX overview](https://huggingface.co/docs/optimum-onnx/onnx/overview)).
- Optimum CLI/API export options including model path, task, opset, dtype, optimization, and local output path (source: [Optimum export guide](https://huggingface.co/docs/optimum-onnx/onnx/usage_guides/export_a_model)).
- ONNX Runtime/Hugging Face export and runtime overview (source: [ONNX Runtime Hugging Face](https://onnxruntime.ai/huggingface)).
- BAAI BGE-M3 ONNX reference with dense/sparse/ColBERT outputs and Optimum export script (source: [aapot/bge-m3-onnx](https://huggingface.co/aapot/bge-m3-onnx)).
- `aapot` output names and normalization implementation (source: [aapot export_onnx.py](https://huggingface.co/aapot/bge-m3-onnx/raw/main/export_onnx.py), [aapot bgem3_model.py](https://huggingface.co/aapot/bge-m3-onnx/raw/main/bgem3_model.py)).
- Cross-language BGE-M3 ONNX implementation and external data handling (source: [yuniko-software/bge-m3-onnx README](https://raw.githubusercontent.com/yuniko-software/bge-m3-onnx/main/README.md)).
- Dense-only BAAI BGE-M3 ONNX serving pattern for Vespa; FP16/INT8 only (source: [hotchpotch/vespa-onnx-BAAI-bge-m3-only-dense](https://huggingface.co/hotchpotch/vespa-onnx-BAAI-bge-m3-only-dense)).
- USER-BGE-M3 INT8 ONNX artifact and cosine benchmark claims, explicitly deferred from FP32 spike (source: [skatzR/USER-BGE-M3-ONNX-INT8](https://huggingface.co/skatzR/USER-BGE-M3-ONNX-INT8)).
