---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Document CI workflow

Update README quality tooling section to mention `.github/workflows/go-quality.yml`, triggers, and remote run pending push.

## Inputs

- `README.md`
- `.github/workflows/go-quality.yml`

## Expected Output

- `README.md updated`

## Verification

README contains workflow path and CI trigger notes.

## Observability Impact

Future maintainers know CI/local parity and remote verification limits.
