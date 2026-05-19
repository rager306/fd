---
id: T02
parent: S01
milestone: M010-84qfzu
key_files:
  - .gsd/milestones/M010-84qfzu/slices/S01/tasks/T02-SUMMARY.md
key_decisions:
  - Rank model-preserving FP32 export from the exact `deepvk/USER-bge-m3` snapshot above community BAAI artifacts.
  - Treat BAAI ONNX repositories as implementation references, not replacement candidates.
  - Defer INT8 `USER-BGE-M3-ONNX-INT8` until after a verified FP32 dense comparator exists.
duration: 
verification_result: mixed
completed_at: 2026-05-19T18:30:57.637Z
blocker_discovered: false
---

# T02: Ranked ONNX artifact candidates and export paths; model-preserving FP32 export from deepvk/USER-bge-m3 is the safest path, while community artifacts are reference material or out-of-scope INT8/BAAI variants.

**Ranked ONNX artifact candidates and export paths; model-preserving FP32 export from deepvk/USER-bge-m3 is the safest path, while community artifacts are reference material or out-of-scope INT8/BAAI variants.**

## What Happened

Researched current ONNX candidates and export paths for the M010 S01 provenance task. The target Hugging Face model `deepvk/USER-bge-m3` is a Russian-focused sentence-transformer model that maps text to 1024-dimensional dense vectors, is initialized from `TatonkaHF/bge-m3_en_ru`, and was trained primarily for Russian. Its own model tree contains safetensors/tokenizer/config/SentenceTransformer files but no ONNX artifact, matching the local TEI cache inspection from T01.

Candidate ranking:

1. **Best path: export FP32 ONNX from the exact local `deepvk/USER-bge-m3` snapshot.** This preserves the model and satisfies M010's no-replacement boundary. Use the local revision `0cc6cfe48e260fb0474c753087a69369e88709ae` and local hashes from T01 as provenance anchors. Export should target dense output only first: CLS pooling plus L2 normalization, matching the model card's direct-transformers usage and SentenceTransformer `normalize_embeddings=True` behavior. Optimum supports ONNX export for Sentence Transformers and XLM-Roberta-family architectures, and its CLI/API can export local model paths. This path still needs S03 proof because the exact SentenceTransformer wrapper/pooling/export behavior must be verified against TEI/API output.

2. **Useful reference: `aapot/bge-m3-onnx`.** This is a BAAI/bge-m3 ONNX conversion using HF Optimum. It includes `export_onnx.py`, `bgem3_model.py`, `model.onnx`, external `model.onnx.data`, tokenizer files, and explicit output names: `dense_vecs`, `sparse_vecs`, `colbert_vecs`. Its export config subclasses `XLMRobertaOnnxConfig`; dense output is `last_hidden_state[:, 0]` normalized. This is valuable for output naming, custom config, and normalization behavior, but it is based on `BAAI/bge-m3`, not `deepvk/USER-bge-m3`, and includes sparse/ColBERT heads that the deepvk model card says were not thoroughly evaluated for USER-bge-m3.

3. **Useful reference only: `yuniko-software/bge-m3-onnx`.** This repository demonstrates complete BGE-M3 ONNX conversion and cross-language use, producing tokenizer and model ONNX files plus external data. It targets `BAAI/bge-m3` and full multi-vector functionality, so it is not a drop-in model-preserving artifact. It is useful for cross-language ONNX Runtime implementation patterns and external-data handling.

4. **Not suitable for current M010 FP32 goal: `hotchpotch/vespa-onnx-BAAI-bge-m3-only-dense`.** This provides dense-only BAAI/bge-m3 ONNX variants for Vespa, but the published artifacts are FP16 and INT8, not FP32, and the base model is BAAI. It is useful as evidence that CLS pooling + normalize is a known dense-only BGE-M3 serving pattern, not as an artifact to adopt.

5. **Out of scope for this milestone: `skatzR/USER-BGE-M3-ONNX-INT8`.** This is model-family relevant because it claims an INT8 ONNX version of `deepvk/USER-bge-m3` with 1024-dimensional embeddings and benchmarked cosine similarity vs FP32. However, M010 is explicitly FP32 dense-only and excludes INT8. It can inform future quantization comparison after a verified FP32 baseline exists, but should not be used as the first feasibility artifact.

Required provenance fields for any future ONNX attempt: source repo/model ID, exact revision/commit, export command, exporter package versions, opset, optimization level, dtype, provider, local file paths, file sizes and SHA256 for ONNX/model external-data/tokenizer/config files, input names/shapes/dtypes, output names/shapes/dtypes, pooling strategy, normalization behavior, max sequence length, tokenizer hash, and comparison metrics against TEI/API dense output on fixed Russian/legal probes.

Dense-output compatibility risks: BGE-M3 examples differ in output shape and included heads; some outputs are already normalized while others require post-processing; SentenceTransformer export may hide pooling in wrapper modules; direct `AutoModel` export returns token hidden states rather than final sentence embeddings unless pooling/normalization is included; community artifacts based on BAAI/bge-m3 are model replacements and cannot be accepted as production-equivalent for fd; INT8 artifacts need a separate Russian/legal quality gate.

## Verification

Verified by reading cited model cards/docs/source pages and ranking candidates against M010 boundaries. Sources checked: https://huggingface.co/deepvk/USER-bge-m3, https://huggingface.co/deepvk/USER-bge-m3/tree/main, https://huggingface.co/docs/optimum-onnx/onnx/overview, https://huggingface.co/docs/optimum-onnx/onnx/usage_guides/export_a_model, https://onnxruntime.ai/huggingface, https://huggingface.co/aapot/bge-m3-onnx, https://huggingface.co/aapot/bge-m3-onnx/tree/main, https://huggingface.co/aapot/bge-m3-onnx/raw/main/export_onnx.py, https://huggingface.co/aapot/bge-m3-onnx/raw/main/bgem3_model.py, https://raw.githubusercontent.com/yuniko-software/bge-m3-onnx/main/README.md, https://huggingface.co/hotchpotch/vespa-onnx-BAAI-bge-m3-only-dense, https://huggingface.co/skatzR/USER-BGE-M3-ONNX-INT8.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Fetched and read `deepvk/USER-bge-m3` model card and file tree; found safetensors/tokenizer/config but no ONNX artifact.` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Fetched and read Optimum ONNX overview/export docs; confirmed local model export via Optimum is a supported path for Sentence Transformers/XLM-Roberta-family models.` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Fetched and read `aapot/bge-m3-onnx` model card, file tree, `export_onnx.py`, and `bgem3_model.py`; recorded output names and normalization behavior.` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Fetched and read `yuniko-software/bge-m3-onnx` README/requirements; classified as BAAI multi-vector implementation reference.` | -1 | unknown (coerced from string) | 0ms |
| 5 | `Fetched and read `hotchpotch/vespa-onnx-BAAI-bge-m3-only-dense` and `skatzR/USER-BGE-M3-ONNX-INT8`; classified as reference/out-of-scope for M010 FP32 dense-only.` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None. No production runtime or code symbols were changed.

## Known Issues

No ready FP32 ONNX artifact for the exact `deepvk/USER-bge-m3` model was found. The preferred path still needs S03 export/load proof. Community BAAI artifacts remain reference-only, and INT8 remains out of scope until FP32 dense baseline/comparator exists.

## Files Created/Modified

- `.gsd/milestones/M010-84qfzu/slices/S01/tasks/T02-SUMMARY.md`
