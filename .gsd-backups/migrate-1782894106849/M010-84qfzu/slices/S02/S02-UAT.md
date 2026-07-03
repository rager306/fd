# S02: Dense comparator baseline — UAT

**Milestone:** M010-84qfzu
**Written:** 2026-05-19T18:38:04.405Z

# S02 UAT — Dense comparator baseline

## Checks

- [x] Comparator script exists at `tools/compare_dense_embeddings.py`.
- [x] Comparator runs with Python 3.13 via uv and exits 0 against the live local API.
- [x] Baseline artifact exists at `benchmark-results/fd-dense-comparator-m010-s02.txt`.
- [x] Artifact records 1024 expected dimensions, PASS verdict, vector hashes, L2 norms, and pairwise cosine similarities.
- [x] Artifact states `raw_probe_texts_logged=false`.
- [x] Validation confirmed no raw probe text strings from script constants appear in the artifact.
- [x] No production runtime default or ONNX adapter was introduced.

## UAT Result

Pass. S03 has a concrete TEI/API baseline artifact and a command to regenerate it before comparing ONNX output.

