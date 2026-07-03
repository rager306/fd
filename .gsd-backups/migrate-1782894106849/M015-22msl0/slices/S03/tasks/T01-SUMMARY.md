---
id: T01
parent: S03
milestone: M015-22msl0
key_files: []
key_decisions:
  - Use ONNX cache namespace `m015-onnx-legal-quality` for the live legal gate.
  - Tagged ONNX RSS remains about 1.68 GiB shortly after startup, similar to M014.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:56:19.229Z
blocker_discovered: false
---

# T01: Started and verified TEI and tagged ONNX runtimes for the live legal gate.

**Started and verified TEI and tagged ONNX runtimes for the live legal gate.**

## What Happened

Verified default TEI API health, confirmed required local ONNX/native artifacts exist, and started tagged ONNX on port 18000 with `hf_tokenizers` and isolated cache namespace `m015-onnx-legal-quality`. `/health` returns ok, logs show `onnx backend ready`, and RSS was captured.

## Verification

TEI health, artifact existence, tagged ONNX readiness, ONNX health, logs, and RSS checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `curl -fsS http://localhost:8000/health` | 0 | ✅ pass — TEI/default API health ok | 0ms |
| 2 | `test -f ONNX/native/tokenizer artifacts` | 0 | ✅ pass — artifacts_ok | 0ms |
| 3 | `bg_shell start/wait_for_ready fd-api-onnx-hf-tokenizer-m015` | 0 | ✅ pass — port 18000 ready | 5000ms |
| 4 | `curl -fsS http://localhost:18000/health` | 0 | ✅ pass — tagged ONNX health ok | 0ms |
| 5 | `ps -o pid,rss,vsz,etime,command -C fd-api -C go` | 0 | ✅ pass — tagged fd-api RSS about 1,766,616 KiB | 0ms |

## Deviations

None.

## Known Issues

GIN debug-mode warning remains a local benchmark caveat. It does not block quality evaluation.

## Files Created/Modified

None.
