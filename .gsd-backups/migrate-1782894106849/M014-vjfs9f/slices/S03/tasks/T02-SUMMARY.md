---
id: T02
parent: S03
milestone: M014-vjfs9f
key_files: []
key_decisions:
  - Tagged ONNX benchmark server runs on port 18000 and must be stopped after artifact verification.
  - Observed resident memory for compiled `fd-api` tagged ONNX process was about 1.69 GiB RSS shortly after startup.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:23:12.196Z
blocker_discovered: false
---

# T02: Started and verified the tagged ONNX benchmark server on port 18000.

**Started and verified the tagged ONNX benchmark server on port 18000.**

## What Happened

Started the tagged ONNX API with `go run -tags hf_tokenizers`, native tokenizer CGO flags, ONNX manifest, ONNX Runtime shared library, tokenizer JSON, max sequence length 128, and isolated Redis namespace `m014-onnx-hf-tokenizer`. Port 18000 became ready, `/health` returned ok, logs show `onnx backend ready`, and process RSS was captured shortly after startup.

## Verification

bg_shell readiness, `/health`, logs, and process RSS checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `bg_shell start/wait_for_ready fd-api-onnx-hf-tokenizer-m014 ready_port=18000` | 0 | ✅ pass — port 18000 ready | 5000ms |
| 2 | `curl -fsS http://localhost:18000/health` | 0 | ✅ pass — {"status":"ok"} | 0ms |
| 3 | `ps -o pid,rss,vsz,etime,command -C fd-api -C go` | 0 | ✅ pass — fd-api RSS about 1,768,196 KiB; go wrapper RSS about 31,056 KiB | 0ms |
| 4 | `bg_shell highlights f738a66d` | 0 | ✅ pass — onnx backend ready; listening on 0.0.0.0:18000 | 0ms |

## Deviations

Startup duration is approximate from bg_shell readiness: port ready after about 5s; process uptime about 10s when health/RSS was captured.

## Known Issues

GIN debug-mode warning appears because `GIN_MODE=release` is not set; this is a benchmark caveat but not a blocker for local comparison. The server logs suggest loading ONNX and initializing runtime took several seconds before listen.

## Files Created/Modified

None.
