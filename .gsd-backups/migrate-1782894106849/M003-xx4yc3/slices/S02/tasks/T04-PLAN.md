---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T04: Negative API tests

Run negative validation tests for invalid JSON, empty input, invalid dimensions, and invalid batch encoding.

## Inputs

- `S02 positive smoke tests`

## Expected Output

- `negative status evidence`

## Verification

curl status codes are 400.

## Observability Impact

Confirms error paths are observable and return 400.
