---
id: T03
parent: S01
milestone: M038-pmw50e
key_files:
  - benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T10:36:20.654Z
blocker_discovered: false
---

# T03: Verified local Go ONNX API health and embedding smoke through the actual target runtime.

**Verified local Go ONNX API health and embedding smoke through the actual target runtime.**

## What Happened

Started the actual Go fd API with `onnx hf_tokenizers`, current ONNX artifact, native tokenizer, ONNX Runtime library, and isolated cache namespace `m038-go-onnx-smoke`. Verified `/health` runtime metadata and `/v1/embeddings` through HTTP. The response produced a 1024-dimensional embedding with L2 norm `0.99999983`, backend `onnx`, artifact and tokenizer verification true, and the expected isolated namespace. Stopped the server and confirmed port 18000 is clean.

## Verification

Go ONNX API smoke passed; outcome checks passed; server stopped.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `bg_shell start Go API on port 18000 with EMBEDDING_BACKEND=onnx and EMBEDDING_CACHE_VERSION=m038-go-onnx-smoke` | 0 | ✅ pass — server ready on port 18000 | 0ms |
| 2 | `HTTP /health and /v1/embeddings smoke via driver script` | 0 | ✅ pass — embedding_length=1024, norm=0.99999983, backend=onnx, artifact/tokenizer verified | 10600ms |
| 3 | `outcome leak/marker checks` | 0 | ✅ pass — required markers present, no raw probe text, no secrets/signed URLs | 9500ms |
| 4 | `stop server and port check` | 0 | ✅ pass — no background processes, port_18000_clean | 0ms |

## Deviations

None.

## Known Issues

This is Go API smoke only. Legal retrieval and performance target-runtime gates still need S02.

## Files Created/Modified

- `benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt`
