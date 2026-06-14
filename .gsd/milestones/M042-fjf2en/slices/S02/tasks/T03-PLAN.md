---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Add async observability header and metrics

Add X-Concurrent-Chunks response header in async mode and metrics/log surfaces that make the chunk count, miss count, and cancellation/error paths observable. Do not introduce noisy globals or unbounded labels. Ensure header/middleware interactions remain compatible with M041 S03 headers work and that new observability code passes gocritic/contextcheck.

## Inputs

- `.gsd/milestones/M041-4tw0w7/M041-4tw0w7-ROADMAP.md`
- `.golangci.yml`

## Expected Output

- `api/handlers/embeddings.go`
- `api/handlers/embeddings_integration_test.go`

## Verification

cd api && go test ./handlers && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./handlers

## Observability Impact

Adds async concurrency signal to responses and testable logs/metrics.
