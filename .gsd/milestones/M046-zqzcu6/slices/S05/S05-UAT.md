# S05: LocalCache correctness — UAT

**Milestone:** M046-zqzcu6
**Written:** 2026-06-14T19:23:30.804Z

# S05: LocalCache correctness — UAT

**Milestone:** M046-zqzcu6
**Written:** 2026-06-14

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: S05 changes internal cache concurrency/accounting/lifecycle behavior. Unit tests, race-enabled tests, static proof, and evidence artifacts directly prove the contract; no external HTTP behavior is expected to change.

## Preconditions

- S05 code and tests are present.
- `benchmark-results/m046-s05-localcache-correctness.md` has been written.

## Smoke Test

Verify LocalCache implementation shape, evidence artifact, R031 status, and shutdown integration.

## Test Cases

### 1. Evidence artifact records gates

1. Read `benchmark-results/m046-s05-localcache-correctness.md`.
2. **Expected:** artifact records P1 #10, targeted tests, race test, full tests, lint, govulncheck, and static proof.

### 2. LocalCache implementation shape

1. Inspect `api/cache/local.go`.
2. **Expected:** no `sync.Map`, no separate `size` counter, one `map[string]l1Entry`, derived `currentSize`, and `Close() error`.

### 3. API lifecycle integration

1. Inspect `api/main.go`.
2. **Expected:** shutdown/error cleanup closes `localCache` through `closeResource`.

### 4. Requirement and plan status

1. Inspect `.gsd/REQUIREMENTS.md` and `documents/issue-3-audit-remediation-plan-m046.md`.
2. **Expected:** R031 is validated and S05 is marked done.

## Edge Cases

### Race-enabled LocalCache tests

1. Run `cd api && go test -race ./cache -run TestLocalCache`.
2. **Expected:** pass.

### Close lifecycle

1. Call `Close()` twice.
2. **Expected:** both calls return nil and do not panic.

## Failure Signals

- LocalCache uses `sync.Map` with independent `size` counter again.
- `Close()` missing or non-idempotent.
- Race detector fails on LocalCache tests.
- API shutdown no longer closes LocalCache.

## Requirements Proved By This UAT

- R031 — LocalCache correctness and lifecycle are deterministic and race-tested.

## Not Proven By This UAT

- S06 residual P1 #6 and P2/P3 closure matrix.

## Notes for Tester

UAT evidence: `9c5fef8f-f26b-4251-ba7c-5251645e39b5`, `35380403-fd5f-4d1a-905f-2dc6b9a55daf`, `870b9370-036b-48d9-9c35-b7f2dd97359e`, `c202aabb-1af1-4f17-bd07-132ebbaae49b`.
Static proof: `f124000a-5996-4c68-888d-1e31237c6d39`.
Completeness proof: `a623b5a9-3cd3-4413-bb88-20ee70f64547`.
