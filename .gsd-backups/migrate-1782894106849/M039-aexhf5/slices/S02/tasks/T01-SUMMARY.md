---
id: T01
parent: S02
milestone: M039-aexhf5
key_files:
  - benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T11:27:05.351Z
blocker_discovered: false
---

# T01: Packaged Docker ONNX legal retrieval gate passed through actual packaged Go endpoint.

**Packaged Docker ONNX legal retrieval gate passed through actual packaged Go endpoint.**

## What Happened

Started packaged image `fd-api:onnx1024-m039` with `ONNX_RUNTIME_SHA256` and namespace `m039-docker-legal`. Health confirmed backend `onnx`, artifact/tokenizer/runtime library verified, CPU provider, dimensions 1024, and the expected namespace. Ran the legal retrieval evaluator against TEI/default Go API on port 8000 and packaged Go ONNX API on port 18000. Gate passed with minimum cross-backend cosine `0.99989883`, top-1 agreement `1.0`, mean overlap@5 `0.997701`, and ONNX recall ratio `1.0`. Artifact leak checks passed; container stopped and port 18000 is clean.

## Verification

Legal evaluator, artifact checks, and cleanup passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker run fd-api:onnx1024-m039 with namespace m039-docker-legal and ONNX_RUNTIME_SHA256` | 0 | ✅ pass — container started | 4000ms |
| 2 | `packaged legal health precheck` | 0 | ✅ pass — artifact/tokenizer/runtime verified, namespace m039-docker-legal | 3900ms |
| 3 | `uv run --python 3.13 --with requests python tools/evaluate_legal_retrieval.py ... --onnx-runtime-label docker-onnx-go-api --onnx-cache-namespace m039-docker-legal` | 0 | ✅ pass — PASS, min cosine 0.99989883, top1 1.0 | 167300ms |
| 4 | `legal artifact leak checks` | 0 | ✅ pass — missing_required=0, leak_markers=0, signed_url_like=0 | 4300ms |
| 5 | `docker rm -f fd-onnx-m039-legal && port check` | 0 | ✅ pass — port_18000_clean | 4300ms |

## Deviations

None.

## Known Issues

Legal gate demonstrates TEI-vs-ONNX parity on selected corpus, not absolute human relevance quality.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt`
