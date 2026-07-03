---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Live dependency health checks

Verify API health, TEI health, and Redis ping against the running stack.

## Inputs

- `S01 healthy stack`

## Expected Output

- `S02 T01 summary`

## Verification

curl API/TEI health and redis-cli ping.

## Observability Impact

Captures live dependency health evidence.
