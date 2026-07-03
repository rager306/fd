---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T03: Ran final M047 gates, wrote issue #6 closure matrix, and validated R034.

Write S04 and milestone closure artifacts, validate R034, run full tests, race cache smoke if relevant, lint, govulncheck, UAT checks, milestone validation, and completion.

## Inputs

- `benchmark-results/m047-s01-contract-cleanup.md`
- `benchmark-results/m047-s02-graceful-listener-shutdown.md`
- `benchmark-results/m047-s03-tei-retry-fast-fail.md`

## Expected Output

- `benchmark-results/m047-s04-warmup-retry-closure.md`
- `benchmark-results/m047-issue-6-closure.md`

## Verification

cd api && go test ./...; golangci-lint; govulncheck; artifact UAT; milestone validation.

## Observability Impact

Closure matrix records all issue #6 findings and evidence IDs.
