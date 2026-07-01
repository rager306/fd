---
milestone_id: M043-dpr0cq
title: Static analysis quality hardening
status: ready-for-planning
predecessors:
  - M040-pbp9z1 (golangci-lint baseline of 7 linters)
  - M041-4tw0w7 (new code in handlers/, middleware/, embed/codec.go)
  - M042-fjf2en (planned: async orchestrator in embed/async.go)
source: docs/static-analysis-recommendation.md (10.4KB analysis, 2026-06-13)
gathered: 2026-06-13
---

# M043-dpr0cq: Static analysis quality hardening

## Source

`docs/static-analysis-recommendation.md` (saved 2026-06-13, 10.4KB) analyzes fd's current 7-linter golangci-lint baseline against the curated `dpolivaev/static-analysis` Go section (~50+ tools) and 2026 Go community consensus. The document maps each candidate tool to fd's actual code patterns and ranks them into three tiers by value-to-noise ratio.

M041 (fd v2 validation+observability) and M042 (TEI perf investigation) both introduce new Go code that should be checked by the expanded lint set. M043 closes the gap between M040's conservative baseline and the de-facto 2026 standard (golangci-lint + staticcheck + targeted security/optim linters + govulncheck).

## Project Description

`fd` is a Go embedding API service with the following linter coverage (`.golangci.yml` v2, set in M040):
- `errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`, `goconst`, `misspell`

The repo (157 LOC middleware, 609 LOC handlers, 850 LOC embed, 527 LOC cache) is small but growing. The new code from M041 (validation middleware, error envelope, recovery, OpenAI-style 17-code registry, encoding_format codec) and the planned M042 S02 async orchestrator (~150 LOC chunked parallel calls) increase the surface for static analysis issues.

M043 hardens the static analysis stack in three phases:
1. **Tier 1** — high value, low risk: `gosec`, `bodyclose`, `prealloc`, `errorlint`, `revive`. Plus govulncheck in CI (always-on security).
2. **Tier 2** — medium value: `gocyclo`, `gocritic`, `durationcheck`, `unparam`, `contextcheck`, `nilnil`.
3. **Tier 3** — opt-in style/deep: `gofumpt`, `structslop`/`maligned`, `dupl`, `nakedret`, `wsl`, `goimports`, `lll`. Defer unless pprof or style audit demands.

Production default remains: lint passes in CI on every push and PR.

## Why This Milestone

The current 7-linter baseline is conservative and predates M041. Without M043:
- **Security**: no `gosec` means G107 (URL from variable), G110 (HTTP server without timeouts), weak-crypto issues aren't caught at lint time.
- **Resource leaks**: no `bodyclose` means future HTTP client code in fd could leak `resp.Body`.
- **Performance**: no `prealloc` means M042 S02 async code could regress to `make([]T, 0)` + `append` patterns.
- **Code-as-documentation**: no `revive` means exported functions/methods from M041 (`WriteError`, `WriteErrorWithRetryAfter`, `HTTPStatusFor`, `AllErrorCodes`, `EncodeEmbedding`, `Float32SliceToBytes`, `BytesToFloat32Slice`) lack enforced godoc comments.
- **Dependencies**: no `govulncheck` means known vulnerabilities in `gin`, `redis/go-redis`, `onnxruntime` (etc) aren't caught at CI time.

M043 closes all five gaps in a phased, low-risk rollout. The first phase is conservative (warn mode, then fix, then fail mode) so that M041's recently-added code can be reviewed in isolation.

## User-Visible Outcome

### When this milestone is complete, the user can:

- Run `golangci-lint run --config .golangci.yml ./...` and see 12+ linters (7 baseline + 5 Tier 1) all active in fail mode with zero false positives.
- See `govulncheck ./...` run as a required CI step on every push and PR; it fails on any known vulnerability in dependencies or stdlib usage.
- Open `docs/static-analysis-recommendation.md` and see the final M043 outcome section: which linters are active, which are deferred, what the noise floor looks like, and the rollout plan for future tiers.
- Use the 2026 Go community standard stack without re-deriving it: golangci-lint + staticcheck + revive + gosec + bodyclose + prealloc + errorlint + govulncheck (CI). This is what Kubernetes, Prometheus, and most production Go services use.

### Entry point / environment

- Entry point: `.golangci.yml` (existing), `.github/workflows/go-quality.yml` (existing), `api/` (Go module).
- Environment: local Go 1.22.2 (matches CI's 1.25.x with version skew already present), no new tooling deps (golangci-lint v2.12.2 already bundles Tier 1 linters).
- Live dependencies: `govulncheck` requires `go install golang.org/x/vuln/cmd/govulncheck@latest` (CI installs via `go run`).

## Current Architecture (lint-relevant)

```
.github/workflows/go-quality.yml
  └─ Run Go tests (go test ./... -short)
  └─ Verify ONNX artifact contract
  └─ Ensure native ONNX binaries are not tracked
  └─ Run GolangCI-Lint with Staticcheck (golangci-lint v2.12.2 run --config ../.golangci.yml ./...)
```

`.golangci.yml` v2 currently enables 7 linters (no per-linter settings, no exclusions, no presets beyond defaults).

`go.mod` has zero dev-tooling deps; all lint tooling is in the `golangci-lint` binary and in CI's `go run` invocations.

CI Go version: 1.25.x (from `actions/setup-go`). Local development Go version: 1.22.2 (from `go version`). Skew is acceptable for golangci-lint v2.12.2 which supports 1.18+.

## Completion Class

- **Phase 1 complete**: 5 Tier 1 linters active in fail mode, all 0 issues in existing fd code, govulncheck integrated in CI.
- **Phase 2 complete**: 6 Tier 2 linters active in fail mode, complexity threshold set (gocyclo ≤ 15-20), M041 CreateEmbedding handler refactored if gocyclo threshold exceeded.
- **Phase 3 complete**: docs/static-analysis-recommendation.md updated with M043 outcome, deferred items documented, Tier 3 explicitly opt-in (not auto-enabled).

## Final Integrated Acceptance

To call this milestone complete, we must prove:

- `golangci-lint run --config .golangci.yml ./...` exits 0 with all 12+ linters active in fail mode.
- `govulncheck ./...` exits 0 (no known vulnerabilities) on a clean main branch; integrated as a required CI step.
- All M041 acceptance tests (45 test cases + 10 behavior scenarios) still pass — no lint-induced refactor breaks existing tests.
- M042 S02 async code (when shipped) passes the expanded lint set from day 1 (no follow-up cleanup needed).
- docs/static-analysis-recommendation.md has a "M043 outcome" section with: linter list, baseline issue count, fixed issue count, deferred items, future work.

## Architectural Decisions

### Phased rollout, Tier 1 first

**Decision:** Tier 1 (gosec, bodyclose, prealloc, errorlint, revive) first; Tier 2 only after Tier 1 is clean; Tier 3 opt-in only.

**Rationale:** Tier 1 has the highest value-to-noise ratio. gosec catches security issues, bodyclose catches resource leaks, prealloc catches allocation regressions, errorlint catches error wrapping bugs, revive catches missing docs. Each is low-noise in fd's existing code patterns (per analysis in `docs/static-analysis-recommendation.md`).

**Alternatives Considered:**
- Enable all 50+ linters at once: rejected. Most have low fd value, several are deprecated, and noise would bury real findings.
- Enable only gosec + govulncheck (skip rest): rejected. bodyclose/prealloc/errorlint have nearly zero false positives in idiomatic Go and catch real issues in fd's HTTP client + slice patterns.

### Each new linter in warn mode for one pass, then fail mode

**Decision:** When enabling a new linter, first run with `severity: warning` (or per-rule), capture false-positive rate, then move to `severity: error` (default).

**Rationale:** Avoids the situation where a new linter generates 50 false positives on first CI run, blocks PRs, and team disables it. Warn → fix → fail is the standard adoption path used by Kubernetes, Terraform, and other large Go codebases.

**Alternatives Considered:**
- Fail mode immediately: rejected. M041 introduced 7 new files; we want to review them under the new lints before locking the build.
- Skip linters that have any false positives: rejected. The whole point of static analysis is to find issues; the noise is the cost.

### gocyclo threshold = 15, not default 30

**Decision:** `gocyclo.min-complexity: 15` for fd. Default is 30.

**Rationale:** M041's `CreateEmbedding` handler is ~150 LOC with nested loops (cache peek + chunk loop + cache.Set). Likely already above default. Setting threshold too high (e.g., 30) means we never catch complexity drift. Setting it too low (e.g., 5) means every handler is over threshold. 15 is calibrated for service code with some unavoidable complexity (e.g., request parsing + dispatch + response shaping).

**Alternatives Considered:**
- Default 30: rejected. fd's `CreateEmbedding` is already at ~12-15.
- Per-function override (//nolint:gocyclo): rejected as default. Use only with explicit justification for legitimately complex code.

### govulncheck runs in CI as a required step, not advisory

**Decision:** `.github/workflows/go-quality.yml` runs `govulncheck` and fails the build on any reported vulnerability. No `[advisory]` or `[skip]` configuration.

**Rationale:** Known vulnerabilities in `gin`, `redis/go-redis`, `onnxruntime`, or stdlib are a real attack surface for fd. Even for a local same-host service, defense in depth means catching them at CI.

**Alternatives Considered:**
- Advisory only: rejected. False sense of security.
- Trivy instead: out of scope. govulncheck is the official Go tool.

### Tier 3 linters NOT auto-enabled; opt-in only

**Decision:** `gofumpt`, `structslop`, `maligned`, `aligncheck`, `dupl`, `nakedret`, `wsl`, `goimports`, `lll` are NOT added to `.golangci.yml`. They stay opt-in for future needs.

**Rationale:** Most are style preferences (gofumpt, wsl) or premature optimization (structslop, maligned). Adding them creates team-wide stylistic debates (tabs vs spaces equivalents) without fixing real bugs. Future needs (e.g., a pprof-driven allocation hotspot) can opt-in specific tools with specific configs.

**Alternatives Considered:**
- Enable all per 2026 best practice: rejected. Best practice is "use golangci-lint with staticcheck and targeted tools", not "enable all 50+".
- Use golangci-lint `presets: [bugs, style, performance, ...]`: considered. These presets include Tier 1+2+3 in known-good combinations. Rejected for fd because the resulting config is less explicit and harder to defend in code review.

### Out of scope (lifted from fd-v2.md + analysis)

- **Custom linters / Semgrep rules**: out of scope. Documented as future work if specific patterns emerge.
- **Pre-commit hooks**: out of scope. CI is the gate. (Project has no pre-commit config; adding one is a separate task.)
- **IDE integration (gopls, vscode-go)**: out of scope. These are developer-side, not repo-side.
- **Auto-fix on save (gofmt -w in pre-commit)**: out of scope. CI fails the build on unformatted code; manual fix is the workflow.

## Error Handling Strategy

- New linter failures: noise floor is captured in `docs/static-analysis-recommendation.md` Phase 1 report. False positives are excluded per-file with `//nolint:<linter>` + justification comment. Genuine issues are fixed in-place with a commit reference.
- govulncheck vulnerabilities: if a dependency has a known vuln, the CI fails. Options: (a) upgrade the dep, (b) pin to last safe version, (c) document why we're stuck and exclude with a per-issue marker (last resort, requires team sign-off).
- No new error envelopes needed — this is a tooling change, not a runtime change.

## Risks and Unknowns

- **Noise from revive** on M041's unexported but used functions (`SetVector`, `SetBase64`, `GetIfPresent`, `Set`). Mitigation: explicit `//nolint:revive` on the methods that don't need godoc.
- **Noise from gocritic** especially in `api/embed/codec.go` (small functions, may trigger "hugeParam" or "rangeValCopy"). Mitigation: per-file `//nolint` with justification; prefer not enabling `gocritic` at all and rely on `staticcheck` for most checks.
- **govulncheck first run** may report a known vuln in one of fd's deps. Mitigation: the S03 task includes a step to upgrade/pin/document before declaring "exit 0" as the baseline.
- **M042 S02 async code** may introduce new patterns (errgroup, semaphore) that one of the new linters flags. Mitigation: M043 S01 runs against current main (no M042 S02 yet); M043 S02 re-runs after M042 S02 merges. Lint regressions caught in PR review.
- **gocyclo false positives** on chunked handler. Mitigation: explicit per-function `//nolint:gocyclo` if handler logic is inherently nested (cache → chunk → sub-chunk). Acceptable if complexity is justified.
- **CI run time** with 12+ linters may exceed the 10-minute `timeout-minutes: 10` in `go-quality.yml`. Mitigation: bump to 15 minutes; golangci-lint caches results so re-runs are fast.

## Slice Plan

- **S01: Tier 1 lint adoption + fix existing issues** (~4h)
- **S02: Tier 2 lint adoption + complexity refactor** (~4h)
- **S03: govulncheck CI integration + docs finalization** (~2h)

## Notes for future slices

- S01 must run BEFORE M042 S02 async code merges, so the new code is checked from day 1. If M042 S02 ships first, S01 must re-run to catch the async code.
- S02 gocyclo threshold may need tuning after first run. Default in M043 is 15; can move to 20 if too many false positives.
- S03 govulncheck may discover vulns that block CI; if so, M043 S03 must resolve them (upgrade, pin, document) before declaring done.

## Hand-off from `docs/static-analysis-recommendation.md`

The recommendation document (saved earlier in planning) is the source of truth for tier classification and tool selection. M043 S03 updates the document with the final M043 outcome, deferred items, and future work — turning it from "recommendation" to "as-implemented + roadmap".

## Communication

Project communication should be in Russian by default (per `KNOWLEDGE.md` rule 1). Linter error messages and CI logs are in English (tooling); we summarize them in Russian in `docs/` and `.gsd/` artifacts.
