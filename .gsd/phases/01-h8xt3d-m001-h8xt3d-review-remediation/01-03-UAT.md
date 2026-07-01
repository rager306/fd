# S03: Local cache semantics — UAT

**Milestone:** M001-h8xt3d
**Written:** 2026-05-19T07:02:25.097Z

# UAT: S03 Local cache semantics

## Verification performed

- `cd api && go test ./cache -run 'TestLocalCache' -count=1` — passed.
- `cd api && go test ./... -short` — passed, 46 tests in 4 packages.

## Acceptance checks

- Existing key updates refresh value and TTL.
- LocalCache does not exceed configured maxSize after inserts.
- Existing TTL/delete/not-found tests still pass.

