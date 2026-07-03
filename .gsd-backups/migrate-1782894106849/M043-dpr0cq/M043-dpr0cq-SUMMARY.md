---
id: M043-dpr0cq
title: "Static analysis quality hardening"
status: complete
completed_at: 2026-06-14T05:16:27.503Z
key_decisions:
  - Keep govulncheck standalone rather than trying to fold it into golangci-lint.
  - Keep Tier 3 style/deep linters opt-in/out of scope for M043.
  - Use one justified test-only gocyclo suppression for a dense table-driven integration matrix while refactoring production complexity.
key_files:
  - .golangci.yml
  - .github/workflows/go-quality.yml
  - docs/static-analysis-recommendation.md
  - docs/static-analysis-phase1-report-m043.md
  - docs/static-analysis-phase2-report-m043.md
  - benchmark-results/m043-s03-final-lint.txt
  - benchmark-results/m043-s03-go-test.txt
  - benchmark-results/m043-s03-govulncheck-final.txt
lessons_learned:
  - govulncheck can report non-reachable imported/required vulnerabilities while still exiting 0; document reachable vs non-reachable clearly.
  - gocyclo is valuable for production handler refactors but should be used carefully on dense table-driven tests.
  - revive:exported creates initial doc churn but becomes useful once the public API surface has meaningful comments.
---

# M043-dpr0cq: Static analysis quality hardening

**fd static analysis hardened from 7 linters to 18 linters in fail mode plus standalone govulncheck CI; all lint/test/vuln gates pass.**

## What Happened

M043 completed across three slices. S01 added Tier 1 linters (gosec, bodyclose, prealloc, errorlint, revive), fixed 11 baseline issues, and documented Phase 1. S02 enabled revive:exported plus Tier 2 linters (gocyclo, gocritic, durationcheck, unparam, contextcheck, nilnil), fixed 44 godoc gaps and 17 Tier 2 issues, and refactored production complexity in ValidateArtifact, CreateEmbedding, ValidateEmbeddingsRequest, and main. S03 added standalone govulncheck to CI, confirmed 0 reachable vulnerabilities locally, and finalized docs/static-analysis-recommendation.md with M043 outcome, false-positive/noise notes, suppressions, and future work. Requirements R023-R025 were validated. A local commit was created before S03 at user request: 3ca1cdd `chore: harden static analysis gates`; S03 changes remain uncommitted unless the user requests a follow-up commit.

## Success Criteria Results

- Tier 1 fail mode: PASS (S01, 12 linters, 0 issues).
- Tier 2 fail mode: PASS (S02/S03, 18 linters, 0 issues).
- govulncheck CI: PASS (workflow step added; local govulncheck 0 reachable vulnerabilities).
- M041 regression protection: PASS (`go test ./...` all packages ok).
- M042 async future gate: PARTIAL/N/A because async code has not shipped; expanded lint gate is now installed.
- Docs updated: PASS (`docs/static-analysis-recommendation.md`).

## Definition of Done Results

- All planned slices complete: PASS.
- Final lint verification: PASS (`benchmark-results/m043-s03-final-lint.txt`).
- Final tests: PASS (`benchmark-results/m043-s03-go-test.txt`).
- Final vulnerability scan: PASS (`benchmark-results/m043-s03-govulncheck-final.txt`).
- Milestone validation: PASS (`M043-dpr0cq-VALIDATION.md`).

## Requirement Outcomes

- R023: validated — Tier 1 linters + fixes + CI.
- R024: validated — Tier 2 linters + fixes.
- R025: validated — govulncheck CI + docs finalization.

## Deviations

The local commit before S03 captured a broad verified working snapshot including accumulated GSD artifacts. `api/report.json` was excluded and remains untracked. Ruby YAML parser was unavailable; PyYAML was used instead.

## Follow-ups

S03 changes were performed after commit 3ca1cdd and remain uncommitted. If desired, create a follow-up commit for govulncheck CI + docs + milestone closure artifacts. Leave `api/report.json` untracked/generated or add it to .gitignore if it recurs.
