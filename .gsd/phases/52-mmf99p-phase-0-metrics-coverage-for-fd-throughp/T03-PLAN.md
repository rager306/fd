---
estimated_steps: 16
estimated_files: 4
skills_used: []
---

# T03: Batch-fill ratio + request duration refinement

Добавить метрики:
1. `fd_tei_batch_fill_ratio` histogram с buckets [0, 0.125, 0.25, 0.375, 0.5, 0.625, 0.75, 0.875, 1.0] — отражает долю заполненности TEI batch (texts_count / 32 cap). Каждое TEI наблюдаемое call record значение.
   - Если текстов 8/32 -> ratio=0.25
   - Если 32/32 (full) -> ratio=1.0
2. `fd_request_duration_seconds` уже существует — расширить buckets с холодными (cold-miss 200ms+):
   `[]float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.2, 0.5, 1.0, 5.0}` — даёт visibility как cache-hot (1-5ms), так и cold-miss (200ms-1s).
   - ⚠️ Risk: изменение buckets — может сломать PromQL queries с `le="X"`. Backward compat: оставить original 4 buckets и add second histogram `fd_request_duration_seconds_v2` отключив — или агрессивно replace.
   - Decision: REPLACE existing buckets (naming change is acceptable since метрика не была широко использована в production dashboards; new buckets give p95/p99 visibility).
3. Apply new bucket set to `fd_request_duration_seconds` in same NewMetrics().

Instrument api/handlers/embeddings.go fillEmbeddingChunk:
- pass len(missTexts) / 32 to metrics.ObserveBatchFillRatio
- Existing ObserveCacheResult calls already in tiered.go (after T01)

Instrument api/embed/tei.go doEmbedRequest:
- existing call records duration; now also records batch_fill_ratio

Files: api/observability/metrics.go (new histogram + bucket change), api/handlers/embeddings.go (record fill), api/embed/tei.go (record fill), tests in metrics_test.go + tiered_test.go.

Verify: existing TestMetricsHandlerExposesPrometheusText should still pass (all metric names present). Add new test that verifies batch_fill_ratio buckets are populated.

## Inputs

- `api/observability/metrics.go`
- `api/handlers/embeddings.go`
- `api/embed/tei.go`

## Expected Output

- `api/observability/metrics.go`
- `api/handlers/embeddings.go`
- `api/embed/tei.go (if not in T02)`
- `api/observability/metrics_test.go`

## Verification

go test ./api/observability/... -race ; new TestBatchFillRatioObserved passes.

## Observability Impact

Batch-fill ratio directly answers "is Phase 1b coalescing worth doing?" - if observed avg ratio >0.6, no much room; if <0.3, lots of room.
