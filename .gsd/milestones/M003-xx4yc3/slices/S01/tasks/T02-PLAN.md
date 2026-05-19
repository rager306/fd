---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Start stack and collect logs

Start the full Compose stack with build, wait for readiness, and capture ps/log/health evidence for TEI, Redis, and API.

## Inputs

- `S01 T01 baseline`

## Expected Output

- `runtime startup evidence`

## Verification

docker compose up -d --build; docker compose ps; docker inspect health states.

## Observability Impact

Captures startup logs and health states for all services.
