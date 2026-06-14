---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T05: Run regression suite for M041 acceptance in async mode

Run the M041 acceptance/regression suite with FD_ASYNC_CHUNKS=false and true. Confirm error envelopes, validation behavior, cache-hit path, headers, and lifecycle assumptions are unchanged. Final completion evidence must include go test ./..., golangci-lint 18 linters, and govulncheck 0 reachable vulnerabilities.

## Inputs

- `.gsd/milestones/M041-4tw0w7/slices/S01/S01-SUMMARY.md`
- `docs/static-analysis-recommendation.md`

## Expected Output

- `benchmark-results/fd-v2-async-regression-m042.md`

## Verification

cd api && go test ./... && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./... && go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## Observability Impact

Regression artifact records async/sync mode and cache behavior evidence.
