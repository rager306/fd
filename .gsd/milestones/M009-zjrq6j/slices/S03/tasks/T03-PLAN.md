---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T03: Verify Redis persistence hardening

Verify Redis restart reuse with persistence enabled: create a cached embedding, force/save if needed, restart Redis, verify key remains, restart API to clear L1, and confirm request succeeds from Redis-backed cache path. Run tests/lint/config and GitNexus detect_changes.

## Inputs

- `docker-compose.yaml`
- `benchmark.py`

## Expected Output

- `benchmark-results/`
- ` .gsd/milestones/M009-zjrq6j/slices/S03/tasks/T03-SUMMARY.md`

## Verification

`docker compose config`; Go tests; lint; Redis restart reuse script; benchmark snapshot parser; GitNexus detect_changes.

## Observability Impact

Produces live evidence that Redis cache survives Redis restart under configured persistence.
