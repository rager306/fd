---
id: T01
parent: S01
milestone: M018-vq2ttb
key_files:
  - tools/evaluate_legal_retrieval.py
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
key_decisions:
  - Use tagged Go ONNX endpoint on port 18000 with `ONNX_MAX_SEQUENCE_LENGTH=1024`.
  - Use isolated ONNX cache namespace `m018-onnx-1024-legal-quality`.
  - Evaluator output path will be `benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt`.
  - Treat evaluator exit code 2 as measured FAIL if an artifact is written.
duration: 
verification_result: passed
completed_at: 2026-05-20T07:36:18.654Z
blocker_discovered: false
---

# T01: Prepared the ONNX 1024 legal gate command plan and verified prerequisites.

**Prepared the ONNX 1024 legal gate command plan and verified prerequisites.**

## What Happened

Prepared the M018 S01 command plan. Required ONNX artifact, native tokenizer library, tokenizer JSON, and legal corpus are present. TEI health is ok on port 8000, and no background process is running. The tagged ONNX 1024 gate will run on port 18000 with isolated cache namespace `m018-onnx-1024-legal-quality`.

## Verification

Prerequisite check passed, TEI health returned ok, and no background processes were running.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `test required ONNX/native tokenizer/tokenizer/corpus paths` | 0 | ✅ pass — m018_prereqs_ok | 0ms |
| 2 | `curl -fsS http://localhost:8000/health` | 0 | ✅ pass — TEI fd API health ok | 0ms |
| 3 | `bg_shell list` | 0 | ✅ pass — No background processes | 0ms |

## Deviations

None.

## Known Issues

The ONNX manifest records export sequence length 128, but runtime inputs have dynamic sequence length and prior Python/Go evidence supports longer runtime lengths. M018 tests the tagged Go path directly.

## Files Created/Modified

- `tools/evaluate_legal_retrieval.py`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
