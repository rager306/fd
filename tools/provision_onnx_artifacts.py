#!/usr/bin/env python3
"""Provision fd opt-in ONNX artifacts from explicit sources.

The helper intentionally has no fake default artifact URLs. In dry-run mode it
prints the required destinations and reports missing sources as blockers. In
normal mode every supplied source is copied/downloaded, then size and sha256 are
verified before the artifact is accepted.
"""

from __future__ import annotations

import argparse
import hashlib
import ipaddress
import json
import os
from pathlib import Path
import shutil
import socket
import sys
import tarfile
import tempfile
from typing import Any
from urllib.error import HTTPError
from urllib.parse import urlparse
from urllib.request import HTTPRedirectHandler, build_opener

SCRIPT_VERSION = 2
DEFAULT_MAX_DOWNLOAD_BYTES = 2 * 1024 * 1024 * 1024
DEFAULT_MAX_ARCHIVE_MEMBER_BYTES = 256 * 1024 * 1024


class NoRedirectHandler(HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):  # type: ignore[no-untyped-def]
        raise HTTPError(req.full_url, code, f"redirects are not allowed for artifact downloads: {source_display(newurl)}", headers, fp)


URL_OPENER = build_opener(NoRedirectHandler)


def csv_values(value: str | None) -> list[str]:
    if not value:
        return []
    return [item.strip().lower() for item in value.split(",") if item.strip()]


def env_bool(key: str, default: bool = False) -> bool:
    value = os.getenv(key)
    if value is None:
        return default
    return value.lower() in {"1", "true", "yes", "on"}


def positive_int_env(key: str, default: int) -> int:
    value = os.getenv(key)
    if value is None or value == "":
        return default
    try:
        parsed = int(value)
    except ValueError as exc:
        raise RuntimeError(f"{key} must be a positive integer") from exc
    if parsed <= 0:
        raise RuntimeError(f"{key} must be a positive integer")
    return parsed


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Provision fd ONNX/native artifacts with checksum verification.")
    parser.add_argument("--onnx-manifest", type=Path, required=True)
    parser.add_argument("--native-tokenizer-manifest", type=Path, required=True)
    parser.add_argument("--repo-root", type=Path, default=Path("."))
    parser.add_argument("--onnx-source", help="Explicit file path or URL for the ONNX model artifact.")
    parser.add_argument("--native-tokenizer-source", help="Explicit file path or URL for libtokenizers.a or an archive containing it.")
    parser.add_argument("--tokenizer-json-source", help="Optional explicit file path or URL for tokenizer.json.")
    parser.add_argument("--tokenizer-json-dest", type=Path, default=Path("tei-models/deepvk--USER-bge-m3/tokenizer.json"))
    parser.add_argument("--onnx-runtime-source", help="Optional explicit file path or URL for libonnxruntime.so.*.")
    parser.add_argument("--onnx-runtime-dest", type=Path, default=Path(".gsd/runtime/onnxruntime/libonnxruntime.so.1.26.0"))
    parser.add_argument("--onnx-runtime-sha256", help="Required when --onnx-runtime-source is provided.")
    parser.add_argument("--allowed-artifact-host", action="append", default=csv_values(os.getenv("FD_ONNX_ALLOWED_ARTIFACT_HOSTS")), help="Allowed HTTPS artifact host. May be repeated. If omitted, any public HTTPS host is allowed.")
    parser.add_argument("--allow-private-artifact-hosts", action="store_true", default=env_bool("FD_ONNX_ALLOW_PRIVATE_ARTIFACT_HOSTS"), help="Allow artifact URLs resolving to private, loopback, link-local, or reserved addresses. Intended only for trusted local testing.")
    parser.add_argument("--max-download-bytes", type=int, default=positive_int_env("FD_ONNX_MAX_DOWNLOAD_BYTES", DEFAULT_MAX_DOWNLOAD_BYTES), help="Maximum bytes to stream from a remote artifact URL before failing.")
    parser.add_argument("--max-archive-member-bytes", type=int, default=positive_int_env("FD_ONNX_MAX_ARCHIVE_MEMBER_BYTES", DEFAULT_MAX_ARCHIVE_MEMBER_BYTES), help="Maximum native tokenizer archive member size before failing.")
    parser.add_argument("--dry-run", action="store_true", help="Report the provisioning plan without copying/downloading.")
    return parser.parse_args()


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


def repo_path(repo_root: Path, value: str | Path) -> Path:
    path = Path(value)
    return path if path.is_absolute() else repo_root / path


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def artifact_expectation(manifest: dict[str, Any], manifest_path: Path, repo_root: Path) -> tuple[str, Path, int | None, str]:
    artifact = manifest.get("artifact")
    if not isinstance(artifact, dict):
        raise RuntimeError(f"manifest missing artifact object: {manifest_path}")
    artifact_id = manifest.get("artifact_id")
    local_path = artifact.get("local_path")
    sha = artifact.get("sha256")
    size = artifact.get("size_bytes")
    if not isinstance(artifact_id, str) or not artifact_id:
        raise RuntimeError(f"manifest missing artifact_id: {manifest_path}")
    if not isinstance(local_path, str) or not local_path:
        raise RuntimeError(f"manifest missing artifact.local_path: {manifest_path}")
    if not isinstance(sha, str) or len(sha) != 64:
        raise RuntimeError(f"manifest artifact.sha256 must be a 64-character hex string: {manifest_path}")
    if size is not None and (not isinstance(size, int) or size <= 0):
        raise RuntimeError(f"manifest artifact.size_bytes must be a positive integer: {manifest_path}")
    return artifact_id, repo_path(repo_root, local_path), size, sha


def tokenizer_json_expectation(onnx_manifest: dict[str, Any]) -> tuple[int | None, str | None]:
    source_files = ((onnx_manifest.get("model") or {}).get("source_files") or {}) if isinstance(onnx_manifest.get("model"), dict) else {}
    tokenizer = source_files.get("tokenizer.json") if isinstance(source_files, dict) else None
    if not isinstance(tokenizer, dict):
        return None, None
    size = tokenizer.get("size_bytes")
    sha = tokenizer.get("sha256")
    return (size if isinstance(size, int) else None), (sha if isinstance(sha, str) else None)


def source_display(source: str | None) -> str | None:
    if not source:
        return None
    parsed = urlparse(source)
    if parsed.scheme in {"http", "https"}:
        host = parsed.netloc
        name = Path(parsed.path).name
        return f"{parsed.scheme}://{host}/.../{name}" if name else f"{parsed.scheme}://{host}/..."
    return source


def is_url(source: str) -> bool:
    return urlparse(source).scheme in {"http", "https"}


def validate_remote_source(source: str, allowed_hosts: set[str], allow_private_hosts: bool) -> None:
    parsed = urlparse(source)
    if parsed.scheme != "https":
        raise RuntimeError(f"artifact URL must use https: {source_display(source)}")
    host = parsed.hostname
    if not host:
        raise RuntimeError(f"artifact URL is missing host: {source_display(source)}")
    normalized_host = host.lower()
    if allowed_hosts and normalized_host not in allowed_hosts:
        raise RuntimeError(f"artifact URL host is not allowed: {normalized_host}")
    if allow_private_hosts:
        return
    try:
        infos = socket.getaddrinfo(host, parsed.port or 443, type=socket.SOCK_STREAM)
    except socket.gaierror as exc:
        raise RuntimeError(f"artifact URL host could not be resolved: {normalized_host}") from exc
    addresses = {item[4][0] for item in infos}
    for address in addresses:
        ip = ipaddress.ip_address(address)
        if ip.is_private or ip.is_loopback or ip.is_link_local or ip.is_reserved or ip.is_multicast or ip.is_unspecified:
            raise RuntimeError(f"artifact URL host resolves to a disallowed address: {normalized_host}")


def copy_bounded(source_handle: Any, dest_handle: Any, max_bytes: int, label: str) -> int:
    total = 0
    while True:
        chunk = source_handle.read(1024 * 1024)
        if not chunk:
            return total
        total += len(chunk)
        if total > max_bytes:
            raise RuntimeError(f"{label} exceeds maximum allowed size of {max_bytes} bytes")
        dest_handle.write(chunk)


def fetch_source(source: str, temp_dir: Path, *, allowed_hosts: set[str], allow_private_hosts: bool, max_download_bytes: int) -> Path:
    if is_url(source):
        validate_remote_source(source, allowed_hosts, allow_private_hosts)
        name = Path(urlparse(source).path).name or "artifact.bin"
        dest = temp_dir / name
        with URL_OPENER.open(source, timeout=300) as response, dest.open("wb") as handle:  # noqa: S310 - policy validated explicit artifact URL.
            length = response.headers.get("Content-Length")
            if length:
                try:
                    content_length = int(length)
                except ValueError as exc:
                    raise RuntimeError(f"artifact URL returned invalid Content-Length: {source_display(source)}") from exc
                if content_length > max_download_bytes:
                    raise RuntimeError(f"artifact URL Content-Length exceeds maximum allowed size: {source_display(source)}")
            copy_bounded(response, handle, max_download_bytes, f"artifact download {source_display(source)}")
        return dest
    path = Path(source).expanduser()
    if not path.exists() or not path.is_file():
        raise RuntimeError(f"source file missing: {source_display(source)}")
    return path


def materialize_source(
    source: str,
    destination: Path,
    *,
    archive_member: str | None = None,
    expected_size: int | None = None,
    max_archive_member_bytes: int = DEFAULT_MAX_ARCHIVE_MEMBER_BYTES,
    allowed_hosts: set[str] | None = None,
    allow_private_hosts: bool = False,
    max_download_bytes: int = DEFAULT_MAX_DOWNLOAD_BYTES,
) -> None:
    destination.parent.mkdir(parents=True, exist_ok=True)
    with tempfile.TemporaryDirectory(prefix="fd-onnx-provision-") as tmp:
        fetched = fetch_source(source, Path(tmp), allowed_hosts=allowed_hosts or set(), allow_private_hosts=allow_private_hosts, max_download_bytes=max_download_bytes)
        if archive_member:
            with tarfile.open(fetched) as archive:
                member = next((item for item in archive.getmembers() if Path(item.name).name == archive_member), None)
                if member is None:
                    raise RuntimeError(f"archive does not contain {archive_member}: {source_display(source)}")
                if not member.isfile():
                    raise RuntimeError(f"archive member is not a regular file: {archive_member}")
                max_member_bytes = expected_size if expected_size is not None else max_archive_member_bytes
                if member.size > max_member_bytes:
                    raise RuntimeError(f"archive member {archive_member} size {member.size} exceeds maximum allowed size {max_member_bytes}")
                extracted = archive.extractfile(member)
                if extracted is None:
                    raise RuntimeError(f"archive member is not readable: {archive_member}")
                with destination.open("wb") as handle:
                    copy_bounded(extracted, handle, max_member_bytes, f"archive member {archive_member}")
        else:
            shutil.copy2(fetched, destination)


def verify_destination(label: str, destination: Path, expected_size: int | None, expected_sha: str | None) -> dict[str, Any]:
    if not destination.exists() or not destination.is_file():
        raise RuntimeError(f"{label} destination missing after provisioning: {destination}")
    actual_size = destination.stat().st_size
    if expected_size is not None and actual_size != expected_size:
        raise RuntimeError(f"{label} size mismatch: expected {expected_size}, got {actual_size}")
    actual_sha = sha256_file(destination)
    if expected_sha and actual_sha != expected_sha:
        raise RuntimeError(f"{label} sha256 mismatch: expected {expected_sha}, got {actual_sha}")
    return {"present": True, "verified": bool(expected_sha), "size_bytes": actual_size, "sha256": actual_sha}


def plan_item(label: str, destination: Path, source: str | None, required: bool, expected_sha: str | None) -> dict[str, Any]:
    return {
        "label": label,
        "destination": str(destination),
        "source_configured": bool(source),
        "source": source_display(source),
        "required": required,
        "expected_sha256": expected_sha,
        "status": "planned" if source else ("blocked_missing_source" if required else "optional_missing_source"),
    }


def main() -> int:
    args = parse_args()
    repo_root = args.repo_root.resolve()
    onnx_manifest = load_manifest(args.onnx_manifest)
    native_manifest = load_manifest(args.native_tokenizer_manifest)
    onnx_id, onnx_dest, onnx_size, onnx_sha = artifact_expectation(onnx_manifest, args.onnx_manifest, repo_root)
    native_id, native_dest, native_size, native_sha = artifact_expectation(native_manifest, args.native_tokenizer_manifest, repo_root)
    tokenizer_size, tokenizer_sha = tokenizer_json_expectation(onnx_manifest)
    tokenizer_dest = repo_path(repo_root, args.tokenizer_json_dest)
    ort_dest = repo_path(repo_root, args.onnx_runtime_dest)

    plan = [
        plan_item("onnx", onnx_dest, args.onnx_source, True, onnx_sha),
        plan_item("native_tokenizer", native_dest, args.native_tokenizer_source, True, native_sha),
        plan_item("tokenizer_json", tokenizer_dest, args.tokenizer_json_source, False, tokenizer_sha),
        plan_item("onnx_runtime", ort_dest, args.onnx_runtime_source, False, args.onnx_runtime_sha256),
    ]

    if args.dry_run:
        print(json.dumps({"script_version": SCRIPT_VERSION, "dry_run": True, "repo_root": str(repo_root), "plan": plan}, indent=2, sort_keys=True))
        return 0

    missing_required = [item["label"] for item in plan if item["required"] and not item["source_configured"]]
    if missing_required:
        raise RuntimeError("missing required artifact sources: " + ", ".join(missing_required))
    if args.onnx_runtime_source and not args.onnx_runtime_sha256:
        raise RuntimeError("--onnx-runtime-sha256 is required when --onnx-runtime-source is provided")
    allowed_hosts = {host.lower() for host in args.allowed_artifact_host}
    if args.max_download_bytes <= 0:
        raise RuntimeError("--max-download-bytes must be a positive integer")
    if args.max_archive_member_bytes <= 0:
        raise RuntimeError("--max-archive-member-bytes must be a positive integer")

    materialize_source(args.onnx_source, onnx_dest, expected_size=onnx_size, allowed_hosts=allowed_hosts, allow_private_hosts=args.allow_private_artifact_hosts, max_download_bytes=args.max_download_bytes)
    materialize_source(args.native_tokenizer_source, native_dest, archive_member="libtokenizers.a", expected_size=native_size, max_archive_member_bytes=args.max_archive_member_bytes, allowed_hosts=allowed_hosts, allow_private_hosts=args.allow_private_artifact_hosts, max_download_bytes=args.max_download_bytes)
    results = [
        {"label": "onnx", "artifact_id": onnx_id, **verify_destination("onnx", onnx_dest, onnx_size, onnx_sha)},
        {"label": "native_tokenizer", "artifact_id": native_id, **verify_destination("native_tokenizer", native_dest, native_size, native_sha)},
    ]
    if args.tokenizer_json_source:
        materialize_source(args.tokenizer_json_source, tokenizer_dest, expected_size=tokenizer_size, allowed_hosts=allowed_hosts, allow_private_hosts=args.allow_private_artifact_hosts, max_download_bytes=args.max_download_bytes)
        results.append({"label": "tokenizer_json", **verify_destination("tokenizer_json", tokenizer_dest, tokenizer_size, tokenizer_sha)})
    if args.onnx_runtime_source:
        materialize_source(args.onnx_runtime_source, ort_dest, allowed_hosts=allowed_hosts, allow_private_hosts=args.allow_private_artifact_hosts, max_download_bytes=args.max_download_bytes)
        results.append({"label": "onnx_runtime", **verify_destination("onnx_runtime", ort_dest, None, args.onnx_runtime_sha256)})

    print(json.dumps({"script_version": SCRIPT_VERSION, "dry_run": False, "repo_root": str(repo_root), "results": results}, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except RuntimeError as exc:
        print(f"ERROR: {exc}", file=sys.stderr)
        raise SystemExit(1)
