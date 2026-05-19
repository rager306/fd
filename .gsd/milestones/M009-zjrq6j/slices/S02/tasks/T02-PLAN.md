---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Implement cache namespace and retention

Implement cache config parsing and Redis key namespace/retention behavior. Preserve current default TTL and key behavior as closely as possible, support explicit namespace fields, and support no-expire mode. Wire from main runtime env and keep benchmark snapshot allowlist aligned.

## Inputs

- `.gsd/milestones/M009-zjrq6j/slices/S02/tasks/T01-SUMMARY.md`

## Expected Output

- `api/cache/redis.go`
- `api/cache/redis_test.go`
- `api/main.go`
- `benchmark.py`
- `.gsd/milestones/M009-zjrq6j/slices/S02/tasks/T02-SUMMARY.md`

## Verification

Targeted Go tests and Python snapshot parser pass.

## Observability Impact

Configurable namespace/retention can be reported and verified without raw text.
