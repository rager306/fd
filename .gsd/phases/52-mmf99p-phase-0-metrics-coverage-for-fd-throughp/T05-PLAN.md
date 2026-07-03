---
estimated_steps: 19
estimated_files: 2
skills_used: []
---

# T05: Integration smoke + README metrics reference

End-to-end verification + documentation.

A. Integration smoke test:
  - Stand up real docker compose stack (fd_api, tei, redis) - same as M051 verification.
  - Make a sequence of POST /v1/embeddings calls: 1st miss (cold), 2nd same input (L1 hit), 3rd different input (L2 hit if Redis was warm).
  - Scrape /metrics endpoint.
  - Assert presence and non-zero values for all new metrics.
  - Save evidence log в .gsd/runtime/M052-mmf99p/.

B. Update README.md:
  - New section "Observability metrics" под Configuration или отдельный раздел перед Performance.
  - Document all fd_*_total, fd_*_seconds, fd_*_bytes metric names with help text.
  - Reference Phase 0 of Issue #9 - "Metrics coverage is the foundation for Phase 1+ tuning; metrics are exposed at /metrics in Prometheus format".
  - Note что /metrics is publicly accessible (FD_MAX_IN_FLIGHT middleware doesn't gate it - защищён auth, но метрика полезна для operators).

C. Update benchmark.py или tools/ если он enumerates метрики - убедиться что naming consistent (если есть script что lists метрики для validation).

Files: README.md (new section), .gsd/runtime/M052-mmf99p/ (evidence logs), maybe benchmark.py if it references metric names.

Verify:
- docker compose up -d всё работает
- curl localhost:8000/metrics содержит fd_cache_hits_total{tier=...}, fd_tei_request_duration_seconds_bucket, fd_tei_batch_fill_ratio, fd_cache_entries{tier=l1|l2}, fd_cache_memory_bytes{tier=l1|l2}, fd_tei_requests_in_flight, fd_tei_errors_total
- README renders markdown correctly
- go test ./... -race -count=1 зелёный

## Inputs

- `README.md (current)`
- `.gsd/runtime/M052-mmf99p/ (evidence directory)`
- `api/observability/metrics.go (T01-T03)`

## Expected Output

- `README.md (+observability section)`
- `.gsd/runtime/M052-mmf99p/metrics-scrape.log`
- `.gsd/runtime/M052-mmf99p/api-logs.jsonl`

## Verification

docker compose up -d ; curl localhost:8000/metrics | grep -E 'fd_cache_hits_total\{|fd_tei_request_duration_seconds_bucket|fd_tei_batch_fill_ratio|fd_cache_entries\{tier="l2'|fd_cache_memory_bytes|fd_tei_errors_total' (all present, non-zero); README renders; full test suite green.

## Observability Impact

Documentation surfaces new metrics для operators.
