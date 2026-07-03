---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Capture TEI dense comparator baseline

Run the comparator against the live Docker stack and save the sanitized baseline artifact. Verify artifact structure, absence of raw probe texts, expected dimensionality, finite values, normalized vectors, and stable summary fields. This artifact becomes S03's baseline for ONNX comparison.

## Inputs

- `tools/compare_dense_embeddings.py`
- `benchmark-results/fd-environment-inxi-m008.txt`

## Expected Output

- `benchmark-results/fd-dense-comparator-m010-s02.txt`

## Verification

Comparator command exits 0; parser/grep checks confirm required sections and no raw probe text leakage.

## Observability Impact

Persists the TEI baseline evidence needed to debug ONNX mismatches in S03.
