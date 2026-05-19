---
id: T04
parent: S06
milestone: M008-6hnowu
key_files:
  - api/embed/types.go
  - api/cache/redis.go
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:02:17.077Z
blocker_discovered: false
---

# T04: BGE-M3 ONNX multi-output is real, but fd should first use dense-only compatibility; sparse/ColBERT are future hybrid retrieval work.

**BGE-M3 ONNX multi-output is real, but fd should first use dense-only compatibility; sparse/ColBERT are future hybrid retrieval work.**

## What Happened

Verified BGE-M3 multi-output implications. BGE documentation describes BGE-M3 as multi-functional: dense retrieval, sparse lexical retrieval, and multi-vector/ColBERT-style retrieval, while also being multilingual and supporting long contexts. Multiple ONNX model cards (`aapot/bge-m3-onnx`, `philipchung/bge-m3-onnx`, and INT8 variants) claim ONNX outputs dense, sparse, and ColBERT representations in a single pass. For fd's current architecture, the dense output is the only directly compatible output because `api/embed/types.go` and the OpenAI-compatible response shape represent a list of float arrays, and `api/cache/redis.go` stores one dense float32 vector per text/dimension key. Sparse token weights and ColBERT token vectors are potentially valuable for future hybrid retrieval on legal text, but they change both API semantics and cache shape. Recommendation: in a future ONNX spike, first consume only dense output for equivalence and latency. Store sparse/ColBERT only in a separate experimental namespace if a hybrid retrieval milestone is explicitly planned, with separate cache keys, memory accounting, and Russian legal retrieval metrics.

## Verification

Read BGE docs and BGE-M3 ONNX model cards; mapped outputs against fd response/cache code.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Search query: BGE-M3 ONNX dense sparse ColBERT outputs aapot bge-m3-onnx FlagEmbedding ONNX Runtime` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Read: https://bge-model.com/bge/bge_m3.html` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Read: https://huggingface.co/aapot/bge-m3-onnx` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Read: https://huggingface.co/philipchung/bge-m3-onnx` | -1 | unknown (coerced from string) | 0ms |
| 5 | `Read: api/embed/types.go and api/cache/redis.go earlier in S06 planning context` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

fd currently exposes dense OpenAI-compatible embeddings only. Sparse and ColBERT outputs would require new internal data structures, cache schemas, search integration, and quality benchmarks; they should not be added to the current `/v1/embeddings` response by default.

## Files Created/Modified

- `api/embed/types.go`
- `api/cache/redis.go`
