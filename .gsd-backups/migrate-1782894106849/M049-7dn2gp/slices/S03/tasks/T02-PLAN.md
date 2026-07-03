---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Rebuilt the API container and passed live health/metrics/cache invalidation smoke.

Rebuild/restart the api container via Docker Compose, wait for health, and run authenticated live HTTP smoke for cache invalidation, health fields, metrics gauges, and embedding behavior.

## Inputs

- `docker-compose.yaml`
- `docker-compose.override.yaml`
- `api/.env`

## Expected Output

- `benchmark-results/m049-s03-live-container-proof.md`

## Verification

docker compose up -d --build api; live HTTP smoke script passes.

## Observability Impact

Proves new surfaces work in the actual rebuilt runtime.
