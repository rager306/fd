# S03: ONNX FP32 CPU feasibility

**Goal:** Attempt a local FP32 dense-only ONNX export/load for the exact `deepvk/USER-bge-m3` snapshot and compare candidate output against the S02 TEI/API baseline, without changing production runtime defaults.
**Demo:** After this, we know whether FP32 BGE-M3 ONNX can be exported/downloaded and loaded on CPU in this environment.

## Must-Haves

- Local ONNX artifact workspace is ignored and does not stage large files.
- Export attempt uses the exact local `deepvk/USER-bge-m3` snapshot, not a BAAI/community replacement.
- Export/load attempt records package versions, command, input/output names, shapes, and hashes.
- If ONNX loads, dense output is compared against S02 TEI baseline for fixed probes.
- Production API/runtime defaults remain unchanged.

## Proof Level

- This slice proves: Actual local export/load attempt with ONNX Runtime CPU EP and comparator evidence; documented blocker if export/load fails.

## Integration Closure

Produces either a loadable FP32 dense ONNX candidate with metadata/comparison evidence, or a documented blocker explaining why the model-preserving ONNX path is not currently feasible. S04 consumes this evidence for the milestone recommendation.

## Verification

- Records export command, dependency versions, source model hashes, ONNX/external-data hashes, input/output names and shapes, provider, and comparison metrics against S02 baseline.

## Tasks

- [x] **T01: Prepare ONNX artifact workspace and provenance** `est:small`
  Prepare `.gsd/runtime/onnx/m010-s03/` as the local ignored artifact workspace and capture source model provenance: model path, revision, source hashes, tokenizer/config hashes, available disk, and dependency plan. Verify no large artifacts are staged.
  - Files: `.gsd/runtime/onnx/m010-s03/`
  - Verify: Workspace exists under ignored `.gsd/runtime`; source hashes recorded; `git status --short` shows no large ONNX/model artifacts staged.

- [x] **T02: Attempt FP32 dense-only ONNX export** `est:large`
  Implement and run a model-preserving FP32 dense-only ONNX export attempt using the local model snapshot. Prefer an explicit wrapper around `AutoModel` that outputs `dense_vecs = normalize(last_hidden_state[:,0])`, with dynamic batch/sequence axes, CPU export, and metadata capture. Store generated ONNX artifacts under `.gsd/runtime/onnx/m010-s03/`.
  - Files: `tools/export_user_bge_m3_dense_onnx.py`, `.gsd/runtime/onnx/m010-s03/`
  - Verify: `uv run --python 3.13 --with torch --with transformers --with onnx --with onnxruntime --with safetensors python tools/export_user_bge_m3_dense_onnx.py --model-path tei-models/deepvk--USER-bge-m3 --output-dir .gsd/runtime/onnx/m010-s03` exits 0 or records a structured failure artifact.

- [x] **T03: Compare ONNX output to TEI baseline or record blocker** `est:medium`
  If export succeeds, load the ONNX with ONNX Runtime CPU EP, run the same fixed probes as S02, and compare dimensions, finite values, L2 norms, vector hashes, and cosine similarity against the TEI baseline. If export failed, synthesize the blocker evidence instead. Save a concise tracked benchmark/result artifact under `benchmark-results/`.
  - Files: `tools/compare_onnx_dense_embeddings.py`, `benchmark-results/fd-onnx-fp32-m010-s03.txt`
  - Verify: Comparison artifact exists and states PASS/FAIL/BLOCKED with output shape/hash/cosine evidence; raw probe texts are not emitted.

## Files Likely Touched

- .gsd/runtime/onnx/m010-s03/
- tools/export_user_bge_m3_dense_onnx.py
- tools/compare_onnx_dense_embeddings.py
- benchmark-results/fd-onnx-fp32-m010-s03.txt
