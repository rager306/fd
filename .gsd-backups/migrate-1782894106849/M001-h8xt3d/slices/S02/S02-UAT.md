# S02: API validation and handler tests — UAT

**Milestone:** M001-h8xt3d
**Written:** 2026-05-19T06:58:36.633Z

# UAT: S02 API validation and handler tests

## Verification performed

- `cd api && go test ./handlers -count=1` — passed.
- `cd api && go test ./... -short` — passed, 44 tests in 4 packages.

## Acceptance checks

- Invalid batch `dimensions` returns 400.
- Invalid batch `encoding_format` returns 400.
- Valid batch defaults and 512d float format still return 200.
- Tests instantiate production handlers through `NewEmbeddingsHandler` and `NewBatchHandler`.

