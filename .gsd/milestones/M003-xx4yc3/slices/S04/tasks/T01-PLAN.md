---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Prepare benchmark environment

Verify Python benchmark dependencies without installing globally.

## Inputs

- `S03 runtime stack`

## Expected Output

- `S04 T01 summary`

## Verification

python3 imports requests and redis.

## Observability Impact

Prevents benchmark failure due to missing local deps.
