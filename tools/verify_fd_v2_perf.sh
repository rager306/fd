#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${FD_BASE_URL:-http://localhost:8000}"
OUT="${FD_PERF_OUT:-benchmark-results/fd-v2-perf-validation-m041-s04.md}"
MODEL="${FD_PERF_MODEL:-deepvk/USER-bge-m3}"

if [[ "${1:-}" == "--help" ]]; then
  cat <<'HELP'
Usage: FD_BASE_URL=http://localhost:8000 tools/verify_fd_v2_perf.sh

Runs fd v2 performance checks:
- batch=1 p95 < 50ms
- batch=10 p95 < 200ms
- batch=32 p95 < 1000ms
- 100 sequential requests: 0 errors
- 4 concurrent workers x 8 requests complete < 2s
- repeated identical input returns X-Cache: HIT with latency < 5ms

Writes markdown results to benchmark-results/fd-v2-perf-validation-m041-s04.md by default.
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


def percentile(values, pct):
    if not values:
        return None
    ordered = sorted(values)
    idx = int((len(ordered) - 1) * pct / 100)
    return ordered[idx]


def request_embedding(batch_size, text_prefix="perf"):
    if batch_size == 1:
        input_value = f"{text_prefix}-single"
    else:
        input_value = [f"{text_prefix}-{i}" for i in range(batch_size)]
    body = json.dumps({"model": model, "input": input_value}).encode()
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
                "error": "",
            }
    except urllib.error.HTTPError as exc:
        elapsed_ms = (time.perf_counter() - start) * 1000
        return {
            "ok": False,
            "status": exc.code,
            "latency_ms": elapsed_ms,
            "headers": dict(exc.headers.items()),
            "bytes": len(exc.read()),
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
            "error": repr(exc),
        }


def run_latency_case(batch_size, count, threshold_ms):
    results = [request_embedding(batch_size, f"perf-b{batch_size}-{i}") for i in range(count)]
    latencies = [r["latency_ms"] for r in results if r["ok"]]
    errors = [r for r in results if not r["ok"]]
    p95 = percentile(latencies, 95)
    passed = not errors and p95 is not None and p95 < threshold_ms
    return {"batch": batch_size, "count": count, "threshold_ms": threshold_ms, "p50": percentile(latencies, 50), "p95": p95, "p99": percentile(latencies, 99), "errors": len(errors), "passed": passed, "error_samples": errors[:3]}


def run_sequential_100():
    results = [request_embedding(1, f"perf-seq-{i}") for i in range(100)]
    errors = [r for r in results if not r["ok"]]
    return {"count": 100, "errors": len(errors), "passed": len(errors) == 0, "error_samples": errors[:3]}


def run_concurrent():
    start = time.perf_counter()
    with concurrent.futures.ThreadPoolExecutor(max_workers=4) as executor:
        futures = [executor.submit(request_embedding, 8, f"perf-concurrent-{i}") for i in range(4)]
        results = [f.result() for f in futures]
    elapsed = time.perf_counter() - start
    errors = [r for r in results if not r["ok"]]
    return {"workers": 4, "batch": 8, "elapsed_s": elapsed, "errors": len(errors), "passed": len(errors) == 0 and elapsed < 2.0, "error_samples": errors[:3]}


def run_cache_check():
    key = f"perf-cache-{time.time_ns()}"
    first = request_embedding(1, key)
    second_start = time.perf_counter()
    second = request_embedding(1, key)
    hit_latency_ms = (time.perf_counter() - second_start) * 1000
    x_cache = {k.lower(): v for k, v in second["headers"].items()}.get("x-cache", "")
    return {"first_status": first["status"], "second_status": second["status"], "x_cache": x_cache, "hit_latency_ms": hit_latency_ms, "passed": first["ok"] and second["ok"] and x_cache.upper() == "HIT" and hit_latency_ms < 5.0}

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
    f"Overall: {'PASS' if all_passed else 'FAIL'}",
    "",
    "## Latency cases",
    "",
    "| batch | count | p50 ms | p95 ms | p99 ms | threshold ms | errors | verdict |",
    "|---:|---:|---:|---:|---:|---:|---:|---|",
]
for c in cases:
    fmt = lambda v: "n/a" if v is None else f"{v:.3f}"
    lines.append(f"| {c['batch']} | {c['count']} | {fmt(c['p50'])} | {fmt(c['p95'])} | {fmt(c['p99'])} | {c['threshold_ms']} | {c['errors']} | {'PASS' if c['passed'] else 'FAIL'} |")
lines += [
    "",
    "## Sequential and concurrent",
    "",
    f"- 100 sequential zero errors: {'PASS' if seq['passed'] else 'FAIL'} ({seq['errors']} errors)",
    f"- 4 concurrent × 8 input < 2s: {'PASS' if concurrent['passed'] else 'FAIL'} ({concurrent['elapsed_s']:.3f}s, {concurrent['errors']} errors)",
    "",
    "## Cache effectiveness",
    "",
    f"- Repeated input X-Cache HIT and <5ms: {'PASS' if cache['passed'] else 'FAIL'} (X-Cache={cache['x_cache']!r}, latency={cache['hit_latency_ms']:.3f}ms)",
    "",
    "## Raw summary",
    "",
    "```json",
    json.dumps({"latency": cases, "sequential": seq, "concurrent": concurrent, "cache": cache}, indent=2),
    "```",
]
Path(out_path).write_text("\n".join(lines) + "\n")
print(f"wrote {out_path}")
print("PASS" if all_passed else "FAIL")
sys.exit(0 if all_passed else 1)
PY
