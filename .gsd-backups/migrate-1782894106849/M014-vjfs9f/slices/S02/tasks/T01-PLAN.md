---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Preflight TEI benchmark runtime

Verify the current default Docker stack and runtime health before running the TEI baseline. Capture compose/health state and ensure no tagged ONNX server is running.

## Inputs

- `docker-compose.yaml`
- `docker-compose.override.yaml`

## Expected Output

- `Task summary with preflight evidence`

## Verification

Docker compose ps and API health show default TEI stack healthy; no background tagged server.

## Observability Impact

Ensures TEI baseline is measured against the intended default runtime.
