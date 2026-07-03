---
id: T02
parent: S01
milestone: M045-d0e5xq
key_files:
  - documents/tei-startup-recon-m045.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T11:30:47.896Z
blocker_discovered: false
---

# T02: Summarized TEI startup log sequence and ONNX/ORT probe pattern.

**Summarized TEI startup log sequence and ONNX/ORT probe pattern.**

## What Happened

Parsed current TEI container logs without restarting the service. The recon artifact records TEI starting with `text-embeddings-router --model-id deepvk/USER-bge-m3`, then attempting ONNX/ORT artifacts before falling back to Candle/safetensors and becoming ready.

## Verification

`documents/tei-startup-recon-m045.md` includes the startup log timeline and observed pattern. It separates current running evidence from the future controlled startup proof still needed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker compose logs --no-color --tail=260 tei parsed for startup, ONNX, ORT, Candle, warmup, ready events` | 0 | ✅ pass: startup warning sequence summarized | 30000ms |

## Deviations

None.

## Known Issues

The current log tail contains startup events but not a fresh controlled restart capture; that remains S03 scope.

## Files Created/Modified

- `documents/tei-startup-recon-m045.md`
