---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T01: Design cache config surface

Inspect current Redis/Tiered cache construction and key generation. Design a small cache config surface with defaults preserving current keys/TTL unless env vars are set, plus validation for TTL/no-expire conflicts.

## Inputs

- `.gsd/milestones/M009-zjrq6j/slices/S01/S01-SUMMARY.md`

## Expected Output

- `.gsd/milestones/M009-zjrq6j/slices/S02/tasks/T01-SUMMARY.md`

## Verification

Summary names defaults, env vars, validation rules, and affected symbols.

## Observability Impact

Defines which cache config fields should be safe to expose in benchmark snapshots.
