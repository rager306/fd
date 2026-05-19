---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Correlate benchmark with resources

Capture docker stats and recent logs after benchmark, then summarize bottlenecks.

## Inputs

- `S04 T02 benchmark output`

## Expected Output

- `S04 summary`

## Verification

docker stats --no-stream and logs captured.

## Observability Impact

Correlates metrics with container resource/log signals.
