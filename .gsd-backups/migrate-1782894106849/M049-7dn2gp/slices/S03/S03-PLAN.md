# S03: Solo scope closure and live verification

**Goal:** Close M049 with final gates, rebuilt-container runtime verification, requirement validation, and issue #8 closure matrix for implemented/deferred items.
**Demo:** After this slice, issue #8 has a closure matrix for implemented/deferred items and the rebuilt container passes live smoke.

## Must-Haves

- Full Go tests, lint, and govulncheck pass after all changes.
- Docker Compose api image is rebuilt and healthy.
- Live HTTP smoke proves auth-protected cache invalidation, MISS->HIT->flush->MISS, health dependencies/capacity, and metrics gauges.
- R040 and R041 are validated; R042 records solo-scope decision for AN-D/E/F.
- Closure artifact maps all issue #8 items to implemented/deferred/scoped outcomes.

## Proof Level

- This slice proves: Final command evidence, live runtime HTTP evidence, artifact UAT, milestone validation.

## Integration Closure

No GitHub issue mutation or push unless user separately authorizes it.

## Verification

- Final runtime evidence proves the new agent-facing surfaces work in the rebuilt container.

## Tasks

- [x] **T01: Final static gates passed after lint fixes.** `est:medium`
  Run full Go tests, golangci-lint, and govulncheck after S01/S02 commits. Fix any regressions before proceeding to runtime verification.
  - Files: `api/**/*`
  - Verify: cd api && go test ./...; cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...; cd api && go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...

- [x] **T02: Rebuilt the API container and passed live health/metrics/cache invalidation smoke.** `est:medium`
  Rebuild/restart the api container via Docker Compose, wait for health, and run authenticated live HTTP smoke for cache invalidation, health fields, metrics gauges, and embedding behavior.
  - Files: `benchmark-results/m049-s03-live-container-proof.md`
  - Verify: docker compose up -d --build api; live HTTP smoke script passes.

- [x] **T03: Wrote issue #8 closure matrix, validated R040-R042, and saved final UAT.** `est:medium`
  Write issue #8 closure matrix, validate R040/R041/R042, save UAT, validate and complete milestone, and commit final artifacts.
  - Files: `benchmark-results/m049-issue-8-closure.md`, `benchmark-results/m049-s03-live-container-proof.md`, `.gsd/REQUIREMENTS.md`
  - Verify: Artifact UAT and milestone validation pass.

## Files Likely Touched

- api/**/*
- benchmark-results/m049-s03-live-container-proof.md
- benchmark-results/m049-issue-8-closure.md
- .gsd/REQUIREMENTS.md
