---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Validate cold versus warm behavior

Measure cold and warm request timing and inspect API logs for cache miss behavior.

## Inputs

- `S03 T02`

## Expected Output

- `S03 T03 summary`

## Verification

Warm request faster than cold and no repeated miss log while L1 warm.

## Observability Impact

Captures cache effect through timing/logs.
