---
title: Static analysis tools for fd — assessment and phased rollout plan
sources:
  - https://github.com/dpolivaev/static-analysis (Go section, ~50+ tools)
  - analysis-tools.dev Go rankings
  - 2026 Go community consensus (Reddit r/golang, blogs)
date: 2026-06-13
author: gsd-orchestrator
related: M041-4tw0w7 (fd v2 quality baseline)
---

# Static analysis tools for fd — assessment and rollout plan

## 1. Current state

`fd` uses `golangci-lint v2.12.2` (pinned in CI) with **7 enabled linters**:

```yaml
# .golangci.yml
linters:
  enable:
    - errcheck
    - govet
    - ineffassign
    - staticcheck
    - unused
    - goconst
    - misspell
```

CI: `.github/workflows/go-quality.yml` runs `go test ./... -short` + `golangci-lint run` on push to master and PRs. Go version: 1.22.2 locally, 1.25.x in CI.

`go.mod` has zero dev-tooling deps (no staticcheck, gosec, etc. as direct imports).

## 2. Gap analysis vs curated awesome list

The dpolivaev/static-analysis repo (Go section) lists ~50+ tools. Mapping to current fd:

| Tier | Tool | Status | fd value | Notes |
|---|---|---|---|---|
| 1 | **gosec** | NOT enabled | **HIGH** | Security. Catches: weak crypto, file perms, hardcoded creds, G107 (URL from variable), G110 (HTTP server without timeouts). fd has HTTP server (M041 S02 lifecycle), request body parsing, env var loading — all relevant. |
| 1 | **govulncheck** | NOT enabled (separate tool) | **HIGH** | Official Go vuln scanner. Checks stdlib + dependencies. Critical for production. Runs separately from golangci-lint (`go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...`). |
| 1 | **bodyclose** | NOT enabled | HIGH | Catches `resp.Body` not closed. fd has 2 outbound HTTP clients (`tei.go`, `onnx.go`); embedded/handlers don't make outbound calls. Risk: future regression in any HTTP client code path. |
| 1 | **prealloc** | NOT enabled | MED-HIGH | Slice preallocation hints. fd already preallocates (`make([][]float32, len(texts))`, `make([]int, 0, len(chunk))`) — prealloc would catch any regression to `make([]T, 0)` + `append` pattern. |
| 1 | **errorlint** | NOT enabled | MED-HIGH | Checks `error` wrapping via `%w` (Go 1.13+). fd already uses `%w` consistently in `api/embed/onnx.go`, `api/embed/hf_tokenizer_native.go`. errorlint catches subtle bugs like `errors.New(fmt.Errorf(...))` (loses wrap), `errors.Wrap` from `pkg/errors` (not in fd, but defensive). |
| 1 | **revive** | NOT enabled (golangci-lint has it in presets) | HIGH | Successor to deprecated `golint`. Has checks staticcheck doesn't: exported comments, package comments, confusing names, dot-imports, early-return, etc. More lints for code-as-documentation quality. |
| 2 | **gocyclo** | NOT enabled | MED | Cyclomatic complexity. M041 embeddings handler is now ~150 LOC with nested loops (cache peek → chunk loop → sequential TEI). Likely above threshold (default 30). May want threshold = 15-20 for fd. |
| 2 | **gocritic** | NOT enabled | MED | Meta-linter with 30+ checks. Significant overlap with staticcheck; revive. Adds value if not already covered. |
| 2 | **durationcheck** | NOT enabled | MED | Catches `time.Duration` conversion bugs like `int * time.Second` (if `int` overflows). fd has `30*time.Second`, `120*time.Second` — these are fine (small literals). Low actual risk but cheap insurance. |
| 2 | **unparam** | NOT enabled | MED | Flags unused function parameters. Helps when refactoring signatures but low direct value. |
| 2 | **nilnil** | NOT enabled | LOW-MED | Catches `return nil, nil` patterns (ambiguous error). fd has `return nil, nil, 0, fmt.Errorf(...)` in onnx tokenize — multi-return, not relevant. |
| 2 | **contextcheck** | NOT enabled | MED | Context propagation in stdlib calls. fd uses `context.WithTimeout` + `c.Request.Context()` correctly. Low risk. |
| 3 | **gofumpt** | NOT enabled | LOW | Stricter formatter than gofmt. Style preference. Not value-add for behavior. |
| 3 | **structslop/maligned/aligncheck** | NOT enabled | LOW | Memory layout optimization. Premature optimization — Go's escape analysis + GC handle this. Skip unless benchmarks show allocations. |
| 3 | **dupl** | NOT enabled | LOW | Duplicated code. fd has 2 `encodeEmbedding`/`float32SliceToBytes` — already extracted to `api/embed/codec.go` (M041 S01 T04). Low remaining dup. |
| 3 | **nakedret** | NOT enabled | LOW | Naked returns. fd has them in error paths of `batch.go` (5 occurrences). Acceptable style; not blocking. |
| 3 | **wsl** | NOT enabled | LOW | Whitespace style. Already covered by gofmt. |
| 3 | **errcheck** | already enabled | n/a | ✅ |
| 3 | **goimports** | NOT enabled | LOW | Grouped imports. Often in golangci-lint default. |
| 3 | **lll** | NOT enabled | LOW | Long lines. fd already mostly within 120 cols. |

**Tier 1 (high value, low risk, recommended for immediate adoption):** gosec, govulncheck, bodyclose, prealloc, errorlint, revive.

**Tier 2 (medium value, add after baseline is clean):** gocyclo, gocritic, durationcheck, unparam, contextcheck, nilnil.

**Tier 3 (specialized, opt-in only):** gofumpt, structslop/maligned, dupl, nakedret, wsl, goimports, lll.

## 3. Phased rollout plan

### Phase 1: Tier 1 (HIGH value, immediate)

**Step 1.1**: Run new linters locally in **warn mode** (not fail) to assess noise level.

```yaml
# .golangci.yml additions
linters:
  enable:
    - existing 7
    - gosec
    - bodyclose
    - prealloc
    - errorlint
    - revive
  settings:
    gosec:
      G107: warn  # URL from variable — fd has TEI_URL, EMBEDDING_BACKEND etc, suppress
      G304: warn  # file inclusion via variable — n/a but defensive
    revive:
      rules:
        - name: exported
        - name: package-comments
        - name: var-naming
        - name: error-strings
        - name: error-naming
        - name: error-returning
        - name: errorf
        - name: blank-imports
        - name: context-as-argument
        - name: dot-imports
        - name: early-return
        - name: exported
        - name: if-return
        - name: increment-decrement
        - name: var-declaration
        - name: range
        - name: receiver-naming
        - name: time-naming
        - name: indent-error-flow
        - name: errorf
issues:
  exclude-rules:
    - linters: [gosec]
      text: "G107"  # TEI_URL, REDIS_HOST etc come from env
```

Run: `golangci-lint run --config .golangci.yml ./...` — capture warnings count, identify false positives.

**Step 1.2**: `govulncheck` standalone. Add to CI:

```yaml
# .github/workflows/go-quality.yml
- name: Run govulncheck
  working-directory: api
  run: |
    go install golang.org/x/vuln/cmd/govulncheck@latest
    $(go env GOPATH)/bin/govulncheck ./...
```

**Step 1.3**: Fix or document each warning. Likely fixes needed:
- `gosec G107` on env-var URL loading (suppress with `//nolint:gosec` and justification)
- `revive exported` on `EmbeddingObj.SetVector`/`SetBase64` (unexported methods, public fields — no action)
- `bodyclose` may find leaks in new HTTP client code
- `prealloc` may find new `make([]T, 0)` patterns (M041 S02 async will introduce these — design with cap from start)
- `errorlint` may find `errors.New(fmt.Errorf(...))` (anti-pattern, loses wrap)

**Step 1.4**: Move from warn to fail mode in CI.

```yaml
# .golangci.yml
linters:
  enable:
    - all the above
  # Remove exclude-rules once baseline is clean
```

### Phase 2: Tier 2 (after Phase 1 baseline clean)

Add gocyclo (with custom threshold), gocritic (selectively — overlap with staticcheck is real), durationcheck, unparam, contextcheck, nilnil.

```yaml
linters:
  enable:
    - all Phase 1
    - gocyclo
    - gocritic
    - durationcheck
    - unparam
    - contextcheck
    - nilnil
  settings:
    gocyclo:
      min-complexity: 15
    gocritic:
      enabled-tags: [diagnostic, performance, style]
      disabled-checks: [hugeParam, rangeValCopy]
```

### Phase 3: Tier 3 (opt-in, style/deep)

`gofumpt` only if team agrees; `structslop` only if pprof shows allocation hot path.

## 4. Tooling dependencies

Adding these linters doesn't require new `go.mod` deps — `golangci-lint v2.12.2` already bundles them. But `govulncheck` needs:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```

Or pin in CI via `go run golang.org/x/vuln/cmd/govulncheck@latest`.

## 5. Risks

- **False positive noise** (especially gosec, gocritic, revive). Mitigation: start in warn mode, exclude rules per file, iterate.
- **Performance regression** in CI. golangci-lint caches results, but adding 10+ linters can slow first run. Mitigation: cache in CI, parallelize.
- **gocyclo false positives** for inherently complex code (e.g., M041 S04 chunked handler). Mitigation: explicit per-function `//nolint:gocyclo` with justification.
- **errorlint + testify**: tests use `assert.Error(t, err)` which is fine, but `require.Error` chains might trigger. Mitigation: check test files separately.

## 6. M041 compatibility

M041 S01 introduced `api/handlers/{errors,recovery,notfound}.go` + `api/middleware/validation.go` + `api/embed/codec.go` — all new code that should pass the new linters cleanly. Phase 1 step 1.3 should fix issues in this new code first (it's the most recently added and most-likely-to-have-style-issues).

## 7. M043 proposal (formal milestone)

If user wants this as a formal GSD milestone, the structure would be:

- **M043: Static analysis quality hardening** (3 slices):
  - **S01: Tier 1 lint adoption** (gosec, bodyclose, prealloc, errorlint, revive + govulncheck). ~4h. Fix existing issues, integrate CI.
  - **S02: Tier 2 lint adoption** (gocyclo, gocritic, durationcheck, unparam, contextcheck, nilnil). ~4h.
  - **S03: Baseline + final `.golangci.yml` consolidation** + per-linter exclusion doc. ~2h.

Total: ~10h. Blocks: none. Predecessor: M041 S01 (for new code review).

Alternative: do this as **quick task** (no milestone) — just update `.golangci.yml`, add govulncheck to CI, fix issues.

## 8. Concrete next steps (recommended)

1. **Quick task approach** (no formal milestone): update `.golangci.yml` with Tier 1, run `golangci-lint run` locally, fix issues (~2h).
2. Add `govulncheck` step to CI (~30min).
3. Document in `docs/static-analysis-recommendation.md` (this file).
4. Re-run mutation testing from M041 S04 to confirm new linters don't suppress critical coverage.

If user prefers formal GSD ceremony: propose M043 with the 3-slice plan above.
