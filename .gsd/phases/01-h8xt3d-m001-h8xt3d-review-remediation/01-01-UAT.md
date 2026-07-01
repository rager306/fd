# S01: Cache correctness and panic safety — UAT

**Milestone:** M001-h8xt3d
**Written:** 2026-05-19T06:52:39.187Z

# UAT: S01 Cache correctness and panic safety

## Verification performed

- `cd api && go test ./cache -run 'Test.*(Tiered|Binary|Redis|Local|Marshal)' -count=1` — passed, 13 tests in 1 package.
- `cd api && go test ./... -short` — passed, 38 tests in 4 packages.

## Acceptance checks

- Same text requested at 512d and 1024d uses separate L1/singleflight keys.
- Short embeddings return explicit errors instead of panicking.
- Existing cache binary round-trip tests still pass.
- Full short suite remains green.

