---
id: S01
parent: M003-xx4yc3
milestone: M003-xx4yc3
provides:
  - Healthy local stack for live endpoint, cache, and benchmark validation.
requires:
  []
affects:
  - S02
  - S03
  - S04
key_files:
  - docker-compose.override.yaml
key_decisions:
  - Keep Redis accessible from host for local benchmark.py but bind it to 127.0.0.1 only.
patterns_established:
  - Runtime validation must inspect logs, not just health checks; Redis logs revealed an exposure issue that config review alone had not fully closed for override mode.
observability_surfaces:
  - Captured Docker Compose build output, service health states, API logs, TEI logs, and Redis logs.
  - Discovered Redis external attack attempts from real logs.
drill_down_paths:
  - .gsd/milestones/M003-xx4yc3/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M003-xx4yc3/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M003-xx4yc3/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T08:12:52.785Z
blocker_discovered: false
---

# S01: Compose startup and logs

**S01 proved the stack starts after fixing a stale API container conflict and localhost-binding Redis.**

## What Happened

S01 validated real Compose startup. The baseline showed Redis and TEI already healthy and API absent. The first API startup failed due to a stale exited `fd_api` container. Logs also exposed a security issue from prior runtime: Redis received external attack attempts while bound to all interfaces. The override was hardened to bind Redis only on localhost, the stale API container was removed, and the stack started successfully. Final evidence shows all services healthy and API `/health` returning ok.

## Verification

All S01 startup gates passed after the local runtime/config fix.

## Requirements Advanced

- Runtime validation startup requirement satisfied. — 

## Requirements Validated

- Compose config renders successfully. — 
- Env-file presence recorded without printing values. — 
- Stack startup attempted and fixed after root-cause diagnosis. — 
- TEI, Redis, and API are healthy. — 
- API health returns ok. — 

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S01 included an immediate config fix to bind Redis to localhost-only after logs showed real external attack attempts against the previously exposed Redis port.

## Known Limitations

Redis host memory overcommit warning remains a host tuning item. TEI logs include historical ONNX fallback warnings but service is healthy on Candle CPU backend.

## Follow-ups

Continue with S02 live health/API smoke tests. Later commit docker-compose.override.yaml and M003 artifacts after milestone or logical checkpoint.

## Files Created/Modified

- `docker-compose.override.yaml` — Changed Redis host port binding from all interfaces to localhost-only.
