# M043-dpr0cq: Static analysis quality hardening

**Vision:** Расширить fd static analysis с 7 linters (M040 baseline) до 2026 Go community standard stack: golangci-lint с 12+ linters (7 baseline + 5 Tier 1: gosec, bodyclose, prealloc, errorlint, revive) + Tier 2 (gocyclo, gocritic, durationcheck, unparam, contextcheck, nilnil) + govulncheck в CI как required step. Phased rollout (Tier 1 warn→fix→fail, потом Tier 2, потом CI govulncheck) с low-risk adoption. docs/static-analysis-recommendation.md обновляется с M043 outcome.

## Success Criteria

- golangci-lint run --config .golangci.yml ./... exit 0 с 12+ linters (7 baseline + 5 Tier 1) в fail mode
- golangci-lint run --config .golangci.yml ./... exit 0 с 18 linters (12 S01 + 6 Tier 2) в fail mode
- govulncheck ./... exit 0 в CI как required step
- All M041 acceptance tests (45 cases + 10 behavior scenarios) pass после всех refactors
- M042 S02 async code (when shipped) passes expanded lint set from day 1 (no follow-up cleanup)
- docs/static-analysis-recommendation.md имеет M043 outcome section: Tier 1+2 active, Tier 3 opt-in, future work
- Final .golangci.yml consolidated, все exclusions documented
- golangci-lint pass в обоих режимах (TEI default, ONNX opt-in) per M042

## Slices

- [x] **S01: Tier 1 lint adoption and fix existing issues** `risk:low` `depends:[]`
  > After this: After this, .golangci.yml имеет 12 linters (7 baseline + 5 Tier 1) в fail mode. golangci-lint run exit 0. CI go-quality.yml runs full lint, fails on any new issue. docs/static-analysis-phase1-report-m043.md фиксирует baseline issue count, fix list, exclusions rationale.

- [x] **S02: Tier 2 lint adoption and complexity refactor** `risk:medium` `depends:[S01]`
  > After this: After this, .golangci.yml имеет 18 linters (12 from S01 + 6 Tier 2: gocyclo, gocritic, durationcheck, unparam, contextcheck, nilnil). gocyclo.min-complexity=15. golangci-lint run exit 0. CreateEmbedding handler refactored если gocyclo exceeds threshold. docs/static-analysis-phase2-report-m043.md с before/after complexity metrics.

- [ ] **S03: govulncheck CI integration and docs finalization** `risk:low` `depends:[S01,S02]`
  > After this: After this, .github/workflows/go-quality.yml имеет step "Run govulncheck" который fails on known vulnerabilities в dependencies или stdlib usage. docs/static-analysis-recommendation.md обновлён с M043 outcome: Tier 1+2 implemented, Tier 3 opt-in, future work. Final artifact: docs/static-analysis-recommendation.md с полной M043 section.

## Boundary Map

### S01 → S02

Produces:
- Tier 1 baseline (12 linters, fail mode, 0 issues)
- 5 security/quality gaps closed
- Phase 1 report with exclusion rationale

Consumes:
- M040 .golangci.yml baseline (7 linters)
- M041 new code (handlers, middleware, codec)

### S02 → S03

Produces:
- Tier 2 baseline (18 linters, fail mode, 0 issues)
- gocyclo threshold established
- Complexity refactor documented

Consumes:
- S01 Tier 1 baseline
- Phase 1 exclusion rationale (don't re-introduce)

### S01 → S03

Produces:
- Dependency vuln scan (0 vulns baseline)
- CI step defined

Consumes:
- S01+S02 lint stack hardened (no regression in security baseline)
- M041 new code passes lints (no exclusions specifically for M041 to worry about)

### S02 → S03

Produces:
- Complexity hotspots documented (gocritic feedback)
- Phase 2 report (refactor + exclusions)

Consumes:
- S02 complexity metrics (for documentation)
