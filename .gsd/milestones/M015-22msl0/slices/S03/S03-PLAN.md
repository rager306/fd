# S03: Run legal retrieval quality gate

**Goal:** Run the evaluator against TEI default and tagged ONNX HF tokenizer path using the 44-FZ corpus.
**Demo:** After this, ONNX has a real first-pass Russian/legal retrieval parity result.

## Must-Haves

- TEI and tagged ONNX endpoints are healthy.
- Redis namespaces are isolated.
- Evaluation artifact is produced.
- Metrics meet or fail stated thresholds with evidence.
- Tagged server cleanup is verified.

## Proof Level

- This slice proves: Runtime evaluation artifact plus hygiene checks.

## Integration Closure

Produces the quality evidence needed before packaging/tuning decisions.

## Verification

- Records endpoint config, corpus hash, metrics, and runtime cleanup evidence.

## Tasks

- [x] **T01: Start legal gate runtimes** `est:small`
  Verify default TEI health and start tagged ONNX API on port 18000 with `EMBEDDING_CACHE_VERSION=m015-onnx-legal-quality`, then capture health/startup evidence.
  - Verify: TEI and tagged ONNX health endpoints return ok; no stale background process.

- [x] **T02: Run live legal retrieval gate** `est:medium`
  Run `tools/evaluate_legal_retrieval.py` live against TEI and tagged ONNX using the 44-ФЗ corpus and write `benchmark-results/fd-legal-retrieval-m015-s03.txt`.
  - Files: `benchmark-results/fd-legal-retrieval-m015-s03.txt`
  - Verify: Evaluator exits 0 for pass or 2 for quality fail and artifact records verdict.

- [x] **T03: Verify legal gate artifact and cleanup** `est:small`
  Verify the legal gate artifact hygiene, record pass/fail evidence, cleanup tagged ONNX server, and run GitNexus scope check.
  - Files: `benchmark-results/fd-legal-retrieval-m015-s03.txt`
  - Verify: Artifact has no raw legal text leaks, runtime cleanup confirmed, GitNexus detect_changes passes.

## Files Likely Touched

- benchmark-results/fd-legal-retrieval-m015-s03.txt
