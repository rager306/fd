---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Authenticated Docker Compose e2e proof passed against real containers.

Подготовить локальный matching API key without printing it, restart/reuse Docker Compose API as needed, then run `cd tests/integration && FD_INTEGRATION_API_KEY=... go test -v .` against real containers. Record command outcomes without secrets.

## Inputs

- `tests/integration/api_test.go`
- `docker-compose.yaml`
- `docker-compose.override.yaml`

## Expected Output

- `benchmark-results/m050-s02-docker-e2e.md`

## Verification

Docker services healthy; authenticated e2e suite passes against `localhost:8000`.

## Observability Impact

Proves auth, embedding runtime, cache and metrics work together in the actual container.
