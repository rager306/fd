#!/usr/bin/env python3
"""Verify local ONNX and native tokenizer artifacts against tracked manifests.

This script is intentionally local/CI-friendly and does not download artifacts.
It fails fast with actionable, sanitized messages when required ignored files are
missing or checksums do not match tracked metadata.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
from typing import Any

SCRIPT_VERSION = 2
APPROVED_ARTIFACT_ROOTS = (
    Path(".gsd/runtime/onnx"),
    Path(".gsd/runtime/tokenizers"),
    Path(".gsd/runtime/onnxruntime"),
    Path("tei-models"),
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Verify fd ONNX/native tokenizer artifacts against manifests.")
    parser.add_argument("--onnx-manifest", type=Path, required=True)
    parser.add_argument("--native-tokenizer-manifest", type=Path, required=True)
    parser.add_argument("--repo-root", type=Path, default=Path("."))
    parser.add_argument("--allow-missing", action="store_true", help="Report missing artifacts without failing. Useful for default CI jobs that do not provision opt-in ONNX artifacts.")
    return parser.parse_args()


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def load_manifest(path: Path) -> dict[str, Any]:
    try:
        data = json.loads(path.read_text(encoding="utf-8"))
    except FileNotFoundError as exc:
        raise RuntimeError(f"manifest not found: {path}") from exc
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"manifest is not valid JSON: {path}: {exc}") from exc
    if not isinstance(data, dict):
        raise RuntimeError(f"manifest must be a JSON object: {path}")
    return data


def git_tracked_paths(repo_root: Path) -> set[str]:
    try:
        result = subprocess.run(
            ["git", "ls-files"],
            cwd=repo_root,
            check=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            timeout=30,
        )
    except (FileNotFoundError, subprocess.CalledProcessError, subprocess.TimeoutExpired):
        return set()
    return set(result.stdout.splitlines())


def require_bool_false(data: dict[str, Any], path: str) -> None:
    current: Any = data
    for part in path.split("."):
        if not isinstance(current, dict) or part not in current:
            raise RuntimeError(f"missing required manifest field: {path}")
        current = current[part]
    if current is not False:
        raise RuntimeError(f"manifest field {path} must be false, got {current!r}")


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


def repo_relative_artifact_path(value: str | Path) -> Path:
    path = Path(value)
    if path.is_absolute():
        raise RuntimeError(f"artifact path must be repo-relative: {safe_path_display(path)}")
    normalized = Path(os.path.normpath(path.as_posix()))
    if normalized == Path(".") or str(normalized).startswith(".."):
        raise RuntimeError(f"artifact path must not traverse outside the repository: {safe_path_display(path)}")
    if not any(normalized == root or root in normalized.parents for root in APPROVED_ARTIFACT_ROOTS):
        roots = ", ".join(root.as_posix() for root in APPROVED_ARTIFACT_ROOTS)
        raise RuntimeError(f"artifact path must be under an approved root ({roots}): {safe_path_display(path)}")
    return normalized


def artifact_local_path(manifest: dict[str, Any], manifest_path: Path, repo_root: Path) -> Path:
    artifact = manifest.get("artifact")
    if not isinstance(artifact, dict):
        raise RuntimeError(f"manifest missing artifact object: {manifest_path}")
    local_path = artifact.get("local_path")
    if not isinstance(local_path, str) or not local_path:
        raise RuntimeError(f"manifest missing artifact.local_path: {manifest_path}")
    path = repo_relative_artifact_path(local_path)
    return repo_root / path


def expected_artifact(manifest: dict[str, Any], manifest_path: Path) -> tuple[str, int | None, str]:
    artifact = manifest.get("artifact")
    if not isinstance(artifact, dict):
        raise RuntimeError(f"manifest missing artifact object: {manifest_path}")
    sha = artifact.get("sha256")
    if not isinstance(sha, str) or len(sha) != 64:
        raise RuntimeError(f"manifest artifact.sha256 must be a 64-character hex string: {manifest_path}")
    size = artifact.get("size_bytes")
    if size is not None and (not isinstance(size, int) or size <= 0):
        raise RuntimeError(f"manifest artifact.size_bytes must be a positive integer: {manifest_path}")
    artifact_id = manifest.get("artifact_id")
    if not isinstance(artifact_id, str) or not artifact_id:
        raise RuntimeError(f"manifest missing artifact_id: {manifest_path}")
    return sha, size, artifact_id


def verify_one(label: str, manifest_path: Path, repo_root: Path, tracked: set[str], allow_missing: bool) -> dict[str, Any]:
    manifest = load_manifest(manifest_path)
    require_bool_false(manifest, "production_default")
    require_bool_false(manifest, "artifact.git_tracked")
    expected_sha, expected_size, artifact_id = expected_artifact(manifest, manifest_path)
    artifact_path = artifact_local_path(manifest, manifest_path, repo_root)
    rel = safe_path_display(artifact_path, repo_root)

    if rel in tracked:
        raise RuntimeError(f"{label} artifact is tracked by git but must remain external: {rel}")

    if not artifact_path.exists():
        if allow_missing:
            return {
                "label": label,
                "artifact_id": artifact_id,
                "artifact_path": rel,
                "present": False,
                "verified": False,
                "reason": "missing_allowed",
            }
        raise RuntimeError(f"{label} artifact missing: {rel}. Provision it outside git, then rerun this verifier.")

    if not artifact_path.is_file():
        raise RuntimeError(f"{label} artifact path is not a file: {rel}")

    actual_size = artifact_path.stat().st_size
    if expected_size is not None and actual_size != expected_size:
        raise RuntimeError(f"{label} artifact size mismatch for {rel}: expected {expected_size}, got {actual_size}")

    actual_sha = sha256_file(artifact_path)
    if actual_sha != expected_sha:
        raise RuntimeError(f"{label} artifact sha256 mismatch for {rel}: expected {expected_sha}, got {actual_sha}")

    return {
        "label": label,
        "artifact_id": artifact_id,
        "artifact_path": rel,
        "present": True,
        "verified": True,
        "size_bytes": actual_size,
        "sha256": actual_sha,
    }


def main() -> int:
    args = parse_args()
    repo_root = args.repo_root.resolve()
    tracked = git_tracked_paths(repo_root)
    results = [
        verify_one("onnx", args.onnx_manifest, repo_root, tracked, args.allow_missing),
        verify_one("native_tokenizer", args.native_tokenizer_manifest, repo_root, tracked, args.allow_missing),
    ]
    print(
        json.dumps(
            {
                "script_version": SCRIPT_VERSION,
                "repo_root": ".",
                "allow_missing": args.allow_missing,
                "results": results,
                "verified_all_present": all(item.get("verified") for item in results),
            },
            ensure_ascii=False,
            indent=2,
            sort_keys=True,
        )
    )
    if args.allow_missing:
        return 0
    return 0 if all(item.get("verified") for item in results) else 1


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except RuntimeError as exc:
        print(f"ERROR: {exc}", file=sys.stderr)
        raise SystemExit(1)
