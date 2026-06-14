---
id: T05
parent: S04
milestone: M041-4tw0w7
key_files:
  - tools/verify_fd_v2_perf.sh
  - docs/fd-v2.md
  - benchmark-results/fd-v2-perf-validation-m041-s04.md
  - benchmark-results/m041-s04-t05-cache-hot-runtime.txt
  - benchmark-results/m041-s04-t05-cache-hot-help.txt
  - benchmark-results/m041-s04-t05-go-test.txt
  - benchmark-results/m041-s04-t05-lint.txt
  - benchmark-results/m041-s04-t05-govulncheck.txt
key_decisions:
  - D045: T-P-1..T-P-5 are cache-hot steady-state validation after prewarm; real cache-miss TEI CPU latency is diagnostic only for M041 S04.
duration: 
verification_result: passed
completed_at: 2026-06-14T07:41:48.720Z
blocker_discovered: false
---

# T05: Validated S04 performance under the accepted D045 cache-hot steady-state contract, with real cache-miss diagnostics preserved.

**Validated S04 performance under the accepted D045 cache-hot steady-state contract, with real cache-miss diagnostics preserved.**

## What Happened

Updated `tools/verify_fd_v2_perf.sh` so T-P latency cases explicitly prewarm each measured payload through real inference and then require `X-Cache: HIT` for measured requests. Updated `docs/fd-v2.md` to clarify R-P0-6 and T-P-1..T-P-5 as warm-service/cache-hot validation after the user explicitly descoped backend remediation. Ran the verifier against current `fd_api` on localhost:8000 backed by real `fd_tei` and Redis. The verifier passed: batch=1 p95 2.236ms (<50ms), batch=10 p95 3.468ms (<200ms), batch=32 p95 7.595ms (<1000ms), 100 sequential cache-hot requests had 0 errors/non-HIT responses, 4x8 concurrent cache-hot requests completed in 0.010s, and repeated input returned `X-Cache: HIT` in 1.870ms. The same artifact keeps non-blocking cache-miss diagnostics visible: batch=1 MISS 235ms, batch=10 MISS 2107ms, batch=32 MISS 6796ms.

## Verification

Fresh runtime verification passed: `FD_BASE_URL=http://localhost:8000 tools/verify_fd_v2_perf.sh` exit 0, writing `benchmark-results/fd-v2-perf-validation-m041-s04.md`. Fresh static gates passed after the script/spec update: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. Script syntax/help checked with `bash -n` and `tools/verify_fd_v2_perf.sh --help`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `bash -n /root/fd/tools/verify_fd_v2_perf.sh && /root/fd/tools/verify_fd_v2_perf.sh --help` | 0 | ✅ pass: verifier syntax/help reflect cache-hot contract | 120000ms |
| 2 | `cd /root/fd && FD_BASE_URL=http://localhost:8000 tools/verify_fd_v2_perf.sh` | 0 | ✅ pass: cache-hot T-P validation passes against current fd with real TEI/Redis | 600000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 4 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 5 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

Original T05 wording required real cache-miss latency after backend remediation. Per explicit user decision D045, T-P targets are now cache-hot steady-state; real cache-miss latency remains diagnostic only.

## Known Issues

Real TEI CPU cache-miss inference remains slow and is documented as non-blocking diagnostic evidence. `api/report.json` remains an unrelated untracked generated file. `.gsd/.../S04-CONTINUE.md` remains an auto-compact untracked artifact not part of the task deliverable.

## Files Created/Modified

- `tools/verify_fd_v2_perf.sh`
- `docs/fd-v2.md`
- `benchmark-results/fd-v2-perf-validation-m041-s04.md`
- `benchmark-results/m041-s04-t05-cache-hot-runtime.txt`
- `benchmark-results/m041-s04-t05-cache-hot-help.txt`
- `benchmark-results/m041-s04-t05-go-test.txt`
- `benchmark-results/m041-s04-t05-lint.txt`
- `benchmark-results/m041-s04-t05-govulncheck.txt`
