# S06: Audit closure — UAT

**Milestone:** M046-zqzcu6
**Written:** 2026-06-15T05:56:01.267Z

# S06: Audit closure — UAT

**Milestone:** M046-zqzcu6
**Written:** 2026-06-14

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: S06's user-visible output is a closure matrix plus code invariants for residual P1 fixes. Full tests, race tests, lint, govulncheck, static proof, and artifact checks prove the closure without needing a browser or external runtime mutation.

## Preconditions

- `benchmark-results/m046-s06-audit-closure.md` exists.
- S01-S05 artifacts exist.
- S06 code changes are present.

## Smoke Test

Verify closure matrix coverage and residual code invariants.

## Test Cases

### 1. Closure matrix covers all findings

1. Inspect `benchmark-results/m046-s06-audit-closure.md`.
2. **Expected:** rows exist for findings #1 through #32 and verdict states M046 closes P0/P1.

### 2. Residual P1 #6 code invariant

1. Inspect `api/handlers/embeddings.go`, `api/cache/tiered.go`, and `api/cache/redis.go`.
2. **Expected:** handler calls `GetManyIfPresent`, TieredCache implements it, and Redis uses MGET.

### 3. Residual P1 #9 code invariant

1. Inspect `api/handlers/errors.go` and `api/handlers/notfound.go`.
2. **Expected:** `CodeMethodNotAllowed` is registered and NoMethod path uses `WriteError`.

### 4. Requirements validated

1. Inspect `.gsd/REQUIREMENTS.md`.
2. **Expected:** R029, R030, R031, and R032 are present and validated.

## Edge Cases

### P2/P3 residuals

1. Inspect closure matrix rows #11-#32.
2. **Expected:** every row is explicitly fixed, mitigated, accepted, or deferred with rationale.

## Failure Signals

- Missing closure row for any issue #3 finding.
- `/v1/embeddings` regresses to per-item cache peeks.
- 405 bypasses the canonical error registry again.
- Requirements do not reflect validated R029-R032.

## Requirements Proved By This UAT

- R032 — `/v1/embeddings` cache peeks are batched per bounded chunk.
- M046 closure: P0/P1 issue #3 findings are fixed and P2/P3 findings are triaged.

## Not Proven By This UAT

- Future deferred P2/P3 remediation implementation.

## Notes for Tester

UAT evidence: `2912d887-d003-4340-ae7b-74c901258f74`, `574e9de5-fbc5-4cbc-bb18-48387e86b1b9`, `961e2e5a-4de9-410e-ab4a-e522bceb0311`, `b7d6f086-2972-4b17-85d9-8cd8449f6419`.
Static proof: `43c16c32-c290-499a-a42a-b8602a0ce6ee`, `ab9ac45b-a646-4b6e-a5ef-22839e715e5c`.
Closure completeness: `2a1a76b7-be59-4cb1-b79c-53aa2dc84ff7`.
