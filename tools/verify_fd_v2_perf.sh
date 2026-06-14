#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${FD_BASE_URL:-http://localhost:8000}"
OUT="${FD_PERF_OUT:-benchmark-results/fd-v2-perf-validation-m041-s04.md}"
MODEL="${FD_PERF_MODEL:-deepvk/USER-bge-m3}"

if [[ "${1:-}" == "--help" ]]; then
  cat <<'HELP'
Usage: FD_BASE_URL=http://localhost:8000 tools/verify_fd_v2_perf.sh

Runs fd v2 cache-hot steady-state performance checks:
- prewarm each measured payload through real inference
- batch=1 cache HIT p95 < 50ms
- batch=10 cache HIT p95 < 200ms
- batch=32 cache HIT p95 < 1000ms
- 100 sequential cache-hot requests: 0 errors
- 4 concurrent callers x 8 cached inputs complete < 2s total
- repeated identical input returns X-Cache: HIT with latency < 5ms

The report also includes non-blocking cache-miss diagnostics for real TEI/ONNX
inference latency. Writes markdown results to
benchmark-results/fd-v2-perf-validation-m041-s04.md by default.
HELP
  exit 0
fi

mkdir -p "$(dirname "$OUT")"

python3 - "$BASE_URL" "$OUT" "$MODEL" <<'PY'
import concurrent.futures
import json
import statistics
import sys
import time
import urllib.error
import urllib.request
from pathlib import Path

base_url, out_path, model = sys.argv[1:4]
run_id = time.time_ns()


def percentile(values, pct):
    if not values:
        return None
    ordered = sorted(values)
    idx = int((len(ordered) - 1) * pct / 100)
    return ordered[idx]


def x_cache(headers):
    return {k.lower(): v for k, v in headers.items()}.get("x-cache", "")


def input_for(batch_size, text_prefix):
    if batch_size == 1:
        return f"{text_prefix}-single"
    return [f"{text_prefix}-{i}" for i in range(batch_size)]


def request_embedding(batch_size, text_prefix="perf"):
    body = json.dumps({"model": model, "input": input_for(batch_size, text_prefix)}).encode()
    req = urllib.request.Request(
        f"{base_url}/v1/embeddings",
        data=body,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    start = time.perf_counter()
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            payload = resp.read()
            elapsed_ms = (time.perf_counter() - start) * 1000
            return {
                "ok": 200 <= resp.status < 300,
                "status": resp.status,
                "latency_ms": elapsed_ms,
                "headers": dict(resp.headers.items()),
                "bytes": len(payload),
                "x_cache": x_cache(dict(resp.headers.items())),
                "error": "",
            }
    except urllib.error.HTTPError as exc:
        elapsed_ms = (time.perf_counter() - start) * 1000
        headers = dict(exc.headers.items())
        return {
            "ok": False,
            "status": exc.code,
            "latency_ms": elapsed_ms,
            "headers": headers,
            "bytes": len(exc.read()),
            "x_cache": x_cache(headers),
            "error": f"HTTPError:{exc.code}",
        }
    except Exception as exc:
        elapsed_ms = (time.perf_counter() - start) * 1000
        return {
            "ok": False,
            "status": 0,
            "latency_ms": elapsed_ms,
            "headers": {},
            "bytes": 0,
            "x_cache": "",
            "error": repr(exc),
        }


def cache_hot_result(result):
    return result["ok"] and result["x_cache"].upper() == "HIT"


def run_latency_case(batch_size, count, threshold_ms):
    prefix = f"perf-hot-{run_id}-b{batch_size}"
    prewarm = request_embedding(batch_size, prefix)
    results = [request_embedding(batch_size, prefix) for _ in range(count)]
    latencies = [r["latency_ms"] for r in results if r["ok"]]
    errors = [r for r in results if not r["ok"]]
    non_hits = [r for r in results if r["ok"] and r["x_cache"].upper() != "HIT"]
    p95 = percentile(latencies, 95)
    passed = prewarm["ok"] and not errors and not non_hits and p95 is not None and p95 < threshold_ms
    return {
        "batch": batch_size,
        "count": count,
        "threshold_ms": threshold_ms,
        "prewarm_status": prewarm["status"],
        "prewarm_x_cache": prewarm["x_cache"],
        "p50": percentile(latencies, 50),
        "p95": p95,
        "p99": percentile(latencies, 99),
        "errors": len(errors),
        "non_hit_responses": len(non_hits),
        "passed": passed,
        "error_samples": errors[:3],
    }


def run_sequential_100():
    prefix = f"perf-hot-{run_id}-seq"
    prewarm = request_embedding(1, prefix)
    results = [request_embedding(1, prefix) for _ in range(100)]
    errors = [r for r in results if not r["ok"]]
    non_hits = [r for r in results if r["ok"] and r["x_cache"].upper() != "HIT"]
    return {
        "count": 100,
        "prewarm_status": prewarm["status"],
        "prewarm_x_cache": prewarm["x_cache"],
        "errors": len(errors),
        "non_hit_responses": len(non_hits),
        "passed": prewarm["ok"] and len(errors) == 0 and len(non_hits) == 0,
        "error_samples": errors[:3],
    }


def run_concurrent():
    prefixes = [f"perf-hot-{run_id}-concurrent-{i}" for i in range(4)]
    prewarms = [request_embedding(8, prefix) for prefix in prefixes]
    start = time.perf_counter()
    with concurrent.futures.ThreadPoolExecutor(max_workers=4) as executor:
        futures = [executor.submit(request_embedding, 8, prefix) for prefix in prefixes]
        results = [f.result() for f in futures]
    elapsed = time.perf_counter() - start
    errors = [r for r in results if not r["ok"]]
    non_hits = [r for r in results if r["ok"] and r["x_cache"].upper() != "HIT"]
    return {
        "workers": 4,
        "batch": 8,
        "elapsed_s": elapsed,
        "prewarm_statuses": [r["status"] for r in prewarms],
        "prewarm_x_cache": [r["x_cache"] for r in prewarms],
        "errors": len(errors),
        "non_hit_responses": len(non_hits),
        "passed": all(r["ok"] for r in prewarms) and len(errors) == 0 and len(non_hits) == 0 and elapsed < 2.0,
        "error_samples": errors[:3],
    }


def run_cache_check():
    key = f"perf-cache-{run_id}"
    first = request_embedding(1, key)
    second_start = time.perf_counter()
    second = request_embedding(1, key)
    hit_latency_ms = (time.perf_counter() - second_start) * 1000
    return {
        "first_status": first["status"],
        "first_x_cache": first["x_cache"],
        "second_status": second["status"],
        "x_cache": second["x_cache"],
        "hit_latency_ms": hit_latency_ms,
        "passed": first["ok"] and second["ok"] and second["x_cache"].upper() == "HIT" and hit_latency_ms < 5.0,
    }


def run_miss_diagnostics():
    diagnostics = []
    for batch_size in (1, 10, 32):
        prefix = f"perf-miss-diag-{run_id}-b{batch_size}"
        result = request_embedding(batch_size, prefix)
        diagnostics.append({
            "batch": batch_size,
            "status": result["status"],
            "x_cache": result["x_cache"],
            "latency_ms": result["latency_ms"],
            "ok": result["ok"],
        })
    return diagnostics

miss_diagnostics = run_miss_diagnostics()
cases = [
    run_latency_case(1, 50, 50),
    run_latency_case(10, 20, 200),
    run_latency_case(32, 10, 1000),
]
seq = run_sequential_100()
concurrent = run_concurrent()
cache = run_cache_check()
all_passed = all(c["passed"] for c in cases) and seq["passed"] and concurrent["passed"] and cache["passed"]

lines = [
    "# fd v2 performance validation — M041 S04",
    "",
    f"Base URL: `{base_url}`",
    "Mode: cache-hot steady-state after explicit prewarm; real cache-miss inference is diagnostic only.",
    f"Overall: {'PASS' if all_passed else 'FAIL'}",
    "",
    "## Cache-hot latency cases",
    "",
    "| batch | count | prewarm X-Cache | p50 ms | p95 ms | p99 ms | threshold ms | errors | non-HIT | verdict |",
    "|---:|---:|---|---:|---:|---:|---:|---:|---:|---|",
]
for c in cases:
    fmt = lambda v: "n/a" if v is None else f"{v:.3f}"
    lines.append(f"| {c['batch']} | {c['count']} | {c['prewarm_x_cache'] or 'n/a'} | {fmt(c['p50'])} | {fmt(c['p95'])} | {fmt(c['p99'])} | {c['threshold_ms']} | {c['errors']} | {c['non_hit_responses']} | {'PASS' if c['passed'] else 'FAIL'} |")
lines += [
    "",
    "## Sequential and concurrent",
    "",
    f"- 100 sequential cache-hot zero errors: {'PASS' if seq['passed'] else 'FAIL'} ({seq['errors']} errors, {seq['non_hit_responses']} non-HIT responses)",
    f"- 4 concurrent × 8 cache-hot inputs < 2s: {'PASS' if concurrent['passed'] else 'FAIL'} ({concurrent['elapsed_s']:.3f}s, {concurrent['errors']} errors, {concurrent['non_hit_responses']} non-HIT responses)",
    "",
    "## Cache effectiveness",
    "",
    f"- Repeated input X-Cache HIT and <5ms: {'PASS' if cache['passed'] else 'FAIL'} (first X-Cache={cache['first_x_cache']!r}, second X-Cache={cache['x_cache']!r}, latency={cache['hit_latency_ms']:.3f}ms)",
    "",
    "## Non-blocking cache-miss diagnostics",
    "",
    "| batch | status | X-Cache | latency ms |",
    "|---:|---:|---|---:|",
]
for diag in miss_diagnostics:
    lines.append(f"| {diag['batch']} | {diag['status']} | {diag['x_cache'] or 'n/a'} | {diag['latency_ms']:.3f} |")
lines += [
    "",
    "## Raw summary",
    "",
    "```json",
    json.dumps({"mode": "cache-hot", "miss_diagnostics": miss_diagnostics, "latency": cases, "sequential": seq, "concurrent": concurrent, "cache": cache}, indent=2),
    "```",
]
Path(out_path).write_text("\n".join(lines) + "\n")
print(f"wrote {out_path}")
print("PASS" if all_passed else "FAIL")
sys.exit(0 if all_passed else 1)
PY
