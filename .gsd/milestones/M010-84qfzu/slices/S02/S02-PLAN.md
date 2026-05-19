# S02: Dense comparator baseline

**Goal:** Create a repeatable dense-output comparator baseline for the current TEI-backed API, using fixed Russian/legal probes and sanitized output, so future ONNX FP32 exports can be compared by shape, normalization, and cosine similarity without changing production runtime.
**Demo:** After this, fd has a repeatable dense-output comparator for current TEI baseline and future ONNX candidates.

## Must-Haves

- A comparator exists that calls the current API/TEI baseline for fixed Russian/legal probes.
- Comparator output records dimensions, finite-value checks, L2 norm checks, stable vector hashes, and pairwise cosine similarities.
- Raw probe texts are not printed in artifacts; labels and char counts are acceptable.
- Comparator artifact is saved under `benchmark-results/` and can be used by S03.
- No ONNX/runtime default is introduced.

## Proof Level

- This slice proves: Executable comparator against live TEI/API baseline plus persisted sanitized benchmark artifact.

## Integration Closure

Produces a baseline comparator script and benchmark artifact consumed by S03 ONNX export/load attempts. S03 must compare ONNX dense output against this baseline before making any adapter recommendation.

## Verification

- Adds explicit comparator output fields for future debugging: model, dimensions, probe labels/char counts, L2 norms, finite-value checks, vector hashes, and cosine similarity metrics without printing raw input texts.

## Tasks

- [x] **T01: Define dense comparator contract and probes** `est:small`
  Define the minimal comparator contract and probe set. Use non-sensitive Russian/legal-style fixed probes, assign stable labels, and document expected output fields: API URL, model, dimensions, probe label/length, finite values, L2 norm, vector hash, and cosine similarities. Ensure raw texts will not be printed in benchmark artifacts.
  - Files: `tools/compare_dense_embeddings.py`, `benchmark-results/fd-dense-comparator-m010-s02.txt`
  - Verify: Review comparator contract/probes; confirm raw texts are not emitted by design.

- [x] **T02: Implement TEI API dense comparator** `est:medium`
  Implement a Python comparator script that queries the current API for the fixed probes, validates output dimensions and finite floats, checks L2 normalization tolerance, computes deterministic vector hashes, and computes pairwise cosine similarity. Keep runtime configurable through env vars like `COMPARE_API_URL`, `COMPARE_MODEL`, and `COMPARE_DIMENSIONS`, with safe defaults matching current fd behavior.
  - Files: `tools/compare_dense_embeddings.py`
  - Verify: `uv run --python 3.13 --with requests python tools/compare_dense_embeddings.py --output benchmark-results/fd-dense-comparator-m010-s02.txt` exits 0 against healthy local stack.

- [x] **T03: Capture TEI dense comparator baseline** `est:small`
  Run the comparator against the live Docker stack and save the sanitized baseline artifact. Verify artifact structure, absence of raw probe texts, expected dimensionality, finite values, normalized vectors, and stable summary fields. This artifact becomes S03's baseline for ONNX comparison.
  - Files: `benchmark-results/fd-dense-comparator-m010-s02.txt`
  - Verify: Comparator command exits 0; parser/grep checks confirm required sections and no raw probe text leakage.

## Files Likely Touched

- tools/compare_dense_embeddings.py
- benchmark-results/fd-dense-comparator-m010-s02.txt
