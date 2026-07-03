---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T02: Run final CI workflow verification

Run final local verification: workflow parse, README snippets, Go tests, GolangCI-Lint, and GitNexus change detection.

## Inputs

- `.github/workflows/go-quality.yml`
- `README.md`

## Expected Output

- `S03 T02 summary`

## Verification

All local checks pass.

## Observability Impact

Final local proof before commit.
