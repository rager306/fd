# M041 Post-M043 Roadmap Assessment

## Verdict

Roadmap adjusted, not structurally replanned.

## Reason

M043 completed after M041 S01 and changed the project-wide definition of done: all future Go work must pass 18-linter golangci-lint, `go test ./...`, and govulncheck with 0 reachable vulnerabilities. M041's remaining feature slices are still valid; the change is a verification contract and acceptance-gate update, not a scope rewrite.

## Adjustments

- S02 lifecycle work must pass M043 gates and keep context propagation clean.
- S03 observability endpoints/headers must include godoc for exported APIs and avoid gocyclo/gocritic issues.
- S04 cache/perf work must avoid new static-analysis suppressions without justification.
- S05 should not reimplement `encoding_format`; M041 S01 already implemented part of that surface. Remaining S05 work should focus on `user`, `priority`, auth/CORS/OpenAPI/rate-limit/traces and compatibility cleanup.

## Required gate for remaining slices

```bash
cd api
go test ./...
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```
