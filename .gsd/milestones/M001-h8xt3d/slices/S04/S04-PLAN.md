# S04: Runtime configuration hardening

**Goal:** Harden Docker/runtime defaults without breaking local development.
**Demo:** Docker Compose healthchecks work with runtime image and config names match application expectations.

## Must-Haves

- API runtime image has the healthcheck dependency or healthcheck no longer requires missing curl.
- Base compose uses REDIS_HOST rather than unused REDIS_ADDR.
- Redis host port exposure is safer by default or clearly isolated to dev override.
- PORT is removed or implemented consistently.

## Proof Level

- This slice proves: config inspection plus Go test suite

## Integration Closure

Compose still starts TEI, API, and Redis; app env names align with main.go.

## Verification

- Healthcheck reports actual API health; Redis is not unnecessarily exposed by default.

## Tasks

- [x] **T01: Assess runtime config blast radius** `est:small`
  Assess config/runtime blast radius for main.go env handling, Dockerfile healthcheck dependency, and compose Redis exposure.
  - Files: `api/main.go`, `api/Dockerfile`, `docker-compose.yaml`, `docker-compose.override.yaml`, `README.md`
  - Verify: No code changes; document findings.

- [x] **T02: Implement runtime hardening** `est:medium`
  Implement PORT env handling in main.go, align compose variable names, ensure API healthcheck command exists in runtime image, and move Redis host port exposure out of base compose into override for local dev.
  - Files: `api/main.go`, `api/Dockerfile`, `docker-compose.yaml`, `docker-compose.override.yaml`, `README.md`
  - Verify: docker compose config && cd api && go test ./...

- [x] **T03: Verify runtime hardening** `est:small`
  Run final full suite and complete milestone if all slices pass.
  - Verify: cd api && go test ./... -short

## Files Likely Touched

- api/main.go
- api/Dockerfile
- docker-compose.yaml
- docker-compose.override.yaml
- README.md
