---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T04: Run regression suite for M041 and M042 async in ONNX mode

Run M041 acceptance tests and M042 S02 async/sync regression checks with FD_BACKEND=onnx where build/runtime prerequisites are available. Validate /health runtime metadata but do not treat HTTP 200 /health as inference readiness; include a smoke embedding request. If S02 async regression artifact exists, consume it; otherwise record dependency status. Final evidence must include 18-linter gate and govulncheck.

## Inputs

- `.gsd/milestones/M041-4tw0w7/slices/S01/S01-SUMMARY.md`
- `docs/static-analysis-recommendation.md`

## Expected Output

- `benchmark-results/fd-v2-onnx-regression-m042.md`

## Verification

cd api && go test ./... && go test -tags onnx ./... && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./... && go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## Observability Impact

Regression artifact records backend, smoke embedding readiness, and health metadata separately.
