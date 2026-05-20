#!/usr/bin/env python3
"""Diagnose ONNX sequence length effects for M016 legal divergence cases."""

from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import json
import math
from pathlib import Path
import re
import struct
from typing import Any

import numpy as np
import onnxruntime as ort
import requests
from transformers import AutoTokenizer

SCRIPT_VERSION = 1
DEFAULT_MODEL = "deepvk/USER-bge-m3"
DEFAULT_DIMENSIONS = 1024


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Compare TEI vectors with local ONNX outputs at multiple sequence lengths.")
    parser.add_argument("--profile", type=Path, required=True, help="M016 S01 divergence profile artifact.")
    parser.add_argument("--corpus", type=Path, required=True)
    parser.add_argument("--onnx-path", type=Path, required=True)
    parser.add_argument("--tokenizer-path", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    parser.add_argument("--api-url", default="http://localhost:8000")
    parser.add_argument("--model", default=DEFAULT_MODEL)
    parser.add_argument("--dimensions", type=int, default=DEFAULT_DIMENSIONS)
    parser.add_argument("--sequence-lengths", default="128,512")
    parser.add_argument("--batch-size", type=int, default=8)
    return parser.parse_args()


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def sha256_text(value: str) -> str:
    return hashlib.sha256(value.encode("utf-8")).hexdigest()


def vector_sha256(vector: list[float]) -> str:
    digest = hashlib.sha256()
    for value in vector:
        digest.update(struct.pack("<f", float(value)))
    return digest.hexdigest()


def stable_number(value: Any, fallback: str) -> str:
    if value is None or value == "":
        return fallback
    return str(value).replace(" ", "_")


def l2_norm(vector: list[float]) -> float:
    return math.sqrt(sum(value * value for value in vector))


def cosine(a: list[float], b: list[float]) -> float:
    denom = l2_norm(a) * l2_norm(b)
    if denom == 0:
        return float("nan")
    return sum(x * y for x, y in zip(a, b)) / denom


def summarize(values: list[float]) -> dict[str, float]:
    if not values:
        return {"count": 0, "mean": 0.0, "min": 0.0, "p50": 0.0, "p95": 0.0}
    ordered = sorted(values)
    return {
        "count": len(values),
        "mean": round(sum(values) / len(values), 8),
        "min": round(min(values), 8),
        "p50": round(ordered[len(ordered) // 2], 8),
        "p95": round(ordered[min(len(ordered) - 1, int(len(ordered) * 0.95))], 8),
    }


def parse_profile(path: Path) -> dict[str, Any]:
    text = path.read_text(encoding="utf-8")
    match = re.search(r"## Result\n\n```json\n(.*?)\n```", text, re.S)
    if not match:
        raise RuntimeError(f"profile JSON block not found in {path}")
    return json.loads(match.group(1))


def load_corpus_texts(path: Path) -> dict[str, dict[str, Any]]:
    records: dict[str, dict[str, Any]] = {}
    for line_no, line in enumerate(path.read_text(encoding="utf-8").splitlines(), 1):
        if not line.strip():
            continue
        obj = json.loads(line)
        article = stable_number(obj.get("article"), f"line{line_no}")
        prefix = stable_number(obj.get("doc_id", "doc"), "doc")
        for part_index, part in enumerate(obj.get("parts") or [], 1):
            part_no = stable_number(part.get("number"), f"idx{part_index}")
            part_text = part.get("text") or ""
            part_id = f"{prefix}:a{article}:p{part_no}"
            if part_text.strip():
                records[part_id] = {"text": part_text, "text_sha256": sha256_text(part_text), "article": article, "kind": "part", "chars": len(part_text)}
            for clause_index, clause in enumerate(part.get("clauses") or [], 1):
                clause_no = stable_number(clause.get("number"), f"idx{clause_index}")
                clause_text = clause.get("text") or ""
                clause_id = f"{prefix}:a{article}:p{part_no}:c{clause_no}"
                if clause_text.strip():
                    records[clause_id] = {"text": clause_text, "text_sha256": sha256_text(clause_text), "article": article, "kind": "clause", "chars": len(clause_text)}
                for subclause_index, subclause in enumerate(clause.get("subclauses") or [], 1):
                    sub_no = stable_number(subclause.get("number"), f"idx{subclause_index}")
                    sub_text = subclause.get("text") or ""
                    sub_id = f"{prefix}:a{article}:p{part_no}:c{clause_no}:s{sub_no}"
                    if sub_text.strip():
                        records[sub_id] = {"text": sub_text, "text_sha256": sha256_text(sub_text), "article": article, "kind": "subclause", "chars": len(sub_text)}
    return records


def request_tei_vectors(api_url: str, model: str, dimensions: int, texts: list[str], batch_size: int) -> list[list[float]]:
    vectors: list[list[float]] = []
    url = f"{api_url.rstrip('/')}/v1/embeddings"
    for start in range(0, len(texts), batch_size):
        batch = texts[start : start + batch_size]
        payload = {"model": model, "input": batch, "dimensions": dimensions}
        response = requests.post(url, json=payload, timeout=180)
        response.raise_for_status()
        body = response.json()
        data = body.get("data")
        if not isinstance(data, list) or len(data) != len(batch):
            raise RuntimeError(f"TEI response count mismatch: got {len(data) if isinstance(data, list) else 'invalid'}, expected {len(batch)}")
        for item in data:
            vector = [float(value) for value in item["embedding"]]
            if len(vector) != dimensions:
                raise RuntimeError(f"TEI dimensions mismatch: {len(vector)} != {dimensions}")
            vectors.append(vector)
    return vectors


def run_onnx_vectors(session: ort.InferenceSession, tokenizer: Any, texts: list[str], max_length: int, batch_size: int) -> list[list[float]]:
    vectors: list[list[float]] = []
    for start in range(0, len(texts), batch_size):
        batch = texts[start : start + batch_size]
        encoded = tokenizer(batch, padding=True, truncation=True, max_length=max_length, return_tensors="np")
        inputs = {
            "input_ids": encoded["input_ids"].astype(np.int64),
            "attention_mask": encoded["attention_mask"].astype(np.int64),
        }
        outputs = session.run(None, inputs)
        if len(outputs) != 1:
            raise RuntimeError(f"expected one ONNX output, got {len(outputs)}")
        dense = outputs[0]
        if dense.ndim != 2:
            raise RuntimeError(f"expected 2D ONNX output, got {dense.shape}")
        vectors.extend(dense.astype(np.float32).tolist())
    return vectors


def token_count(tokenizer: Any, text: str) -> int:
    ids = tokenizer(text, add_special_tokens=True, truncation=False)["input_ids"]
    if ids and isinstance(ids[0], list):
        ids = ids[0]
    return len(ids)


def build_result(args: argparse.Namespace) -> dict[str, Any]:
    profile = parse_profile(args.profile)
    corpus = load_corpus_texts(args.corpus)
    sequence_lengths = [int(item.strip()) for item in args.sequence_lengths.split(",") if item.strip()]
    cases = []
    for case in profile["cases"]:
        record = corpus.get(case["id"])
        if record is None:
            raise RuntimeError(f"case not found in corpus: {case['id']}")
        if record["text_sha256"] != case["text_sha256"]:
            raise RuntimeError(f"case hash mismatch: {case['id']}")
        cases.append({**case, "text": record["text"]})

    texts = [case["text"] for case in cases]
    tokenizer = AutoTokenizer.from_pretrained(str(args.tokenizer_path), local_files_only=True)
    session = ort.InferenceSession(str(args.onnx_path), providers=["CPUExecutionProvider"])
    tei_vectors = request_tei_vectors(args.api_url, args.model, args.dimensions, texts, args.batch_size)

    per_length: dict[str, Any] = {}
    for length in sequence_lengths:
        onnx_vectors = run_onnx_vectors(session, tokenizer, texts, length, args.batch_size)
        cosines = [cosine(tei, onnx) for tei, onnx in zip(tei_vectors, onnx_vectors)]
        rows = []
        for case, tei, onnx, similarity in zip(cases, tei_vectors, onnx_vectors, cosines):
            tokens = token_count(tokenizer, case["text"])
            rows.append(
                {
                    "id": case["id"],
                    "article": case["article"],
                    "kind": case["kind"],
                    "chars": case["chars"],
                    "tokens_with_specials": tokens,
                    "truncated_at_length": tokens > length,
                    "tokens_dropped_at_length": max(0, tokens - length),
                    "m015_cosine": case["m015_cosine"],
                    "tei_onnx_cosine": round(similarity, 8),
                    "tei_norm": round(l2_norm(tei), 8),
                    "onnx_norm": round(l2_norm(onnx), 8),
                    "onnx_sha256_float32": vector_sha256(onnx),
                    "text_sha256": case["text_sha256"],
                }
            )
        rows.sort(key=lambda item: item["tei_onnx_cosine"])
        per_length[str(length)] = {
            "summary": summarize(cosines),
            "truncated_cases": sum(1 for row in rows if row["truncated_at_length"]),
            "worst_cases": rows[:10],
        }

    return {
        "script_version": SCRIPT_VERSION,
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "raw_text_logged": False,
        "config": {
            "profile": str(args.profile),
            "profile_sha256": sha256_file(args.profile),
            "corpus": str(args.corpus),
            "corpus_sha256": sha256_file(args.corpus),
            "onnx_path": str(args.onnx_path),
            "onnx_sha256": sha256_file(args.onnx_path),
            "tokenizer_path": str(args.tokenizer_path),
            "tokenizer_json_sha256": sha256_file(args.tokenizer_path / "tokenizer.json") if (args.tokenizer_path / "tokenizer.json").exists() else None,
            "api_url": args.api_url,
            "model": args.model,
            "dimensions": args.dimensions,
            "sequence_lengths": sequence_lengths,
            "case_count": len(cases),
        },
        "onnxruntime": {
            "providers": session.get_providers(),
            "inputs": [{"name": item.name, "shape": item.shape, "type": item.type} for item in session.get_inputs()],
            "outputs": [{"name": item.name, "shape": item.shape, "type": item.type} for item in session.get_outputs()],
        },
        "results_by_sequence_length": per_length,
        "caveat": "This compares TEI API vectors with local Python ONNX outputs at selected tokenizer max lengths. It isolates sequence length effects from the Go service path but does not change production runtime.",
    }


def render_markdown(result: dict[str, Any]) -> str:
    return "\n".join(
        [
            "# M016 S02 ONNX Sequence Length Diagnostics",
            "",
            "## Result",
            "",
            "```json",
            json.dumps(result, ensure_ascii=False, indent=2, sort_keys=True),
            "```",
            "",
            "Raw legal corpus text is intentionally excluded from this artifact.",
            "",
        ]
    )


def main() -> int:
    args = parse_args()
    result = build_result(args)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(render_markdown(result), encoding="utf-8")
    print(render_markdown(result))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
