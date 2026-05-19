#!/usr/bin/env python3
"""Export deepvk/USER-bge-m3 to a dense-only FP32 ONNX candidate.

This is a local spike tool. It does not change production runtime defaults.
Artifacts are written under an ignored runtime directory by default.
"""

from __future__ import annotations

import argparse
import hashlib
import importlib.metadata
import inspect
import json
from pathlib import Path
import platform
import subprocess
import sys
import time
import traceback
from typing import Any


def sha256_file(path: Path) -> str | None:
    if not path.exists() or not path.is_file():
        return None
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def package_version(name: str) -> str | None:
    try:
        return importlib.metadata.version(name)
    except importlib.metadata.PackageNotFoundError:
        return None


def run_metadata_command(args: list[str], timeout: int = 10) -> str | None:
    try:
        result = subprocess.run(
            args,
            check=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            timeout=timeout,
        )
        return result.stdout.strip() or None
    except Exception:
        return None


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Export USER-bge-m3 dense FP32 ONNX.")
    parser.add_argument("--model-path", type=Path, required=True, help="Local model snapshot path.")
    parser.add_argument("--output-dir", type=Path, required=True, help="Ignored output directory for ONNX artifacts.")
    parser.add_argument("--output-name", default="user-bge-m3-dense.onnx", help="ONNX output filename.")
    parser.add_argument("--opset", type=int, default=17, help="ONNX opset version.")
    parser.add_argument("--sequence-length", type=int, default=128, help="Dummy export sequence length.")
    parser.add_argument("--metadata-name", default="export-metadata.json", help="Metadata JSON filename.")
    return parser.parse_args()


def collect_source_artifacts(model_path: Path) -> list[dict[str, Any]]:
    rels = [
        "model.safetensors",
        "tokenizer.json",
        "config.json",
        "modules.json",
        "sentence_bert_config.json",
        "1_Pooling/config.json",
        "tokenizer_config.json",
        "special_tokens_map.json",
        "sentencepiece.bpe.model",
    ]
    artifacts = []
    for rel in rels:
        path = model_path / rel
        artifacts.append(
            {
                "path": str(path),
                "exists": path.exists(),
                "size_bytes": path.stat().st_size if path.exists() else None,
                "sha256": sha256_file(path) if path.exists() else None,
            }
        )
    return artifacts


def collect_output_artifacts(output_dir: Path) -> list[dict[str, Any]]:
    artifacts = []
    if not output_dir.exists():
        return artifacts
    for path in sorted(output_dir.iterdir()):
        if path.is_file() and path.name != "export-metadata.json":
            artifacts.append(
                {
                    "path": str(path),
                    "size_bytes": path.stat().st_size,
                    "sha256": sha256_file(path),
                }
            )
    return artifacts


def write_metadata(path: Path, metadata: dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(metadata, ensure_ascii=False, indent=2, sort_keys=True), encoding="utf-8")


def main() -> int:
    args = parse_args()
    started = time.time()
    args.output_dir.mkdir(parents=True, exist_ok=True)
    output_path = args.output_dir / args.output_name
    metadata_path = args.output_dir / args.metadata_name

    metadata: dict[str, Any] = {
        "script": "tools/export_user_bge_m3_dense_onnx.py",
        "status": "started",
        "production_runtime_changed": False,
        "model_path": str(args.model_path),
        "output_dir": str(args.output_dir),
        "output_path": str(output_path),
        "opset": args.opset,
        "sequence_length": args.sequence_length,
        "python": sys.version,
        "platform": platform.platform(),
        "git": {
            "commit": run_metadata_command(["git", "rev-parse", "HEAD"]),
            "branch": run_metadata_command(["git", "branch", "--show-current"]),
        },
        "packages": {
            "torch": package_version("torch"),
            "transformers": package_version("transformers"),
            "onnx": package_version("onnx"),
            "onnxruntime": package_version("onnxruntime"),
            "safetensors": package_version("safetensors"),
        },
        "source_artifacts": collect_source_artifacts(args.model_path),
        "export": {
            "input_names": ["input_ids", "attention_mask"],
            "output_names": ["dense_vecs"],
            "dynamic_axes": {
                "input_ids": {"0": "batch_size", "1": "sequence_length"},
                "attention_mask": {"0": "batch_size", "1": "sequence_length"},
                "dense_vecs": {"0": "batch_size", "1": "embedding"},
            },
            "pooling": "cls_token_last_hidden_state[:,0]",
            "normalization": "torch.nn.functional.normalize(..., p=2, dim=1)",
            "dtype": "fp32",
        },
    }

    try:
        import torch
        import torch.nn.functional as F
        from transformers import AutoModel, AutoTokenizer

        class DenseOnlyModel(torch.nn.Module):
            def __init__(self, model_path: Path) -> None:
                super().__init__()
                self.model = AutoModel.from_pretrained(str(model_path), local_files_only=True)
                self.model.eval()

            def forward(self, input_ids: torch.Tensor, attention_mask: torch.Tensor) -> torch.Tensor:
                outputs = self.model(input_ids=input_ids, attention_mask=attention_mask, return_dict=True)
                dense_vecs = outputs.last_hidden_state[:, 0]
                return F.normalize(dense_vecs, p=2, dim=1)

        tokenizer = AutoTokenizer.from_pretrained(str(args.model_path), local_files_only=True)
        model = DenseOnlyModel(args.model_path)
        model.eval()

        dummy_texts = ["Проверка плотного эмбеддинга для ONNX экспорта."]
        encoded = tokenizer(
            dummy_texts,
            padding="max_length",
            truncation=True,
            max_length=args.sequence_length,
            return_tensors="pt",
        )
        input_ids = encoded["input_ids"]
        attention_mask = encoded["attention_mask"]

        export_kwargs: dict[str, Any] = {
            "model": model,
            "args": (input_ids, attention_mask),
            "f": str(output_path),
            "export_params": True,
            "opset_version": args.opset,
            "do_constant_folding": True,
            "input_names": ["input_ids", "attention_mask"],
            "output_names": ["dense_vecs"],
            "dynamic_axes": {
                "input_ids": {0: "batch_size", 1: "sequence_length"},
                "attention_mask": {0: "batch_size", 1: "sequence_length"},
                "dense_vecs": {0: "batch_size", 1: "embedding"},
            },
        }
        signature = inspect.signature(torch.onnx.export)
        if "external_data" in signature.parameters:
            export_kwargs["external_data"] = True
        if "dynamo" in signature.parameters:
            export_kwargs["dynamo"] = False

        with torch.no_grad():
            expected = model(input_ids, attention_mask)
            metadata["export"]["dummy_output_shape"] = list(expected.shape)
            metadata["export"]["dummy_output_l2_norm"] = float(torch.linalg.vector_norm(expected[0]).item())
            torch.onnx.export(**export_kwargs)

        # Validate ONNX structure and CPU EP load.
        import onnx
        import onnxruntime as ort

        onnx_model = onnx.load(str(output_path), load_external_data=True)
        onnx.checker.check_model(onnx_model)
        session = ort.InferenceSession(str(output_path), providers=["CPUExecutionProvider"])
        ort_outputs = session.run(
            None,
            {
                "input_ids": input_ids.cpu().numpy(),
                "attention_mask": attention_mask.cpu().numpy(),
            },
        )
        metadata["onnxruntime"] = {
            "providers": session.get_providers(),
            "inputs": [
                {"name": item.name, "shape": item.shape, "type": item.type}
                for item in session.get_inputs()
            ],
            "outputs": [
                {"name": item.name, "shape": item.shape, "type": item.type}
                for item in session.get_outputs()
            ],
            "dummy_output_shape": list(ort_outputs[0].shape),
        }
        metadata["output_artifacts"] = collect_output_artifacts(args.output_dir)
        metadata["duration_seconds"] = round(time.time() - started, 3)
        metadata["status"] = "success"
        write_metadata(metadata_path, metadata)
        print(json.dumps({"status": "success", "metadata": str(metadata_path), "output": str(output_path)}, ensure_ascii=False))
        return 0
    except Exception as exc:  # noqa: BLE001 - spike tool must persist full failure mode.
        metadata["duration_seconds"] = round(time.time() - started, 3)
        metadata["status"] = "failed"
        metadata["error"] = {
            "type": type(exc).__name__,
            "message": str(exc),
            "traceback": traceback.format_exc(limit=20),
        }
        metadata["output_artifacts"] = collect_output_artifacts(args.output_dir)
        write_metadata(metadata_path, metadata)
        print(json.dumps({"status": "failed", "metadata": str(metadata_path), "error_type": type(exc).__name__}, ensure_ascii=False), file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
