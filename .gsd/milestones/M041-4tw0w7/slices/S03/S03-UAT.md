# S03: Observability surface endpoints headers and deep health — UAT

**Milestone:** M041-4tw0w7
**Written:** 2026-06-14T06:57:27.537Z

# UAT — M041 S03 Observability surface endpoints headers and deep health

## Verdict
PASS

## Checks

- [x] `/version`, `/info`, `/metrics`, `/v1/healthcheck`, and `/warmup` exist and return expected status/content type.
  - Evidence: `benchmark-results/m041-s03-t07-observability-integration.txt`.
- [x] `/health`, `/live`, `/ready`, and `/version` health/version checks pass.
  - Evidence: `benchmark-results/m041-s03-t07-observability-integration.txt`.
- [x] Deep `/health` reports lifecycle fields and `last_inference_at` after successful `/v1/embeddings`.
  - Evidence: `benchmark-results/m041-s03-t03-deep-health.txt` and `benchmark-results/m041-s03-t07-observability-integration.txt`.
- [x] Prometheus `/metrics` exposes request, duration, batch, error, model-loaded, and cache hit/miss metrics in text/plain.
  - Evidence: `benchmark-results/m041-s03-t04-metrics.txt`.
- [x] Headers are present: `Server`, `X-Request-Id` echo/generation, `X-Model-Id`, `X-Dimensions`, `Retry-After`, `Connection`.
  - Evidence: `benchmark-results/m041-s03-t05-headers.txt` and `benchmark-results/m041-s03-t07-observability-integration.txt`.
- [x] Static quality gates: full `go test ./...`, golangci-lint v2.12.2, and govulncheck pass with 0 reachable vulnerabilities.
  - Evidence: `benchmark-results/m041-s03-t07-go-test.txt`, `benchmark-results/m041-s03-t07-lint.txt`, `benchmark-results/m041-s03-t07-govulncheck.txt`.

## Notes
Root-level `tests/integration` remains non-executable with the current Go module layout, so executable integration tests live under `api/` and run with `cd api && go test . -run TestFdV2Observability -v`.

