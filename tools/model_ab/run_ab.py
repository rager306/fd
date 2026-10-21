#!/usr/bin/env python3
"""Reproducible A/B rig for embedding models behind an OpenAI-compatible API.

The rig embeds one JSONL pair file with a baseline and a candidate endpoint,
scores retrieval over the whole document pool, and writes a sanitized JSON
report. Raw texts are never copied into the report or into error messages.
"""

from __future__ import annotations

import argparse
import json
import math
import sys
import time
from pathlib import Path
from typing import Any, Sequence

import requests

from ab_common import (
    InputError,
    METRIC_KS,
    cosine,
    load_pairs,
    pairs_sha256,
    percentile,
    rank_of,
    summarize_ranks,
    unique_in_order,
    vector_set_sha256,
)

SCHEMA_VERSION = 1
DEFAULT_BATCH_SIZE = 16
DEFAULT_TOP_K = 10
DEFAULT_TIMEOUT = 120.0


def parse_args(argv: Sequence[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Compare two OpenAI-compatible embedding endpoints on one JSONL pair file.",
        epilog="The report contains counts, hashes, metrics and latency; it never contains raw texts.",
    )
    parser.add_argument("--baseline-url", required=True, help="Baseline base URL, e.g. http://localhost:30080.")
    parser.add_argument("--candidate-url", required=True, help="Candidate base URL, e.g. http://localhost:30081.")
    parser.add_argument("--pairs", required=True, type=Path, help="JSONL file with {id, query, positive, negatives} rows.")
    parser.add_argument("--out", required=True, type=Path, help="Sanitized JSON report path.")
    parser.add_argument("--api-key", default="", help="Bearer token sent to both endpoints (optional).")
    parser.add_argument(
        "--dimensions",
        type=int,
        default=None,
        help="Requested output dimensions. Must be supported by both models; fd accepts 512 or 1024 only.",
    )
    parser.add_argument(
        "--top-k",
        type=int,
        default=DEFAULT_TOP_K,
        help="Ranking depth cap. Must be at least 10 because nDCG@10 is reported.",
    )
    parser.add_argument(
        "--batch-size",
        type=int,
        default=DEFAULT_BATCH_SIZE,
        help="Inputs per /v1/embeddings request. Use the same value on both sides.",
    )
    parser.add_argument(
        "--timeout", type=float, default=DEFAULT_TIMEOUT, help="Per-request timeout in seconds."
    )
    parser.add_argument(
        "--warmup",
        action="store_true",
        help="Send one extra request per endpoint and exclude it from latency.",
    )
    parser.add_argument("--baseline-model", default="", help="Optional model id sent to the baseline endpoint.")
    parser.add_argument("--candidate-model", default="", help="Optional model id sent to the candidate endpoint.")
    parser.add_argument(
        "--baseline-query-prefix",
        default="",
        help="Prefix applied to baseline queries. Use the model card's instruction, if any.",
    )
    parser.add_argument(
        "--candidate-query-prefix",
        default="",
        help="Prefix applied to candidate queries. Use the model card's instruction, if any.",
    )
    parser.add_argument(
        "--baseline-document-prefix",
        default="",
        help="Prefix applied to baseline documents. Agreement is skipped when the two document prefixes differ.",
    )
    parser.add_argument(
        "--candidate-document-prefix",
        default="",
        help="Prefix applied to candidate documents. Agreement is skipped when the two document prefixes differ.",
    )
    args = parser.parse_args(argv)
    if args.dimensions is not None and args.dimensions <= 0:
        parser.error("--dimensions must be positive")
    if args.top_k < max(METRIC_KS):
        parser.error(f"--top-k must be at least {max(METRIC_KS)} because nDCG@10 is reported")
    if args.batch_size <= 0:
        parser.error("--batch-size must be positive")
    if args.timeout <= 0:
        parser.error("--timeout must be positive")
    return args


class EmbeddingClient:
    """Batched OpenAI-compatible embedding client with per-batch latency accounting."""

    def __init__(
        self,
        name: str,
        base_url: str,
        *,
        model: str = "",
        dimensions: int | None = None,
        api_key: str = "",
        batch_size: int = DEFAULT_BATCH_SIZE,
        timeout: float = DEFAULT_TIMEOUT,
    ) -> None:
        self.name = name
        self.url = base_url.rstrip("/") + "/v1/embeddings"
        self.model = model
        self.dimensions = dimensions
        self.api_key = api_key
        self.batch_size = batch_size
        self.timeout = timeout
        self.session = requests.Session()
        self.latency_ms: list[float] = []
        self.dimensions_seen: int | None = None

    def embed(self, texts: Sequence[str]) -> list[list[float]]:
        vectors: list[list[float]] = []
        for start in range(0, len(texts), self.batch_size):
            vectors.extend(self._request(list(texts[start : start + self.batch_size])))
        return vectors

    def _request(self, batch: list[str]) -> list[list[float]]:
        payload: dict[str, Any] = {"input": batch}
        if self.model:
            payload["model"] = self.model
        if self.dimensions is not None:
            payload["dimensions"] = self.dimensions
        headers = {"Content-Type": "application/json"}
        if self.api_key:
            headers["Authorization"] = f"Bearer {self.api_key}"
        started = time.perf_counter()
        try:
            response = self.session.post(self.url, json=payload, headers=headers, timeout=self.timeout)
        except requests.RequestException as exc:
            raise RuntimeError(f"{self.name}: embedding request failed ({type(exc).__name__})") from None
        self.latency_ms.append((time.perf_counter() - started) * 1000.0)
        if response.status_code != 200:
            raise RuntimeError(f"{self.name}: HTTP {response.status_code}{self._error_code(response)}")
        return self._parse(response, len(batch))

    @staticmethod
    def _error_code(response: requests.Response) -> str:
        """Extract only the OpenAI error code/type; never echo messages or bodies."""
        try:
            body = response.json()
        except ValueError:
            return ""
        if not isinstance(body, dict):
            return ""
        error = body.get("error")
        if not isinstance(error, dict):
            return ""
        for key in ("code", "type"):
            value = error.get(key)
            if isinstance(value, str) and 0 < len(value) <= 64:
                return f" ({key}={value})"
        return ""

    def _parse(self, response: requests.Response, expected: int) -> list[list[float]]:
        try:
            body = response.json()
        except ValueError:
            raise RuntimeError(f"{self.name}: response is not valid JSON") from None
        data = body.get("data") if isinstance(body, dict) else None
        if not isinstance(data, list) or len(data) != expected:
            got = len(data) if isinstance(data, list) else "none"
            raise RuntimeError(f"{self.name}: response item count {got}, expected {expected}")
        indexed = [isinstance(item, dict) and isinstance(item.get("index"), int) for item in data]
        if any(indexed):
            if not all(indexed):
                raise RuntimeError(f"{self.name}: response mixes indexed and unindexed items")
            indexes = [item["index"] for item in data]
            if sorted(indexes) != list(range(expected)):
                raise RuntimeError(
                    f"{self.name}: response indexes are not a permutation of 0..{expected - 1}"
                )
            data = [item for _, item in sorted(zip(indexes, data), key=lambda pair: pair[0])]
        vectors: list[list[float]] = []
        for item in data:
            embedding = item.get("embedding") if isinstance(item, dict) else None
            if not isinstance(embedding, list) or not embedding:
                raise RuntimeError(f"{self.name}: response item without a numeric embedding")
            try:
                vector = [float(value) for value in embedding]
            except (TypeError, ValueError):
                raise RuntimeError(f"{self.name}: response embedding contains non-numeric values") from None
            if not all(math.isfinite(value) for value in vector):
                raise RuntimeError(f"{self.name}: response embedding contains non-finite values")
            if self.dimensions is not None and len(vector) != self.dimensions:
                raise RuntimeError(
                    f"{self.name}: returned {len(vector)} dimensions, requested {self.dimensions}"
                )
            if self.dimensions_seen is None:
                self.dimensions_seen = len(vector)
            elif self.dimensions_seen != len(vector):
                raise RuntimeError(
                    f"{self.name}: mixed dimensions {self.dimensions_seen} and {len(vector)}"
                )
            vectors.append(vector)
        return vectors


def evaluate_model(
    client: EmbeddingClient,
    pairs: Sequence[dict[str, Any]],
    documents: Sequence[str],
    top_k: int,
    query_prefix: str,
    document_prefix: str,
) -> tuple[dict[str, Any], list[list[float]]]:
    query_texts = unique_in_order(query_prefix + pair["query"] for pair in pairs)
    document_texts = unique_in_order(document_prefix + document for document in documents)
    query_vectors = dict(zip(query_texts, client.embed(query_texts)))
    document_vectors = dict(zip(document_texts, client.embed(document_texts)))
    document_matrix = [document_vectors[document_prefix + document] for document in documents]
    document_index = {document: index for index, document in enumerate(documents)}

    ranks = [
        rank_of(
            query_vectors[query_prefix + pair["query"]],
            document_matrix,
            document_index[pair["positive"]],
            top_k,
        )
        for pair in pairs
    ]

    retrieval = summarize_ranks(ranks)
    retrieval.update({"queries": len(pairs), "documents": len(documents), "top_k": top_k})
    rounded = {
        key: round(value, 6) if isinstance(value, float) else value for key, value in retrieval.items()
    }
    report = {
        "retrieval": rounded,
        "latency_ms": {
            "p50": round(percentile(client.latency_ms, 0.50), 3),
            "p95": round(percentile(client.latency_ms, 0.95), 3),
            "batches": len(client.latency_ms),
        },
        "texts": {"queries": len(query_texts), "documents": len(documents)},
        "dimensions": client.dimensions_seen,
        "ranks": [
            {"id": pair["id"], "rank": rank} for pair, rank in zip(pairs, ranks)
        ],
        "embeddings_sha256": vector_set_sha256(
            [query_vectors[text] for text in query_texts]
            + [document_vectors[text] for text in document_texts]
        ),
    }
    return report, document_matrix


def compute_agreement(
    baseline_documents: Sequence[Sequence[float]],
    candidate_documents: Sequence[Sequence[float]],
    comparable: bool = True,
) -> dict[str, Any]:
    if not comparable:
        return {"status": "skipped", "reason": "document_prefix_mismatch"}
    if not baseline_documents or len(baseline_documents) != len(candidate_documents):
        return {"status": "skipped", "reason": "document_count_mismatch"}
    baseline_dimensions = len(baseline_documents[0])
    candidate_dimensions = len(candidate_documents[0])
    if baseline_dimensions != candidate_dimensions:
        return {
            "status": "skipped",
            "reason": "dimension_mismatch",
            "baseline_dimensions": baseline_dimensions,
            "candidate_dimensions": candidate_dimensions,
        }
    values = [cosine(left, right) for left, right in zip(baseline_documents, candidate_documents)]
    if not all(math.isfinite(value) for value in values):
        return {"status": "skipped", "reason": "non_finite_cosine"}
    return {
        "status": "ok",
        "scope": "documents",
        "count": len(values),
        "dimensions": baseline_dimensions,
        "cosine": {
            "mean": round(sum(values) / len(values), 8),
            "min": round(min(values), 8),
            "p1": round(percentile(values, 0.01), 8),
        },
    }


def run_ab(args: argparse.Namespace) -> dict[str, Any]:
    pairs = load_pairs(args.pairs)
    documents = unique_in_order(
        text for pair in pairs for text in [pair["positive"], *pair["negatives"]]
    )
    baseline = EmbeddingClient(
        "baseline",
        args.baseline_url,
        model=args.baseline_model,
        dimensions=args.dimensions,
        api_key=args.api_key,
        batch_size=args.batch_size,
        timeout=args.timeout,
    )
    candidate = EmbeddingClient(
        "candidate",
        args.candidate_url,
        model=args.candidate_model,
        dimensions=args.dimensions,
        api_key=args.api_key,
        batch_size=args.batch_size,
        timeout=args.timeout,
    )

    if args.warmup:
        for client in (baseline, candidate):
            client.embed(["warmup"])
            client.latency_ms.clear()

    baseline_report, baseline_documents = evaluate_model(
        baseline,
        pairs,
        documents,
        args.top_k,
        args.baseline_query_prefix,
        args.baseline_document_prefix,
    )
    candidate_report, candidate_documents = evaluate_model(
        candidate,
        pairs,
        documents,
        args.top_k,
        args.candidate_query_prefix,
        args.candidate_document_prefix,
    )

    baseline_metrics = baseline_report["retrieval"]
    candidate_metrics = candidate_report["retrieval"]
    delta = {
        f"recall_at_{k}_pp": round(
            (candidate_metrics[f"recall_at_{k}"] - baseline_metrics[f"recall_at_{k}"]) * 100.0, 4
        )
        for k in METRIC_KS
    }
    delta["mrr"] = round(candidate_metrics["mrr"] - baseline_metrics["mrr"], 6)
    delta["ndcg_at_10"] = round(candidate_metrics["ndcg_at_10"] - baseline_metrics["ndcg_at_10"], 6)

    return {
        "schema_version": SCHEMA_VERSION,
        "pairs": {
            "file": args.pairs.name,
            "count": len(pairs),
            "documents": len(documents),
            "sha256": pairs_sha256(pairs),
        },
        "config": {
            "baseline_url": args.baseline_url,
            "candidate_url": args.candidate_url,
            "baseline_model": args.baseline_model or None,
            "candidate_model": args.candidate_model or None,
            "dimensions_requested": args.dimensions,
            "batch_size": args.batch_size,
            "top_k": args.top_k,
            "timeout_seconds": args.timeout,
            "warmup_excluded": args.warmup,
            "query_prefixes_applied": bool(args.baseline_query_prefix or args.candidate_query_prefix),
            "document_prefixes_applied": bool(
                args.baseline_document_prefix or args.candidate_document_prefix
            ),
        },
        "baseline": baseline_report,
        "candidate": candidate_report,
        "delta_candidate_minus_baseline": delta,
        "agreement": compute_agreement(
            baseline_documents,
            candidate_documents,
            args.baseline_document_prefix == args.candidate_document_prefix,
        ),
        "raw_texts_in_report": False,
    }


def main(argv: Sequence[str] | None = None) -> int:
    args = parse_args(argv)
    try:
        report = run_ab(args)
    except (InputError, RuntimeError, OSError) as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 2
    try:
        args.out.parent.mkdir(parents=True, exist_ok=True)
        args.out.write_text(
            json.dumps(report, ensure_ascii=False, indent=2, sort_keys=True) + "\n", encoding="utf-8"
        )
    except OSError as exc:
        print(f"error: cannot write report ({exc.strerror or type(exc).__name__})", file=sys.stderr)
        return 2
    print(f"report written to {args.out}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
