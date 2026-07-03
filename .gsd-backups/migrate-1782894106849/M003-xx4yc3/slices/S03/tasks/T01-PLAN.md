---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Validate Redis key and 1024d payload

Flush local Redis, make a 1024d request, and verify Redis key suffix and binary payload size.

## Inputs

- `S02 live stack`

## Expected Output

- `S03 T01 summary`

## Verification

Redis scan and STRLEN show d1024 key and 4098 bytes.

## Observability Impact

Proves L2 key/payload format in runtime.
