---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Start legal gate runtimes

Verify default TEI health and start tagged ONNX API on port 18000 with `EMBEDDING_CACHE_VERSION=m015-onnx-legal-quality`, then capture health/startup evidence.

## Inputs

- `tools/evaluate_legal_retrieval.py`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`

## Expected Output

- `Task summary with runtime evidence`

## Verification

TEI and tagged ONNX health endpoints return ok; no stale background process.

## Observability Impact

Ensures evaluation targets the intended runtimes with isolated ONNX namespace.
