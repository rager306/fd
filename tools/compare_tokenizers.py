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
import subprocess
import tempfile
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
        choices=["baseline", "go-current"],
        default="baseline",
        help="Operation mode. Writes the HF baseline or compares the current Go tokenizer against it.",
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
        "--baseline",
        type=Path,
        default=Path("benchmark-results/fd-tokenizer-baseline-m012-s01.txt"),
        help="Hugging Face baseline artifact to compare against in go-current mode.",
    )
    parser.add_argument(
        "--output",
        type=Path,
        default=None,
        help="Markdown artifact path to write. Defaults depend on mode.",
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


def extract_token_evidence(artifact_path: Path) -> list[dict[str, Any]]:
    text = artifact_path.read_text(encoding="utf-8")
    marker = "## Token Evidence"
    if marker not in text:
        raise RuntimeError(f"token evidence section not found in {artifact_path}")
    after_marker = text.split(marker, 1)[1]
    try:
        json_block = after_marker.split("```json", 1)[1].split("```", 1)[0]
    except IndexError as exc:
        raise RuntimeError(f"token evidence JSON block not found in {artifact_path}") from exc
    value = json.loads(json_block)
    if not isinstance(value, list):
        raise RuntimeError("token evidence JSON must be a list")
    return value


def first_mismatch(left: list[int], right: list[int]) -> int | None:
    for idx, (left_value, right_value) in enumerate(zip(left, right)):
        if left_value != right_value:
            return idx
    if len(left) != len(right):
        return min(len(left), len(right))
    return None


def list_window(values: list[int], mismatch: int | None, radius: int = 3) -> list[int]:
    if mismatch is None:
        return values[: min(len(values), 6)]
    start = max(0, mismatch - radius)
    end = min(len(values), mismatch + radius + 1)
    return values[start:end]


def run_go_current_tokenizer(args: argparse.Namespace, probes: list[dict[str, str]]) -> list[dict[str, Any]]:
    go_tokenizer_path = args.tokenizer_path / "tokenizer.json" if args.tokenizer_path.is_dir() else args.tokenizer_path
    payload = {
        "tokenizer_path": str(go_tokenizer_path.resolve()),
        "probes": probes,
    }
    go_source = r'''
package main

import (
  "encoding/json"
  "fmt"
  "os"

  "github.com/sugarme/tokenizer/pretrained"
)

type probe struct {
  Label string `json:"label"`
  Text string `json:"text"`
}

type payload struct {
  TokenizerPath string `json:"tokenizer_path"`
  Probes []probe `json:"probes"`
}

type result struct {
  Label string `json:"label"`
  InputIDs []int `json:"input_ids"`
  AttentionMask []int `json:"attention_mask"`
}

func main() {
  if len(os.Args) != 2 {
    fmt.Fprintln(os.Stderr, "usage: go-tokenizer <payload.json>")
    os.Exit(2)
  }
  data, err := os.ReadFile(os.Args[1])
  if err != nil { panic(err) }
  var in payload
  if err := json.Unmarshal(data, &in); err != nil { panic(err) }
  tk, err := pretrained.FromFile(in.TokenizerPath)
  if err != nil { panic(err) }
  out := make([]result, 0, len(in.Probes))
  for _, p := range in.Probes {
    enc, err := tk.EncodeSingle(p.Text, true)
    if err != nil { panic(err) }
    out = append(out, result{Label: p.Label, InputIDs: enc.Ids, AttentionMask: enc.AttentionMask})
  }
  encoded, err := json.Marshal(out)
  if err != nil { panic(err) }
  fmt.Println(string(encoded))
}
'''
    with tempfile.TemporaryDirectory(prefix="fd-tokenizer-") as temp_dir:
        temp = Path(temp_dir)
        payload_path = temp / "payload.json"
        source_path = temp / "go_tokenizer_probe.go"
        payload_path.write_text(json.dumps(payload, ensure_ascii=False), encoding="utf-8")
        source_path.write_text(go_source, encoding="utf-8")
        completed = subprocess.run(
            ["go", "run", str(source_path), str(payload_path)],
            cwd="api",
            text=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            check=False,
            timeout=120,
        )
    if completed.returncode != 0:
        raise RuntimeError(f"go tokenizer probe failed with exit {completed.returncode}: {completed.stderr.strip()}")
    value = json.loads(completed.stdout)
    if not isinstance(value, list):
        raise RuntimeError("go tokenizer output must be a list")
    return value


def build_go_current_comparison(args: argparse.Namespace) -> dict[str, Any]:
    probes_input = load_probes(args.probes_source)
    baseline_items = extract_token_evidence(args.baseline)
    go_items = run_go_current_tokenizer(args, probes_input)
    baseline_by_label = {item["label"]: item for item in baseline_items}
    all_passed = True
    comparisons: list[dict[str, Any]] = []
    for probe, go_item in zip(probes_input, go_items):
        label = probe["label"]
        expected = baseline_by_label[label]
        go_ids = [int(value) for value in go_item["input_ids"]]
        go_mask = [int(value) for value in go_item["attention_mask"]]
        hf_ids = [int(value) for value in expected["input_ids"]]
        hf_mask = [int(value) for value in expected["attention_mask"]]
        ids_equal = go_ids == hf_ids
        mask_equal = go_mask == hf_mask
        passed = ids_equal and mask_equal
        all_passed = all_passed and passed
        ids_mismatch = first_mismatch(hf_ids, go_ids)
        mask_mismatch = first_mismatch(hf_mask, go_mask)
        comparisons.append(
            {
                "label": label,
                "chars": len(probe["text"]),
                "hf_token_count": len(hf_ids),
                "go_token_count": len(go_ids),
                "hf_input_ids_sha256": sha256_json(hf_ids),
                "go_input_ids_sha256": sha256_json(go_ids),
                "hf_attention_mask_sha256": sha256_json(hf_mask),
                "go_attention_mask_sha256": sha256_json(go_mask),
                "input_ids_equal": ids_equal,
                "attention_mask_equal": mask_equal,
                "first_input_ids_mismatch_index": ids_mismatch,
                "first_attention_mask_mismatch_index": mask_mismatch,
                "hf_ids_window": list_window(hf_ids, ids_mismatch),
                "go_ids_window": list_window(go_ids, ids_mismatch),
                "passed": passed,
            }
        )
    return {
        "script_version": SCRIPT_VERSION,
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "mode": args.mode,
        "model": args.model,
        "baseline": str(args.baseline),
        "probes_source": str(args.probes_source),
        "probes_source_sha256": sha256_file(args.probes_source),
        "tokenizer_path": str(args.tokenizer_path),
        "go_tokenizer_library": "github.com/sugarme/tokenizer/pretrained.FromFile + EncodeSingle(addSpecialTokens=true)",
        "raw_probe_texts_logged": False,
        "probe_count": len(probes_input),
        "comparisons": comparisons,
        "passed": all_passed,
    }


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


def render_comparison_markdown(result: dict[str, Any]) -> str:
    config = {key: value for key, value in result.items() if key not in {"comparisons", "passed"}}
    lines = [
        "# Go Tokenizer Comparison — M012 S02",
        "",
        "## Effective Configuration",
        "",
        "```json",
        json.dumps(config, ensure_ascii=False, indent=2, sort_keys=True),
        "```",
        "",
        "## HF Baseline vs Current Go Tokenizer",
        "",
        "| Label | Chars | HF Tokens | Go Tokens | IDs Equal | Mask Equal | First IDs Mismatch | Passed |",
        "|---|---:|---:|---:|---|---|---:|---|",
    ]
    for item in result["comparisons"]:
        mismatch = item["first_input_ids_mismatch_index"]
        mismatch_text = "" if mismatch is None else str(mismatch)
        lines.append(
            "| {label} | {chars} | {hf_count} | {go_count} | {ids_equal} | {mask_equal} | {mismatch} | {passed} |".format(
                label=item["label"],
                chars=item["chars"],
                hf_count=item["hf_token_count"],
                go_count=item["go_token_count"],
                ids_equal="yes" if item["input_ids_equal"] else "no",
                mask_equal="yes" if item["attention_mask_equal"] else "no",
                mismatch=mismatch_text,
                passed="yes" if item["passed"] else "no",
            )
        )
    lines.extend([
        "",
        "## Mismatch Evidence",
        "",
        "```json",
        json.dumps(result["comparisons"], ensure_ascii=False, indent=2, sort_keys=True),
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
    if args.output is None:
        args.output = DEFAULT_OUTPUT if args.mode == "baseline" else Path("benchmark-results/fd-tokenizer-go-current-m012-s02.txt")
    if args.mode == "baseline":
        result = build_baseline(args)
        artifact = render_markdown(result)
    else:
        result = build_go_current_comparison(args)
        artifact = render_comparison_markdown(result)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(artifact, encoding="utf-8")
    print(artifact)
    return 0 if result["passed"] else 2


if __name__ == "__main__":
    raise SystemExit(main())
