# S01 Research: Issue #3 Audit Validation Map

## Source

- GitHub issue #3: `[audit] 3 P0 / 7 P1 / 22 other — default-open auth, batch-endpoint DoS, N+1 TEI, LocalCache races`
- Current HEAD at S01 start: `682871c` plus PR #1/#2 merge history already present.
- Static validation evidence: `gsd_exec:cafb3f84-b852-4d21-b71f-c13c9f5afd77`.

## P0 Inventory

| ID | Status | Issue claim | Current-code evidence | Root decision or assumption | Remediation wave |
|---|---|---|---|---|---|
| #1 | Confirmed policy risk | Auth disabled by default because empty `FD_API_KEY` disables `APIKeyAuth`; compose publishes API port. | `api/middleware/auth.go` allows all protected endpoints when `apiKey == ""`; compose exposes `8000:8000`. | Local/dev parity and same-host assumptions leaked into default compose exposure posture. | S04 Exposure posture policy |
| #2 | Confirmed defect | `/embeddings/batch` mounted as bare handler with no middleware guardrails. | `api/main.go` route is `r.POST("/embeddings/batch", batchHandler.CreateBatchEmbeddings)`; handler performs only partial inline validation. | Legacy batch endpoint was kept separate from `/v1/embeddings` validation chain and never received the hardened request boundary. | S02 Batch endpoint guardrails |
| #3 | Confirmed defect | `/v1/batch` skips validation, rate limit, and input caps. | `api/main.go` route only uses `LifecycleGateWithCapacity`; request shape has inner batch validation but no body cap or rate limiter. | Batch endpoint was added as a parallel surface rather than sharing a common endpoint guardrail policy. | S02 Batch endpoint guardrails |

## P1 Inventory

| ID | Status | Issue claim | Current-code evidence | Root decision or assumption | Remediation wave |
|---|---|---|---|---|---|
| #4 | Confirmed performance and resilience defect | Legacy batch endpoint does N+1 TEI calls on misses. | `api/handlers/batch.go` loops over inputs and loader calls `Embed(ctx, []string{text})` per input. | Cache abstraction optimized single-input semantics; batch handler did not collect misses. | S03 Batch backend work shaping |
| #5 | Confirmed performance and resilience defect | `/v1/batch` does per-input `Embed` across inner batches. | `api/handlers/v1batch.go` loops per text and loader calls `Embed(ctx, []string{text})`. | Same as #4; cache helper lacks batch-miss API. | S03 Batch backend work shaping |
| #6 | Confirmed performance concern | `/v1/embeddings` cache peek performs sequential L2 Redis checks up to request max. | Handler peeks item-by-item to avoid TEI work. It is bounded by `/v1/embeddings` max input count, so P1 performance concern, not P0 DoS. | Correctness and cache-hit preservation were prioritized over batched Redis lookup. | S06 residual or future optimization unless S03 creates reusable batch cache API |
| #7 | Confirmed policy risk | `/metrics` is auth-exempt and can leak telemetry. | `api/middleware/auth.go` includes `publicMetrics`. | Observability was made operator-friendly without explicit diagnostics exposure policy. | S04 Exposure posture policy |
| #8 | Confirmed policy and reliability risk | Rate limiter can be spoofed if trusted proxy configuration is wrong; limiter map can grow. | Rate limiter keys off client identity; no explicit trusted proxy setup observed in route bootstrap. | Local same-host assumptions reduced attention to internet-edge proxy semantics. | S04 or S06 depending on exposure policy |
| #9 | Confirmed contract cleanup | 405 code handling bypasses central registry shape. | `NoMethod` path has custom handler; issue claims code is absent from registry. Needs exact confirmation during S06. | Error contract evolved over time with separate middleware paths. | S06 residual closure |
| #10 | Confirmed correctness risk | LocalCache size counter and lifecycle can drift; no Close for eviction goroutine. | `LocalCache` uses `sync.Map` plus separate mutex-protected `size`; `evictLoop` runs forever and lacks `Close`. | Simplicity and hot-path concurrency were prioritized over deterministic lifecycle/accounting. | S05 LocalCache correctness |

## Safe Probe Result

`gsd_exec:cafb3f84-b852-4d21-b71f-c13c9f5afd77` confirmed these current-code signals:

```text
v1_batch_route_only_lifecycle: CONFIRMED
legacy_batch_bare_handler: CONFIRMED
auth_empty_key_allows_all: CONFIRMED
metrics_public: CONFIRMED
localcache_no_close: CONFIRMED
localcache_size_counter: CONFIRMED
legacy_batch_per_input_embed: CONFIRMED
v1batch_per_input_embed: CONFIRMED
```

## Wave Ordering

1. S02 fixes P0 batch guardrails before performance work.
2. S03 fixes N+1 backend work after request bounds exist.
3. S04 fixes exposure posture while preserving PR #2 health/readiness probe behavior.
4. S05 fixes LocalCache correctness with race/capacity evidence.
5. S06 rechecks all P0/P1 and triages P2/P3 residuals.
