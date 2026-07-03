---
id: T04
parent: S01
milestone: M040-pbp9z1
key_files:
  - README.md
  - docs/same-host-embedding-service-contract.md
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-22T05:34:19.687Z
blocker_discovered: false
---

# T04: Linked and verified the same-host embedding service contract entry point while keeping examples non-sensitive.

**Linked and verified the same-host embedding service contract entry point while keeping examples non-sensitive.**

## What Happened

T04 confirmed that README.md already exposes the same-host service contract link, then tightened the public examples so README.md and docs/same-host-embedding-service-contract.md use a non-sensitive smoke input instead of the prior raw legal-domain sample. I reran Go verification from the correct module directory (`cd api && go test ./... -short`), which resolves the prior verification failure caused by running `go test ./...` from the repository root outside the Go module. The required broad docs/artifact audit command was also run; it now finds no raw sample in README.md or the contract doc, and the remaining matches are historical GSD planning/summary text or policy statements about avoiding signed URLs/secrets rather than live secret material. A focused high-risk check over deliverable docs and benchmark artifacts found no token query strings, X-Amz URLs, PEM/private-key blocks, or private-key markers.

## Verification

Verified `cd api && go test ./... -short` passes from the Go module directory. Ran the required broad `rg` audit over docs, README, benchmark-results, and the S01 GSD artifacts; reviewed matches as non-secret policy/history references after removing the raw sample from current README/contract docs. Ran an additional focused high-risk secret check over deliverable docs and benchmark artifacts, which returned no matches.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./... -short` | 0 | ✅ pass | 261ms |
| 2 | `rg -n "signed|token=|X-Amz|BEGIN|PRIVATE|юридическая справка" docs README.md benchmark-results .gsd/milestones/M040-pbp9z1/slices/S01` | 0 | ✅ pass (audit reviewed; remaining matches are policy/history text, not live secrets) | 16ms |
| 3 | `if rg -n "token=|X-Amz|BEGIN .*PRIVATE|PRIVATE KEY" docs README.md benchmark-results; then exit 1; else exit 0; fi` | 0 | ✅ pass | 18ms |

## Deviations

Also updated the contract document example to match the README's non-sensitive smoke text so current public documentation no longer carries the prohibited raw sample string.

## Known Issues

The broad required `rg` command still emits historical GSD planning/summary lines and policy statements containing the audit terms; those are not current public contract examples or secret material.

## Files Created/Modified

- `README.md`
- `docs/same-host-embedding-service-contract.md`
