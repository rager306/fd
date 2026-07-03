---
id: T01
parent: S04
milestone: M050-rfqm1p
key_files:
  - README.md
  - benchmark-results/m050-s04-test-gates-closure.md
key_decisions:
  - Document heavy Docker e2e and mutation as manual/local gates until CI runtime and secret handling are explicitly configured.
duration: 
verification_result: passed
completed_at: 2026-06-15T14:57:21.509Z
blocker_discovered: false
---

# T01: README updated with current test levels, e2e modes, and mutation policy.

**README updated with current test levels, e2e modes, and mutation policy.**

## What Happened

Development documentation now lists regular API tests, short CI tests, lint, govulncheck, Docker Compose startup, standalone integration no-key mode, authenticated Docker e2e mode with `FD_INTEGRATION_API_KEY`, and bounded mutation baseline command. It also records that mutation is local/manual informational, not mandatory CI.

## Verification

Final commands listed in the README were run in S04 and passed; closure artifact `benchmark-results/m050-s04-test-gates-closure.md` records results.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./...` | 0 | ✅ pass | 9200ms |
| 2 | `cd api && go test ./... -short` | 0 | ✅ pass | 9100ms |
| 3 | `cd tests/integration && go test -v .` | 0 | ✅ pass | 9100ms |

## Deviations

None.

## Known Issues

Authenticated e2e and mutation remain manual/local policy, not CI hard gates.

## Files Created/Modified

- `README.md`
- `benchmark-results/m050-s04-test-gates-closure.md`
