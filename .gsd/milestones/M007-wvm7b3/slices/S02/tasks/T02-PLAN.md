---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Verify workflow command parity

Check workflow syntax structurally and verify command parity against README/local quality commands.

## Inputs

- `.github/workflows/go-quality.yml`
- `README.md`

## Expected Output

- `S02 T02 summary`

## Verification

YAML parse/snippet checks pass.

## Observability Impact

Avoids divergence between local and CI gates.
