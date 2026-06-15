---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: README updated with current test levels, e2e modes, and mutation policy.

Update README with concise testing section: regular api suite, lint/govulncheck, root integration no-key mode, authenticated Docker e2e mode, bounded mutation baseline. Do not print secrets.

## Inputs

- `benchmark-results/m050-s01-test-actuality.md`
- `benchmark-results/m050-s02-docker-e2e.md`
- `benchmark-results/m050-s03-mutation-baseline.md`

## Expected Output

- `README.md`

## Verification

README contains commands and secret-handling notes for each test level.

## Observability Impact

Future agents can run the right tests without re-deriving commands.
