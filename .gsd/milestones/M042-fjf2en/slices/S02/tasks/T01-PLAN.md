---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Implement lint-aware async chunked orchestrator

Implement bounded parallel TEI chunk orchestration for batches larger than TEI's per-request limit. Keep production functions below gocyclo 15 by extracting small helpers for chunk planning, worker execution, ordered result assembly, and error aggregation. Propagate request context into every goroutine/call so contextcheck stays clean; ensure goroutines stop on cancellation and do not leak. Avoid exported APIs unless required; any exported helper/type needs meaningful godoc because revive:exported is now enforced.

## Inputs

- `.golangci.yml`
- `docs/static-analysis-phase2-report-m043.md`
- `docs/fd-v2.md`

## Expected Output

- `api/handlers/embeddings.go`
- `api/handlers/embeddings_integration_test.go`

## Verification

cd api && go test ./handlers ./middleware && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./handlers ./middleware

## Observability Impact

Worker errors, cancellation, and chunk count should be visible in logs/tests without high-cardinality metrics.
