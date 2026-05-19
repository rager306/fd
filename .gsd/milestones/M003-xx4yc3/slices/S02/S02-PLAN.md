# S02: Live API smoke tests

**Goal:** Run live health and endpoint smoke tests.
**Demo:** Live API returns valid embeddings and rejects invalid requests against real TEI/Redis.

## Must-Haves

- API health, TEI health, Redis ping pass.
- `/v1/embeddings` returns correct 1024d and 512d shapes.
- `/embeddings/batch` returns valid base64/float responses.
- Negative validation cases return 400.

## Proof Level

- This slice proves: curl/jq/http status evidence

## Integration Closure

Confirms public runtime endpoints work against real dependencies.

## Verification

- Captures HTTP status/body summaries and service logs around real calls.

## Tasks

- [x] **T01: Live dependency health checks** `est:small`
  Verify API health, TEI health, and Redis ping against the running stack.
  - Verify: curl API/TEI health and redis-cli ping.

- [x] **T02: Smoke test embeddings endpoint** `est:medium`
  Smoke test `/v1/embeddings` with single 1024d input and array 512d input, recording response shapes.
  - Verify: curl /v1/embeddings and jq summaries show expected dimensions and lengths.

- [x] **T03: Smoke test batch endpoint** `est:medium`
  Smoke test `/embeddings/batch` for base64 and float formats.
  - Verify: curl /embeddings/batch summaries show expected count/dimensions/payloads.

- [x] **T04: Negative API tests** `est:small`
  Run negative validation tests for invalid JSON, empty input, invalid dimensions, and invalid batch encoding.
  - Verify: curl status codes are 400.
