#!/usr/bin/env python3
"""Profile M015 legal divergence worst cases without raw text leakage."""

from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import json
from pathlib import Path
import re
from typing import Any

from transformers import AutoTokenizer

SCRIPT_VERSION = 1
DEFAULT_SEQUENCE_LENGTHS = [128, 256, 512, 1024]


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Profile legal divergence cases from M015 artifacts.")
    parser.add_argument("--corpus", type=Path, required=True)
    parser.add_argument("--m015-artifact", type=Path, required=True)
    parser.add_argument("--tokenizer-path", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    parser.add_argument("--sequence-lengths", default=",".join(str(v) for v in DEFAULT_SEQUENCE_LENGTHS))
    parser.add_argument("--limit", type=int, default=20)
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


def parse_metrics(path: Path) -> dict[str, Any]:
    text = path.read_text(encoding="utf-8")
    match = re.search(r"## Metrics\n\n```json\n(.*?)\n```", text, re.S)
    if not match:
        raise RuntimeError(f"metrics JSON block not found in {path}")
    return json.loads(match.group(1))


def extract_worst_case_rows(metrics: dict[str, Any], limit: int) -> list[dict[str, Any]]:
    cosine = metrics.get("cross_backend_cosine", {})
    rows: list[dict[str, Any]] = []
    for source, kind in ((cosine.get("worst_documents", []), "document"), (cosine.get("worst_queries", []), "query")):
        for row in source:
            document_id = row.get("id", "")
            if document_id.startswith("self:"):
                document_id = document_id[5:]
            rows.append(
                {
                    "source_kind": kind,
                    "id": document_id,
                    "original_id": row.get("id"),
                    "article": row.get("article"),
                    "kind": row.get("kind"),
                    "chars": row.get("chars"),
                    "cosine": row.get("cosine"),
                    "text_sha256": row.get("text_sha256"),
                }
            )
    dedup: dict[str, dict[str, Any]] = {}
    for row in sorted(rows, key=lambda item: (float(item.get("cosine") or 1.0), item["id"])):
        existing = dedup.get(row["id"])
        if existing is None or float(row.get("cosine") or 1.0) < float(existing.get("cosine") or 1.0):
            dedup[row["id"]] = row
    return list(dedup.values())[:limit]


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
                records[part_id] = {
                    "id": part_id,
                    "article": article,
                    "kind": "part",
                    "text": part_text,
                    "chars": len(part_text),
                    "text_sha256": sha256_text(part_text),
                    "invalid": bool(part.get("invalid")),
                }
            for clause_index, clause in enumerate(part.get("clauses") or [], 1):
                clause_no = stable_number(clause.get("number"), f"idx{clause_index}")
                clause_text = clause.get("text") or ""
                clause_id = f"{prefix}:a{article}:p{part_no}:c{clause_no}"
                if clause_text.strip():
                    records[clause_id] = {
                        "id": clause_id,
                        "article": article,
                        "kind": "clause",
                        "text": clause_text,
                        "chars": len(clause_text),
                        "text_sha256": sha256_text(clause_text),
                        "invalid": bool(part.get("invalid")) or bool(clause.get("invalid")),
                    }
                for subclause_index, subclause in enumerate(clause.get("subclauses") or [], 1):
                    sub_no = stable_number(subclause.get("number"), f"idx{subclause_index}")
                    sub_text = subclause.get("text") or ""
                    sub_id = f"{prefix}:a{article}:p{part_no}:c{clause_no}:s{sub_no}"
                    if sub_text.strip():
                        records[sub_id] = {
                            "id": sub_id,
                            "article": article,
                            "kind": "subclause",
                            "text": sub_text,
                            "chars": len(sub_text),
                            "text_sha256": sha256_text(sub_text),
                            "invalid": bool(part.get("invalid")) or bool(clause.get("invalid")) or bool(subclause.get("invalid")),
                        }
    return records


def token_diagnostics(tokenizer: Any, text: str, sequence_lengths: list[int]) -> dict[str, Any]:
    encoded_without_specials = tokenizer(text, add_special_tokens=False, truncation=False)
    encoded_with_specials = tokenizer(text, add_special_tokens=True, truncation=False)
    input_ids = encoded_with_specials["input_ids"]
    if input_ids and isinstance(input_ids[0], list):
        input_ids = input_ids[0]
    base = {
        "tokens_without_specials": len(encoded_without_specials["input_ids"]),
        "tokens_with_specials": len(input_ids),
        "sequence_lengths": {},
    }
    for length in sequence_lengths:
        encoded = tokenizer(text, add_special_tokens=True, truncation=True, max_length=length)
        ids = encoded["input_ids"]
        if ids and isinstance(ids[0], list):
            ids = ids[0]
        base["sequence_lengths"][str(length)] = {
            "encoded_tokens": len(ids),
            "truncated": len(input_ids) > length,
            "tokens_dropped": max(0, len(input_ids) - length),
            "retained_ratio": round(min(len(input_ids), length) / len(input_ids), 6) if input_ids else 0.0,
        }
    return base


def build_result(args: argparse.Namespace) -> dict[str, Any]:
    sequence_lengths = [int(item.strip()) for item in args.sequence_lengths.split(",") if item.strip()]
    metrics = parse_metrics(args.m015_artifact)
    worst_rows = extract_worst_case_rows(metrics, args.limit)
    corpus = load_corpus_texts(args.corpus)
    tokenizer = AutoTokenizer.from_pretrained(str(args.tokenizer_path), local_files_only=True)

    cases = []
    for row in worst_rows:
        record = corpus.get(row["id"])
        if record is None:
            cases.append({**row, "resolved": False})
            continue
        hash_matches = record["text_sha256"] == row.get("text_sha256")
        diagnostics = token_diagnostics(tokenizer, record["text"], sequence_lengths)
        cases.append(
            {
                "id": row["id"],
                "source_kind": row["source_kind"],
                "article": record["article"],
                "kind": record["kind"],
                "chars": record["chars"],
                "text_sha256": record["text_sha256"],
                "artifact_hash_matches": hash_matches,
                "invalid": record["invalid"],
                "m015_cosine": row.get("cosine"),
                "tokenizer": diagnostics,
                "resolved": True,
            }
        )

    resolved_cases = [case for case in cases if case.get("resolved")]
    return {
        "script_version": SCRIPT_VERSION,
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "raw_text_logged": False,
        "inputs": {
            "corpus_path": str(args.corpus),
            "corpus_sha256": sha256_file(args.corpus),
            "m015_artifact": str(args.m015_artifact),
            "m015_artifact_sha256": sha256_file(args.m015_artifact),
            "tokenizer_path": str(args.tokenizer_path),
            "tokenizer_json_sha256": sha256_file(args.tokenizer_path / "tokenizer.json") if (args.tokenizer_path / "tokenizer.json").exists() else None,
            "sequence_lengths": sequence_lengths,
        },
        "summary": {
            "requested_limit": args.limit,
            "cases": len(cases),
            "resolved_cases": len(resolved_cases),
            "all_hashes_match": all(case.get("artifact_hash_matches") for case in resolved_cases),
            "min_m015_cosine": min((float(case["m015_cosine"]) for case in resolved_cases), default=None),
            "max_tokens_with_specials": max((case["tokenizer"]["tokens_with_specials"] for case in resolved_cases), default=0),
            "cases_truncated_at_128": sum(1 for case in resolved_cases if case["tokenizer"]["sequence_lengths"].get("128", {}).get("truncated")),
            "cases_truncated_at_512": sum(1 for case in resolved_cases if case["tokenizer"]["sequence_lengths"].get("512", {}).get("truncated")),
        },
        "cases": cases,
        "caveat": "Token diagnostics use the local Hugging Face tokenizer. They identify truncation risk but do not prove TEI internal truncation behavior by themselves.",
    }


def render_markdown(result: dict[str, Any]) -> str:
    return "\n".join(
        [
            "# M016 S01 Legal Divergence Profile",
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
    markdown = render_markdown(result)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(markdown, encoding="utf-8")
    print(markdown)
    if result["summary"]["resolved_cases"] != result["summary"]["cases"]:
        return 2
    if not result["summary"]["all_hashes_match"]:
        return 3
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
