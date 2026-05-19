---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Smoke test embeddings endpoint

Smoke test `/v1/embeddings` with single 1024d input and array 512d input, recording response shapes.

## Inputs

- `S02 T01 health`

## Expected Output

- `embedding response summaries`

## Verification

curl /v1/embeddings and jq summaries show expected dimensions and lengths.

## Observability Impact

Validates real TEI/API response shape.
