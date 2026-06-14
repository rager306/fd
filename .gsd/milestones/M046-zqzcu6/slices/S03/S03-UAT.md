# S03: Batch backend work shaping — UAT

**Milestone:** M046-zqzcu6
**Written:** 2026-06-14T16:26:41.630Z

# S03: Batch backend work shaping — UAT

**Milestone:** M046-zqzcu6
**Written:** 2026-06-14

## UAT Type

- UAT mode: runtime-executable
- Why this mode is sufficient: Unit/static tests prove call-count reduction. Runtime HTTP checks against the rebuilt local API prove behavior and S02 regression safety.

## Preconditions

- API container has been rebuilt with S03 changes using `docker compose up -d --build api`.
- TEI and Redis are healthy.

## Smoke Test

Request `/ready` and `/health`; then send valid requests to `/v1/batch` and `/embeddings/batch` and verify success.

## Test Cases

### 1. API ready after rebuild

1. GET `/ready`.
2. GET `/health`.
3. Confirm runtime backend is `tei` and dimensions are `1024`.
4. **Expected:** API is ready and runtime metadata is intact.

### 2. `/v1/batch` valid nested shape

1. POST two inner batches to `/v1/batch`.
2. **Expected:** response contains two nested batches, lengths match input `[3,2]`, and all vectors are 1024-dimensional.

### 3. `/embeddings/batch` valid legacy shape

1. POST three inputs to `/embeddings/batch`.
2. **Expected:** response has `count=3`, `dimensions=1024`, and three encoded embeddings.

### 4. S02 rejection regression safety

1. POST a 2049-char input to `/v1/batch` and `/embeddings/batch`.
2. **Expected:** both return HTTP 413 with error code `input_too_long` and the expected parameter path.

## Edge Cases

### Cache hit short-circuit

1. Unit tests repeat the same batch request after first miss fill.
2. **Expected:** second request does not increment embedder call count.

## Failure Signals

- Valid batch shape changes.
- Too-long inputs stop returning `input_too_long`.
- Unit tests show per-input embedder calls return.
- Backend wrong-count responses do not fail closed.

## Requirements Proved By This UAT

- R029 — batch endpoints now enforce bounded request posture and bounded backend work.

## Not Proven By This UAT

- S04 auth/exposure posture.
- S05 LocalCache concurrency/lifecycle correctness.
- P1 #6 `/v1/embeddings` cache-peek sequencing.

## Notes for Tester

Runtime evidence: `cb8a0f47-c9f4-4daa-ba02-b68152bb85ac`, `db1bbf65-af81-41ea-b67c-bcb1f74c6efc`, `e48a4774-94ed-4d02-9ab2-cffba624437e`, `ba43c8b5-9581-4d65-99ee-6d00f568b4b0`.
Static proof: `6591611c-d4d4-4485-b17e-ac2be3aa5d6d`.
