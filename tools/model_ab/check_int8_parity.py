#!/usr/bin/env python3
"""Numerical parity check between a dense ONNX artifact and its INT8 variant.

The reference artifact may be FP32 or FP16, but both artifacts must be dense
pooled exports with dynamic sequence axes, as produced by
tools/export_user_bge_m3_dense_onnx.py. The check runs ONNX Runtime on the CPU
and is a numerics gate, not a latency benchmark.

Run it in the ONNX export environment (numpy, onnxruntime, transformers), not
in the fd Go environment. Model files and artifacts are never committed.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import math
import sys
from pathlib import Path
from typing import Any, Sequence

from ab_common import (
    InputError,
    cosine,
    load_pairs,
    pairs_sha256,
    percentile,
    rank_of,
    summarize_ranks,
    unique_in_order,
    vector_set_sha256,
)


def parse_args(argv: Sequence[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Compare dense ONNX embeddings of a reference artifact and an INT8 S8S8-QDQ artifact.",
        epilog="The report contains hashes, metrics and acceptance flags; it never contains raw texts.",
    )
    parser.add_argument("--reference-onnx", type=Path, required=True, help="FP32 or FP16 dense ONNX artifact.")
    parser.add_argument("--int8-onnx", type=Path, required=True, help="INT8 S8S8-QDQ dense ONNX artifact.")
    parser.add_argument("--model-path", type=Path, required=True, help="Local tokenizer snapshot path.")
    parser.add_argument("--pairs", type=Path, required=True, help="JSONL pair file in the run_ab schema.")
    parser.add_argument("--output", type=Path, required=True, help="Sanitized JSON report path.")
    parser.add_argument("--max-length", type=int, default=512, help="Tokenizer truncation length.")
    parser.add_argument("--batch-size", type=int, default=16, help="Sequences per ONNX run.")
    parser.add_argument("--top-k", type=int, default=10, help="Ranking depth cap; must be at least 10.")
    parser.add_argument("--min-mean-cosine", type=float, default=0.999, help="Acceptance floor for mean cosine.")
    parser.add_argument("--min-cosine", type=float, default=0.995, help="Acceptance floor for the worst cosine.")
    parser.add_argument(
        "--max-recall-drop-pp",
        type=float,
        default=1.0,
        help="Largest allowed Recall@10 drop in percentage points.",
    )
    args = parser.parse_args(argv)
    if args.max_length <= 0:
        parser.error("--max-length must be positive")
    if args.batch_size <= 0:
        parser.error("--batch-size must be positive")
    if args.top_k < 10:
        parser.error("--top-k must be at least 10 because nDCG@10 is reported")
    return args


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def artifact_info(path: Path) -> dict[str, Any]:
    return {"path": str(path), "size_bytes": path.stat().st_size, "sha256": sha256_file(path)}


def load_sessions(args: argparse.Namespace) -> tuple[Any, Any, tuple[Any, Any]]:
    try:
        import numpy as np
        import onnxruntime as ort
        from transformers import AutoTokenizer
    except ImportError:
        raise InputError(
            "this check needs numpy, onnxruntime and transformers; run it in the ONNX export environment"
        ) from None
    tokenizer = AutoTokenizer.from_pretrained(str(args.model_path), local_files_only=True)
    reference = ort.InferenceSession(str(args.reference_onnx), providers=["CPUExecutionProvider"])
    int8 = ort.InferenceSession(str(args.int8_onnx), providers=["CPUExecutionProvider"])
    return np, tokenizer, (reference, int8)


def embed_all(
    np: Any,
    session: Any,
    tokenizer: Any,
    texts: Sequence[str],
    max_length: int,
    batch_size: int,
) -> list[list[float]]:
    input_names = {item.name for item in session.get_inputs()}
    missing = {"input_ids", "attention_mask"} - input_names
    if missing:
        raise InputError(f"ONNX model is missing inputs: {sorted(missing)}")
    vectors: list[list[float]] = []
    for start in range(0, len(texts), batch_size):
        batch = list(texts[start : start + batch_size])
        encoded = tokenizer(
            batch, padding=True, truncation=True, max_length=max_length, return_tensors="np"
        )
        feeds = {
            "input_ids": encoded["input_ids"].astype(np.int64),
            "attention_mask": encoded["attention_mask"].astype(np.int64),
        }
        if "token_type_ids" in input_names:
            feeds["token_type_ids"] = np.zeros_like(feeds["input_ids"])
        outputs = session.run(None, feeds)
        if len(outputs) != 1:
            raise InputError(f"expected one ONNX output, got {len(outputs)}")
        dense = outputs[0]
        if dense.ndim != 2:
            raise InputError(f"expected a 2D pooled output, got shape {tuple(dense.shape)}")
        vectors.extend(dense.astype(np.float32).tolist())
    return vectors


def retrieval_metrics(
    vectors_by_text: dict[str, list[float]],
    pairs: Sequence[dict[str, Any]],
    documents: Sequence[str],
    top_k: int,
) -> dict[str, Any]:
    document_matrix = [vectors_by_text[document] for document in documents]
    document_index = {document: index for index, document in enumerate(documents)}
    ranks = [
        rank_of(
            vectors_by_text[pair["query"]],
            document_matrix,
            document_index[pair["positive"]],
            top_k,
        )
        for pair in pairs
    ]
    summary = summarize_ranks(ranks)
    summary.update({"queries": len(pairs), "documents": len(documents), "top_k": top_k})
    return {
        key: round(value, 6) if isinstance(value, float) else value for key, value in summary.items()
    }


def run_check(args: argparse.Namespace) -> dict[str, Any]:
    pairs = load_pairs(args.pairs)
    texts = unique_in_order(
        text for pair in pairs for text in [pair["query"], pair["positive"], *pair["negatives"]]
    )
    documents = unique_in_order(
        text for pair in pairs for text in [pair["positive"], *pair["negatives"]]
    )
    np, tokenizer, (reference_session, int8_session) = load_sessions(args)
    reference_vectors = embed_all(
        np, reference_session, tokenizer, texts, args.max_length, args.batch_size
    )
    int8_vectors = embed_all(np, int8_session, tokenizer, texts, args.max_length, args.batch_size)

    reference_dimensions = len(reference_vectors[0])
    int8_dimensions = len(int8_vectors[0])
    if reference_dimensions != int8_dimensions:
        raise InputError(
            f"artifact dimensions differ: reference {reference_dimensions}, int8 {int8_dimensions}"
        )

    similarities = [cosine(left, right) for left, right in zip(reference_vectors, int8_vectors)]
    if not all(math.isfinite(value) for value in similarities):
        raise RuntimeError("parity produced a non-finite cosine; check for zero vectors")
    mean_cosine = sum(similarities) / len(similarities)
    min_cosine = min(similarities)

    reference_by_text = dict(zip(texts, reference_vectors))
    int8_by_text = dict(zip(texts, int8_vectors))
    reference_retrieval = retrieval_metrics(reference_by_text, pairs, documents, args.top_k)
    int8_retrieval = retrieval_metrics(int8_by_text, pairs, documents, args.top_k)
    delta = {
        f"{key}_pp": round((int8_retrieval[key] - reference_retrieval[key]) * 100.0, 4)
        for key in ("recall_at_1", "recall_at_5", "recall_at_10")
    }
    delta["mrr"] = round(int8_retrieval["mrr"] - reference_retrieval["mrr"], 6)
    delta["ndcg_at_10"] = round(
        int8_retrieval["ndcg_at_10"] - reference_retrieval["ndcg_at_10"], 6
    )

    checks = {
        "mean_cosine": mean_cosine >= args.min_mean_cosine,
        "min_cosine": min_cosine >= args.min_cosine,
        "recall_at_10_drop": delta["recall_at_10_pp"] >= -args.max_recall_drop_pp,
    }
    return {
        "schema_version": 1,
        "reference_artifact": artifact_info(args.reference_onnx),
        "int8_artifact": artifact_info(args.int8_onnx),
        "model_path": str(args.model_path),
        "pairs": {"count": len(pairs), "documents": len(documents), "sha256": pairs_sha256(pairs)},
        "texts": len(texts),
        "dimensions": reference_dimensions,
        "max_length": args.max_length,
        "batch_size": args.batch_size,
        "embeddings_sha256": {
            "reference": vector_set_sha256(reference_vectors),
            "int8": vector_set_sha256(int8_vectors),
        },
        "cosine": {
            "mean": round(mean_cosine, 8),
            "min": round(min_cosine, 8),
            "p1": round(percentile(similarities, 0.01), 8),
            "count": len(similarities),
        },
        "retrieval": {"reference": reference_retrieval, "int8": int8_retrieval},
        "delta_int8_minus_reference": delta,
        "acceptance": {
            "min_mean_cosine": args.min_mean_cosine,
            "min_cosine": args.min_cosine,
            "max_recall_drop_pp": args.max_recall_drop_pp,
            "checks": checks,
            "passed": all(checks.values()),
        },
        "raw_texts_in_report": False,
    }


def main(argv: Sequence[str] | None = None) -> int:
    args = parse_args(argv)
    try:
        report = run_check(args)
    except (InputError, RuntimeError, OSError) as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 2
    try:
        args.output.parent.mkdir(parents=True, exist_ok=True)
        args.output.write_text(
            json.dumps(report, ensure_ascii=False, indent=2, sort_keys=True) + "\n", encoding="utf-8"
        )
    except OSError as exc:
        print(f"error: cannot write report ({exc.strerror or type(exc).__name__})", file=sys.stderr)
        return 2
    print(f"report written to {args.output}", file=sys.stderr)
    return 0 if report["acceptance"]["passed"] else 2


if __name__ == "__main__":
    raise SystemExit(main())
