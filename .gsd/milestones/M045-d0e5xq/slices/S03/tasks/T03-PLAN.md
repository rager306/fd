---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verified fd and direct TEI runtime after local-path startup.

After TEI is healthy, run fd `/health`, fd `/ready`, fd `/v1/embeddings`, and direct TEI `/embeddings`. Confirm runtime backend/model/dimensions unchanged from fd client perspective.

## Inputs

- `benchmark-results/m045-tei-local-path-startup-proof.md`

## Expected Output

- `benchmark-results/m045-tei-local-path-startup-proof.md`

## Verification

Smoke results are captured and pass with 1024-dimensional embeddings; fd runtime still reports TEI and model deepvk/USER-bge-m3.

## Observability Impact

Confirms service recovers after local-path startup proof.
