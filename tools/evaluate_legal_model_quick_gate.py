#!/usr/bin/env python3
"""Bounded legal embedding-model quick gate.

This tool compares replacement-candidate retrieval metrics against the current
same-host embedding service without changing the service API contract. Artifacts
are sanitized: raw corpus text and request payloads are never rendered.
"""

from __future__ import annotations

import argparse
from dataclasses import dataclass
from datetime import datetime, timezone
import json
from pathlib import Path
import sys
from typing import Any

import requests

from evaluate_legal_retrieval import (  # reuse existing corpus/evaluator patterns
    Document,
    Query,
    build_queries,
    load_corpus,
    ranking,
    recall_at_k,
    reciprocal_rank,
    select_documents,
    sha256_file,
    summarize,
)

SCRIPT_VERSION = 1
DEFAULT_TOP_K = 5
ALLOWED_OUTCOMES = {"keep_current", "reject_candidate", "defer_candidate"}
SMOKE_TEXT = "sanitized legal retrieval smoke text"


@dataclass(frozen=True)
class EndpointSpec:
    role: str
    api_url: str | None
    model: str
    runtime_label: str
    dimensions: int
    cache_namespace: str


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Bounded quick gate for Russian/legal embedding model candidates.")
    parser.add_argument("--corpus", type=Path, required=True, help="Structured legal JSONL corpus path.")
    parser.add_argument("--output", type=Path, required=True, help="Sanitized markdown artifact output path.")
    parser.add_argument("--dry-run", action="store_true", help="Render availability-only artifact; do not call endpoints.")

    parser.add_argument("--baseline-api-url", help="Current same-host embedding API URL for live runs.")
    parser.add_argument("--baseline-model", required=True, help="Runtime model expected from baseline /health.")
    parser.add_argument("--baseline-runtime-label", required=True)
    parser.add_argument("--baseline-dimensions", type=int, required=True)
    parser.add_argument("--baseline-cache-namespace", default="default-or-current")

    parser.add_argument("--candidate-api-url", action="append", default=[], help="Candidate embedding API URL; repeat per candidate for live runs.")
    parser.add_argument("--candidate-model", action="append", required=True, help="Candidate model ID; repeat per candidate, max 2.")
    parser.add_argument("--candidate-runtime-label", action="append", required=True)
    parser.add_argument("--candidate-dimensions", action="append", type=int, required=True)
    parser.add_argument("--candidate-cache-namespace", action="append", required=True)
    parser.add_argument("--candidate-stop-reason", action="append", default=[], help="Optional dry-run availability stop reason per candidate.")

    parser.add_argument("--max-docs", type=int, default=128)
    parser.add_argument("--max-title-queries", type=int, default=32)
    parser.add_argument("--max-self-queries", type=int, default=32)
    parser.add_argument("--batch-size", type=int, default=16)
    parser.add_argument("--top-k", type=int, default=DEFAULT_TOP_K)
    parser.add_argument("--timeout-seconds", type=float, default=60.0)
    parser.add_argument("--min-candidate-recall-ratio", type=float, default=0.98)
    parser.add_argument("--min-candidate-mrr-ratio", type=float, default=0.98)
    return parser.parse_args()


def sanitize_error(exc: BaseException) -> str:
    text = f"{type(exc).__name__}: {exc}"
    for marker in ("Bearer ", "token=", "api_key=", "password="):
        if marker.lower() in text.lower():
            return f"{type(exc).__name__}: redacted sensitive error detail"
    return text[:240]


def validate_candidate_args(args: argparse.Namespace) -> None:
    count = len(args.candidate_model)
    if count > 2:
        raise ValueError(f"at most two candidates are allowed; got {count}")
    if count < 1:
        raise ValueError("at least one candidate is required")
    required_lists = {
        "candidate-runtime-label": args.candidate_runtime_label,
        "candidate-dimensions": args.candidate_dimensions,
        "candidate-cache-namespace": args.candidate_cache_namespace,
    }
    if not args.dry_run:
        required_lists["candidate-api-url"] = args.candidate_api_url
        if not args.baseline_api_url:
            raise ValueError("--baseline-api-url is required for live runs")
    for name, values in required_lists.items():
        if len(values) != count:
            raise ValueError(f"--{name} count ({len(values)}) must match --candidate-model count ({count})")
    if args.candidate_stop_reason and len(args.candidate_stop_reason) not in {0, count}:
        raise ValueError("--candidate-stop-reason count must be omitted or match candidate count")


def candidate_specs(args: argparse.Namespace) -> list[EndpointSpec]:
    return [
        EndpointSpec(
            role=f"candidate_{idx + 1}",
            api_url=args.candidate_api_url[idx] if idx < len(args.candidate_api_url) else None,
            model=args.candidate_model[idx],
            runtime_label=args.candidate_runtime_label[idx],
            dimensions=args.candidate_dimensions[idx],
            cache_namespace=args.candidate_cache_namespace[idx],
        )
        for idx in range(len(args.candidate_model))
    ]


def selected_summary(docs: list[Document], queries: list[Query]) -> dict[str, Any]:
    return {
        "documents": len(docs),
        "queries": len(queries),
        "title_queries": sum(1 for query in queries if query.kind == "title"),
        "self_document_queries": sum(1 for query in queries if query.kind == "self_document"),
        "document_chars": summarize([float(doc.chars) for doc in docs]),
        "query_chars": summarize([float(query.chars) for query in queries]),
        "document_hashes_sample": [doc.text_sha256 for doc in docs[:5]],
        "query_hashes_sample": [query.text_sha256 for query in queries[:5]],
    }


def health_probe(spec: EndpointSpec, timeout: float) -> dict[str, Any]:
    if not spec.api_url:
        raise RuntimeError("missing api_url")
    response = requests.get(f"{spec.api_url.rstrip('/')}/health", timeout=timeout)
    response.raise_for_status()
    body = response.json()
    runtime = body.get("runtime") if isinstance(body, dict) else None
    if not isinstance(runtime, dict):
        raise RuntimeError("/health missing runtime object")
    observed_model = runtime.get("model")
    observed_dimensions = runtime.get("dimensions")
    observed_cache = runtime.get("cache_namespace")
    if observed_model != spec.model:
        raise RuntimeError(f"/health model mismatch: observed {observed_model!r}, expected {spec.model!r}")
    if int(observed_dimensions) != spec.dimensions:
        raise RuntimeError(f"/health dimensions mismatch: observed {observed_dimensions!r}, expected {spec.dimensions}")
    return {
        "ok": True,
        "runtime": {
            "backend": runtime.get("backend"),
            "model": observed_model,
            "dimensions": observed_dimensions,
            "production_default": runtime.get("production_default"),
            "cache_namespace": observed_cache,
        },
    }


def request_embeddings(spec: EndpointSpec, texts: list[str], batch_size: int, timeout: float) -> list[list[float]]:
    if not spec.api_url:
        raise RuntimeError("missing api_url")
    vectors: list[list[float]] = []
    url = f"{spec.api_url.rstrip('/')}/v1/embeddings"
    for start in range(0, len(texts), batch_size):
        batch = texts[start : start + batch_size]
        payload = {"model": spec.model, "input": batch, "dimensions": spec.dimensions}
        response = requests.post(url, json=payload, timeout=timeout)
        response.raise_for_status()
        body = response.json()
        data = body.get("data") if isinstance(body, dict) else None
        if not isinstance(data, list) or len(data) != len(batch):
            raise RuntimeError("embedding response count mismatch")
        for item in data:
            embedding = item.get("embedding") if isinstance(item, dict) else None
            if not isinstance(embedding, list):
                raise RuntimeError("embedding item missing vector")
            vector = [float(value) for value in embedding]
            if len(vector) != spec.dimensions:
                raise RuntimeError(f"embedding dimension mismatch: got {len(vector)}, expected {spec.dimensions}")
            vectors.append(vector)
    return vectors


def smoke_probe(spec: EndpointSpec, timeout: float) -> dict[str, Any]:
    vector = request_embeddings(spec, [SMOKE_TEXT], batch_size=1, timeout=timeout)[0]
    return {"ok": True, "embedding_dimensions": len(vector)}


def retrieval_metrics(docs: list[Document], queries: list[Query], doc_vectors: list[list[float]], query_vectors: list[list[float]], top_k: int) -> dict[str, Any]:
    doc_ids = [doc.doc_id for doc in docs]
    recalls: list[float] = []
    mrrs: list[float] = []
    worst: list[dict[str, Any]] = []
    for query, query_vector in zip(queries, query_vectors):
        ranked_indexes = ranking(query_vector, doc_vectors)
        ranked_ids = [doc_ids[idx] for idx in ranked_indexes]
        recall = recall_at_k(ranked_ids, query.relevant_doc_ids, top_k)
        mrr = reciprocal_rank(ranked_ids, query.relevant_doc_ids)
        recalls.append(recall)
        mrrs.append(mrr)
        worst.append(
            {
                "query_id": query.query_id,
                "query_kind": query.kind,
                "query_chars": query.chars,
                "query_text_sha256": query.text_sha256,
                f"recall_at_{top_k}": round(recall, 4),
                "mrr": round(mrr, 4),
                "top1_doc_id": ranked_ids[0] if ranked_ids else None,
            }
        )
    worst.sort(key=lambda row: (row["mrr"], row[f"recall_at_{top_k}"]))
    return {f"mean_recall_at_{top_k}": round(sum(recalls) / len(recalls), 6) if recalls else 0.0, "mrr": round(sum(mrrs) / len(mrrs), 6) if mrrs else 0.0, "worst_queries": worst[:10]}


def metric_ratios(candidate: dict[str, Any], baseline: dict[str, Any], top_k: int) -> dict[str, float]:
    recall_key = f"mean_recall_at_{top_k}"
    baseline_recall = float(baseline.get(recall_key) or 0.0)
    baseline_mrr = float(baseline.get("mrr") or 0.0)
    return {
        "recall_ratio": round((float(candidate.get(recall_key) or 0.0) / baseline_recall) if baseline_recall else 1.0, 6),
        "mrr_ratio": round((float(candidate.get("mrr") or 0.0) / baseline_mrr) if baseline_mrr else 1.0, 6),
    }


def build_config(args: argparse.Namespace, corpus_hash: str, corpus_stats: dict[str, Any], docs: list[Document], queries: list[Query]) -> dict[str, Any]:
    return {
        "script_version": SCRIPT_VERSION,
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "corpus_path": str(args.corpus),
        "corpus_sha256": corpus_hash,
        "raw_text_logged": False,
        "redaction_status": "sanitized_hashes_only_no_raw_legal_text",
        "dry_run": bool(args.dry_run),
        "top_k": args.top_k,
        "max_docs": args.max_docs,
        "max_title_queries": args.max_title_queries,
        "max_self_queries": args.max_self_queries,
        "batch_size": args.batch_size,
        "timeout_seconds": args.timeout_seconds,
        "candidate_count": len(args.candidate_model),
        "thresholds": {
            "min_candidate_recall_ratio": args.min_candidate_recall_ratio,
            "min_candidate_mrr_ratio": args.min_candidate_mrr_ratio,
        },
        "baseline": {
            "api_url_label": args.baseline_api_url or "not-called-dry-run",
            "model": args.baseline_model,
            "runtime_label": args.baseline_runtime_label,
            "dimensions": args.baseline_dimensions,
            "cache_namespace": args.baseline_cache_namespace,
        },
        "corpus_stats": corpus_stats,
        "selected": selected_summary(docs, queries),
    }


def build_dry_run_candidate(spec: EndpointSpec, provided_reason: str | None) -> dict[str, Any]:
    return {
        "role": spec.role,
        "api_url_label": spec.api_url or "not-called-dry-run",
        "model": spec.model,
        "runtime_label": spec.runtime_label,
        "expected_dimensions": spec.dimensions,
        "cache_namespace": spec.cache_namespace,
        "health": {"ok": False, "phase": "dry_run", "stop_reason": provided_reason or "dry_run_availability_only_no_endpoint_calls"},
        "smoke_embedding": {"ok": False, "phase": "dry_run", "stop_reason": provided_reason or "dry_run_availability_only_no_endpoint_calls"},
        "metrics": None,
        "operational_compatibility": "deferred",
        "outcome": "defer_candidate",
        "stop_reason": provided_reason or "dry_run_availability_only_no_endpoint_calls",
    }


def fail_candidate(spec: EndpointSpec, phase: str, exc: BaseException, outcome: str = "defer_candidate") -> dict[str, Any]:
    return {
        "role": spec.role,
        "api_url_label": spec.api_url or "missing",
        "model": spec.model,
        "runtime_label": spec.runtime_label,
        "expected_dimensions": spec.dimensions,
        "cache_namespace": spec.cache_namespace,
        "health": {"ok": False, "phase": phase} if phase == "health" else {"ok": None, "phase": "not_reached"},
        "smoke_embedding": {"ok": False, "phase": phase} if phase == "smoke" else {"ok": None, "phase": "not_reached"},
        "metrics": None,
        "operational_compatibility": "failed_closed",
        "outcome": outcome,
        "stop_reason": f"{phase}: {sanitize_error(exc)}",
    }


def build_result(args: argparse.Namespace) -> dict[str, Any]:
    validate_candidate_args(args)
    documents_all, title_queries, corpus_stats = load_corpus(args.corpus)
    docs = select_documents(documents_all, args.max_docs)
    queries = build_queries(title_queries, docs, args.max_title_queries, args.max_self_queries)
    if not docs or not queries:
        raise ValueError("corpus selection produced no documents or queries")
    config = build_config(args, sha256_file(args.corpus), corpus_stats, docs, queries)
    specs = candidate_specs(args)

    result: dict[str, Any] = {
        "config": config,
        "baseline_status": {
            "api_url_label": args.baseline_api_url or "not-called-dry-run",
            "model": args.baseline_model,
            "runtime_label": args.baseline_runtime_label,
            "expected_dimensions": args.baseline_dimensions,
            "cache_namespace": args.baseline_cache_namespace,
            "health": {"ok": False, "phase": "dry_run", "stop_reason": "dry_run_availability_only_no_endpoint_calls"} if args.dry_run else None,
            "smoke_embedding": {"ok": False, "phase": "dry_run", "stop_reason": "dry_run_availability_only_no_endpoint_calls"} if args.dry_run else None,
            "metrics": None,
        },
        "candidates": [],
        "cross_model_cosine_parity": {
            "applicable": False,
            "reason": "Different embedding models are compared by retrieval metrics and operational compatibility; cross-model cosine/top-1 parity is not an acceptance metric.",
        },
        "final_outcome": "defer_candidate",
        "redaction": {"raw_text_logged": False, "statement": "Raw legal corpus text and smoke payload text are intentionally excluded; only counts, IDs, dimensions, metrics, and SHA-256 hashes are rendered."},
    }

    if args.dry_run:
        reasons = args.candidate_stop_reason or []
        result["candidates"] = [build_dry_run_candidate(spec, reasons[idx] if idx < len(reasons) else None) for idx, spec in enumerate(specs)]
        result["final_outcome"] = "defer_candidate"
        return result

    baseline = EndpointSpec("baseline", args.baseline_api_url, args.baseline_model, args.baseline_runtime_label, args.baseline_dimensions, args.baseline_cache_namespace)
    try:
        result["baseline_status"]["health"] = health_probe(baseline, args.timeout_seconds)
        result["baseline_status"]["smoke_embedding"] = smoke_probe(baseline, args.timeout_seconds)
        baseline_doc_vectors = request_embeddings(baseline, [doc.text for doc in docs], args.batch_size, args.timeout_seconds)
        baseline_query_vectors = request_embeddings(baseline, [query.text for query in queries], args.batch_size, args.timeout_seconds)
        baseline_metrics = retrieval_metrics(docs, queries, baseline_doc_vectors, baseline_query_vectors, args.top_k)
        result["baseline_status"]["metrics"] = baseline_metrics
    except Exception as exc:  # noqa: BLE001 - gate must fail closed with sanitized artifact.
        result["baseline_status"]["stop_reason"] = f"baseline: {sanitize_error(exc)}"
        result["final_outcome"] = "defer_candidate"
        result["candidates"] = [fail_candidate(spec, "baseline_unavailable", exc) for spec in specs]
        return result

    any_candidate_kept = False
    any_candidate_deferred = False
    for spec in specs:
        candidate_row: dict[str, Any] = {
            "role": spec.role,
            "api_url_label": spec.api_url or "missing",
            "model": spec.model,
            "runtime_label": spec.runtime_label,
            "expected_dimensions": spec.dimensions,
            "cache_namespace": spec.cache_namespace,
        }
        try:
            candidate_row["health"] = health_probe(spec, args.timeout_seconds)
            candidate_row["smoke_embedding"] = smoke_probe(spec, args.timeout_seconds)
            candidate_doc_vectors = request_embeddings(spec, [doc.text for doc in docs], args.batch_size, args.timeout_seconds)
            candidate_query_vectors = request_embeddings(spec, [query.text for query in queries], args.batch_size, args.timeout_seconds)
            metrics = retrieval_metrics(docs, queries, candidate_doc_vectors, candidate_query_vectors, args.top_k)
            ratios = metric_ratios(metrics, result["baseline_status"]["metrics"], args.top_k)
            metrics["against_baseline"] = ratios
            candidate_row["metrics"] = metrics
            candidate_row["operational_compatibility"] = "compatible"
            if ratios["recall_ratio"] >= args.min_candidate_recall_ratio and ratios["mrr_ratio"] >= args.min_candidate_mrr_ratio:
                candidate_row["outcome"] = "keep_current"
                candidate_row["stop_reason"] = "candidate_not_rejected_by_quick_gate; human/legal-domain review still required before replacement"
                any_candidate_kept = True
            else:
                candidate_row["outcome"] = "reject_candidate"
                candidate_row["stop_reason"] = "retrieval_metrics_below_threshold"
        except requests.HTTPError as exc:
            candidate_row = fail_candidate(spec, "health_or_embedding_http", exc, outcome="defer_candidate")
            any_candidate_deferred = True
        except Exception as exc:  # noqa: BLE001 - fail closed per threat model.
            phase = "candidate_probe_or_metrics"
            candidate_row = fail_candidate(spec, phase, exc, outcome="defer_candidate")
            any_candidate_deferred = True
        result["candidates"].append(candidate_row)

    if any_candidate_deferred:
        result["final_outcome"] = "defer_candidate"
    elif any_candidate_kept:
        result["final_outcome"] = "keep_current"
    else:
        result["final_outcome"] = "reject_candidate"
    return result


def render_markdown(result: dict[str, Any]) -> str:
    return "\n".join(
        [
            "# M040 S03 Legal Model Quick Gate",
            "",
            "## Effective Configuration",
            "",
            "```json",
            json.dumps(result["config"], ensure_ascii=False, indent=2, sort_keys=True),
            "```",
            "",
            "## Baseline Status",
            "",
            "```json",
            json.dumps(result["baseline_status"], ensure_ascii=False, indent=2, sort_keys=True),
            "```",
            "",
            "## Candidate Results",
            "",
            "```json",
            json.dumps(result["candidates"], ensure_ascii=False, indent=2, sort_keys=True),
            "```",
            "",
            "## Cross-Model Cosine/Parity",
            "",
            "```json",
            json.dumps(result["cross_model_cosine_parity"], ensure_ascii=False, indent=2, sort_keys=True),
            "```",
            "",
            "## Verdict",
            "",
            result["final_outcome"],
            "",
            "## Redaction",
            "",
            "```json",
            json.dumps(result["redaction"], ensure_ascii=False, indent=2, sort_keys=True),
            "```",
            "",
            "Raw legal corpus text is intentionally excluded from this artifact.",
            "",
        ]
    )


def main() -> int:
    args = parse_args()
    try:
        result = build_result(args)
        markdown = render_markdown(result)
        args.output.parent.mkdir(parents=True, exist_ok=True)
        args.output.write_text(markdown, encoding="utf-8")
        print(markdown)
        if result["final_outcome"] == "reject_candidate":
            return 2
        return 0
    except Exception as exc:  # noqa: BLE001 - produce sanitized blocked artifact when possible.
        blocked = "\n".join(
            [
                "# M040 S03 Legal Model Quick Gate",
                "",
                "## Verdict",
                "",
                "defer_candidate",
                "",
                "## Stop Reason",
                "",
                f"setup: {sanitize_error(exc)}",
                "",
                "Raw legal corpus text is intentionally excluded from this artifact.",
                "",
            ]
        )
        args.output.parent.mkdir(parents=True, exist_ok=True)
        args.output.write_text(blocked, encoding="utf-8")
        print(blocked, file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
