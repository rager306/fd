---
id: T02
parent: S03
milestone: M045-d0e5xq
key_files:
  - docker-compose.yaml
  - docker-compose.override.yaml
  - benchmark-results/m045-tei-local-path-startup-proof.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T12:18:58.870Z
blocker_discovered: false
---

# T02: Applied local snapshot TEI command and measured successful startup.

**Applied local snapshot TEI command and measured successful startup.**

## What Happened

Updated compose/override so TEI uses the cached local USER-bge-m3 snapshot path as `--model-id`, applied `docker compose up -d tei`, and polled Docker health. The container command now includes the local snapshot path and `--max-batch-tokens 8192`. TEI reached healthy at 2026-06-14T12:15:15 after starting at 2026-06-14T12:12:14.

## Verification

`benchmark-results/m045-tei-local-path-startup-proof.md` records compose apply, polling, healthy outcome, start/healthy timestamps, post-state, and notable TEI logs. Current `docker inspect` shows fd_tei healthy with the local snapshot command.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker compose up -d tei; poll docker health until healthy; capture logs` | 0 | ✅ pass: TEI reached healthy with local snapshot command | 193300ms |

## Deviations

None after S03 replan.

## Known Issues

TEI still logs an immediate local ORT error for missing ONNX, but it no longer blocks for remote Hub probes and falls through to a healthy Candle/safetensors startup.

## Files Created/Modified

- `docker-compose.yaml`
- `docker-compose.override.yaml`
- `benchmark-results/m045-tei-local-path-startup-proof.md`
