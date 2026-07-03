---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Validate dimension-specific cache keys

Request same text at 1024d and 512d, then verify separate Redis keys and expected payload sizes.

## Inputs

- `S03 T01`

## Expected Output

- `S03 T02 summary`

## Verification

Redis scan and STRLEN show d1024=4098 and d512=2050.

## Observability Impact

Proves dimension-scoped cache behavior in runtime.
