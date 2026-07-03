# S02: Cache observability and log noise — UAT

**Milestone:** M004-9886ht
**Written:** 2026-05-19T10:25:53.921Z

# UAT: S02 Cache observability and log noise

## Verification performed

- `gofmt` over modified Go files — passed.
- `cd api && go test ./cache ./handlers -short` — passed.
- `cd api && go test ./... -short` — passed, 49 tests in 4 packages.
- `docker compose up -d --build api` — rebuilt and started API.
- Runtime smoke:
  - `GET /health` — success.
  - `POST /v1/embeddings` — success.
  - `POST /embeddings/batch` — success.
- Log grep confirmed no old success INFO messages:
  - `embeddings generated`
  - `batch embeddings generated`
  - `cache miss, calling TEI`
  - `batch cache miss, calling TEI`
- `gitnexus_detect_changes(repo=fd, scope=all)` — low risk, no affected processes.

## Result

Default API logs now keep startup/connect/listen INFO lines but do not spam one success log per embedding request. Cache-path diagnostics are available at debug level and do not log raw input text.

