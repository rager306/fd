---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Wire FD_ASYNC_CHUNKS env into handler and main config

Add FD_ASYNC_CHUNKS configuration with default false and wire it into EmbeddingsHandler without changing default TEI behavior. Keep config parsing small and testable; use contextual errors for invalid env values if any. Ensure the handler path is explicit enough for future agents to see sync vs async behavior, and keep public config comments/godoc aligned with revive:exported requirements if new exported symbols are introduced.

## Inputs

- `.golangci.yml`
- `docs/static-analysis-recommendation.md`

## Expected Output

- `api/main.go`
- `api/handlers/embeddings.go`
- `api/main_test.go`

## Verification

cd api && go test ./... && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...

## Observability Impact

Config should be visible in startup logs without logging secrets.
