---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Final static gates passed after lint fixes.

Run full Go tests, golangci-lint, and govulncheck after S01/S02 commits. Fix any regressions before proceeding to runtime verification.

## Inputs

- `api/go.mod`
- `api/go.sum`

## Expected Output

- `api/**/*`

## Verification

cd api && go test ./...; cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...; cd api && go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...

## Observability Impact

Confirms codebase-wide safety before container proof.
