---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Verify runtime hardening

Run final full suite and complete milestone if all slices pass.

## Inputs

- `S04 T02 changes`

## Expected Output

- `S04 summaries`
- `milestone validation`

## Verification

cd api && go test ./... -short

## Observability Impact

Confirms project is ready after remediation.
