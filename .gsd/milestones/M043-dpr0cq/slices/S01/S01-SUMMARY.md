---
id: S01
parent: M043-dpr0cq
milestone: M043-dpr0cq
provides:
  - (none)
requires:
  []
affects:
  []
key_files:
  - .golangci.yml (12 linters, fail mode, per-linter settings, exclude-rules, package-comments documentation)
  - .github/workflows/go-quality.yml (timeout 15min, step name + description)
  - docs/static-analysis-phase1-report-m043.md (10KB before/after report)
  - benchmark-results/m043-tier1-baseline.txt (raw lint output)
  - api/cache/redis.go (errors.Is, maxUint16, errors import)
  - api/cache/local.go (package comment)
  - api/embed/codec.go (package comment)
  - api/embed/onnx_manifest.go (2x //nolint:gosec G304)
  - api/handlers/batch.go (early-return, package comment)
  - api/handlers/constants.go (deleted, empty file)
  - api/handlers/errors_test.go (paramInput, paramDimensions, paramEncodingFormat consts)
  - api/main.go (package comment, defaultValue, ReadHeaderTimeout, //nolint:gosec G304)
  - api/middleware/validation.go (deleted unused teiSubBatchSize const)
  - api/middleware/validation_test.go (paramInput const)
key_decisions: []
patterns_established:
  - (none)
observability_surfaces:
  - none
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-14T03:52:16.491Z
blocker_discovered: false
---

# S01: Tier 1 lint adoption and fix existing issues

**12 linters в fail mode (7 baseline + 5 Tier 1: gosec, bodyclose, prealloc, errorlint, revive). 11 baseline issues → 0 (16 fix changes). CI timeout 10→15min. Phase 1 report saved.**

## What Happened

S01 complete. Tier 1 linters adopted: gosec (security), bodyclose (HTTP leaks), prealloc (slice allocations), errorlint (%w wrapping), revive (golint successor, 19 rules; `exported` deferred to S02). 11 baseline issues fixed across 16 changes: real bugs (errors.Is redis.Nil, G112 ReadHeaderTimeout, G115 maxUint16 bounds), suppressions with justification (G304 for env-controlled paths), style refactors (early-return, default_ → defaultValue, package comments), dead code removal (errorKey, teiSubBatchSize consts), test fixture consts (paramInput, paramDimensions, paramEncodingFormat). .github/workflows/go-quality.yml timeout 10→15min for 12-linter run. docs/static-analysis-phase1-report-m043.md (10KB) documents все changes с rationale. Out-of-phase items: revive:exported → S02, govulncheck → S03, Tier 2 → S02, Tier 3 opt-in.

## Verification

Live local reproduction: `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` exits 0, 0 issues. go test ./api/... all 5 packages pass. docs/static-analysis-phase1-report-m043.md (10068 bytes) с before/after metrics, fix list (16 changes), exclusions rationale, out-of-phase items. .github/workflows/go-quality.yml timeout bumped 10→15min.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

User requested Phase 1 manual execution with focus on Tier 1 (security/quality gaps). Three intentional deferrals: (1) revive:exported rule excluded due to 50+ missing godocs (deferred to S02 godoc pass), (2) goconst path-pattern в golangci-lint v2.12.2 не работал — workaround через consts в test files, (3) gosec G107/G304 excluded at config level instead of per-call //nolint (3 G304 spots still per-call because of explicit env-source justification; G107 fully env-var-driven).

## Known Limitations

None.

## Follow-ups

None.

## Files Created/Modified

None.
