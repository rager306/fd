#!/usr/bin/env python3
"""Stdlib-only helpers shared by the model A/B rig and the ONNX parity check.

Keeping these helpers free of HTTP and ONNX imports lets the parity check run in
the ONNX export environment without pulling in the API client.
"""

from __future__ import annotations

import hashlib
import json
import math
import struct
from pathlib import Path
from typing import Any, Iterable, Sequence

METRIC_KS = (1, 5, 10)


class InputError(ValueError):
    """Invalid pair file or CLI input. The message never contains input texts."""


def _require_text(item: dict[str, Any], key: str, where: str) -> str:
    value = item.get(key)
    if not isinstance(value, str) or not value.strip():
        raise InputError(f"{where}: {key!r} must be a non-empty string")
    return value


def load_pairs(path: Path) -> list[dict[str, Any]]:
    """Load and validate pair rows. Error messages name the row, never its text."""
    pairs: list[dict[str, Any]] = []
    seen_ids: set[str] = set()
    try:
        handle = path.open(encoding="utf-8")
    except OSError as exc:
        raise InputError(f"cannot read pair file: {exc.strerror or type(exc).__name__}") from None
    with handle:
        for lineno, line in enumerate(handle, 1):
            stripped = line.strip()
            if not stripped:
                continue
            where = f"{path.name}:{lineno}"
            try:
                item = json.loads(stripped)
            except json.JSONDecodeError as exc:
                raise InputError(f"{where}: invalid JSON ({exc.msg})") from None
            if not isinstance(item, dict):
                raise InputError(f"{where}: row must be a JSON object")
            pair_id = _require_text(item, "id", where)
            if pair_id in seen_ids:
                raise InputError(f"{where}: duplicate id {pair_id!r}")
            seen_ids.add(pair_id)
            query = _require_text(item, "query", where)
            positive = _require_text(item, "positive", where)
            negatives = item.get("negatives")
            if not isinstance(negatives, list) or not negatives:
                raise InputError(f"{where}: 'negatives' must be a non-empty array")
            for index, negative in enumerate(negatives):
                if not isinstance(negative, str) or not negative.strip():
                    raise InputError(f"{where}: negatives[{index}] must be a non-empty string")
            if positive in negatives:
                raise InputError(f"{where}: positive text is also listed as a negative")
            pairs.append(
                {"id": pair_id, "query": query, "positive": positive, "negatives": list(negatives)}
            )
    if not pairs:
        raise InputError("pair file contains no pairs")
    return pairs


def pairs_sha256(pairs: Sequence[dict[str, Any]]) -> str:
    """Fingerprint the pair file without publishing its texts."""
    canonical = json.dumps(pairs, ensure_ascii=False, sort_keys=True, separators=(",", ":"))
    return hashlib.sha256(canonical.encode("utf-8")).hexdigest()


def unique_in_order(values: Iterable[str]) -> list[str]:
    seen: set[str] = set()
    ordered: list[str] = []
    for value in values:
        if value not in seen:
            seen.add(value)
            ordered.append(value)
    return ordered


def cosine(left: Sequence[float], right: Sequence[float]) -> float:
    if len(left) != len(right):
        return float("nan")
    dot = sum(a * b for a, b in zip(left, right))
    norm_left = math.sqrt(sum(a * a for a in left))
    norm_right = math.sqrt(sum(b * b for b in right))
    if norm_left == 0.0 or norm_right == 0.0:
        return float("nan")
    return dot / (norm_left * norm_right)


def percentile(values: Sequence[float], quantile: float) -> float:
    """Nearest-rank percentile; deterministic for small samples."""
    ordered = sorted(values)
    index = math.ceil(quantile * len(ordered)) - 1
    index = max(0, min(len(ordered) - 1, index))
    return ordered[index]


def vector_set_sha256(vectors: Sequence[Sequence[float]]) -> str:
    digest = hashlib.sha256()
    for vector in vectors:
        for value in vector:
            digest.update(struct.pack("<f", value))
    return digest.hexdigest()


def rank_of(
    query: Sequence[float],
    documents: Sequence[Sequence[float]],
    relevant: int,
    top_k: int,
) -> int:
    """Rank of the relevant document; ranks beyond top_k collapse to top_k + 1."""
    scores = [cosine(query, document) for document in documents]
    if not all(math.isfinite(score) for score in scores):
        raise RuntimeError("ranking produced a non-finite cosine; check for zero vectors")
    order = sorted(range(len(scores)), key=lambda index: (-scores[index], index))
    rank = order.index(relevant) + 1
    return min(rank, top_k + 1)


def summarize_ranks(ranks: Sequence[int]) -> dict[str, float]:
    count = len(ranks)
    if count == 0:
        raise InputError("no ranks to summarize")
    summary = {f"recall_at_{k}": sum(1 for rank in ranks if rank <= k) / count for k in METRIC_KS}
    summary["mrr"] = sum(1.0 / rank for rank in ranks) / count
    summary["ndcg_at_10"] = sum(
        1.0 / math.log2(rank + 1) if rank <= 10 else 0.0 for rank in ranks
    ) / count
    return summary
