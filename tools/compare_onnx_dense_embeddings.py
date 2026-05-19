#!/usr/bin/env python3
"""Compare local dense ONNX output against the fd TEI/API dense baseline."""

from __future__ import annotations

import argparse
import json
import math
from pathlib import Path
import re
import sys
from typing import Any

import numpy as np
import onnxruntime as ort
from transformers import AutoTokenizer

from compare_dense_embeddings import (
    DEFAULT_API_URL,
    DEFAULT_DIMENSIONS,
    DEFAULT_MODEL,
    DEFAULT_NORM_TOLERANCE,
    PROBES,
    cosine,
    extract_vectors,
    finite_values,
    l2_norm,
    request_embeddings,
    vector_sha256,
)

DEFAULT_COSINE_THRESHOLD = 0.999


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Compare dense ONNX embeddings to current TEI/API baseline.")
    parser.add_argument("--onnx-path", type=Path, required=True, help="Local dense ONNX model path.")
    parser.add_argument("--model-path", type=Path, required=True, help="Local tokenizer/model snapshot path.")
    parser.add_argument("--baseline-artifact", type=Path, required=True, help="S02 baseline artifact path.")
    parser.add_argument("--output", type=Path, required=True, help="Markdown result artifact path.")
    parser.add_argument("--api-url", default=DEFAULT_API_URL, help="fd API URL for live TEI comparison.")
    parser.add_argument("--model", default=DEFAULT_MODEL, help="fd model name for live TEI comparison.")
    parser.add_argument("--dimensions", type=int, default=DEFAULT_DIMENSIONS, help="Expected dense dimensions.")
    parser.add_argument("--norm-tolerance", type=float, default=DEFAULT_NORM_TOLERANCE)
    parser.add_argument("--cosine-threshold", type=float, default=DEFAULT_COSINE_THRESHOLD)
    return parser.parse_args()


def parse_baseline_hashes(path: Path) -> dict[str, str]:
    text = path.read_text(encoding="utf-8")
    pattern = re.compile(r"^\| ([a-z0-9_]+) \| \d+ \| \d+ \| yes \| [0-9.]+ \| yes \| `([0-9a-f]{64})` \| yes \|$", re.MULTILINE)
    return {label: digest for label, digest in pattern.findall(text)}


def run_onnx(onnx_path: Path, model_path: Path) -> list[list[float]]:
    tokenizer = AutoTokenizer.from_pretrained(str(model_path), local_files_only=True)
    encoded = tokenizer(
        [probe["text"] for probe in PROBES],
        padding=True,
        truncation=True,
        return_tensors="np",
    )
    session = ort.InferenceSession(str(onnx_path), providers=["CPUExecutionProvider"])
    inputs = {
        "input_ids": encoded["input_ids"].astype(np.int64),
        "attention_mask": encoded["attention_mask"].astype(np.int64),
    }
    outputs = session.run(None, inputs)
    if len(outputs) != 1:
        raise RuntimeError(f"expected one ONNX output, got {len(outputs)}")
    dense = outputs[0]
    if dense.ndim != 2:
        raise RuntimeError(f"expected 2D dense output, got shape {dense.shape}")
    return dense.astype(np.float32).tolist()


def check_vector(vector: list[float], dimensions: int, norm_tolerance: float) -> dict[str, Any]:
    norm = l2_norm(vector)
    return {
        "dimensions": len(vector),
        "dimensions_ok": len(vector) == dimensions,
        "finite_values": finite_values(vector),
        "l2_norm": round(norm, 8),
        "normalized_within_tolerance": abs(norm - 1.0) <= norm_tolerance,
        "vector_sha256_float32": vector_sha256(vector),
    }


def build_result(args: argparse.Namespace) -> dict[str, Any]:
    baseline_hashes = parse_baseline_hashes(args.baseline_artifact)
    live_body = request_embeddings(args.api_url, args.model, args.dimensions)
    tei_vectors = extract_vectors(live_body, len(PROBES))
    onnx_vectors = run_onnx(args.onnx_path, args.model_path)

    rows = []
    all_passed = True
    for probe, tei_vec, onnx_vec in zip(PROBES, tei_vectors, onnx_vectors):
        label = probe["label"]
        tei_check = check_vector(tei_vec, args.dimensions, args.norm_tolerance)
        onnx_check = check_vector(onnx_vec, args.dimensions, args.norm_tolerance)
        tei_hash_matches_baseline = baseline_hashes.get(label) == tei_check["vector_sha256_float32"]
        similarity = cosine(tei_vec, onnx_vec)
        max_abs_diff = max(abs(a - b) for a, b in zip(tei_vec, onnx_vec))
        passed = (
            tei_check["dimensions_ok"]
            and tei_check["finite_values"]
            and tei_check["normalized_within_tolerance"]
            and onnx_check["dimensions_ok"]
            and onnx_check["finite_values"]
            and onnx_check["normalized_within_tolerance"]
            and tei_hash_matches_baseline
            and math.isfinite(similarity)
            and similarity >= args.cosine_threshold
        )
        all_passed = all_passed and passed
        rows.append(
            {
                "label": label,
                "chars": len(probe["text"]),
                "tei": tei_check,
                "onnx": onnx_check,
                "tei_hash_matches_s02_baseline": tei_hash_matches_baseline,
                "tei_onnx_cosine": round(similarity, 8),
                "max_abs_diff": round(max_abs_diff, 10),
                "passed": passed,
            }
        )

    session = ort.InferenceSession(str(args.onnx_path), providers=["CPUExecutionProvider"])
    return {
        "api_url": args.api_url,
        "model": args.model,
        "dimensions": args.dimensions,
        "norm_tolerance": args.norm_tolerance,
        "cosine_threshold": args.cosine_threshold,
        "onnx_path": str(args.onnx_path),
        "model_path": str(args.model_path),
        "baseline_artifact": str(args.baseline_artifact),
        "raw_probe_texts_logged": False,
        "onnxruntime": {
            "providers": session.get_providers(),
            "inputs": [{"name": item.name, "shape": item.shape, "type": item.type} for item in session.get_inputs()],
            "outputs": [{"name": item.name, "shape": item.shape, "type": item.type} for item in session.get_outputs()],
        },
        "rows": rows,
        "passed": all_passed,
    }


def render_markdown(result: dict[str, Any]) -> str:
    lines = [
        "# ONNX FP32 Dense Comparison — M010 S03",
        "",
        "## Effective Configuration",
        "",
        "```json",
        json.dumps({k: result[k] for k in ["api_url", "model", "dimensions", "norm_tolerance", "cosine_threshold", "onnx_path", "model_path", "baseline_artifact", "raw_probe_texts_logged"]}, ensure_ascii=False, indent=2, sort_keys=True),
        "```",
        "",
        "## ONNX Runtime Metadata",
        "",
        "```json",
        json.dumps(result["onnxruntime"], ensure_ascii=False, indent=2, sort_keys=True),
        "```",
        "",
        "## TEI vs ONNX Probe Comparison",
        "",
        "| Label | Chars | TEI Hash Matches S02 | TEI Norm | ONNX Norm | Cosine | Max Abs Diff | ONNX SHA256 (float32) | Passed |",
        "|---|---:|---|---:|---:|---:|---:|---|---|",
    ]
    for row in result["rows"]:
        lines.append(
            "| {label} | {chars} | {hash_match} | {tei_norm:.8f} | {onnx_norm:.8f} | {cosine:.8f} | {diff:.10f} | `{onnx_hash}` | {passed} |".format(
                label=row["label"],
                chars=row["chars"],
                hash_match="yes" if row["tei_hash_matches_s02_baseline"] else "no",
                tei_norm=row["tei"]["l2_norm"],
                onnx_norm=row["onnx"]["l2_norm"],
                cosine=row["tei_onnx_cosine"],
                diff=row["max_abs_diff"],
                onnx_hash=row["onnx"]["vector_sha256_float32"],
                passed="yes" if row["passed"] else "no",
            )
        )
    lines.extend(
        [
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
        result = build_result(args)
        markdown = render_markdown(result)
        args.output.parent.mkdir(parents=True, exist_ok=True)
        args.output.write_text(markdown, encoding="utf-8")
        print(markdown)
        return 0 if result["passed"] else 2
    except Exception as exc:  # noqa: BLE001 - spike comparator should write a blocked artifact.
        args.output.parent.mkdir(parents=True, exist_ok=True)
        blocked = "\n".join(
            [
                "# ONNX FP32 Dense Comparison — M010 S03",
                "",
                "## Verdict",
                "",
                "BLOCKED",
                "",
                f"Error type: `{type(exc).__name__}`",
                "",
                f"Error message: `{exc}`",
                "",
                "Raw probe texts are intentionally excluded from this artifact.",
                "",
            ]
        )
        args.output.write_text(blocked, encoding="utf-8")
        print(blocked, file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
