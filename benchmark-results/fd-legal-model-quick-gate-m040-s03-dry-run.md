# M040 S03 Legal Model Quick Gate

## Effective Configuration

```json
{
  "baseline": {
    "api_url_label": "not-called-dry-run",
    "cache_namespace": "default-or-current",
    "dimensions": 1024,
    "model": "deepvk/USER-bge-m3",
    "runtime_label": "tei-default"
  },
  "batch_size": 16,
  "candidate_count": 1,
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
  "dry_run": true,
  "generated_at": "2026-05-22T08:16:46.097277+00:00",
  "max_docs": 32,
  "max_self_queries": 8,
  "max_title_queries": 8,
  "raw_text_logged": false,
  "redaction_status": "sanitized_hashes_only_no_raw_legal_text",
  "script_version": 1,
  "selected": {
    "document_chars": {
      "count": 32,
      "mean": 547.125,
      "min": 35.0,
      "p50": 449.0,
      "p95": 1577.0
    },
    "document_hashes_sample": [
      "8861126648f5eaa9e4b94159283e5d3ee63e7864270fa882f7dec3d75ecd13a8",
      "028887bdd3a97a69cd86eaa9e09f36df013ee0bfa9253e744f49a6485ecf82d5",
      "c044b6c23d73c0e0457c400095a0177e70f7c1dafc68c4fc839f4df814300056",
      "164210d205da67271a3cb2e087f0ad4190a686070650e2debbed282d2f69e37a",
      "a0b2f919487baa17ce23dc1bfae5955cee83d934bb5c59b1aca38713dfdcaebf"
    ],
    "documents": 32,
    "queries": 11,
    "query_chars": {
      "count": 11,
      "mean": 212.090909,
      "min": 35.0,
      "p50": 52.0,
      "p95": 1577.0
    },
    "query_hashes_sample": [
      "55103533a8ea0b149e38b10b89806ee13294cc66dd5d24aeef9fa5db50fb8297",
      "9901c5d33147ccc88ae097c26708c472656b341c28576ea6dabe37a9e6dd8f33",
      "657c4561100c121ca92dd64616a2bf5db9352e14e05e490b20a505929d40a04c",
      "8861126648f5eaa9e4b94159283e5d3ee63e7864270fa882f7dec3d75ecd13a8",
      "028887bdd3a97a69cd86eaa9e09f36df013ee0bfa9253e744f49a6485ecf82d5"
    ],
    "self_document_queries": 8,
    "title_queries": 3
  },
  "thresholds": {
    "min_candidate_mrr_ratio": 0.98,
    "min_candidate_recall_ratio": 0.98
  },
  "timeout_seconds": 60.0,
  "top_k": 5
}
```

## Baseline Status

```json
{
  "api_url_label": "not-called-dry-run",
  "cache_namespace": "default-or-current",
  "expected_dimensions": 1024,
  "health": {
    "ok": false,
    "phase": "dry_run",
    "stop_reason": "dry_run_availability_only_no_endpoint_calls"
  },
  "metrics": null,
  "model": "deepvk/USER-bge-m3",
  "runtime_label": "tei-default",
  "smoke_embedding": {
    "ok": false,
    "phase": "dry_run",
    "stop_reason": "dry_run_availability_only_no_endpoint_calls"
  }
}
```

## Candidate Results

```json
[
  {
    "api_url_label": "not-called-dry-run",
    "cache_namespace": "m040-s03-candidate-bge-m3",
    "expected_dimensions": 1024,
    "health": {
      "ok": false,
      "phase": "dry_run",
      "stop_reason": "dry_run_availability_only_no_endpoint_calls"
    },
    "metrics": null,
    "model": "BAAI/bge-m3",
    "operational_compatibility": "deferred",
    "outcome": "defer_candidate",
    "role": "candidate_1",
    "runtime_label": "candidate-bge-m3",
    "smoke_embedding": {
      "ok": false,
      "phase": "dry_run",
      "stop_reason": "dry_run_availability_only_no_endpoint_calls"
    },
    "stop_reason": "dry_run_availability_only_no_endpoint_calls"
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
