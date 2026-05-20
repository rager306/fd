---
id: T02
parent: S03
milestone: M015-22msl0
key_files:
  - tools/evaluate_legal_retrieval.py
  - benchmark-results/fd-legal-retrieval-m015-s03.txt
key_decisions:
  - Treat the live gate as a quality FAIL because cross-backend cosine minimum is far below threshold even though ranking parity metrics mostly pass.
  - Worst cosine diagnostics are required in legal quality artifacts to make failures actionable.
duration: 
verification_result: mixed
completed_at: 2026-05-20T05:01:54.576Z
blocker_discovered: false
---

# T02: Ran the live Russian/legal retrieval gate; tagged ONNX fails the strict quality gate due severe cross-backend cosine outliers.

**Ran the live Russian/legal retrieval gate; tagged ONNX fails the strict quality gate due severe cross-backend cosine outliers.**

## What Happened

Ran the live legal retrieval gate against TEI and tagged ONNX on 256 selected legal documents and 87 derived queries. The first result revealed an evaluator ID defect for unnumbered subclauses, which was fixed and rerun. The final result is a quality FAIL: ranking parity is mostly strong, but cross-backend cosine has severe outliers on legal texts, especially longer clauses. The artifact includes worst document/query IDs, character counts, hashes, and metrics without raw legal text.

## Verification

Evaluator exited 2 as expected for a quality FAIL and wrote a sanitized artifact with verdict FAIL and worst-case diagnostics.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests python tools/evaluate_legal_retrieval.py ... --max-docs 256 --max-title-queries 32 --max-self-queries 64` | 2 | ⚠️ quality fail — artifact written with verdict FAIL | 6900ms |
| 2 | `python3 -m py_compile tools/evaluate_legal_retrieval.py` | 0 | ✅ pass — evaluator compiles after ID/worst-cosine diagnostics changes | 0ms |

## Deviations

The first live run exposed duplicate `sNone` IDs for unnumbered subclauses. The evaluator was fixed to generate fallback IDs (`sidxN`) and the live gate was rerun. The final artifact still returns FAIL, now with diagnostic worst cross-backend cosine IDs.

## Known Issues

The gate fails primarily on cross-backend cosine for longer legal clauses/subclauses. Ranking parity is near/pass on several metrics: top1 agreement 0.977, mean overlap@5 0.908, ONNX recall ratio 0.991. Cross-backend cosine minimum is 0.369 and query minimum is 0.656, far below threshold 0.999.

## Files Created/Modified

- `tools/evaluate_legal_retrieval.py`
- `benchmark-results/fd-legal-retrieval-m015-s03.txt`
