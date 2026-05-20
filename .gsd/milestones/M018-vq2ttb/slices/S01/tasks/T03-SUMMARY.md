---
id: T03
parent: S01
milestone: M018-vq2ttb
key_files:
  - benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt
key_decisions:
  - Tagged Go ONNX 1024 passes the full legal retrieval gate under the current strict thresholds.
  - Chunking is not the immediate next requirement for this selected corpus gate; next gate should be performance and packaging/operational validation for the 1024 path, while chunking remains a future policy for unbounded documents.
duration: 
verification_result: passed
completed_at: 2026-05-20T07:38:33.609Z
blocker_discovered: false
---

# T03: Ran the full ONNX 1024 legal quality gate and it passed strict thresholds.

**Ran the full ONNX 1024 legal quality gate and it passed strict thresholds.**

## What Happened

Ran the full legal retrieval evaluator against TEI and tagged Go ONNX configured with max sequence length 1024. The gate passed strict thresholds: minimum cross-backend cosine was 0.99989883, top-1 agreement was 1.0, mean overlap@5 was 0.997701, and ONNX recall ratio was 1.0. The artifact uses isolated namespace `m018-onnx-1024-legal-quality` and contains no raw legal text.

## Verification

Evaluator exited 0 with PASS verdict and artifact hygiene check passed with no raw legal text leaks.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests python tools/evaluate_legal_retrieval.py ... --onnx-runtime-label tagged-onnx-hf-1024 --onnx-cache-namespace m018-onnx-1024-legal-quality` | 0 | ✅ pass — strict legal gate PASS; artifact written | 75400ms |
| 2 | `python artifact hygiene check for benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt` | 0 | ✅ pass — m018_s01_artifact_hygiene=pass; raw_legal_text_leaks=0 | 0ms |

## Deviations

None.

## Known Issues

This is a quality gate only. Performance, memory, Docker/CI packaging, artifact distribution, and production rollout remain unvalidated.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt`
