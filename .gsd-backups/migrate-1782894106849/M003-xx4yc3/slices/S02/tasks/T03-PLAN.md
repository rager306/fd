---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Smoke test batch endpoint

Smoke test `/embeddings/batch` for base64 and float formats.

## Inputs

- `S02 T02 embeddings smoke`

## Expected Output

- `batch response summaries`

## Verification

curl /embeddings/batch summaries show expected count/dimensions/payloads.

## Observability Impact

Validates internal batch endpoint against real stack.
