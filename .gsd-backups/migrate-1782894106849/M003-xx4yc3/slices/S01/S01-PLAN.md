# S01: Compose startup and logs

**Goal:** Validate full Compose startup and dependency readiness.
**Demo:** Docker Compose stack is started or root-cause-diagnosed with logs and container health evidence.

## Must-Haves

- `docker compose config` renders successfully.
- Existing volumes/env-file state recorded.
- `docker compose up -d --build` attempted.
- Container status and logs captured.
- Any startup failure has a classified root cause.

## Proof Level

- This slice proves: compose config/up/ps/logs/inspect evidence

## Integration Closure

Establishes whether Redis, TEI, and API can run together before endpoint tests.

## Verification

- Captures logs and health state for startup diagnostics.

## Tasks

- [x] **T01: Record runtime baseline** `est:small`
  Record baseline repo/runtime state: git status, compose config, env-file presence without printing secrets, and relevant Docker volumes/containers.
  - Files: `docker-compose.yaml`, `docker-compose.override.yaml`, `api/Dockerfile`
  - Verify: docker compose config and docker compose ps/volume listing succeed.

- [x] **T02: Start stack and collect logs** `est:medium`
  Start the full Compose stack with build, wait for readiness, and capture ps/log/health evidence for TEI, Redis, and API.
  - Verify: docker compose up -d --build; docker compose ps; docker inspect health states.

- [x] **T03: Classify startup health** `est:small`
  Classify startup result. If healthy, complete S01; if failed, capture root cause and create fix plan before editing.
  - Verify: All services healthy or root cause documented.

## Files Likely Touched

- docker-compose.yaml
- docker-compose.override.yaml
- api/Dockerfile
