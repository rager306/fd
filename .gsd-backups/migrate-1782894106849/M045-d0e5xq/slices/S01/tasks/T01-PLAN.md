---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Captured current TEI runtime state without restarting containers.

Collect read-only evidence from docker compose config, docker inspect, fd `/health`, fd `/ready`, fd embedding smoke, direct TEI embedding smoke, and container health state. Do not restart or recreate containers. Save a compact evidence artifact under benchmark-results or documents.

## Inputs

- `docker-compose.yaml`
- `docker-compose.override.yaml`
- `M042 S01/S02 summaries`

## Expected Output

- `documents/tei-startup-recon-m045.md`

## Verification

Read-only commands succeed and artifact contains image, command, env subset, health state, fd runtime block, and smoke results.

## Observability Impact

Creates current known-good runtime baseline before any mitigation.
