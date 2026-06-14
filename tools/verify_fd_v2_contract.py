#!/usr/bin/env python3
"""fd v2 contract verifier for M041.

Runs 45 black-box checks against a running fd service and writes a markdown
artifact with pass/fail details. Uses only the Python standard library.
"""

from __future__ import annotations

import base64
import json
import statistics
import sys
import time
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

BASE_URL = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8000"
OUT = Path("benchmark-results/fd-v2-validation-m041.md")
MODEL = "deepvk/USER-bge-m3"


@dataclass
class HTTPResult:
    status: int
    body: bytes
    headers: dict[str, str]
    latency_ms: float
    error: str = ""

    def json(self) -> Any:
        return json.loads(self.body.decode("utf-8"))


@dataclass
class Check:
    id: str
    description: str
    passed: bool
    detail: str = ""
    latency_ms: float | None = None


RESULTS: list[Check] = []


def request(method: str, path: str, body: Any = None, headers: dict[str, str] | None = None, timeout: int = 30) -> HTTPResult:
    data = None
    merged_headers = dict(headers or {})
    if body is not None:
        data = json.dumps(body).encode("utf-8")
        merged_headers.setdefault("Content-Type", "application/json")
    req = urllib.request.Request(BASE_URL + path, data=data, headers=merged_headers, method=method)
    started = time.perf_counter()
    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            payload = resp.read()
            return HTTPResult(resp.status, payload, dict(resp.headers.items()), (time.perf_counter() - started) * 1000)
    except urllib.error.HTTPError as exc:
        payload = exc.read()
        return HTTPResult(exc.code, payload, dict(exc.headers.items()), (time.perf_counter() - started) * 1000, f"HTTPError:{exc.code}")
    except Exception as exc:  # pragma: no cover - surfaced in artifact
        return HTTPResult(0, b"", {}, (time.perf_counter() - started) * 1000, repr(exc))


def add(check_id: str, description: str, passed: bool, detail: str = "", latency_ms: float | None = None) -> None:
    RESULTS.append(Check(check_id, description, passed, detail, latency_ms))


def header(result: HTTPResult, name: str) -> str:
    name_lower = name.lower()
    for key, value in result.headers.items():
        if key.lower() == name_lower:
            return value
    return ""


def expect_status(check_id: str, description: str, result: HTTPResult, status: int) -> None:
    add(check_id, description, result.status == status, f"status={result.status}, want={status}, body={result.body[:160]!r}", result.latency_ms)


def expect_error(check_id: str, description: str, result: HTTPResult, status: int, code: str) -> None:
    try:
        payload = result.json()
        got_code = payload.get("error", {}).get("code")
    except Exception as exc:
        add(check_id, description, False, f"unparseable error body: {exc}; status={result.status}", result.latency_ms)
        return
    add(check_id, description, result.status == status and got_code == code, f"status={result.status}, code={got_code}, want={status}/{code}", result.latency_ms)


def embedding_body(input_value: Any, **extra: Any) -> dict[str, Any]:
    body = {"model": MODEL, "input": input_value}
    body.update(extra)
    return body


def percentile(values: list[float], pct: int) -> float:
    ordered = sorted(values)
    idx = int((len(ordered) - 1) * pct / 100)
    return ordered[idx]


def run_checks() -> dict[str, Any]:
    # Endpoint existence and shape.
    live = request("GET", "/live")
    expect_status("T001", "GET /live returns 200", live, 200)
    ready = request("GET", "/ready")
    expect_status("T002", "GET /ready returns 200 after warmup", ready, 200)
    health = request("GET", "/health")
    add("T003", "GET /health exposes runtime", health.status == 200 and "runtime" in health.json(), f"status={health.status}", health.latency_ms)
    healthcheck = request("GET", "/v1/healthcheck")
    add("T004", "GET /v1/healthcheck exposes runtime", healthcheck.status == 200 and "runtime" in healthcheck.json(), f"status={healthcheck.status}", healthcheck.latency_ms)
    version = request("GET", "/version")
    add("T005", "GET /version exposes service/model", version.status == 200 and "service" in version.json() and "model" in version.json(), f"status={version.status}", version.latency_ms)
    info = request("GET", "/info")
    add("T006", "GET /info exposes service model metadata", info.status == 200 and "models" in info.json(), f"status={info.status}", info.latency_ms)
    metrics = request("GET", "/metrics")
    add("T007", "GET /metrics exposes Prometheus text", metrics.status == 200 and b"fd_requests_total" in metrics.body, f"status={metrics.status}", metrics.latency_ms)
    warmup_get = request("GET", "/warmup")
    expect_status("T008", "GET /warmup returns 200", warmup_get, 200)
    warmup_post = request("POST", "/warmup")
    add("T009", "POST /warmup returns 200 or 202", warmup_post.status in (200, 202), f"status={warmup_post.status}", warmup_post.latency_ms)
    openapi = request("GET", "/openapi.json")
    add("T010", "GET /openapi.json returns OpenAPI 3.1", openapi.status == 200 and openapi.json().get("openapi") == "3.1.0", f"status={openapi.status}", openapi.latency_ms)
    docs = request("GET", "/docs")
    add("T011", "GET /docs returns Swagger UI HTML", docs.status == 200 and b"swagger-ui" in docs.body, f"status={docs.status}", docs.latency_ms)
    traces = request("GET", "/v1/traces")
    add("T012", "GET /v1/traces returns JSON array", traces.status == 200 and isinstance(traces.json(), list), f"status={traces.status}", traces.latency_ms)

    # Embeddings happy path and headers.
    unique = f"contract-{time.time_ns()}"
    emb = request("POST", "/v1/embeddings", embedding_body(unique))
    emb_payload = emb.json() if emb.status == 200 else {}
    add("T013", "POST /v1/embeddings string input returns one embedding", emb.status == 200 and len(emb_payload.get("data", [])) == 1, f"status={emb.status}", emb.latency_ms)
    add("T014", "Embedding response includes authoritative model", emb.status == 200 and emb_payload.get("model") == MODEL, f"model={emb_payload.get('model')}", emb.latency_ms)
    add("T015", "Embedding response includes Server and request id headers", header(emb, "Server").startswith("fd/") and header(emb, "X-Request-Id") != "", str({k: header(emb, k) for k in ['Server','X-Request-Id']}), emb.latency_ms)
    add("T016", "Embedding response includes model/dim/cache headers", all(header(emb, h) for h in ["X-Model-Id", "X-Dimensions", "X-Cache"]), str({k: header(emb, k) for k in ['X-Model-Id','X-Dimensions','X-Cache']}), emb.latency_ms)
    arr = request("POST", "/v1/embeddings", embedding_body([unique + "-a", unique + "-b"]))
    add("T017", "Array input returns two embeddings", arr.status == 200 and len(arr.json().get("data", [])) == 2, f"status={arr.status}", arr.latency_ms)
    dims = request("POST", "/v1/embeddings", embedding_body(unique + "-512", dimensions=512))
    add("T018", "dimensions=512 returns 512 dimensions", dims.status == 200 and dims.json()["data"][0]["dimensions"] == 512 and header(dims, "X-Dimensions") == "512", f"status={dims.status}", dims.latency_ms)
    b64 = request("POST", "/v1/embeddings", embedding_body(unique + "-b64", encoding_format="base64"))
    is_b64 = False
    if b64.status == 200:
        value = b64.json()["data"][0]["embedding"]
        is_b64 = isinstance(value, str) and len(base64.b64decode(value)) > 0
    add("T019", "encoding_format=base64 returns base64 string", is_b64, f"status={b64.status}", b64.latency_ms)
    priority = request("POST", "/v1/embeddings", embedding_body(unique + "-priority", priority="high"))
    expect_status("T020", "priority=high is accepted", priority, 200)
    user = request("POST", "/v1/embeddings", embedding_body(unique + "-user", user="caller-1"))
    expect_status("T021", "user field is accepted", user, 200)

    # Cache / ETag / CORS.
    cache_key = unique + "-cache"
    miss = request("POST", "/v1/embeddings", embedding_body(cache_key))
    hit = request("POST", "/v1/embeddings", embedding_body(cache_key))
    add("T022", "First unique embedding request is cache MISS", miss.status == 200 and header(miss, "X-Cache") == "MISS", f"X-Cache={header(miss, 'X-Cache')}", miss.latency_ms)
    add("T023", "Repeated embedding request is cache HIT", hit.status == 200 and header(hit, "X-Cache") == "HIT", f"X-Cache={header(hit, 'X-Cache')}", hit.latency_ms)
    add("T024", "Embedding response includes ETag and Cache-Control", header(hit, "ETag") != "" and header(hit, "Cache-Control") == "public, max-age=86400", str({k: header(hit, k) for k in ['ETag','Cache-Control']}), hit.latency_ms)
    not_modified = request("POST", "/v1/embeddings", embedding_body(cache_key), headers={"If-None-Match": header(hit, "ETag")})
    expect_status("T025", "Embedding If-None-Match returns 304", not_modified, 304)
    info_again = request("GET", "/info", headers={"If-None-Match": header(info, "ETag")})
    expect_status("T026", "Info If-None-Match returns 304", info_again, 304)
    cors = request("OPTIONS", "/v1/embeddings", headers={"Origin": "https://example.test", "Access-Control-Request-Method": "POST"})
    add("T027", "CORS preflight returns 204 with headers", cors.status == 204 and header(cors, "Access-Control-Allow-Methods") == "GET,POST,OPTIONS", f"status={cors.status}", cors.latency_ms)

    # Error paths.
    expect_error("T028", "Malformed JSON returns invalid_json", request("POST", "/v1/embeddings", headers={"Content-Type": "application/json"}, body=None), 400, "invalid_json")
    expect_error("T029", "Missing input returns input_required", request("POST", "/v1/embeddings", {"model": MODEL}), 400, "input_required")
    expect_error("T030", "Empty input returns input_required", request("POST", "/v1/embeddings", embedding_body([])), 400, "input_required")
    expect_error("T031", "Oversized input batch returns batch_too_large", request("POST", "/v1/embeddings", embedding_body(["x"] * 129)), 413, "batch_too_large")
    expect_error("T032", "Invalid dimensions returns dimensions_invalid", request("POST", "/v1/embeddings", embedding_body("x", dimensions=256)), 400, "dimensions_invalid")
    expect_error("T033", "Invalid encoding_format returns encoding_format_invalid", request("POST", "/v1/embeddings", embedding_body("x", encoding_format="hex")), 400, "encoding_format_invalid")
    expect_error("T034", "Invalid priority returns priority_invalid", request("POST", "/v1/embeddings", embedding_body("x", priority="urgent")), 400, "priority_invalid")
    expect_error("T035", "Too long input returns input_too_long", request("POST", "/v1/embeddings", embedding_body("x" * 2049)), 413, "input_too_long")
    expect_error("T036", "GET /v1/embeddings returns method_not_allowed envelope", request("GET", "/v1/embeddings"), 405, "method_not_allowed")
    expect_error("T037", "Unknown route returns not_found", request("GET", "/nope"), 404, "not_found")

    # Batch endpoints.
    v1batch = request("POST", "/v1/batch", {"batches": [["a", "b", "c", "d"], ["e", "f", "g", "h"]]})
    add("T038", "/v1/batch returns 2x4 embeddings", v1batch.status == 200 and [len(b) for b in v1batch.json().get("batches", [])] == [4, 4], f"status={v1batch.status}", v1batch.latency_ms)
    expect_error("T039", "/v1/batch empty batches returns input_required", request("POST", "/v1/batch", {"batches": []}), 400, "input_required")
    expect_error("T040", "/v1/batch oversized inner batch returns batch_too_large", request("POST", "/v1/batch", {"batches": [["x"] * 33]}), 413, "batch_too_large")
    legacy = request("POST", "/embeddings/batch", {"inputs": ["legacy"]})
    add("T041", "Legacy /embeddings/batch returns base64 embeddings", legacy.status == 200 and isinstance(legacy.json().get("embeddings", [None])[0], str), f"status={legacy.status}", legacy.latency_ms)

    # Trace/openapi and cache-hot performance checks.
    traces_after = request("GET", "/v1/traces")
    add("T042", "Traces contain recorded request entries", traces_after.status == 200 and len(traces_after.json()) >= 5, f"count={len(traces_after.json()) if traces_after.status == 200 else 0}", traces_after.latency_ms)
    spec = openapi.json() if openapi.status == 200 else {}
    paths = spec.get("paths", {})
    add("T043", "OpenAPI includes /v1/embeddings, /v1/batch, /v1/traces", all(p in paths for p in ["/v1/embeddings", "/v1/batch", "/v1/traces"]), f"paths={len(paths)}", openapi.latency_ms)

    perf_key = unique + "-perf"
    prewarm = request("POST", "/v1/embeddings", embedding_body(perf_key))
    latencies = [request("POST", "/v1/embeddings", embedding_body(perf_key)).latency_ms for _ in range(20)]
    p95 = percentile(latencies, 95)
    add("T044", "Cache-hot single input p95 < 50ms", prewarm.status == 200 and p95 < 50, f"p95={p95:.3f}ms", p95)
    batch_key = [f"{unique}-perf-batch-{i}" for i in range(10)]
    prewarm_batch = request("POST", "/v1/embeddings", embedding_body(batch_key))
    batch_latencies = [request("POST", "/v1/embeddings", embedding_body(batch_key)).latency_ms for _ in range(10)]
    batch_p95 = percentile(batch_latencies, 95)
    add("T045", "Cache-hot 10 input p95 < 200ms", prewarm_batch.status == 200 and batch_p95 < 200, f"p95={batch_p95:.3f}ms", batch_p95)

    return {"single_p95_ms": p95, "batch10_p95_ms": batch_p95}


def write_artifact(perf: dict[str, Any]) -> None:
    OUT.parent.mkdir(parents=True, exist_ok=True)
    passed = sum(1 for r in RESULTS if r.passed)
    failed = len(RESULTS) - passed
    lines = [
        "# fd v2 contract validation — M041",
        "",
        f"Base URL: `{BASE_URL}`",
        f"Summary: {passed}/{len(RESULTS)} passed, {failed} failed",
        f"Perf: single p95={perf.get('single_p95_ms'):.3f}ms, batch10 p95={perf.get('batch10_p95_ms'):.3f}ms",
        "",
        "| ID | Result | Description | Detail | Latency ms |",
        "|---|---|---|---|---:|",
    ]
    for r in RESULTS:
        verdict = "PASS" if r.passed else "FAIL"
        latency = "" if r.latency_ms is None else f"{r.latency_ms:.3f}"
        detail = r.detail.replace("\n", " ").replace("|", "\\|")
        lines.append(f"| {r.id} | {verdict} | {r.description} | {detail} | {latency} |")
    lines += ["", "## JSON", "", "```json", json.dumps([r.__dict__ for r in RESULTS], indent=2, default=str), "```", ""]
    OUT.write_text("\n".join(lines), encoding="utf-8")


def main() -> int:
    perf = run_checks()
    write_artifact(perf)
    failed = [r for r in RESULTS if not r.passed]
    print(f"wrote {OUT}")
    print(f"passed={len(RESULTS)-len(failed)} failed={len(failed)} total={len(RESULTS)}")
    if failed:
        for item in failed:
            print(f"FAIL {item.id}: {item.description} — {item.detail}")
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
