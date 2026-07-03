---
id: T02
parent: S03
milestone: M050-rfqm1p
key_files:
  - benchmark-results/m050-s03-mutation-baseline.md
key_decisions:
  - Do not make mutation a mandatory CI hard gate yet; use local/manual bounded baseline until runner/toolchain cost is pinned.
duration: 
verification_result: passed
completed_at: 2026-06-15T14:54:36.925Z
blocker_discovered: false
---

# T02: Bounded mutation baseline прошёл на критичных cache, handlers и lifecycle файлах.

**Bounded mutation baseline прошёл на критичных cache, handlers и lifecycle файлах.**

## What Happened

Mutation baseline запущен на `cache/hash.go`, `cache/keys.go`, `handlers/cache.go`, `handlers/health.go`, `lifecycle/state.go` с per-mutant command `go test ./cache ./handlers ./lifecycle`. Runner reported score 1.000000 with 143 killed mutants, 0 survivors, 4 duplicates, 0 skipped.

## Verification

Critical baseline command exited 0. Summary: `The mutation score is 1.000000 (143 passed, 0 failed, 4 duplicated, 0 skipped, total is 143)`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go run github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest --exec 'go test ./cache ./handlers ./lifecycle' --exec-timeout 45 ./cache/hash.go ./cache/keys.go ./handlers/cache.go ./handlers/health.go ./lifecycle/state.go` | 0 | ✅ pass: mutation score 1.000000 (143 killed, 0 survived) | 50600ms |

## Deviations

None.

## Known Issues

Full repo mutation remains out of scope; baseline is bounded and informational.

## Files Created/Modified

- `benchmark-results/m050-s03-mutation-baseline.md`
