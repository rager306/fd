# S04: Runtime configuration hardening — UAT

**Milestone:** M001-h8xt3d
**Written:** 2026-05-19T07:07:45.141Z

# UAT: S04 Runtime configuration hardening

## Verification performed

- `docker compose config` — passed with non-blocking obsolete-version warning.
- `docker compose -f docker-compose.yaml config` — passed.
- Base compose Redis host port check — passed: no `published: "6379"`.
- `cd api && go test ./...` — passed.
- `cd api && go test ./... -short` — passed, 46 tests in 4 packages.

## Acceptance checks

- API runtime image installs curl for healthcheck.
- Base compose uses `REDIS_HOST`, not `REDIS_ADDR`.
- Base compose does not expose Redis on host port 6379.
- Local override keeps Redis host port available for development.
- App honors `PORT` and README documents it.

