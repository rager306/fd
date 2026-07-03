---
estimated_steps: 1
estimated_files: 5
skills_used: []
---

# T02: Implement runtime hardening

Implement PORT env handling in main.go, align compose variable names, ensure API healthcheck command exists in runtime image, and move Redis host port exposure out of base compose into override for local dev.

## Inputs

- `S04 T01`

## Expected Output

- `api/main.go`
- `api/Dockerfile`
- `docker-compose.yaml`
- `docker-compose.override.yaml`
- `README.md`

## Verification

docker compose config && cd api && go test ./...

## Observability Impact

Healthcheck and runtime env config become trustworthy.
