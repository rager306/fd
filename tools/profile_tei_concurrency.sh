#!/usr/bin/env bash
set -euo pipefail

OUT=${1:-benchmark-results/te-concurrency-profile-m042-s01.md}
TEI_BASE=${TEI_BASE:-http://localhost:30080}
MODEL=${MODEL:-deepvk/USER-bge-m3}
CONTAINER=${TEI_CONTAINER:-fd_tei}
COMPOSE_SERVICE=${TEI_COMPOSE_SERVICE:-tei}

mkdir -p "$(dirname "$OUT")"

python3 - "$OUT" "$TEI_BASE" "$MODEL" "$CONTAINER" "$COMPOSE_SERVICE" <<'PY'
from __future__ import annotations

import concurrent.futures
import datetime as dt
import json
import re
import statistics
import subprocess
import sys
import time
import urllib.error
import urllib.request
from dataclasses import dataclass
from pathlib import Path

out_path = Path(sys.argv[1])
tei_base = sys.argv[2].rstrip("/")
model = sys.argv[3]
container = sys.argv[4]
compose_service = sys.argv[5]

ANSI = re.compile(r"\x1b\[[0-9;]*m")
TIMING_RE = re.compile(
    r'total_time="([^"]+)".*?tokenization_time="([^"]+)".*?queue_time="([^"]+)".*?inference_time="([^"]+)"'
)


def utc_now() -> str:
    return dt.datetime.now(dt.timezone.utc).isoformat(timespec="seconds")


def to_ms(value: str) -> float:
    if value.endswith("µs"):
        return float(value[:-2]) / 1000.0
    if value.endswith("ms"):
        return float(value[:-2])
    if value.endswith("s"):
        return float(value[:-1]) * 1000.0
    return float("nan")


def percentile(values: list[float], pct: float) -> float:
    if not values:
        return float("nan")
    vals = sorted(values)
    idx = int(len(vals) * pct / 100.0) - 1
    idx = max(0, min(idx, len(vals) - 1))
    return vals[idx]


def fmt_ms(value: float) -> str:
    if value != value:
        return "n/a"
    return f"{value:.3f}ms"


def stats(values: list[float]) -> dict[str, float]:
    if not values:
        return {"min": float("nan"), "p50": float("nan"), "p95": float("nan"), "max": float("nan")}
    return {
        "min": min(values),
        "p50": statistics.median(values),
        "p95": percentile(values, 95),
        "max": max(values),
    }


@dataclass
class ClientResult:
    scenario: str
    request_id: int
    batch_size: int
    status: int
    wall_ms: float
    error: str = ""


@dataclass
class ScenarioResult:
    name: str
    description: str
    batch_size: int
    request_count: int
    concurrency: int
    started_at: str
    clients: list[ClientResult]
    server_rows: list[tuple[float, float, float, float]]


def request_embedding(batch_size: int, request_id: int, scenario: str) -> ClientResult:
    inputs = [f"m042 {scenario} request {request_id} item {i}" for i in range(batch_size)]
    payload = json.dumps({"model": model, "input": inputs}).encode("utf-8")
    req = urllib.request.Request(
        f"{tei_base}/v1/embeddings",
        data=payload,
        headers={"content-type": "application/json"},
        method="POST",
    )
    start = time.perf_counter()
    try:
        with urllib.request.urlopen(req, timeout=120) as resp:
            resp.read(200)
            status = resp.status
            error = ""
    except urllib.error.HTTPError as exc:
        status = exc.code
        error = exc.read(500).decode("utf-8", "replace")
    except Exception as exc:  # noqa: BLE001 - profiler records diagnostics
        status = 0
        error = f"{type(exc).__name__}: {exc}"
    return ClientResult(scenario, request_id, batch_size, status, (time.perf_counter() - start) * 1000.0, error)


def parse_logs_since(started_at: str) -> list[tuple[float, float, float, float]]:
    logs = subprocess.check_output(
        ["docker", "logs", "--since", started_at, container],
        stderr=subprocess.STDOUT,
        text=True,
        errors="replace",
    )
    logs = ANSI.sub("", logs)
    return [tuple(to_ms(part) for part in match.groups()) for match in TIMING_RE.finditer(logs)]


def wait_for_tei(timeout_s: int = 240) -> None:
    deadline = time.time() + timeout_s
    last_error = ""
    while time.time() < deadline:
        try:
            with urllib.request.urlopen(f"{tei_base}/health", timeout=3) as resp:
                if resp.status == 200:
                    return
        except Exception as exc:  # noqa: BLE001 - readiness diagnostics
            last_error = f"{type(exc).__name__}: {exc}"
        time.sleep(2)
    raise RuntimeError(f"TEI did not become healthy within {timeout_s}s; last_error={last_error}")


def run_scenario(name: str, description: str, batch_size: int, request_count: int, concurrency: int) -> ScenarioResult:
    started_at = utc_now()
    print(f"SCENARIO {name} started_at={started_at} batch={batch_size} requests={request_count} concurrency={concurrency}", flush=True)
    clients: list[ClientResult] = []
    with concurrent.futures.ThreadPoolExecutor(max_workers=concurrency) as executor:
        futures = [executor.submit(request_embedding, batch_size, i + 1, name) for i in range(request_count)]
        for future in concurrent.futures.as_completed(futures):
            clients.append(future.result())
    clients.sort(key=lambda row: row.request_id)
    time.sleep(0.75)
    server_rows = parse_logs_since(started_at)
    # If a health/warmup request or previous line leaked into the timestamp boundary, keep the most recent rows matching client count.
    if len(server_rows) > request_count:
        server_rows = server_rows[-request_count:]
    return ScenarioResult(name, description, batch_size, request_count, concurrency, started_at, clients, server_rows)


run_started = utc_now()
wait_for_tei()
info = json.loads(urllib.request.urlopen(f"{tei_base}/info", timeout=10).read().decode("utf-8"))

results: list[ScenarioResult] = []
results.append(run_scenario("sequential_batch32_control", "Four sequential direct TEI batch=32 requests as the control for the parallel batch=32 test.", 32, 4, 1))
results.append(run_scenario("parallel_4_batch32", "Four simultaneous direct TEI batch=32 requests.", 32, 4, 4))
results.append(run_scenario("parallel_16_batch1", "Sixteen simultaneous direct TEI batch=1 requests.", 1, 16, 16))
print("Sleeping 30s before idle single batch=32 scenario", flush=True)
time.sleep(30)
results.append(run_scenario("idle_30s_batch32", "Single direct TEI batch=32 request after 30s idle.", 32, 1, 1))

restart_started = utc_now()
print(f"Restarting compose service {compose_service} at {restart_started}", flush=True)
subprocess.check_call(["docker", "compose", "restart", compose_service])
wait_for_tei()
results.append(run_scenario("restart_then_batch32", "Single direct TEI batch=32 request after docker compose restart of TEI.", 32, 1, 1))

lines: list[str] = []
lines.extend([
    "---",
    "milestone: M042-fjf2en",
    "slice: S01",
    "task: T02",
    f"captured: {run_started}",
    "---",
    "",
    "# M042 S01 T02 — TEI concurrency profile",
    "",
    "This profile varies direct TEI request concurrency and batch size. It bypasses fd and Redis so the measurements describe TEI scheduling/inference behavior, not fd cache-hot latency.",
    "",
    "## Runtime limits",
    "",
    "| Field | Value |",
    "|---|---:|",
])
for key in ["model_id", "model_dtype", "version", "max_concurrent_requests", "max_batch_requests", "max_client_batch_size", "max_batch_tokens", "max_input_length", "tokenization_workers"]:
    lines.append(f"| {key} | `{info.get(key)}` |")

lines.extend([
    "",
    "## Scenario summary",
    "",
    "| Scenario | Requests | Batch | Concurrency | OK | Client wall p50 | Client wall p95 | TEI total p50 | TEI total p95 | Queue p50 | Queue p95 | Inference p50 | Inference p95 | Server rows |",
    "|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|",
])
for result in results:
    wall = stats([row.wall_ms for row in result.clients])
    total = stats([row[0] for row in result.server_rows])
    queue = stats([row[2] for row in result.server_rows])
    infer = stats([row[3] for row in result.server_rows])
    ok = sum(1 for row in result.clients if row.status == 200)
    lines.append(
        f"| `{result.name}` | {result.request_count} | {result.batch_size} | {result.concurrency} | {ok} | "
        f"{fmt_ms(wall['p50'])} | {fmt_ms(wall['p95'])} | {fmt_ms(total['p50'])} | {fmt_ms(total['p95'])} | "
        f"{fmt_ms(queue['p50'])} | {fmt_ms(queue['p95'])} | {fmt_ms(infer['p50'])} | {fmt_ms(infer['p95'])} | {len(result.server_rows)} |"
    )

lines.extend([
    "",
    "## Observations",
    "",
    "- `parallel_4_batch32` tests whether multiple large client batches improve throughput or amplify TEI queueing.",
    "- `parallel_16_batch1` tests whether the high `max_concurrent_requests` value translates to low-latency small-request concurrency.",
    "- `idle_30s_batch32` checks whether the batch=32 queue pattern disappears after a short idle period.",
    "- `restart_then_batch32` approximates cold process behavior after TEI restart using the already-downloaded model artifacts.",
    "",
    "## Raw client rows",
    "",
    "| Scenario | request | batch | status | client_wall_ms | error |",
    "|---|---:|---:|---:|---:|---|",
])
for result in results:
    for row in result.clients:
        error = row.error.replace("|", "\\|").replace("\n", " ")[:180]
        lines.append(f"| `{row.scenario}` | {row.request_id} | {row.batch_size} | {row.status} | {row.wall_ms:.3f} | {error} |")

lines.extend([
    "",
    "## Raw TEI timing rows",
    "",
    "| Scenario | row | total_ms | token_ms | queue_ms | infer_ms |",
    "|---|---:|---:|---:|---:|---:|",
])
for result in results:
    for idx, (total, token, queue, infer) in enumerate(result.server_rows, 1):
        lines.append(f"| `{result.name}` | {idx} | {total:.3f} | {token:.3f} | {queue:.3f} | {infer:.3f} |")

lines.extend([
    "",
    "## T02 interpretation",
    "",
    "Use this artifact together with `documents/te-perf-snapshot-m042-s01.md` to decide whether S02 async chunking can reduce fd caller latency enough, or whether S03 ONNX remains necessary for true cold-path targets.",
])

out_path.write_text("\n".join(lines) + "\n", encoding="utf-8")
print(f"wrote {out_path}")
PY
