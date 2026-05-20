#!/usr/bin/env python3
"""
Tokenizer parity comparator for fd ONNX work.

Baseline mode records Hugging Face tokenizer output for the fixed Russian/legal
probes used by the dense embedding comparator. Artifacts intentionally exclude
raw probe text and use labels, character counts, token lengths, hashes, and token
ID/mask evidence for reproducible parity debugging.
"""

from __future__ import annotations

import argparse
import ast
import hashlib
import importlib.metadata
import json
import os
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

from transformers import AutoTokenizer

SCRIPT_VERSION = 1
DEFAULT_MODEL = "deepvk/USER-bge-m3"
DEFAULT_PROBES_SOURCE = Path("tools/compare_dense_embeddings.py")
DEFAULT_TOKENIZER_PATH = Path("tei-models/deepvk--USER-bge-m3")
DEFAULT_OUTPUT = Path("benchmark-results/fd-tokenizer-baseline-m012-s01.txt")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Compare tokenizer output for fd ONNX parity work.")
    parser.add_argument(
        "--mode",
        choices=["baseline"],
        default="baseline",
        help="Operation mode. Currently only writes the Hugging Face baseline.",
    )
    parser.add_argument(
        "--tokenizer-path",
        type=Path,
        default=Path(os.getenv("TOKENIZER_BASELINE_PATH", str(DEFAULT_TOKENIZER_PATH))),
        help="Local Hugging Face tokenizer directory.",
    )
    parser.add_argument(
        "--model",
        default=os.getenv("TOKENIZER_BASELINE_MODEL", DEFAULT_MODEL),
        help="Model identifier recorded in the artifact.",
    )
    parser.add_argument(
        "--probes-source",
        type=Path,
        default=DEFAULT_PROBES_SOURCE,
        help="Python source file containing the fixed PROBES literal.",
    )
    parser.add_argument(
        "--output",
        type=Path,
        default=DEFAULT_OUTPUT,
        help="Markdown artifact path to write.",
    )
    return parser.parse_args()


def sha256_file(path: Path) -> str | None:
    if not path.exists() or not path.is_file():
        return None
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def sha256_json(value: Any) -> str:
    payload = json.dumps(value, ensure_ascii=False, separators=(",", ":"), sort_keys=True).encode("utf-8")
    return hashlib.sha256(payload).hexdigest()


def package_version(name: str) -> str:
    try:
        return importlib.metadata.version(name)
    except importlib.metadata.PackageNotFoundError:
        return "not-installed"


def load_probes(source_path: Path) -> list[dict[str, str]]:
    try:
        module = ast.parse(source_path.read_text(encoding="utf-8"), filename=str(source_path))
    except OSError as exc:
        raise RuntimeError(f"failed to read probes source: {source_path}") from exc
    for node in module.body:
        if isinstance(node, ast.AnnAssign) and isinstance(node.target, ast.Name) and node.target.id == "PROBES":
            value = ast.literal_eval(node.value)
            if not isinstance(value, list):
                break
            probes: list[dict[str, str]] = []
            for idx, item in enumerate(value):
                if not isinstance(item, dict) or not isinstance(item.get("label"), str) or not isinstance(item.get("text"), str):
                    raise RuntimeError(f"invalid probe item at index {idx}")
                probes.append({"label": item["label"], "text": item["text"]})
            return probes
    raise RuntimeError(f"PROBES literal not found in {source_path}")


def load_tokenizer(tokenizer_path: Path) -> Any:
    if not tokenizer_path.exists():
        raise RuntimeError(f"tokenizer path does not exist: {tokenizer_path}")
    try:
        return AutoTokenizer.from_pretrained(tokenizer_path, local_files_only=True)
    except Exception as exc:  # noqa: BLE001 - surface package-specific load failures with context.
        raise RuntimeError(f"failed to load Hugging Face tokenizer from {tokenizer_path}") from exc


def encode_probe(tokenizer: Any, text: str) -> dict[str, Any]:
    encoded = tokenizer(text, add_special_tokens=True, return_attention_mask=True, return_token_type_ids=False)
    input_ids = [int(value) for value in encoded["input_ids"]]
    attention_mask = [int(value) for value in encoded["attention_mask"]]
    if len(input_ids) != len(attention_mask):
        raise RuntimeError("tokenizer produced input_ids and attention_mask with different lengths")
    return {
        "input_ids": input_ids,
        "attention_mask": attention_mask,
        "input_ids_sha256": sha256_json(input_ids),
        "attention_mask_sha256": sha256_json(attention_mask),
        "token_count": len(input_ids),
        "attention_count": sum(attention_mask),
        "first_token_id": input_ids[0] if input_ids else None,
        "last_token_id": input_ids[-1] if input_ids else None,
    }


def build_baseline(args: argparse.Namespace) -> dict[str, Any]:
    tokenizer = load_tokenizer(args.tokenizer_path)
    probes_source = args.probes_source
    probes_input = load_probes(probes_source)
    tokenizer_json = args.tokenizer_path / "tokenizer.json"
    config_json = args.tokenizer_path / "config.json"
    revision_file = args.tokenizer_path / "refs" / "main"

    probes: list[dict[str, Any]] = []
    for probe in probes_input:
        encoded = encode_probe(tokenizer, probe["text"])
        probes.append(
            {
                "label": probe["label"],
                "chars": len(probe["text"]),
                **encoded,
            }
        )

    return {
        "script_version": SCRIPT_VERSION,
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "mode": args.mode,
        "model": args.model,
        "tokenizer_path": str(args.tokenizer_path),
        "tokenizer_class": tokenizer.__class__.__name__,
        "tokenizer_json_sha256": sha256_file(tokenizer_json),
        "config_json_sha256": sha256_file(config_json),
        "local_revision": revision_file.read_text(encoding="utf-8").strip() if revision_file.exists() else None,
        "package_versions": {
            "transformers": package_version("transformers"),
            "tokenizers": package_version("tokenizers"),
            "sentencepiece": package_version("sentencepiece"),
        },
        "probes_source": str(probes_source),
        "probes_source_sha256": sha256_file(probes_source),
        "raw_probe_texts_logged": False,
        "probe_count": len(probes_input),
        "probes": probes,
        "passed": True,
    }


def render_markdown(result: dict[str, Any]) -> str:
    config = {key: value for key, value in result.items() if key not in {"probes", "passed"}}
    lines = [
        "# Hugging Face Tokenizer Baseline — M012 S01",
        "",
        "## Effective Configuration",
        "",
        "```json",
        json.dumps(config, ensure_ascii=False, indent=2, sort_keys=True),
        "```",
        "",
        "## Probe Tokenization Summary",
        "",
        "| Label | Chars | Tokens | Attention Count | First ID | Last ID | IDs SHA256 | Mask SHA256 |",
        "|---|---:|---:|---:|---:|---:|---|---|",
    ]

    for probe in result["probes"]:
        lines.append(
            "| {label} | {chars} | {token_count} | {attention_count} | {first_token_id} | {last_token_id} | `{ids_hash}` | `{mask_hash}` |".format(
                label=probe["label"],
                chars=probe["chars"],
                token_count=probe["token_count"],
                attention_count=probe["attention_count"],
                first_token_id=probe["first_token_id"],
                last_token_id=probe["last_token_id"],
                ids_hash=probe["input_ids_sha256"],
                mask_hash=probe["attention_mask_sha256"],
            )
        )

    lines.extend([
        "",
        "## Token Evidence",
        "",
        "```json",
        json.dumps(result["probes"], ensure_ascii=False, indent=2, sort_keys=True),
        "```",
        "",
        "## Verdict",
        "",
        "PASS" if result["passed"] else "FAIL",
        "",
        "Raw probe texts are intentionally excluded from this artifact.",
        "",
    ])
    return "\n".join(lines)


def main() -> int:
    args = parse_args()
    result = build_baseline(args)
    artifact = render_markdown(result)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(artifact, encoding="utf-8")
    print(artifact)
    return 0 if result["passed"] else 2


if __name__ == "__main__":
    raise SystemExit(main())
