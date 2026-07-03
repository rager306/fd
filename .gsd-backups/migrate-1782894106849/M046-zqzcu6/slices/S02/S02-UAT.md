# S02: Batch endpoint guardrails — UAT

**Milestone:** M046-zqzcu6
**Written:** 2026-06-14T16:11:20.671Z

# S02: Batch endpoint guardrails — UAT

**Milestone:** M046-zqzcu6
**Written:** 2026-06-14

## UAT Type

- UAT mode: runtime-executable
- Why this mode is sufficient: S02 changes backend HTTP behavior. Runtime HTTP checks against the rebuilt local API prove both rejection paths and valid batch paths.

## Preconditions

- API container has been rebuilt with S02 changes using `docker compose up -d --build api`.
- TEI and Redis are healthy.

## Smoke Test

Request `/ready` and `/health`; then send valid requests to `/v1/batch` and `/embeddings/batch` and verify success.

## Test Cases

### 1. API ready after rebuild

1. GET `/ready`.
2. GET `/health`.
3. Confirm runtime backend is `tei` and dimensions are `1024`.
4. **Expected:** API is ready and runtime metadata is intact.

### 2. `/v1/batch` rejects too-long input

1. POST `{"batches":[["x" repeated 2049 times]]}` to `/v1/batch`.
2. **Expected:** HTTP 413 with error code `input_too_long` and param `batches[0][0]`.

### 3. `/embeddings/batch` rejects too-long input

1. POST `{"inputs":["x" repeated 2049 times]}` to `/embeddings/batch`.
2. **Expected:** HTTP 413 with error code `input_too_long` and param `inputs[0]`.

### 4. Valid batch requests still work

1. POST a small valid request to `/v1/batch`.
2. POST a small valid request to `/embeddings/batch`.
3. **Expected:** `/v1/batch` returns a 1024-dimensional vector and `/embeddings/batch` returns count `1`, dimensions `1024`.

## Edge Cases

### Body cap route middleware

1. Use package tests to send declared oversized content-length.
2. **Expected:** request returns `payload_too_large` before downstream handler runs.

## Failure Signals

- Either batch endpoint reaches embedder/cache for rejected input.
- Either batch endpoint returns generic invalid JSON for body cap failures.
- Valid batch requests no longer work.
- Health/readiness regress after rebuild.

## Requirements Proved By This UAT

- R029 — Batch endpoints enforce bounded request and lifecycle posture before backend work.

## Not Proven By This UAT

- P1 N+1 backend call shaping remains for S03.
- Exposure posture and LocalCache correctness remain for S04/S05.

## Notes for Tester

Evidence: `0e10d0e6-0ad5-4824-8456-778945b19345`, `0d532f91-4886-4a4b-9866-cb788e9347d2`, `9b990efc-8758-4440-a8d7-2d9ec4973c53`, `86b73f80-3859-4671-9691-3c1a97da5b1b`.
