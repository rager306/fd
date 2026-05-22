# M040 S03 Legal Model Quick Gate

## Effective Configuration

```json
{
  "baseline": {
    "api_url_label": "http://127.0.0.1:8000",
    "cache_namespace": "m040-s03-baseline-deepvk",
    "dimensions": 1024,
    "model": "deepvk/USER-bge-m3",
    "runtime_label": "current-same-host-api-tei"
  },
  "batch_size": 16,
  "candidate_count": 2,
  "corpus_path": "tests/44-FZ-2026-articles.jsonl",
  "corpus_sha256": "de03cda6b266085a9b1f2376afcb9dffbb00fec922dee1f1553cadcfb6d03869",
  "corpus_stats": {
    "articles": 94,
    "candidate_documents_invalid": 197,
    "candidate_documents_non_invalid": 1655,
    "candidate_documents_total": 1852,
    "clauses": 912,
    "parts": 668,
    "subclauses": 272,
    "title_queries_total": 76
  },
  "dry_run": false,
  "generated_at": "2026-05-22T08:12:43.813473+00:00",
  "max_docs": 128,
  "max_self_queries": 32,
  "max_title_queries": 32,
  "raw_text_logged": false,
  "redaction_status": "sanitized_hashes_only_no_raw_legal_text",
  "script_version": 1,
  "selected": {
    "document_chars": {
      "count": 128,
      "mean": 421.03125,
      "min": 15.0,
      "p50": 241.0,
      "p95": 1577.0
    },
    "document_hashes_sample": [
      "8861126648f5eaa9e4b94159283e5d3ee63e7864270fa882f7dec3d75ecd13a8",
      "028887bdd3a97a69cd86eaa9e09f36df013ee0bfa9253e744f49a6485ecf82d5",
      "c044b6c23d73c0e0457c400095a0177e70f7c1dafc68c4fc839f4df814300056",
      "164210d205da67271a3cb2e087f0ad4190a686070650e2debbed282d2f69e37a",
      "a0b2f919487baa17ce23dc1bfae5955cee83d934bb5c59b1aca38713dfdcaebf"
    ],
    "documents": 128,
    "queries": 42,
    "query_chars": {
      "count": 42,
      "mean": 433.5,
      "min": 20.0,
      "p50": 244.0,
      "p95": 1310.0
    },
    "query_hashes_sample": [
      "55103533a8ea0b149e38b10b89806ee13294cc66dd5d24aeef9fa5db50fb8297",
      "9901c5d33147ccc88ae097c26708c472656b341c28576ea6dabe37a9e6dd8f33",
      "657c4561100c121ca92dd64616a2bf5db9352e14e05e490b20a505929d40a04c",
      "3fe823d6b6b5ed2f9b618c004f7682d765907018f05a00782988353c19d65ed6",
      "cf1d23848bc57acd72a8e515511e7790b8b0b0f124d0d708f4ef9f90aae72014"
    ],
    "self_document_queries": 32,
    "title_queries": 10
  },
  "thresholds": {
    "min_candidate_mrr_ratio": 0.98,
    "min_candidate_recall_ratio": 0.98
  },
  "timeout_seconds": 10.0,
  "top_k": 5
}
```

## Baseline Status

```json
{
  "api_url_label": "http://127.0.0.1:8000",
  "cache_namespace": "m040-s03-baseline-deepvk",
  "expected_dimensions": 1024,
  "health": null,
  "metrics": null,
  "model": "deepvk/USER-bge-m3",
  "runtime_label": "current-same-host-api-tei",
  "smoke_embedding": null,
  "stop_reason": "baseline: RuntimeError: /health missing runtime object"
}
```

## Candidate Results

```json
[
  {
    "api_url_label": "http://127.0.0.1:18001",
    "cache_namespace": "m040-s03-candidate-bge-m3",
    "expected_dimensions": 1024,
    "health": {
      "ok": null,
      "phase": "not_reached"
    },
    "metrics": null,
    "model": "BAAI/bge-m3",
    "operational_compatibility": "failed_closed",
    "outcome": "defer_candidate",
    "role": "candidate_1",
    "runtime_label": "candidate-bge-m3-separate-endpoint-required",
    "smoke_embedding": {
      "ok": null,
      "phase": "not_reached"
    },
    "stop_reason": "baseline_unavailable: RuntimeError: /health missing runtime object"
  },
  {
    "api_url_label": "http://127.0.0.1:18002",
    "cache_namespace": "m040-s03-candidate-multilingual-e5-large",
    "expected_dimensions": 1024,
    "health": {
      "ok": null,
      "phase": "not_reached"
    },
    "metrics": null,
    "model": "intfloat/multilingual-e5-large",
    "operational_compatibility": "failed_closed",
    "outcome": "defer_candidate",
    "role": "candidate_2",
    "runtime_label": "candidate-multilingual-e5-large-separate-endpoint-required",
    "smoke_embedding": {
      "ok": null,
      "phase": "not_reached"
    },
    "stop_reason": "baseline_unavailable: RuntimeError: /health missing runtime object"
  }
]
```

## Cross-Model Cosine/Parity

```json
{
  "applicable": false,
  "reason": "Different embedding models are compared by retrieval metrics and operational compatibility; cross-model cosine/top-1 parity is not an acceptance metric."
}
```

## Verdict

defer_candidate

## Redaction

```json
{
  "raw_text_logged": false,
  "statement": "Raw legal corpus text and smoke payload text are intentionally excluded; only counts, IDs, dimensions, metrics, and SHA-256 hashes are rendered."
}
```

Raw legal corpus text is intentionally excluded from this artifact.
