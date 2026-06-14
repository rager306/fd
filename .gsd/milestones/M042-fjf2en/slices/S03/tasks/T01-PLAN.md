---
estimated_steps: 1
estimated_files: 5
skills_used: []
---

# T01: Audit ONNX implementation and build matrix under M043 gates

Audit existing ONNX implementation, build tags, native libraries, manifest validation, tokenizer/runtime paths, and current tests before changing behavior. Identify any M043 lint/security risks: exported API godoc, gosec G304 operator-path handling, context propagation, gocyclo hotspots, and govulncheck dependency exposure. Record whether the GitNexus/index view is stale before edits if code changes are needed.

## Inputs

- `docs/static-analysis-recommendation.md`
- `docs/onnx-artifacts/README.md`
- `.golangci.yml`

## Expected Output

- `benchmark-results/fd-v2-onnx-audit-m042.md`

## Verification

cd api && go test ./embed && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./embed

## Observability Impact

Audit records which runtime metadata is safe to expose and which checks are startup-only.
