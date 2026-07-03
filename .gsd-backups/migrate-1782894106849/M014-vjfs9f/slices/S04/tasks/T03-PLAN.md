---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Run final M014 verification gates

Run final verification gates: Go tests, pinned lint, tagged tests, artifact hygiene, tracked binary checks, GitNexus detect, milestone validation.

## Inputs

- `benchmark.py`
- `api`

## Expected Output

- `Task summary with final verification evidence`

## Verification

All commands pass and no tagged server/background process remains.

## Observability Impact

Prevents closing the milestone without fresh evidence.
