#!/usr/bin/env python3
"""Verify M040 S02 Docker restart/cache proof artifacts.

The verifier intentionally checks artifact semantics rather than benchmark timing
thresholds: the proof is valid only when restart hooks were configured, Redis L2
restart sections actually ran, effective config was sanitized, health metadata
identifies the ONNX runtime/cache namespace, and artifacts do not contain known
secret or raw-probe leak patterns.
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path
from typing import Any

EXPECTED_CACHE_VERSION = "m040-s02-onnx-restart"
EXPECTED_API_URL = "http://127.0.0.1:18000"
EXPECTED_MODEL = "deepvk/USER-bge-m3"
EXPECTED_DIMENSIONS = 1024

PROHIBITED_PATTERNS = [
    # Secret material and auth-bearing URLs.
    re.compile(r"-----BEGIN [A-Z ]*PRIVATE KEY-----"),
    re.compile(r"\bBearer\s+[A-Za-z0-9._~+/=-]{12,}", re.IGNORECASE),
    re.compile(r"(?i)(api[_-]?key|secret|password|credential|auth[_-]?token)\s*[:=]\s*[^\s,;}]{6,}"),
    re.compile(r"https?://[^\s]+\?(?:[^\s]*\b(?:X-Amz-Signature|Signature|token|access_token)=)", re.IGNORECASE),
    # Raw benchmark/legal probe texts that should not be needed in proof artifacts.
    re.compile(r"Привет, как дела\?"),
    re.compile(r"Москва — столица России"),
    re.compile(r"Искусственный интеллект — это область компьютерных наук"),
    re.compile(r"В современном мире технологии машинного обучения"),
    re.compile(r"Redis L2 persistence diagnostic text"),
    re.compile(r"batch-cache-item-\d+"),
    re.compile(r"research-chunk-\d+"),
]


def fail(message: str) -> None:
    raise AssertionError(message)


def read_text(path: Path) -> str:
    if not path.exists():
        fail(f"missing artifact: {path}")
    if not path.is_file():
        fail(f"artifact is not a file: {path}")
    return path.read_text(encoding="utf-8", errors="replace")


def extract_marked_json(text: str, name: str) -> Any:
    pattern = re.compile(rf"^{name}_JSON_BEGIN\n(.*?)\n{name}_JSON_END$", re.MULTILINE | re.DOTALL)
    match = pattern.search(text)
    if not match:
        fail(f"missing marked JSON block: {name}")
    try:
        return json.loads(match.group(1))
    except json.JSONDecodeError as err:
        fail(f"invalid {name} JSON: {err}")


def extract_effective_config(benchmark_text: str) -> dict[str, Any]:
    heading = "## 0. Effective Configuration Snapshot (sanitized)"
    heading_index = benchmark_text.find(heading)
    if heading_index < 0:
        fail("missing effective configuration snapshot section")
    json_start = benchmark_text.find("{", heading_index)
    if json_start < 0:
        fail("effective configuration snapshot does not contain JSON")
    decoder = json.JSONDecoder()
    try:
        value, _ = decoder.raw_decode(benchmark_text[json_start:])
    except json.JSONDecodeError as err:
        fail(f"effective configuration snapshot JSON is invalid: {err}")
    if not isinstance(value, dict):
        fail("effective configuration snapshot is not an object")
    return value


def section_text(text: str, section_number: int) -> str:
    start_pattern = re.compile(rf"^##\s+{section_number}\.\s+.*$", re.MULTILINE)
    start = start_pattern.search(text)
    if not start:
        fail(f"missing benchmark section {section_number}")
    next_section = re.search(r"^##\s+\d+\.\s+.*$", text[start.end() :], re.MULTILINE)
    end = start.end() + next_section.start() if next_section else len(text)
    return text[start.start() : end]


def require_path(data: dict[str, Any], dotted_path: str) -> Any:
    current: Any = data
    for part in dotted_path.split("."):
        if not isinstance(current, dict) or part not in current:
            fail(f"missing config field: {dotted_path}")
        current = current[part]
    return current


def assert_no_leaks(path: Path, text: str) -> None:
    for pattern in PROHIBITED_PATTERNS:
        if pattern.search(text):
            fail(f"prohibited leak pattern {pattern.pattern!r} found in {path}")


def verify_preflight(preflight: Path, text: str) -> None:
    if "BLOCKER:" in text:
        fail("preflight artifact contains BLOCKER")
    for required in [
        "phase=preflight",
        "phase=health status=ready",
        "phase=smoke status=requesting_embedding",
        "phase=benchmark status=running",
    ]:
        if required not in text:
            fail(f"preflight artifact missing phase marker: {required}")

    health = extract_marked_json(text, "HEALTH")
    runtime = health.get("runtime") if isinstance(health, dict) else None
    if not isinstance(runtime, dict):
        fail("health JSON missing runtime object")
    expected_runtime = {
        "backend": "onnx",
        "model": EXPECTED_MODEL,
        "dimensions": EXPECTED_DIMENSIONS,
        "artifact_id": "user-bge-m3-dense-fp32",
        "provider": "CPUExecutionProvider",
        "artifact_verified": True,
        "tokenizer_verified": True,
        "production_default": False,
    }
    for key, expected in expected_runtime.items():
        if runtime.get(key) != expected:
            fail(f"health.runtime.{key}={runtime.get(key)!r}, expected {expected!r}")
    cache_namespace = runtime.get("cache_namespace")
    if not isinstance(cache_namespace, str) or EXPECTED_CACHE_VERSION not in cache_namespace:
        fail(f"health.runtime.cache_namespace does not include {EXPECTED_CACHE_VERSION!r}: {cache_namespace!r}")
    for forbidden in ["manifest_path", "runtime_library_path", "tokenizer_path", "tei_url", "onnx_runtime_sha256"]:
        if forbidden in runtime:
            fail(f"health.runtime exposes forbidden field: {forbidden}")

    smoke = extract_marked_json(text, "SMOKE")
    if smoke.get("object") != "list" or smoke.get("model") != EXPECTED_MODEL:
        fail("smoke embedding response has unexpected object/model fields")
    data = smoke.get("data")
    if not isinstance(data, list) or not data or data[0].get("object") != "embedding":
        fail("smoke embedding response missing embedding data")
    if not isinstance(data[0].get("embedding"), str) or not data[0].get("embedding"):
        fail("smoke embedding is not a non-empty base64 string")

    container = extract_marked_json(text, "CONTAINER")
    if not isinstance(container, list) or not container:
        fail("container inspect JSON is empty")
    host_config = container[0].get("HostConfig", {})
    port_bindings = host_config.get("PortBindings", {})
    binding = port_bindings.get("18000/tcp", [{}])[0]
    if binding.get("HostIp") != "127.0.0.1" or binding.get("HostPort") != "18000":
        fail(f"API container is not bound to 127.0.0.1:18000: {binding}")


def verify_benchmark(benchmark: Path, text: str) -> None:
    config = extract_effective_config(text)
    if require_path(config, "benchmark.api_url") != EXPECTED_API_URL:
        fail("benchmark effective config has unexpected API URL")
    if require_path(config, "benchmark.model") != EXPECTED_MODEL:
        fail("benchmark effective config has unexpected model")
    if require_path(config, "benchmark.dimensions") != EXPECTED_DIMENSIONS:
        fail("benchmark effective config has unexpected dimensions")
    if require_path(config, "benchmark.input_texts_logged") is not False:
        fail("benchmark effective config must set input_texts_logged=false")
    if require_path(config, "benchmark.api_restart_command_configured") is not True:
        fail("benchmark effective config must set api_restart_command_configured=true")

    env = require_path(config, "environment")
    if not isinstance(env, dict):
        fail("effective environment snapshot is not an object")
    expected_env = {
        "BENCHMARK_API_URL": EXPECTED_API_URL,
        "BENCHMARK_REDIS_HOST": "127.0.0.1",
        "BENCHMARK_REDIS_PORT": "16379",
        "BENCHMARK_RUNTIME_LABEL": "onnx-docker-m040-s02",
        "BENCHMARK_API_RESTART_COMMAND": "docker restart fd-m040-s02-onnx-api",
        "EMBEDDING_BACKEND": "onnx",
        "EMBEDDING_CACHE_VERSION": EXPECTED_CACHE_VERSION,
    }
    for key, expected in expected_env.items():
        if env.get(key) != expected:
            fail(f"environment.{key}={env.get(key)!r}, expected {expected!r}")
    forbidden_env_keys = [key for key in env if re.search(r"(?i)(secret|password|credential|auth|cookie)", key)]
    if forbidden_env_keys:
        fail(f"secret-like keys were included in sanitized environment: {forbidden_env_keys}")

    redaction = require_path(config, "redaction_policy")
    if redaction.get("raw_benchmark_texts_excluded") is not True:
        fail("redaction policy must exclude raw benchmark texts")

    runtime = require_path(config, "runtime")
    if runtime.get("runtime_label") != "onnx-docker-m040-s02":
        fail("runtime label mismatch")
    if runtime.get("onnx_artifact", {}).get("artifact_id") != "user-bge-m3-dense-fp32":
        fail("runtime artifact metadata missing expected ONNX artifact")
    if runtime.get("native_tokenizer_artifact", {}).get("configured") is not True:
        fail("native tokenizer artifact metadata must be configured")
    if runtime.get("onnx_runtime_library", {}).get("configured") is not True:
        fail("ONNX runtime library metadata must be configured")

    section5 = section_text(text, 5)
    if "skipped" in section5.lower():
        fail("Section 5 Redis L2 restart proof was skipped")
    for needle in ["after API restart:", "expectation:", "redis_delta/l2_after_api_restart"]:
        if needle not in section5:
            fail(f"Section 5 missing required evidence: {needle}")

    section6 = section_text(text, 6)
    if "redis_l2_batch_after_api_restart: skipped" in section6.lower():
        fail("Section 6 batch Redis L2 restart proof was skipped")
    for needle in ["redis_l2_batch_after_api_restart:", "redis_delta/batch_l2_after_api_restart"]:
        if needle not in section6:
            fail(f"Section 6 missing required evidence: {needle}")

    summary_required = ["Redis L2 restart:", "Batch L2 p95:"]
    for needle in summary_required:
        if needle not in text:
            fail(f"summary missing {needle}")
    for forbidden_summary in ["Redis L2 restart:   skipped", "Batch L2 p95:       skipped"]:
        if forbidden_summary in text:
            fail(f"summary reports skipped proof: {forbidden_summary}")


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--preflight", type=Path, default=Path("benchmark-results/fd-m040-s02-onnx-docker-preflight.txt"))
    parser.add_argument("--benchmark", type=Path, default=Path("benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt"))
    args = parser.parse_args(argv)

    preflight_text = read_text(args.preflight)
    benchmark_text = read_text(args.benchmark)
    assert_no_leaks(args.preflight, preflight_text)
    assert_no_leaks(args.benchmark, benchmark_text)
    verify_preflight(args.preflight, preflight_text)
    verify_benchmark(args.benchmark, benchmark_text)
    print("M040 S02 artifact verification: PASS")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main(sys.argv[1:]))
    except AssertionError as err:
        print(f"M040 S02 artifact verification: FAIL: {err}", file=sys.stderr)
        raise SystemExit(1)
