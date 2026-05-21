#!/usr/bin/env python3
"""Verify the local USER-bge-m3 ONNX export contract.

This verifier checks the existing local export artifact against tracked manifest,
source provenance, and export metadata. It does not regenerate the ONNX binary
and must not be used as byte-for-byte reproducibility evidence for a fresh
export. It is a contract/provenance verifier for the current ignored artifact.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import sys
from typing import Any

SCRIPT_VERSION = 1
APPROVED_ARTIFACT_ROOTS = (
    Path(".gsd/runtime/onnx"),
    Path("tei-models"),
)
REQUIRED_PYTHON_PREFIX = "3.13.12"
REQUIRED_EXPORT_PACKAGES = {
    "torch": "2.12.0",
    "transformers": "4.51.3",
    "onnx": "1.21.0",
    "onnxruntime": "1.26.0",
    "safetensors": "0.7.0",
}
REQUIRED_EXPORT = {
    "script": "tools/export_user_bge_m3_dense_onnx.py",
    "opset": 17,
    "sequence_length": 128,
    "pooling": "cls_token_last_hidden_state[:,0]",
    "normalization": "torch.nn.functional.normalize(..., p=2, dim=1)",
    "output_name": "dense_vecs",
    "expected_dimensions": 1024,
}


class ContractError(RuntimeError):
    def __init__(self, label: str, message: str) -> None:
        super().__init__(message)
        self.label = label


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Verify fd USER-bge-m3 ONNX export contract.")
    parser.add_argument("--manifest", type=Path, default=Path("docs/onnx-artifacts/user-bge-m3-dense-fp32.json"))
    parser.add_argument("--provenance", type=Path, default=Path(".gsd/runtime/onnx/m010-s03/source-provenance.json"))
    parser.add_argument("--export-metadata", type=Path, default=Path(".gsd/runtime/onnx/m010-s03/export-metadata.json"))
    parser.add_argument("--repo-root", type=Path, default=Path("."))
    parser.add_argument("--allow-missing-artifact", action="store_true", help="Report a missing ONNX binary without failing. This is metadata-only and not runtime proof.")
    return parser.parse_args()


def safe_path_display(path: str | Path, repo_root: Path | None = None) -> str:
    candidate = Path(path)
    if repo_root is not None:
        try:
            return candidate.resolve().relative_to(repo_root.resolve()).as_posix()
        except (OSError, ValueError):
            pass
    if candidate.is_absolute():
        return f".../{candidate.name}" if candidate.name else "..."
    return candidate.as_posix()


def repo_relative_path(value: str | Path, roots: tuple[Path, ...]) -> Path:
    path = Path(value)
    if path.is_absolute():
        raise ContractError("path_policy", f"path must be repo-relative: {safe_path_display(path)}")
    normalized = Path(os.path.normpath(path.as_posix()))
    if normalized == Path(".") or str(normalized).startswith(".."):
        raise ContractError("path_policy", f"path must not traverse outside the repository: {safe_path_display(path)}")
    if not any(normalized == root or root in normalized.parents for root in roots):
        allowed = ", ".join(root.as_posix() for root in roots)
        raise ContractError("path_policy", f"path must be under approved roots ({allowed}): {safe_path_display(path)}")
    return normalized


def load_json(path: Path, label: str) -> dict[str, Any]:
    try:
        data = json.loads(path.read_text(encoding="utf-8"))
    except FileNotFoundError as exc:
        raise ContractError(label, f"metadata file not found: {safe_path_display(path)}") from exc
    except json.JSONDecodeError as exc:
        raise ContractError(label, f"metadata file is invalid JSON: {safe_path_display(path)}: {exc}") from exc
    if not isinstance(data, dict):
        raise ContractError(label, f"metadata file must contain a JSON object: {safe_path_display(path)}")
    return data


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def require(condition: bool, label: str, message: str) -> None:
    if not condition:
        raise ContractError(label, message)


def source_artifact_map(items: Any) -> dict[str, dict[str, Any]]:
    result: dict[str, dict[str, Any]] = {}
    if not isinstance(items, list):
        return result
    for item in items:
        if not isinstance(item, dict):
            continue
        path = item.get("path")
        if isinstance(path, str):
            result[Path(path).as_posix()] = item
    return result


def verify_manifest_shape(manifest: dict[str, Any]) -> None:
    require(manifest.get("production_default") is False, "production_default", "manifest production_default must be false")
    artifact = manifest.get("artifact")
    require(isinstance(artifact, dict), "manifest_artifact", "manifest missing artifact object")
    require(artifact.get("git_tracked") is False, "git_tracking", "ONNX artifact must remain untracked")
    require(artifact.get("format") == "onnx", "manifest_artifact", "artifact format must be onnx")
    require(artifact.get("dtype") == "fp32", "manifest_artifact", "artifact dtype must be fp32")
    require(isinstance(artifact.get("sha256"), str) and len(artifact["sha256"]) == 64, "manifest_artifact", "artifact sha256 must be a 64-character string")
    require(isinstance(artifact.get("size_bytes"), int) and artifact["size_bytes"] > 0, "manifest_artifact", "artifact size_bytes must be positive")

    export = manifest.get("export")
    require(isinstance(export, dict), "manifest_export", "manifest missing export object")
    for key in ("script", "opset", "sequence_length", "pooling", "normalization"):
        require(export.get(key) == REQUIRED_EXPORT[key], "manifest_export", f"manifest export.{key} mismatch")
    packages = export.get("packages")
    require(isinstance(packages, dict), "manifest_export", "manifest export.packages must be an object")
    require(packages.get("python") == REQUIRED_PYTHON_PREFIX, "manifest_export", f"manifest export package python must be {REQUIRED_PYTHON_PREFIX}")
    for name, expected in REQUIRED_EXPORT_PACKAGES.items():
        require(packages.get(name) == expected, "manifest_export", f"manifest export package {name} must be {expected}")

    runtime = manifest.get("runtime")
    require(isinstance(runtime, dict), "manifest_runtime", "manifest missing runtime object")
    require(runtime.get("provider") == "CPUExecutionProvider", "manifest_runtime", "runtime provider must be CPUExecutionProvider")
    require(runtime.get("expected_dimensions") == REQUIRED_EXPORT["expected_dimensions"], "manifest_runtime", "runtime expected_dimensions mismatch")
    require(runtime.get("expected_normalized") is True, "manifest_runtime", "runtime expected_normalized must be true")
    require(runtime.get("validated_max_sequence_length") == 1024, "manifest_runtime", "validated max sequence length must remain 1024")


def verify_artifact(manifest: dict[str, Any], repo_root: Path, allow_missing: bool) -> dict[str, Any]:
    artifact = manifest["artifact"]
    rel = repo_relative_path(artifact["local_path"], (Path(".gsd/runtime/onnx"),))
    path = repo_root / rel
    result: dict[str, Any] = {"path": rel.as_posix(), "expected_size_bytes": artifact["size_bytes"], "expected_sha256": artifact["sha256"]}
    if not path.exists():
        result.update({"present": False, "verified": False, "reason": "missing_allowed" if allow_missing else "missing"})
        if allow_missing:
            return result
        raise ContractError("artifact_missing", f"ONNX artifact missing: {safe_path_display(path, repo_root)}")
    require(path.is_file(), "artifact_path", f"ONNX artifact path is not a file: {safe_path_display(path, repo_root)}")
    size = path.stat().st_size
    require(size == artifact["size_bytes"], "artifact_size", f"ONNX artifact size mismatch for {safe_path_display(path, repo_root)}")
    sha = sha256_file(path)
    require(sha == artifact["sha256"], "artifact_sha256", f"ONNX artifact sha256 mismatch for {safe_path_display(path, repo_root)}")
    result.update({"present": True, "verified": True, "size_bytes": size, "sha256": sha})
    return result


def verify_source_files(manifest: dict[str, Any], provenance: dict[str, Any]) -> list[dict[str, Any]]:
    model = manifest.get("model")
    require(isinstance(model, dict), "manifest_model", "manifest missing model object")
    require(provenance.get("model_revision") == model.get("revision"), "model_revision", "provenance model revision does not match manifest")
    require(provenance.get("production_runtime_changed") is False, "production_default", "provenance must not change production runtime")

    expected_files = model.get("source_files")
    require(isinstance(expected_files, dict), "manifest_model", "manifest model.source_files must be an object")
    observed = source_artifact_map(provenance.get("source_artifacts"))
    results = []
    local_source = model.get("local_source_path")
    require(isinstance(local_source, str) and local_source, "manifest_model", "manifest model.local_source_path missing")
    source_root = repo_relative_path(local_source, (Path("tei-models"),))

    for rel, expected in expected_files.items():
        require(isinstance(expected, dict), "manifest_model", f"source_files.{rel} must be an object")
        source_path = (source_root / rel).as_posix()
        item = observed.get(source_path)
        require(isinstance(item, dict), "source_file", f"source provenance missing {source_path}")
        require(item.get("exists") is True, "source_file", f"source artifact not present in provenance: {source_path}")
        require(item.get("size_bytes") == expected.get("size_bytes"), "source_file", f"source artifact size mismatch: {source_path}")
        require(item.get("sha256") == expected.get("sha256"), "source_file", f"source artifact sha256 mismatch: {source_path}")
        results.append({"path": source_path, "size_bytes": item.get("size_bytes"), "sha256": item.get("sha256"), "verified": True})
    return results


def verify_export_metadata(manifest: dict[str, Any], metadata: dict[str, Any]) -> dict[str, Any]:
    require(metadata.get("status") == "success", "export_metadata", "export metadata status must be success")
    require(metadata.get("production_runtime_changed") is False, "production_default", "export metadata must not change production runtime")
    require(metadata.get("script") == REQUIRED_EXPORT["script"], "export_metadata", "export script mismatch")
    require(metadata.get("opset") == REQUIRED_EXPORT["opset"], "export_metadata", "export opset mismatch")
    require(metadata.get("sequence_length") == REQUIRED_EXPORT["sequence_length"], "export_metadata", "export sequence length mismatch")

    packages = metadata.get("packages")
    require(isinstance(packages, dict), "export_metadata", "export metadata packages must be an object")
    python_version = metadata.get("python")
    require(isinstance(python_version, str) and python_version.startswith(REQUIRED_PYTHON_PREFIX), "export_metadata", f"export metadata python must start with {REQUIRED_PYTHON_PREFIX}")
    for name, expected in REQUIRED_EXPORT_PACKAGES.items():
        require(packages.get(name) == expected, "export_metadata", f"export metadata package {name} must be {expected}")

    export = metadata.get("export")
    require(isinstance(export, dict), "export_metadata", "export metadata export object missing")
    require(export.get("pooling") == REQUIRED_EXPORT["pooling"], "export_metadata", "export pooling mismatch")
    require(export.get("normalization") == REQUIRED_EXPORT["normalization"], "export_metadata", "export normalization mismatch")
    require(export.get("dtype") == "fp32", "export_metadata", "export dtype mismatch")
    require(export.get("output_names") == [REQUIRED_EXPORT["output_name"]], "export_metadata", "export output names mismatch")
    require(export.get("dummy_output_shape") == [1, REQUIRED_EXPORT["expected_dimensions"]], "export_metadata", "dummy output shape mismatch")

    ort = metadata.get("onnxruntime")
    require(isinstance(ort, dict), "export_metadata", "onnxruntime metadata missing")
    require("CPUExecutionProvider" in (ort.get("providers") or []), "export_metadata", "CPUExecutionProvider missing from export metadata")
    require(ort.get("dummy_output_shape") == [1, REQUIRED_EXPORT["expected_dimensions"]], "export_metadata", "ONNX Runtime dummy output shape mismatch")

    output_artifacts = metadata.get("output_artifacts")
    require(isinstance(output_artifacts, list), "export_metadata", "output_artifacts must be a list")
    expected_sha = manifest["artifact"]["sha256"]
    expected_size = manifest["artifact"]["size_bytes"]
    matched = [item for item in output_artifacts if isinstance(item, dict) and item.get("sha256") == expected_sha and item.get("size_bytes") == expected_size]
    require(bool(matched), "export_metadata", "export metadata does not include the manifest ONNX artifact checksum")
    return {
        "script": metadata.get("script"),
        "opset": metadata.get("opset"),
        "sequence_length": metadata.get("sequence_length"),
        "packages": {"python": REQUIRED_PYTHON_PREFIX, **{name: packages.get(name) for name in REQUIRED_EXPORT_PACKAGES}},
        "output_artifact_recorded": True,
    }


def main() -> int:
    args = parse_args()
    repo_root = args.repo_root.resolve()
    manifest = load_json(args.manifest, "manifest")
    provenance = load_json(args.provenance, "provenance")
    metadata = load_json(args.export_metadata, "export_metadata")

    verify_manifest_shape(manifest)
    artifact = verify_artifact(manifest, repo_root, args.allow_missing_artifact)
    source_files = verify_source_files(manifest, provenance)
    export = verify_export_metadata(manifest, metadata)

    print(
        json.dumps(
            {
                "script_version": SCRIPT_VERSION,
                "repo_root": ".",
                "verdict": "pass",
                "claim_scope": "existing_artifact_contract_verification_not_regenerated_export",
                "production_default": False,
                "artifact": artifact,
                "source_files_verified": len(source_files),
                "export": export,
                "remaining_blocker": "exact ONNX binary still needs immutable external source or separately validated reproducible export before hosted proof",
            },
            ensure_ascii=False,
            indent=2,
            sort_keys=True,
        )
    )
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except ContractError as exc:
        print(json.dumps({"script_version": SCRIPT_VERSION, "verdict": "fail", "error_label": exc.label, "error": str(exc)}, ensure_ascii=False, sort_keys=True), file=sys.stderr)
        raise SystemExit(1)
