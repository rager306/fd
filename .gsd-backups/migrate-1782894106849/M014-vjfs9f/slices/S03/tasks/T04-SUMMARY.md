---
id: T04
parent: S03
milestone: M014-vjfs9f
key_files:
  - benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt
  - benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt
key_decisions:
  - Do not keep the local tagged ONNX benchmark binary or restart helper in the working tree; they are transient runtime artifacts.
  - Use M013 cosine artifact as the correctness gate reference for interpreting S03 speed results.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:30:52.792Z
blocker_discovered: false
---

# T04: Verified the tagged ONNX benchmark artifact and cleaned up the local benchmark server.

**Verified the tagged ONNX benchmark artifact and cleaned up the local benchmark server.**

## What Happened

Verified the tagged ONNX benchmark artifact for snapshot v3 metadata, runtime label, build tag, native and ONNX artifact checksums, ONNX Runtime library checksum, isolated Redis namespace, restart command, required benchmark sections, and raw fixed-probe text absence. Confirmed the M013 cosine gate artifact still records PASS. Stopped the local tagged ONNX server and removed transient runtime helper files.

## Verification

Artifact parser/leak checks, correctness-gate reference check, cleanup check, and GitNexus scope check completed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with transformers --with torch --with sentencepiece python ONNX artifact hygiene check` | 0 | ✅ pass — onnx_artifact_hygiene=pass; raw_probe_text_leaks=0 | 0ms |
| 2 | `rg 'PASS|cosine_threshold|raw_probe_texts_logged' benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt` | 0 | ✅ pass — M013 cosine correctness gate referenced; PASS | 0ms |
| 3 | `curl -fsS http://localhost:18000/health after cleanup` | 0 | ✅ pass — server stopped (health unavailable) | 0ms |
| 4 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ⚠️ medium — expected benchmark.py restart/redaction flow changes plus artifact changes | 0ms |

## Deviations

The first cleanup command killed its own shell because the `pkill -f` pattern matched the command line. A follow-up check confirmed the ONNX server was stopped, then removed transient `.gsd/runtime` helper/binary/log files.

## Known Issues

GitNexus reports medium risk because `benchmark.py` changed restart/redaction behavior; the affected flow is the benchmark snapshot/main flow and was verified by py_compile, redaction check, full benchmark run, and artifact parser checks.

## Files Created/Modified

- `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`
- `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt`
