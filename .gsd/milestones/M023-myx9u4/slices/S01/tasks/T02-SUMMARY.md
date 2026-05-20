---
id: T02
parent: S01
milestone: M023-myx9u4
key_files:
  - benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt
key_decisions:
  - Packaged ONNX legal gate uses runtime label `packaged-onnx1024-docker` and cache namespace `m023-onnx-docker-legal`.
  - Strict M018-equivalent cross-backend cosine threshold `0.999` remains in force.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:56:35.922Z
blocker_discovered: false
---

# T02: Packaged ONNX Docker 1024 passed the Russian/legal retrieval quality gate.

**Packaged ONNX Docker 1024 passed the Russian/legal retrieval quality gate.**

## What Happened

Ran the legal retrieval evaluator against TEI baseline at port 8000 and the packaged ONNX Docker image at port 18000. The packaged ONNX image passed the strict legal quality gate with minimum cross-backend cosine `0.99989883`, top-1 agreement `1.0`, mean overlap@5 `0.997701`, and ONNX recall ratio `1.0`. The artifact records sanitized configuration and metrics with `raw_text_logged=false`.

## Verification

Evaluator exited 0 and wrote `benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt` with verdict PASS.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests python tools/evaluate_legal_retrieval.py --corpus tests/44-FZ-2026-articles.jsonl --output benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt --tei-api-url http://localhost:8000 --onnx-api-url http://localhost:18000 --tei-runtime-label tei-default-compose --onnx-runtime-label packaged-onnx1024-docker --tei-cache-namespace default-or-current --onnx-cache-namespace m023-onnx-docker-legal --min-cross-backend-cosine 0.999` | 0 | ✅ pass — Verdict PASS; minimum_overall=0.99989883; top1_agreement=1.0; mean_overlap_at_5=0.997701; onnx_recall_ratio=1.0 | 154000ms |

## Deviations

None.

## Known Issues

The corpus remains unlabeled; the artifact proves TEI-vs-ONNX parity and synthetic known-item behavior, not absolute human relevance quality.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt`
