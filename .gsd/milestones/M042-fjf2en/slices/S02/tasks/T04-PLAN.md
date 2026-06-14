---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T04: Benchmark async versus sync cold and warm paths

Create or update perf benchmark tooling to compare FD_ASYNC_CHUNKS=false versus true for cold batch=32, cold batch=128, and cache-hit path. Record before/after numbers in benchmark-results/fd-v2-async-perf-m042.md. Include the expanded M043 gate in the verification transcript so perf work cannot regress lint/test/security gates.

## Inputs

- `docs/fd-v2.md`
- `benchmark-results/m043-s03-final-lint.txt`
- `benchmark-results/m043-s03-govulncheck-final.txt`

## Expected Output

- `tools/verify_fd_async_perf.sh`
- `benchmark-results/fd-v2-async-perf-m042.md`

## Verification

cd api && go test ./... && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./... && go run golang.org/x/vuln/cmd/govulncheck@latest ./... && ../tools/verify_fd_async_perf.sh

## Observability Impact

Benchmark artifact records async header/metrics evidence alongside latency.
