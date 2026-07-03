---
id: T03
parent: S03
milestone: M014-vjfs9f
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt
key_decisions:
  - Use `BENCHMARK_API_RESTART_COMMAND=.gsd/runtime/restart-fd-api-onnx-m014.sh` for tagged ONNX benchmark so L2 restart tests the actual port 18000 server.
  - Use snapshot_version 3 for restart-aware benchmark artifacts.
  - Final ONNX benchmark artifact is `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:29:42.711Z
blocker_discovered: false
---

# T03: Ran the tagged ONNX benchmark with restart-aware snapshot v3 metadata and isolated Redis namespace.

**Ran the tagged ONNX benchmark with restart-aware snapshot v3 metadata and isolated Redis namespace.**

## What Happened

Updated the benchmark harness to support an explicit API restart command and tightened safe env redaction so tokenizer manifest path is not falsely omitted. Built a tagged ONNX binary in `.gsd/runtime`, created a local restart helper, and ran the full benchmark against `http://localhost:18000` with build tag, native tokenizer manifest, ONNX manifest, ONNX Runtime library hash, and isolated cache namespace metadata. The final ONNX artifact records snapshot_version 3 and summary metrics: best cold latency 10.2ms, warm mean 1.63ms, max throughput ~891 req/s at 4 concurrent, Redis L2 restart 2.70ms, batch L1 p95 5.62ms, batch L2 p95 4.41ms, chunk reuse warm p95 9.00ms.

## Verification

py_compile and snapshot redaction checks passed; ONNX benchmark command exited 0; artifact metadata and summary fields were verified.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile benchmark.py` | 0 | ✅ pass | 0ms |
| 2 | `BENCHMARK_NATIVE_TOKENIZER_MANIFEST=... SECRET_TOKEN=redacted uv run --python 3.13 --with requests --with redis python redaction check` | 0 | ✅ pass — tokenizer manifest retained; SECRET_TOKEN omitted | 0ms |
| 3 | `BENCHMARK_RUNTIME_LABEL=tagged-onnx-hf ... BENCHMARK_API_RESTART_COMMAND=.gsd/runtime/restart-fd-api-onnx-m014.sh uv run --python 3.13 --with requests --with redis python benchmark.py > benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt` | 0 | ✅ pass — ONNX benchmark artifact written | 33500ms |
| 4 | `artifact metadata grep/check` | 0 | ✅ pass — snapshot_version=3; tagged metadata present; ORT sha256 present | 0ms |

## Deviations

A benchmark harness issue was discovered and fixed before accepting the ONNX artifact: Redis L2 restart now uses configurable `BENCHMARK_API_RESTART_COMMAND` instead of always restarting Compose `api`. Snapshot version is now 3 for the ONNX artifact. A broad secret-redaction false positive for `TOKENIZER` was also fixed so tokenizer manifest env paths are retained while real token/secret keys remain omitted.

## Known Issues

ONNX artifact uses a local binary and restart helper under `.gsd/runtime`; this remains local benchmark infrastructure, not production packaging. TEI baseline artifact is snapshot v2 while ONNX artifact is snapshot v3 because the restart bug was discovered during S03.

## Files Created/Modified

- `benchmark.py`
- `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`
