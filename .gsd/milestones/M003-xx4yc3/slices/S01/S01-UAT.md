# S01: Compose startup and logs — UAT

**Milestone:** M003-xx4yc3
**Written:** 2026-05-19T08:12:52.785Z

# UAT: S01 Compose startup and logs

## Verification performed

- `docker compose config` — passed.
- `docker compose up -d --build` — initial failure due to stale `fd_api` container; retry passed after `docker rm fd_api`.
- `docker compose ps` — `fd_api`, `fd_redis`, `fd_tei` healthy.
- `curl -fsS http://localhost:8000/health` — returned `{"status":"ok",...}`.
- Logs inspected for api/tei/redis.

## Findings

- Fixed: stale exited `fd_api` container name conflict.
- Fixed: Redis host exposure in override changed from `0.0.0.0:6379` to `127.0.0.1:6379` after logs showed external attack attempts.
- Non-blocking: Redis warns host memory overcommit is disabled.

