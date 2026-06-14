---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T05: Document ONNX mode, legal quality gate deferral, and static-analysis constraints

Update ONNX mode docs to state FD_BACKEND=onnx is opt-in, TEI remains production default, legal quality gate is deferred/reference-only, cache namespace must isolate TEI/ONNX comparisons, and M043 static-analysis gates are mandatory for future ONNX changes. Include gosec/govulncheck notes for operator-controlled artifacts and dependency changes.

## Inputs

- `docs/static-analysis-recommendation.md`
- `docs/onnx-artifacts/README.md`

## Expected Output

- `docs/onnx-mode-m042.md`

## Verification

test -f docs/onnx-mode-m042.md && cd api && go test ./... && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./... && go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## Observability Impact

Docs explain which runtime metadata is observable and which checks still require smoke embedding requests.
