---
id: T08
parent: S05
milestone: M041-4tw0w7
key_files:
  - tools/verify_fd_v2_contract.py
  - benchmark-results/fd-v2-validation-m041.md
  - benchmark-results/m041-s05-t08-contract-run.txt
  - benchmark-results/m041-s05-t08-go-test.txt
  - benchmark-results/m041-s05-t08-lint.txt
  - benchmark-results/m041-s05-t08-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T08:33:53.287Z
blocker_discovered: false
---

# T08: Added and ran the final M041 45-check fd v2 contract acceptance suite, passing 45/45 against current running fd.

**Added and ran the final M041 45-check fd v2 contract acceptance suite, passing 45/45 against current running fd.**

## What Happened

Implemented `tools/verify_fd_v2_contract.py`, a standard-library Python black-box verifier for the fd v2 HTTP contract. The suite runs 45 checks against a running fd service covering probes, health, version/info/metrics/warmup, OpenAPI/docs, traces, embeddings success paths, headers, cache HIT/MISS, ETag/If-None-Match, CORS preflight, validation/error envelopes, `/v1/batch`, legacy batch, OpenAPI path coverage, and cache-hot performance p95 checks. Rebuilt `fd_api` from the current working tree, waited for warmup, and ran the verifier against `http://localhost:8000`; it wrote `benchmark-results/fd-v2-validation-m041.md` and passed 45/45.

## Verification

Runtime acceptance passed: `tools/verify_fd_v2_contract.py http://localhost:8000` exit 0, 45/45 checks passed. Fresh full M043 gate passed after the script: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. Requirements R017, R018, and R019 were updated to validated with this evidence.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd && tools/verify_fd_v2_contract.py http://localhost:8000` | 0 | ✅ pass: 45/45 fd v2 contract checks passed | 900000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 4 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

The suite is a pragmatic black-box acceptance suite for implemented M041 surfaces; external OpenAPI validator evidence remains in T07. Auth with FD_API_KEY is covered by unit evidence rather than toggling the running acceptance service env mid-run.

## Known Issues

`api/report.json` and `.gsd/.../S04-CONTINUE.md` remain unrelated untracked files. Real cache-miss TEI CPU latency remains diagnostic per D045; T08 cache-hot checks follow the accepted contract.

## Files Created/Modified

- `tools/verify_fd_v2_contract.py`
- `benchmark-results/fd-v2-validation-m041.md`
- `benchmark-results/m041-s05-t08-contract-run.txt`
- `benchmark-results/m041-s05-t08-go-test.txt`
- `benchmark-results/m041-s05-t08-lint.txt`
- `benchmark-results/m041-s05-t08-govulncheck.txt`
