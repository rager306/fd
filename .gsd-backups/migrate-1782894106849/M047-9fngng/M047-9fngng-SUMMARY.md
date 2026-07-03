---
id: M047-9fngng
title: "Audit followup resilience contract"
status: complete
completed_at: 2026-06-15T08:39:11.162Z
key_decisions:
  - Remove un-emitted error codes instead of marking them reserved.
  - Use synthetic `server_error` os.Signal to reuse existing lifecycle shutdown path.
  - Use lightweight process-local TEI circuit state rather than external breaker dependency.
  - Use deterministic retry policies in tests with bounded production defaults.
key_files:
  - documents/issue-6-current-m047.md
  - benchmark-results/m047-s01-contract-cleanup.md
  - benchmark-results/m047-s02-graceful-listener-shutdown.md
  - benchmark-results/m047-s03-tei-retry-fast-fail.md
  - benchmark-results/m047-s04-warmup-retry-closure.md
  - benchmark-results/m047-issue-6-closure.md
  - api/main.go
  - api/main_test.go
  - api/embed/tei.go
  - api/embed/tei_test.go
  - api/handlers/errors.go
  - api/handlers/errors_test.go
lessons_learned:
  - Reliability follow-ups are easiest to prove with deterministic injected retry policies and fake HTTP servers.
  - Process-control failures should route through lifecycle state instead of exiting from background goroutines.
  - Public error registries need tests that ensure registered codes correspond to real emitters.
---

# M047-9fngng: Audit followup resilience contract

**M047 resolved issue #6 by fixing TEI retry/fast-fail, warmup retry, graceful listener shutdown, safe env parsing, and un-emitted error codes.**

## What Happened

M047 followed up M046 by addressing the remaining G4/G5 reliability and contract findings from GitHub issue #6. S01 persisted the issue input, fixed `getEnvInt` overflow/negative parsing with `strconv.Atoi`, removed un-emitted public error codes, and added a registry emitter test. S02 replaced listener-goroutine `os.Exit(1)` with controlled shutdown signalling and `errors.Is` matching for wrapped `http.ErrServerClosed`. S03 added bounded retry/backoff and a lightweight circuit fast-fail to the TEI client for transient network/502/503/504 failures while preserving non-retriable behavior. S04 added bounded warmup retry/backoff, readiness recovery after later success, final gate evidence, and a full issue #6 closure matrix. Final gates passed: `go test ./...` 290 tests, lint 0 issues, and govulncheck 0 reachable vulnerabilities.

## Success Criteria Results

- ✅ Issue #6 findings #11, #14, #13, #32, #25, and #15 revalidated and fixed.
- ✅ Safe env integer parsing implemented and tested.
- ✅ Fatal listener errors routed into lifecycle shutdown with `errors.Is`.
- ✅ TEI transient failures retry and repeated outages fast-fail.
- ✅ Warmup transient failures retry and later success clears error/readies service.
- ✅ Dead public error codes removed and future emitter coverage enforced.
- ✅ Final tests/lint/vulnerability gates passed.

## Definition of Done Results

- ✅ S01-S04 complete.
- ✅ Requirements R033-R036 validated.
- ✅ Issue #6 closure matrix written.
- ✅ Full Go tests, lint, govulncheck, and UAT passed.
- ✅ Local commits created through S03; S04 closure artifacts ready for final commit.
- ✅ No GitHub issue mutation performed.

## Requirement Outcomes

- R033 validated by S03 TEI retry/circuit tests and artifact proof.
- R034 validated by S04 warmup retry tests and artifact proof.
- R035 validated by S02 listener shutdown tests and artifact proof.
- R036 validated by S01 env parsing and error registry tests and artifact proof.

## Deviations

Live TEI container restart testing was not performed; deterministic tests prove retry/circuit code behavior without mutating external runtime. Final lint required minor cleanup before closeout.

## Follow-ups

Optional outward action: after explicit confirmation, push local commits and comment on or close GitHub issue #6 with the closure matrix.
