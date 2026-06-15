# S04: Warmup retry and closure — UAT

**Milestone:** M047-9fngng
**Written:** 2026-06-15T08:38:12.101Z

# S04: Warmup retry and closure — UAT

**Milestone:** M047-9fngng
**Written:** 2026-06-15

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: S04 changes backend warmup retry behavior and writes closure artifacts. The observable contract is covered by deterministic unit tests, final gates, static proof, and artifact checks; no browser surface is involved.

## Preconditions

- `benchmark-results/m047-s04-warmup-retry-closure.md` exists.
- `benchmark-results/m047-issue-6-closure.md` exists.

## Smoke Test

Verify warmup retry code and issue #6 closure matrix.

## Test Cases

### 1. Warmup retry invariant

1. Inspect `api/main.go`.
2. **Expected:** warmup retry policy/helper exists, default max attempts is three, attempt failures are logged, and success calls `state.MarkWarmupDone()`.

### 2. Issue #6 closure matrix

1. Inspect `benchmark-results/m047-issue-6-closure.md`.
2. **Expected:** rows exist for #11, #14, #13, #32, #25, and #15, all marked fixed.

### 3. Requirements validated

1. Inspect `.gsd/REQUIREMENTS.md`.
2. **Expected:** R033-R036 are present and validated.

### 4. Final gates recorded

1. Inspect `benchmark-results/m047-s04-warmup-retry-closure.md`.
2. **Expected:** final gate evidence records 290 passing tests, lint 0 issues, and govulncheck 0 reachable vulnerabilities.

## Edge Cases

- Warmup failure then later success marks ready and clears prior error.
- Terminal warmup failure records last error and does not mark ready.
- Closure matrix does not leave any issue #6 finding unclassified.

## Failure Signals

- Warmup failure returns after one attempt again.
- R034 is not validated.
- Final gate evidence is missing.

## Requirements Proved By This UAT

- R034: bounded warmup retry and readiness recovery.
- M047 aggregate: issue #6 findings are closed and requirements R033-R036 are validated.

## Not Proven By This UAT

- Live TEI container restart behavior. The retry/circuit behavior is verified with deterministic tests and artifact checks.

## Notes for Tester

UAT evidence: `19a0bc3e-9186-40ba-98ce-c456535bae8d`, `81114e5d-c0a5-423c-8d73-0e40fceed740`, `83eeafc9-ed6a-46ab-a427-7863a5f3fb51`, `0644ab40-ee81-4eb7-9481-22ab68329ab0`.
