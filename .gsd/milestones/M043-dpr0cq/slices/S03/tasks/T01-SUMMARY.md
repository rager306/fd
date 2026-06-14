---
id: T01
parent: S03
milestone: M043-dpr0cq
key_files:
  - benchmark-results/m043-s03-govulncheck-baseline.txt
  - benchmark-results/m043-govulncheck-baseline.txt
  - benchmark-results/m043-s03-govulncheck-final.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T05:14:38.734Z
blocker_discovered: false
---

# T01: Initial govulncheck scan completed: 0 reachable vulnerabilities, exit 0

**Initial govulncheck scan completed: 0 reachable vulnerabilities, exit 0**

## What Happened

Ran standalone `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` from api module. The scan found no reachable vulnerabilities. govulncheck also reported 2 vulnerabilities in imported packages and 19 in required modules that fd does not appear to call; this is passing/expected govulncheck behavior because the vulnerable symbols are not reachable from fd code. Saved outputs to benchmark-results/m043-s03-govulncheck-baseline.txt, benchmark-results/m043-govulncheck-baseline.txt, and final rerun benchmark-results/m043-s03-govulncheck-final.txt.

## Verification

`cd /root/fd/api && go run golang.org/x/vuln/cmd/govulncheck@latest ./...` → No vulnerabilities found. Your code is affected by 0 vulnerabilities. govulncheck_exit=0.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `benchmark-results/m043-s03-govulncheck-baseline.txt`
- `benchmark-results/m043-govulncheck-baseline.txt`
- `benchmark-results/m043-s03-govulncheck-final.txt`
