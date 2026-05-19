---
id: T02
parent: S01
milestone: M003-xx4yc3
key_files:
  - docker-compose.override.yaml
key_decisions:
  - Bind Redis host port to localhost in docker-compose.override.yaml instead of 0.0.0.0 to keep local benchmark access without exposing Redis externally.
duration: 
verification_result: mixed
completed_at: 2026-05-19T08:12:03.471Z
blocker_discovered: false
---

# T02: Started the full stack after fixing stale API container conflict and localhost-binding Redis.

**Started the full stack after fixing stale API container conflict and localhost-binding Redis.**

## What Happened

Started the full Compose stack with build. The first attempt built the API image successfully but failed to create `fd_api` because an exited container with the same name already existed. Logs also revealed a serious prior Redis exposure: external clients attempted cross-protocol/replication attacks while Redis was published on 0.0.0.0:6379. After GitNexus impact analysis for docker-compose.override.yaml returned LOW risk, the Redis override port was changed to `127.0.0.1:6379:6379`, the stale fd_api container was removed, and `docker compose up -d --build` was rerun successfully. Final ps/inspect showed fd_tei, fd_redis, and fd_api all healthy; API /health returned JSON status ok.

## Verification

`docker compose up -d --build` succeeded on retry. `docker compose ps` showed fd_api, fd_redis, fd_tei healthy. API health returned `{"status":"ok"...}`. Redis port is now bound to `127.0.0.1:6379`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker compose up -d --build` | 1 | ❌ initial fail: stale fd_api name conflict | 0ms |
| 2 | `gitnexus_impact docker-compose.override.yaml repo=fd` | 0 | ✅ LOW risk, no upstream code impact | 0ms |
| 3 | `docker rm fd_api && docker compose config | grep host_ip` | 0 | ✅ pass: stale container removed and Redis host_ip is 127.0.0.1 | 0ms |
| 4 | `docker compose up -d --build && docker compose ps && curl -fsS http://localhost:8000/health` | 0 | ✅ pass: all services healthy and API health OK | 0ms |

## Deviations

S01 discovered and fixed an additional runtime exposure issue: docker-compose.override.yaml published Redis on all interfaces, and Redis logs showed actual external attack attempts. The override was changed to bind Redis to 127.0.0.1 only. A stale exited fd_api container was removed to resolve a Compose name conflict.

## Known Issues

Redis logs still warn that host memory overcommit is disabled; this is host configuration, not an app/container startup blocker.

## Files Created/Modified

- `docker-compose.override.yaml`
