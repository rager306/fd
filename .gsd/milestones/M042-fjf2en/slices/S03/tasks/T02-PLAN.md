---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Wire FD_BACKEND env selection with safe ONNX runtime config

Wire or refine FD_BACKEND=onnx selection in main.go while preserving default TEI. Treat ONNX manifest/tokenizer/runtime paths as operator-controlled env/config values with explicit validation and gosec G304 comments only where justified. Ensure runtime identity in responses/health uses authoritative backend metadata, not request model. Any new exported config/type requires godoc.

## Inputs

- `.golangci.yml`
- `docs/fd-v2.md`
- `docs/static-analysis-phase1-report-m043.md`

## Expected Output

- `api/main.go`
- `api/main_test.go`
- `api/handlers/health.go`

## Verification

cd api && go test ./... && go test -tags onnx ./... && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...

## Observability Impact

Startup logs and /health expose backend/runtime identity without secrets or full local paths.
