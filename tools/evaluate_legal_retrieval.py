#!/usr/bin/env python3
"""Russian/legal retrieval parity evaluator for fd embeddings.

The evaluator intentionally writes sanitized artifacts: IDs, counts, hashes,
lengths, and metrics only. Raw legal corpus text is never rendered.
"""

from __future__ import annotations

import argparse
from dataclasses import dataclass
from datetime import datetime, timezone
import hashlib
import json
import math
from pathlib import Path
import statistics
import sys
from typing import Any

import requests

SCRIPT_VERSION = 1
DEFAULT_MODEL = "deepvk/USER-bge-m3"
DEFAULT_DIMENSIONS = 1024


@dataclass(frozen=True)
class Document:
    doc_id: str
    article: str
    kind: str
    text: str
    text_sha256: str
    chars: int
    invalid: bool


@dataclass(frozen=True)
class Query:
    query_id: str
    article: str
    kind: str
    text: str
    text_sha256: str
    chars: int
    relevant_doc_ids: frozenset[str]


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Evaluate TEI-vs-ONNX legal retrieval parity on structured JSONL corpus.")
    parser.add_argument("--corpus", type=Path, required=True, help="Structured legal JSONL corpus path.")
    parser.add_argument("--output", type=Path, required=True, help="Markdown artifact output path.")
    parser.add_argument("--tei-api-url", default="http://localhost:8000", help="Baseline TEI fd API URL.")
    parser.add_argument("--onnx-api-url", default="http://localhost:18000", help="Tagged ONNX fd API URL.")
    parser.add_argument("--model", default=DEFAULT_MODEL)
    parser.add_argument("--dimensions", type=int, default=DEFAULT_DIMENSIONS)
    parser.add_argument("--max-docs", type=int, default=256, help="Maximum non-invalid clause/part docs to evaluate.")
    parser.add_argument("--max-title-queries", type=int, default=64)
    parser.add_argument("--max-self-queries", type=int, default=64)
    parser.add_argument("--batch-size", type=int, default=16)
    parser.add_argument("--top-k", type=int, default=5)
    parser.add_argument("--min-top1-agreement", type=float, default=0.90)
    parser.add_argument("--min-mean-overlap-at-k", type=float, default=0.90)
    parser.add_argument("--min-onnx-recall-ratio", type=float, default=0.98)
    parser.add_argument("--min-cross-backend-cosine", type=float, default=0.999)
    parser.add_argument("--tei-runtime-label", default="tei-default")
    parser.add_argument("--onnx-runtime-label", default="tagged-onnx-hf")
    parser.add_argument("--tei-cache-namespace", default="default-or-current")
    parser.add_argument("--onnx-cache-namespace", default="m015-onnx-legal-quality")
    parser.add_argument("--dry-run", action="store_true", help="Only parse/profile corpus and render evaluator plan; no API calls.")
    return parser.parse_args()


def sha256_text(value: str) -> str:
    return hashlib.sha256(value.encode("utf-8")).hexdigest()


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def stable_number(value: Any, fallback: str) -> str:
    if value is None or value == "":
        return fallback
    return str(value).replace(" ", "_")


def load_corpus(path: Path) -> tuple[list[Document], list[Query], dict[str, Any]]:
    documents: list[Document] = []
    title_queries: list[Query] = []
    articles = 0
    parts = clauses = subclauses = 0
    invalid_docs = 0

    for line in path.read_text(encoding="utf-8").splitlines():
        if not line.strip():
            continue
        obj = json.loads(line)
        articles += 1
        article = stable_number(obj.get("article"), f"line{articles}")
        doc_prefix = stable_number(obj.get("doc_id", "doc"), "doc")
        article_doc_ids: list[str] = []

        for part_index, part in enumerate(obj.get("parts") or [], 1):
            parts += 1
            part_no = stable_number(part.get("number"), f"idx{part_index}")
            part_invalid = bool(part.get("invalid"))
            part_text = part.get("text") or ""
            part_id = f"{doc_prefix}:a{article}:p{part_no}"
            if part_text.strip():
                document = Document(
                    doc_id=part_id,
                    article=article,
                    kind="part",
                    text=part_text,
                    text_sha256=sha256_text(part_text),
                    chars=len(part_text),
                    invalid=part_invalid,
                )
                documents.append(document)
                article_doc_ids.append(document.doc_id)
                invalid_docs += int(part_invalid)

            for clause_index, clause in enumerate(part.get("clauses") or [], 1):
                clauses += 1
                clause_no = stable_number(clause.get("number"), f"idx{clause_index}")
                clause_invalid = part_invalid or bool(clause.get("invalid"))
                clause_text = clause.get("text") or ""
                clause_id = f"{doc_prefix}:a{article}:p{part_no}:c{clause_no}"
                if clause_text.strip():
                    document = Document(
                        doc_id=clause_id,
                        article=article,
                        kind="clause",
                        text=clause_text,
                        text_sha256=sha256_text(clause_text),
                        chars=len(clause_text),
                        invalid=clause_invalid,
                    )
                    documents.append(document)
                    article_doc_ids.append(document.doc_id)
                    invalid_docs += int(clause_invalid)

                for subclause_index, subclause in enumerate(clause.get("subclauses") or [], 1):
                    subclauses += 1
                    sub_no = stable_number(subclause.get("number"), f"idx{subclause_index}")
                    sub_invalid = clause_invalid or bool(subclause.get("invalid"))
                    sub_text = subclause.get("text") or ""
                    sub_id = f"{doc_prefix}:a{article}:p{part_no}:c{clause_no}:s{sub_no}"
                    if sub_text.strip():
                        document = Document(
                            doc_id=sub_id,
                            article=article,
                            kind="subclause",
                            text=sub_text,
                            text_sha256=sha256_text(sub_text),
                            chars=len(sub_text),
                            invalid=sub_invalid,
                        )
                        documents.append(document)
                        article_doc_ids.append(document.doc_id)
                        invalid_docs += int(sub_invalid)

        title = obj.get("title") or ""
        if title.strip() and article_doc_ids:
            title_queries.append(
                Query(
                    query_id=f"{doc_prefix}:a{article}:title",
                    article=article,
                    kind="title",
                    text=title,
                    text_sha256=sha256_text(title),
                    chars=len(title),
                    relevant_doc_ids=frozenset(article_doc_ids),
                )
            )

    non_invalid_documents = [doc for doc in documents if not doc.invalid]
    stats = {
        "articles": articles,
        "parts": parts,
        "clauses": clauses,
        "subclauses": subclauses,
        "candidate_documents_total": len(documents),
        "candidate_documents_non_invalid": len(non_invalid_documents),
        "candidate_documents_invalid": invalid_docs,
        "title_queries_total": len(title_queries),
    }
    return non_invalid_documents, title_queries, stats


def select_documents(documents: list[Document], max_docs: int) -> list[Document]:
    # Prefer clause/subclause granularity for legal retrieval; include parts only when needed.
    preferred = [doc for doc in documents if doc.kind in {"clause", "subclause"}]
    fallback = [doc for doc in documents if doc.kind == "part"]
    selected = (preferred + fallback)[:max_docs]
    return selected


def build_queries(title_queries: list[Query], documents: list[Document], max_title: int, max_self: int) -> list[Query]:
    selected_doc_ids = {doc.doc_id for doc in documents}
    queries: list[Query] = []
    for query in title_queries:
        relevant = frozenset(doc_id for doc_id in query.relevant_doc_ids if doc_id in selected_doc_ids)
        if relevant:
            queries.append(
                Query(
                    query_id=query.query_id,
                    article=query.article,
                    kind=query.kind,
                    text=query.text,
                    text_sha256=query.text_sha256,
                    chars=query.chars,
                    relevant_doc_ids=relevant,
                )
            )
        if len([item for item in queries if item.kind == "title"]) >= max_title:
            break

    for doc in documents[:max_self]:
        queries.append(
            Query(
                query_id=f"self:{doc.doc_id}",
                article=doc.article,
                kind="self_document",
                text=doc.text,
                text_sha256=doc.text_sha256,
                chars=doc.chars,
                relevant_doc_ids=frozenset({doc.doc_id}),
            )
        )
    return queries


def request_embeddings(api_url: str, model: str, dimensions: int, texts: list[str], batch_size: int) -> list[list[float]]:
    vectors: list[list[float]] = []
    url = f"{api_url.rstrip('/')}/v1/embeddings"
    for start in range(0, len(texts), batch_size):
        batch = texts[start : start + batch_size]
        payload = {"model": model, "input": batch, "dimensions": dimensions}
        try:
            response = requests.post(url, json=payload, timeout=180)
            response.raise_for_status()
            body = response.json()
        except requests.RequestException as exc:
            raise RuntimeError(f"embedding request failed for {api_url}: {type(exc).__name__}") from exc
        except ValueError as exc:
            raise RuntimeError(f"embedding response was not JSON for {api_url}") from exc
        data = body.get("data") if isinstance(body, dict) else None
        if not isinstance(data, list) or len(data) != len(batch):
            raise RuntimeError(f"embedding response count mismatch for {api_url}: got {len(data) if isinstance(data, list) else 'invalid'}, expected {len(batch)}")
        for item in data:
            embedding = item.get("embedding") if isinstance(item, dict) else None
            if not isinstance(embedding, list):
                raise RuntimeError("embedding item missing vector")
            vector = [float(value) for value in embedding]
            if len(vector) != dimensions:
                raise RuntimeError(f"embedding dimensions mismatch: got {len(vector)}, expected {dimensions}")
            vectors.append(vector)
    return vectors


def l2_norm(vector: list[float]) -> float:
    return math.sqrt(sum(value * value for value in vector))


def cosine(a: list[float], b: list[float]) -> float:
    denom = l2_norm(a) * l2_norm(b)
    if denom == 0:
        return float("nan")
    return sum(x * y for x, y in zip(a, b)) / denom


def ranking(query: list[float], docs: list[list[float]]) -> list[int]:
    scored = [(idx, cosine(query, doc)) for idx, doc in enumerate(docs)]
    scored.sort(key=lambda item: item[1], reverse=True)
    return [idx for idx, _ in scored]


def recall_at_k(ranked_doc_ids: list[str], relevant: frozenset[str], k: int) -> float:
    if not relevant:
        return 0.0
    return len(set(ranked_doc_ids[:k]) & relevant) / min(len(relevant), k)


def reciprocal_rank(ranked_doc_ids: list[str], relevant: frozenset[str]) -> float:
    for idx, doc_id in enumerate(ranked_doc_ids, 1):
        if doc_id in relevant:
            return 1.0 / idx
    return 0.0


def mean(values: list[float]) -> float:
    return statistics.mean(values) if values else 0.0


def summarize(values: list[float]) -> dict[str, float]:
    if not values:
        return {"count": 0, "mean": 0.0, "min": 0.0, "p50": 0.0, "p95": 0.0}
    ordered = sorted(values)
    return {
        "count": len(values),
        "mean": round(mean(values), 6),
        "min": round(min(values), 6),
        "p50": round(ordered[len(ordered) // 2], 6),
        "p95": round(ordered[min(len(ordered) - 1, int(len(ordered) * 0.95))], 6),
    }


def evaluate_rankings(
    docs: list[Document],
    queries: list[Query],
    tei_doc_vectors: list[list[float]],
    tei_query_vectors: list[list[float]],
    onnx_doc_vectors: list[list[float]],
    onnx_query_vectors: list[list[float]],
    top_k: int,
) -> dict[str, Any]:
    doc_ids = [doc.doc_id for doc in docs]
    top1_agreements: list[float] = []
    overlaps_at_k: list[float] = []
    tei_recalls: list[float] = []
    onnx_recalls: list[float] = []
    tei_mrrs: list[float] = []
    onnx_mrrs: list[float] = []
    worst: list[dict[str, Any]] = []

    for query, tei_query, onnx_query in zip(queries, tei_query_vectors, onnx_query_vectors):
        tei_rank = ranking(tei_query, tei_doc_vectors)
        onnx_rank = ranking(onnx_query, onnx_doc_vectors)
        tei_ranked_ids = [doc_ids[idx] for idx in tei_rank]
        onnx_ranked_ids = [doc_ids[idx] for idx in onnx_rank]
        top1_same = tei_ranked_ids[0] == onnx_ranked_ids[0]
        overlap = len(set(tei_ranked_ids[:top_k]) & set(onnx_ranked_ids[:top_k])) / top_k
        top1_agreements.append(1.0 if top1_same else 0.0)
        overlaps_at_k.append(overlap)
        tei_recall = recall_at_k(tei_ranked_ids, query.relevant_doc_ids, top_k)
        onnx_recall = recall_at_k(onnx_ranked_ids, query.relevant_doc_ids, top_k)
        tei_mrr = reciprocal_rank(tei_ranked_ids, query.relevant_doc_ids)
        onnx_mrr = reciprocal_rank(onnx_ranked_ids, query.relevant_doc_ids)
        tei_recalls.append(tei_recall)
        onnx_recalls.append(onnx_recall)
        tei_mrrs.append(tei_mrr)
        onnx_mrrs.append(onnx_mrr)
        worst.append(
            {
                "query_id": query.query_id,
                "query_kind": query.kind,
                "query_chars": query.chars,
                "tei_top1": tei_ranked_ids[0],
                "onnx_top1": onnx_ranked_ids[0],
                "top1_agreement": top1_same,
                f"overlap_at_{top_k}": round(overlap, 4),
                f"tei_recall_at_{top_k}": round(tei_recall, 4),
                f"onnx_recall_at_{top_k}": round(onnx_recall, 4),
                "tei_mrr": round(tei_mrr, 4),
                "onnx_mrr": round(onnx_mrr, 4),
            }
        )

    worst.sort(key=lambda row: (row[f"overlap_at_{top_k}"], row["onnx_mrr"] - row["tei_mrr"]))
    return {
        "top1_agreement": round(mean(top1_agreements), 6),
        f"mean_overlap_at_{top_k}": round(mean(overlaps_at_k), 6),
        f"tei_mean_recall_at_{top_k}": round(mean(tei_recalls), 6),
        f"onnx_mean_recall_at_{top_k}": round(mean(onnx_recalls), 6),
        "tei_mrr": round(mean(tei_mrrs), 6),
        "onnx_mrr": round(mean(onnx_mrrs), 6),
        "worst_queries": worst[:10],
    }


def build_config(args: argparse.Namespace, corpus_hash: str, corpus_stats: dict[str, Any], docs: list[Document], queries: list[Query]) -> dict[str, Any]:
    return {
        "script_version": SCRIPT_VERSION,
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "corpus_path": str(args.corpus),
        "corpus_sha256": corpus_hash,
        "raw_text_logged": False,
        "model": args.model,
        "dimensions": args.dimensions,
        "top_k": args.top_k,
        "max_docs": args.max_docs,
        "max_title_queries": args.max_title_queries,
        "max_self_queries": args.max_self_queries,
        "batch_size": args.batch_size,
        "tei": {"api_url": args.tei_api_url, "runtime_label": args.tei_runtime_label, "cache_namespace": args.tei_cache_namespace},
        "onnx": {"api_url": args.onnx_api_url, "runtime_label": args.onnx_runtime_label, "cache_namespace": args.onnx_cache_namespace},
        "thresholds": {
            "min_top1_agreement": args.min_top1_agreement,
            f"min_mean_overlap_at_{args.top_k}": args.min_mean_overlap_at_k,
            "min_onnx_recall_ratio": args.min_onnx_recall_ratio,
            "min_cross_backend_cosine": args.min_cross_backend_cosine,
        },
        "corpus_stats": corpus_stats,
        "selected": {
            "documents": len(docs),
            "queries": len(queries),
            "title_queries": sum(1 for query in queries if query.kind == "title"),
            "self_document_queries": sum(1 for query in queries if query.kind == "self_document"),
            "document_chars": summarize([float(doc.chars) for doc in docs]),
            "query_chars": summarize([float(query.chars) for query in queries]),
        },
    }


def worst_cross_backend_cosines(items: list[Any], cosines: list[float], limit: int = 10) -> list[dict[str, Any]]:
    rows = []
    for item, value in zip(items, cosines):
        rows.append(
            {
                "id": getattr(item, "doc_id", getattr(item, "query_id", "unknown")),
                "kind": item.kind,
                "article": item.article,
                "chars": item.chars,
                "cosine": round(value, 8),
                "text_sha256": item.text_sha256,
            }
        )
    rows.sort(key=lambda row: row["cosine"])
    return rows[:limit]


def build_result(args: argparse.Namespace) -> dict[str, Any]:
    documents_all, title_queries, corpus_stats = load_corpus(args.corpus)
    docs = select_documents(documents_all, args.max_docs)
    queries = build_queries(title_queries, docs, args.max_title_queries, args.max_self_queries)
    corpus_hash = sha256_file(args.corpus)
    config = build_config(args, corpus_hash, corpus_stats, docs, queries)

    if args.dry_run:
        return {
            "config": config,
            "dry_run": True,
            "passed": True,
            "verdict": "DRY_RUN",
            "metrics": {},
            "caveat": "No API calls were made; this validates corpus parsing and sanitized artifact shape only.",
        }

    texts_docs = [doc.text for doc in docs]
    texts_queries = [query.text for query in queries]
    tei_doc_vectors = request_embeddings(args.tei_api_url, args.model, args.dimensions, texts_docs, args.batch_size)
    tei_query_vectors = request_embeddings(args.tei_api_url, args.model, args.dimensions, texts_queries, args.batch_size)
    onnx_doc_vectors = request_embeddings(args.onnx_api_url, args.model, args.dimensions, texts_docs, args.batch_size)
    onnx_query_vectors = request_embeddings(args.onnx_api_url, args.model, args.dimensions, texts_queries, args.batch_size)

    doc_cosines = [cosine(left, right) for left, right in zip(tei_doc_vectors, onnx_doc_vectors)]
    query_cosines = [cosine(left, right) for left, right in zip(tei_query_vectors, onnx_query_vectors)]
    ranking_metrics = evaluate_rankings(docs, queries, tei_doc_vectors, tei_query_vectors, onnx_doc_vectors, onnx_query_vectors, args.top_k)
    onnx_recall = ranking_metrics[f"onnx_mean_recall_at_{args.top_k}"]
    tei_recall = ranking_metrics[f"tei_mean_recall_at_{args.top_k}"]
    recall_ratio = onnx_recall / tei_recall if tei_recall else 1.0
    cross_backend_cosine = min(doc_cosines + query_cosines)
    passed = (
        ranking_metrics["top1_agreement"] >= args.min_top1_agreement
        and ranking_metrics[f"mean_overlap_at_{args.top_k}"] >= args.min_mean_overlap_at_k
        and recall_ratio >= args.min_onnx_recall_ratio
        and cross_backend_cosine >= args.min_cross_backend_cosine
    )
    return {
        "config": config,
        "dry_run": False,
        "passed": passed,
        "verdict": "PASS" if passed else "FAIL",
        "metrics": {
            "cross_backend_cosine": {
                "documents": summarize(doc_cosines),
                "queries": summarize(query_cosines),
                "minimum_overall": round(cross_backend_cosine, 8),
                "worst_documents": worst_cross_backend_cosines(docs, doc_cosines),
                "worst_queries": worst_cross_backend_cosines(queries, query_cosines),
            },
            "ranking": ranking_metrics,
            "onnx_recall_ratio": round(recall_ratio, 6),
        },
        "caveat": "Unlabeled corpus: metrics show TEI-vs-ONNX parity and synthetic known-item behavior, not absolute human relevance quality.",
    }


def render_markdown(result: dict[str, Any]) -> str:
    lines = [
        "# M015 Russian Legal Retrieval Gate",
        "",
        "## Effective Configuration",
        "",
        "```json",
        json.dumps(result["config"], ensure_ascii=False, indent=2, sort_keys=True),
        "```",
        "",
        "## Metrics",
        "",
        "```json",
        json.dumps(result["metrics"], ensure_ascii=False, indent=2, sort_keys=True),
        "```",
        "",
        "## Verdict",
        "",
        str(result["verdict"]),
        "",
        "## Caveat",
        "",
        result["caveat"],
        "",
        "Raw legal corpus text is intentionally excluded from this artifact.",
        "",
    ]
    return "\n".join(lines)


def main() -> int:
    args = parse_args()
    try:
        result = build_result(args)
        markdown = render_markdown(result)
        args.output.parent.mkdir(parents=True, exist_ok=True)
        args.output.write_text(markdown, encoding="utf-8")
        print(markdown)
        if result["dry_run"]:
            return 0
        return 0 if result["passed"] else 2
    except Exception as exc:  # noqa: BLE001 - gate should write a blocked artifact.
        blocked = "\n".join(
            [
                "# M015 Russian Legal Retrieval Gate",
                "",
                "## Verdict",
                "",
                "BLOCKED",
                "",
                f"Error type: `{type(exc).__name__}`",
                "",
                f"Error message: `{exc}`",
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
