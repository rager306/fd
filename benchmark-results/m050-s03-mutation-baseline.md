# M050 S03 Mutation Baseline

Date: 2026-06-15
Milestone: M050-rfqm1p
Slice: S03

## Goal

Create an initial bounded mutation-testing baseline for critical Go backend packages, or document a reproducible blocker if tooling is not viable.

## Tooling Probe

Candidates checked:

- `github.com/zimmski/go-mutesting/cmd/go-mutesting@latest`
  - Available, but old upstream lineage.
- `github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest`
  - Available and selected for baseline.
  - Requires Go >= 1.25.5 and uses automatic toolchain switch to `go1.25.11` in this environment.

Selected runner:

```bash
go run github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
```

Reason: more current fork, supports config/html output, and ran successfully against this Go 1.25 module.

## Smoke Result

Command:

```bash
cd api && go run github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest --exec 'go test ./cache' --exec-timeout 30 ./cache/hash.go
```

Result:

```text
ok fd-api/cache
PASS mutation 0
PASS mutation 1
mutation score 1.000000 (2 passed, 0 failed, 0 duplicated, 0 skipped, total is 2)
```

Interpretation: for this runner output, `PASS` means the test command detected/killed the mutation and the mutant did not survive.

## Bounded Baseline Scope

Initial critical-file scope:

- `api/cache/hash.go`
- `api/cache/keys.go`
- `api/handlers/cache.go`
- `api/handlers/health.go`
- `api/lifecycle/state.go`

Test command per mutation:

```bash
go test ./cache ./handlers ./lifecycle
```

## Results

### Critical-file baseline

Command:

```bash
cd api && go run github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest \
  --exec 'go test ./cache ./handlers ./lifecycle' \
  --exec-timeout 45 \
  ./cache/hash.go ./cache/keys.go ./handlers/cache.go ./handlers/health.go ./lifecycle/state.go
```

Result:

```text
The mutation score is 1.000000 (143 passed, 0 failed, 4 duplicated, 0 skipped, total is 143)
```

Interpretation:

- 143 mutants were killed by the bounded package tests.
- 0 mutants survived in the selected critical-file scope.
- 4 duplicate mutants were reported by the runner.
- Duration observed by the agent harness: about 47-51 seconds per run after dependency/toolchain setup.

Full local log path for this session: `/tmp/m050-s03-mutation-critical.log`.

## Policy

This is an informational baseline, not yet a mandatory CI hard gate. It is suitable for targeted local verification on critical backend files. Future CI use should either pin the runner version and cache the Go 1.25.11 toolchain or run as a manual workflow because each mutation reruns package tests.

## Verdict

S03 established a working bounded mutation baseline for selected critical cache, handler, and lifecycle files. The baseline score is 1.0 with no surviving mutants in scope.
