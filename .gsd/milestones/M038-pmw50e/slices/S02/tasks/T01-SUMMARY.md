---
id: T01
parent: S02
milestone: M038-pmw50e
key_files:
  - benchmark-results/fd-legal-retrieval-m038-go-onnx-target-runtime.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T10:42:57.727Z
blocker_discovered: false
---

# T01: Ran the Go target-runtime legal retrieval gate successfully through actual Go endpoints.

**Ran the Go target-runtime legal retrieval gate successfully through actual Go endpoints.**

## What Happened

Started the Go ONNX API with namespace `m038-go-onnx-legal` and ran the legal retrieval evaluator against actual endpoints: TEI/default Go API on port 8000 and Go ONNX API on port 18000. The gate passed with minimum cross-backend cosine `0.99989883`, top-1 agreement `1.0`, mean overlap@5 `0.997701`, and ONNX recall ratio `1.0`. Raw legal text remained excluded. Stopped the server and confirmed port 18000 is clean.

## Verification

Legal evaluator passed and artifact checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests python tools/evaluate_legal_retrieval.py ... --tei-api-url http://localhost:8000 --onnx-api-url http://localhost:18000 --onnx-cache-namespace m038-go-onnx-legal` | 0 | ✅ pass — legal gate PASS, min cosine 0.99989883, top1 1.0 | 179300ms |
| 2 | `legal artifact marker/leak checks` | 0 | ✅ pass — PASS marker and namespace present, no secret markers or signed URLs | 6600ms |
| 3 | `stop server and port check` | 0 | ✅ pass — no background processes, port_18000_clean | 0ms |

## Deviations

None.

## Known Issues

Legal gate demonstrates parity on the selected corpus, not absolute human relevance quality. Performance gate remains pending.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m038-go-onnx-target-runtime.txt`
