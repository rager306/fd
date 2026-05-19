#!/usr/bin/env python3
"""
Dense embedding comparator baseline for fd.

Calls the current TEI-backed API with fixed Russian/legal-style probes and writes a
sanitized artifact: probe labels, character counts, dimensions, norm checks,
vector hashes, and pairwise cosine similarities. Raw probe texts are never
printed to stdout or the artifact.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import math
import os
from datetime import datetime, timezone
from pathlib import Path
import struct
import sys
from typing import Any

import requests

DEFAULT_API_URL = "http://localhost:8000"
DEFAULT_MODEL = "deepvk/USER-bge-m3"
DEFAULT_DIMENSIONS = 1024
DEFAULT_NORM_TOLERANCE = 0.02
SCRIPT_VERSION = 1

# Non-sensitive fixed probes. Do not print raw text in artifacts; use labels and
# character counts only. These stay in source so future ONNX comparisons are
# reproducible without a separate data dependency.
PROBES: list[dict[str, str]] = [
    {
        "label": "contract_question_jurisdiction",
        "text": "Как определяется подсудность спора по договору поставки между российскими юридическими лицами?",
    },
    {
        "label": "contract_clause_delivery",
        "text": "Поставщик обязан передать товар в срок, согласованный сторонами, а покупатель обязан принять товар и оплатить его по цене договора.",
    },
    {
        "label": "labor_question_termination",
        "text": "Какие гарантии предоставляются работнику при расторжении трудового договора по инициативе работодателя?",
    },
    {
        "label": "civil_clause_damages",
        "text": "Лицо, право которого нарушено, может требовать полного возмещения причинённых ему убытков, если законом или договором не предусмотрено иное.",
    },
    {
        "label": "neutral_russian_reference",
        "text": "Москва — крупный научный, культурный и транспортный центр России с развитой городской инфраструктурой.",
    },
]


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Compare dense embeddings from the fd API baseline.")
    parser.add_argument(
        "--api-url",
        default=os.getenv("COMPARE_API_URL", DEFAULT_API_URL),
        help="fd API base URL (default: COMPARE_API_URL or http://localhost:8000)",
    )
    parser.add_argument(
        "--model",
        default=os.getenv("COMPARE_MODEL", DEFAULT_MODEL),
        help="Embedding model name to request and record.",
    )
    parser.add_argument(
        "--dimensions",
        type=int,
        default=int(os.getenv("COMPARE_DIMENSIONS", str(DEFAULT_DIMENSIONS))),
        help="Expected embedding dimensions.",
    )
    parser.add_argument(
        "--norm-tolerance",
        type=float,
        default=float(os.getenv("COMPARE_NORM_TOLERANCE", str(DEFAULT_NORM_TOLERANCE))),
        help="Allowed absolute deviation from L2 norm 1.0.",
    )
    parser.add_argument(
        "--output",
        type=Path,
        default=None,
        help="Optional artifact path to write markdown output.",
    )
    return parser.parse_args()


def request_embeddings(api_url: str, model: str, dimensions: int) -> dict[str, Any]:
    url = f"{api_url.rstrip('/')}/v1/embeddings"
    payload = {
        "model": model,
        "input": [probe["text"] for probe in PROBES],
        "dimensions": dimensions,
    }
    try:
        response = requests.post(url, json=payload, timeout=120)
    except requests.RequestException as exc:
        raise RuntimeError(f"embedding request failed: {type(exc).__name__}") from exc

    if response.status_code != 200:
        raise RuntimeError(f"embedding request returned HTTP {response.status_code}")

    try:
        body = response.json()
    except ValueError as exc:
        raise RuntimeError("embedding response was not valid JSON") from exc

    if not isinstance(body, dict):
        raise RuntimeError("embedding response must be a JSON object")
    return body


def extract_vectors(body: dict[str, Any], expected_count: int) -> list[list[float]]:
    data = body.get("data")
    if not isinstance(data, list):
        raise RuntimeError("embedding response missing data array")
    if len(data) != expected_count:
        raise RuntimeError(f"embedding response count mismatch: got {len(data)}, expected {expected_count}")

    vectors: list[list[float]] = []
    for idx, item in enumerate(data):
        if not isinstance(item, dict):
            raise RuntimeError(f"embedding item {idx} is not an object")
        embedding = item.get("embedding")
        if not isinstance(embedding, list):
            raise RuntimeError(f"embedding item {idx} missing embedding array")
        try:
            vector = [float(value) for value in embedding]
        except (TypeError, ValueError) as exc:
            raise RuntimeError(f"embedding item {idx} contains non-numeric values") from exc
        vectors.append(vector)
    return vectors


def l2_norm(vector: list[float]) -> float:
    return math.sqrt(sum(value * value for value in vector))


def cosine(a: list[float], b: list[float]) -> float:
    denom = l2_norm(a) * l2_norm(b)
    if denom == 0:
        return float("nan")
    return sum(x * y for x, y in zip(a, b)) / denom


def finite_values(vector: list[float]) -> bool:
    return all(math.isfinite(value) for value in vector)


def vector_sha256(vector: list[float]) -> str:
    digest = hashlib.sha256()
    for value in vector:
        # Hash as float32 little-endian to match the service/cache precision rather than Python repr.
        digest.update(struct.pack("<f", float(value)))
    return digest.hexdigest()


def build_result(args: argparse.Namespace, body: dict[str, Any], vectors: list[list[float]]) -> dict[str, Any]:
    vector_checks = []
    all_passed = True

    for probe, vector in zip(PROBES, vectors):
        norm = l2_norm(vector)
        dimensions_ok = len(vector) == args.dimensions
        finite_ok = finite_values(vector)
        normalized_ok = abs(norm - 1.0) <= args.norm_tolerance
        passed = dimensions_ok and finite_ok and normalized_ok
        all_passed = all_passed and passed
        vector_checks.append(
            {
                "label": probe["label"],
                "chars": len(probe["text"]),
                "dimensions": len(vector),
                "dimensions_ok": dimensions_ok,
                "finite_values": finite_ok,
                "l2_norm": round(norm, 8),
                "normalized_within_tolerance": normalized_ok,
                "vector_sha256_float32": vector_sha256(vector),
                "passed": passed,
            }
        )

    similarities = []
    for i, left in enumerate(PROBES):
        for j in range(i + 1, len(PROBES)):
            right = PROBES[j]
            similarities.append(
                {
                    "left": left["label"],
                    "right": right["label"],
                    "cosine": round(cosine(vectors[i], vectors[j]), 8),
                }
            )

    return {
        "script_version": SCRIPT_VERSION,
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "api_url": args.api_url,
        "requested_model": args.model,
        "response_model": body.get("model"),
        "expected_dimensions": args.dimensions,
        "norm_tolerance": args.norm_tolerance,
        "raw_probe_texts_logged": False,
        "probe_count": len(PROBES),
        "vector_checks": vector_checks,
        "pairwise_cosine_similarities": similarities,
        "usage": body.get("usage"),
        "passed": all_passed,
    }


def render_markdown(result: dict[str, Any]) -> str:
    lines = [
        "# Dense Embedding Comparator Baseline",
        "",
        "## Effective Configuration",
        "",
        "```json",
        json.dumps(
            {
                "script_version": result["script_version"],
                "generated_at": result["generated_at"],
                "api_url": result["api_url"],
                "requested_model": result["requested_model"],
                "response_model": result["response_model"],
                "expected_dimensions": result["expected_dimensions"],
                "norm_tolerance": result["norm_tolerance"],
                "raw_probe_texts_logged": result["raw_probe_texts_logged"],
                "probe_count": result["probe_count"],
            },
            ensure_ascii=False,
            indent=2,
            sort_keys=True,
        ),
        "```",
        "",
        "## Probe Summary",
        "",
        "| Label | Chars | Dimensions | Finite | L2 Norm | Normalized | Vector SHA256 (float32) | Passed |",
        "|---|---:|---:|---|---:|---|---|---|",
    ]

    for check in result["vector_checks"]:
        lines.append(
            "| {label} | {chars} | {dimensions} | {finite} | {norm:.8f} | {normalized} | `{sha}` | {passed} |".format(
                label=check["label"],
                chars=check["chars"],
                dimensions=check["dimensions"],
                finite="yes" if check["finite_values"] else "no",
                norm=check["l2_norm"],
                normalized="yes" if check["normalized_within_tolerance"] else "no",
                sha=check["vector_sha256_float32"],
                passed="yes" if check["passed"] else "no",
            )
        )

    lines.extend(
        [
            "",
            "## Pairwise Cosine Similarities",
            "",
            "| Left | Right | Cosine |",
            "|---|---|---:|",
        ]
    )
    for similarity in result["pairwise_cosine_similarities"]:
        lines.append(
            f"| {similarity['left']} | {similarity['right']} | {similarity['cosine']:.8f} |"
        )

    lines.extend(
        [
            "",
            "## Usage",
            "",
            "```json",
            json.dumps(result["usage"], ensure_ascii=False, indent=2, sort_keys=True),
            "```",
            "",
            "## Verdict",
            "",
            "PASS" if result["passed"] else "FAIL",
            "",
            "Raw probe texts are intentionally excluded from this artifact.",
            "",
        ]
    )
    return "\n".join(lines)


def main() -> int:
    args = parse_args()
    try:
        body = request_embeddings(args.api_url, args.model, args.dimensions)
        vectors = extract_vectors(body, len(PROBES))
        result = build_result(args, body, vectors)
        markdown = render_markdown(result)
    except RuntimeError as exc:
        print(f"ERROR: {exc}", file=sys.stderr)
        return 1

    if args.output:
        args.output.parent.mkdir(parents=True, exist_ok=True)
        args.output.write_text(markdown, encoding="utf-8")
    print(markdown)
    return 0 if result["passed"] else 2


if __name__ == "__main__":
    raise SystemExit(main())
