---
id: T02
parent: S02
milestone: M010-84qfzu
key_files:
  - tools/compare_dense_embeddings.py
  - benchmark-results/fd-dense-comparator-m010-s02.txt
key_decisions:
  - Created a dedicated `tools/compare_dense_embeddings.py` script rather than extending `benchmark.py`, keeping benchmark and equivalence-comparator concerns separate.
  - Hashed vectors as little-endian float32 bytes to align with service/cache precision and avoid Python float repr instability.
  - Made normalization tolerance configurable via `COMPARE_NORM_TOLERANCE`, defaulting to `0.02` for current 1024-dimensional baseline checks.
duration: 
verification_result: passed
completed_at: 2026-05-19T18:36:54.846Z
blocker_discovered: false
---

# T02: Implemented and ran the TEI API dense comparator script with sanitized output and baseline vector checks.

**Implemented and ran the TEI API dense comparator script with sanitized output and baseline vector checks.**

## What Happened

Implemented `tools/compare_dense_embeddings.py`. The script calls `/v1/embeddings` with fixed non-sensitive Russian/legal-style probes, validates response count, dimensions, finite float values, and L2 normalization, computes deterministic float32 vector hashes, and prints pairwise cosine similarities. It supports `COMPARE_API_URL`, `COMPARE_MODEL`, `COMPARE_DIMENSIONS`, and `COMPARE_NORM_TOLERANCE`, plus CLI overrides and `--output`. The generated markdown artifact intentionally excludes raw probe texts and records labels/character counts only.

## Verification

`uv run --python 3.13 --with requests python -m py_compile tools/compare_dense_embeddings.py` passed. `uv run --python 3.13 --with requests python tools/compare_dense_embeddings.py --output benchmark-results/fd-dense-comparator-m010-s02.txt` exited 0 against the live local API and produced a PASS artifact with 5 probes, 1024 dimensions, finite values, L2 norms within tolerance, vector hashes, and pairwise cosine similarities.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests python -m py_compile tools/compare_dense_embeddings.py` | 0 | ✅ pass | 0ms |
| 2 | `uv run --python 3.13 --with requests python tools/compare_dense_embeddings.py --output benchmark-results/fd-dense-comparator-m010-s02.txt` | 0 | ✅ pass | 0ms |

## Deviations

The T02 implementation command also generated the baseline artifact planned for T03; T03 will still verify and record that artifact formally.

## Known Issues

The comparator currently targets the current API baseline only. S03 will need either an ONNX-runner input mode or a separate ONNX comparison wrapper to compare candidate vectors against this baseline.

## Files Created/Modified

- `tools/compare_dense_embeddings.py`
- `benchmark-results/fd-dense-comparator-m010-s02.txt`
