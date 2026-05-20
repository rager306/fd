---
id: T01
parent: S01
milestone: M017-j10hmp
key_files:
  - tools/evaluate_legal_retrieval.py
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
key_decisions:
  - Use tagged Go ONNX endpoint on port 18000 with `ONNX_MAX_SEQUENCE_LENGTH=512`.
  - Use isolated ONNX cache namespace `m017-onnx-512-legal-quality`.
  - Evaluator output path will be `benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt`.
  - Treat evaluator exit code 2 as a measured FAIL artifact, not an execution failure, because the script returns 2 when quality thresholds fail.
duration: 
verification_result: passed
completed_at: 2026-05-20T07:25:13.080Z
blocker_discovered: false
---

# T01: Prepared the ONNX 512 legal gate command plan and verified prerequisites.

**Prepared the ONNX 512 legal gate command plan and verified prerequisites.**

## What Happened

Prepared the M017 S01 command plan. Existing evaluator supports TEI and ONNX API URLs, cache namespace labels, runtime labels, and thresholded PASS/FAIL. Required local ONNX artifact, native tokenizer library, tokenizer JSON, and legal corpus are present. The tagged ONNX service will be launched on port 18000 with max sequence length 512 and an isolated Redis namespace before running the evaluator.

## Verification

Prerequisite check passed, TEI health returned ok, and no background processes were running.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `test required ONNX/native tokenizer/tokenizer/corpus paths` | 0 | ✅ pass — m017_prereqs_ok | 0ms |
| 2 | `curl -fsS http://localhost:8000/health` | 0 | ✅ pass — TEI fd API health ok | 0ms |
| 3 | `bg_shell list` | 0 | ✅ pass — No background processes | 0ms |

## Deviations

None.

## Known Issues

TEI health is currently ok on port 8000. No tagged ONNX server is running yet.

## Files Created/Modified

- `tools/evaluate_legal_retrieval.py`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
