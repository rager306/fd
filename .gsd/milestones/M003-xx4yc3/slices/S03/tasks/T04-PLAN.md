---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T04: Validate Redis L2 after API restart

Restart API and validate Redis L2 hit for an already cached key.

## Inputs

- `S03 T03`

## Expected Output

- `S03 T04 summary`

## Verification

Second request after API restart succeeds without TEI miss for same key.

## Observability Impact

Proves cache survives API process restart via Redis.
