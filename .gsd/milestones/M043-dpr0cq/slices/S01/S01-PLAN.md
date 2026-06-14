# S01: Tier 1 lint adoption and fix existing issues

**Goal:** Расширить .golangci.yml с Tier 1 linters (gosec, bodyclose, prealloc, errorlint, revive). Warn mode первом проходе, fix issues в fd коде, move to fail mode. Документировать baseline noise floor в Phase 1 report.
**Demo:** After this, .golangci.yml имеет 12 linters (7 baseline + 5 Tier 1) в fail mode. golangci-lint run exit 0. CI go-quality.yml runs full lint, fails on any new issue. docs/static-analysis-phase1-report-m043.md фиксирует baseline issue count, fix list, exclusions rationale.

## Must-Haves

- .golangci.yml содержит 12 enabled linters (7 baseline + 5 Tier 1) в fail mode. golangci-lint run --config .golangci.yml ./... exit 0 на текущем main. 0 false positives без //nolint explanations. Все M041 acceptance tests (45 cases) pass. docs/static-analysis-phase1-report-m043.md фиксирует baseline, fix list, exclusions rationale.

## Proof Level

- This slice proves: integration + contract

## Integration Closure

Builds on M040 .golangci.yml baseline и M041 new code. Tier 1 linters adopted in fail mode. CI integration complete. Closes 5 quality gaps (gosec, bodyclose, prealloc, errorlint, revive).

## Verification

- No new runtime observability. CI annotations: golangci-lint shows new linter issues in PR diff. CI timeout bumped 10min → 15min for 12+ linter run.

## Tasks

- [x] **T01: Baseline noise floor captured: 11 issues (1 errorlint, 3 goconst, 5 gosec, 2 unused) перед fixes** `est:1h`
  .golangci.yml: добавить Tier 1 linters (gosec, bodyclose, prealloc, errorlint, revive) с per-linter settings (gosec: G107 warn для URL-from-env, revive: explicit rules list, errorlint: check errorf formatting). Запустить `golangci-lint run --config .golangci.yml ./...` локально. Собрать output: per-linter issue count, false positive classification, fix vs exclude decision. Это baseline для Phase 1 report.
  - Files: `.golangci.yml`
  - Verify: golangci-lint run exit code ≠ 0 (issues captured), но без падений compile. Per-linter counts recorded.

- [x] **T02: Fixed all 11 genuine issues: errorlint redis.Nil, gosec G112/G115/G304, revive package-comments + var-naming + early-return, unused consts, test fixture consts** `est:2h`
  Fix genuine issues found в T01: (a) gosec false positives — добавить //nolint:gosec с justification для legit G107 (URL from env) или fix actual issue. (b) bodyclose — добавить defer resp.Body.Close() если missing. (c) prealloc — preallocate slices с cap. (d) errorlint — заменить errors.New(fmt.Errorf(...)) на fmt.Errorf(...). (e) revive — добавить godoc на exported funcs (WriteError, HTTPStatusFor, AllErrorCodes, EncodeEmbedding, Float32SliceToBytes, BytesToFloat32Slice). Также: убрать underscore в имени parameter если revive flags (`_ = strconv.Itoa` style). Не делать refactor to fix every issue — для genuinely difficult issues (e.g., exported method without doc) добавить //nolint:revive с explicit comment.
  - Files: `api/handlers/embeddings.go`, `api/handlers/batch.go`, `api/handlers/errors.go`, `api/handlers/recovery.go`, `api/handlers/notfound.go`, `api/middleware/validation.go`, `api/embed/codec.go`, `api/main.go`
  - Verify: golangci-lint run exit 0. Каждое fix с явным commit reference в fix list. Никаких //nolint:gosec без комментария justification.

- [x] **T03: Tier 1 linters moved to fail mode (12 linters, 0 issues, CI timeout 15min) + Phase 1 report saved** `est:1h`
  .golangci.yml: move Tier 1 linters from warn to fail mode (severity: error default). Verify `golangci-lint run` exit 0. .github/workflows/go-quality.yml: bump timeout to 15 min (12+ linters медленнее). Проверить CI step не использует cached results иначе новые linter issues могут быть masked. docs/static-analysis-phase1-report-m043.md: baseline noise count, fix list, exclusions rationale, false positive rate per linter, before/after .golangci.yml diff. Также: confirm M041 acceptance tests (45 cases) все ещё pass.
  - Files: `.golangci.yml`, `.github/workflows/go-quality.yml`, `docs/static-analysis-phase1-report-m043.md`
  - Verify: golangci-lint run exit 0. CI workflow YAML valid (parseable). docs/static-analysis-phase1-report-m043.md существует, ≥2KB. go test ./api/... -short exit 0 (all M041 tests pass).

## Files Likely Touched

- .golangci.yml
- api/handlers/embeddings.go
- api/handlers/batch.go
- api/handlers/errors.go
- api/handlers/recovery.go
- api/handlers/notfound.go
- api/middleware/validation.go
- api/embed/codec.go
- api/main.go
- .github/workflows/go-quality.yml
- docs/static-analysis-phase1-report-m043.md
