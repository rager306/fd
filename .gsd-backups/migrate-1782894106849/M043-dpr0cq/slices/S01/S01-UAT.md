# S01: Tier 1 lint adoption and fix existing issues — UAT

**Milestone:** M043-dpr0cq
**Written:** 2026-06-14T03:52:16.491Z

# S01 UAT: Tier 1 lint adoption

## Test method

Local reproduction matches CI: `cd /root/fd/api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...`

## Results

| Check | Expected | Actual | Status |
|---|---|---|---|
| `golangci-lint run` exit code | 0 | 0 | **PASS** |
| Issue count | 0 | 0 | **PASS** |
| Linters enabled | 12 (7 baseline + 5 Tier 1) | 12 | **PASS** |
| Linter config: gosec | G107 + G304 excluded (env-controlled) | excluded | **PASS** |
| Linter config: errorlint | `checks: all` | set | **PASS** |
| Linter config: revive | 19 rules; `exported` excluded | configured | **PASS** |
| Linter config: goconst | test files excluded | consts extracted in test files (path-pattern не работал в golangci-lint v2.12.2) | **PASS** (functional) |
| CI: `.github/workflows/go-quality.yml` timeout | 15 minutes | 15 | **PASS** |
| Unit tests | все 5 packages pass | ok fd-api, fd-api/cache, fd-api/embed, fd-api/handlers, fd-api/middleware | **PASS** |

## Bug closure

| Original bug | Status |
|---|---|
| B11 `err == redis.Nil` (loses wrap) | **FIXED** → `errors.Is(err, redis.Nil)` (errorlint) |
| B12 `http.Server` без `ReadHeaderTimeout` (Slowloris) | **FIXED** → 10s timeout (gosec G112) |
| B13 `uint16(dim)` без bounds check (silent overflow potential) | **FIXED** → explicit `dim > maxUint16` + `//nolint:gosec` G115 |
| B14 G304 file inclusion через variable (3 spots) | **FIXED via suppression** → `//nolint:gosec` G304 + config-level exclusion (env-controlled paths) |
| B15 `default_` underscore naming | **FIXED** → `defaultValue` (revive var-naming) |
| B16 `if/else` early-return anti-pattern | **FIXED** → inverted (revive early-return) |
| B17 missing package comments (4 packages) | **FIXED** → added (revive package-comments) |
| B18 dead consts (`errorKey`, `teiSubBatchSize`) | **FIXED** → deleted (unused) |
| B19 goconst test fixtures (3 strings) | **FIXED** → extracted paramInput, paramDimensions, paramEncodingFormat consts |

## Out-of-phase (carries to S02/S03)

- **`revive:exported` rule** (50+ godoc comments needed): deferred to M043 S02.
- **govulncheck** (standalone tool): M043 S03.
- **Tier 2 linters** (gocyclo, gocritic, durationcheck, unparam, contextcheck, nilnil): M043 S02.
- **Tier 3 linters** (gofumpt, structslop, maligned, aligncheck, dupl, nakedret, wsl, goimports, lll): opt-in only.

## Known limitations

- **Path-pattern for goconst in golangci-lint v2.12.2** did not work; we worked around by extracting consts in test files. If golangci-lint v2.13+ fixes this, we can revert to exclude-rules.
- **Per-call `//nolint`** added in 3 places (G304 on env-controlled paths). Future code review should check that new G304 candidates get the same suppression pattern.

## S01 verdict

**PASS.** 12 linters active in fail mode with 0 issues. All M041 acceptance tests pass. CI gate defined. Phase 1 report documents all changes with rationale.
