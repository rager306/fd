# M049 Issue #8 Closure Matrix

Captured: 2026-06-15

Issue: https://github.com/rager306/fd/issues/8
Input artifact: `documents/issue-8-current-m049.md`

## Final Verification Summary

```text
go test ./...: 295 passed in 10 packages
golangci-lint: 0 issues
govulncheck: 0 reachable vulnerabilities
live container smoke: 5 passed, 0 failed
Docker Compose: api, redis, tei healthy
```

## Closure Matrix

| Issue #8 item | Outcome | Evidence |
|---|---|---|
| AN-A cache flush / invalidation surface | Implemented | `LocalCache.Flush/Size`, `RedisCache.Delete/FlushNamespace`, `TieredCache.Delete/Flush/LocalSize`, `POST /v1/cache/flush`, `POST /v1/cache/delete`. Evidence: `benchmark-results/m049-s01-cache-invalidation.md`, live proof `benchmark-results/m049-s03-live-container-proof.md`, R040 validated. |
| AN-B health hides last_error and dependencies | Implemented | `/health` can emit `last_error`, `dependencies.tei`, `dependencies.redis`, and dependency latency/reachability. Evidence: `benchmark-results/m049-s02-health-metrics-context.md`, live proof, R041 validated. |
| AN-C metrics lack capacity/cache occupancy | Implemented for solo-safe cheap signals | `/health` exposes `in_flight_capacity`; `/metrics` exposes `fd_in_flight_requests`, `fd_in_flight_capacity`, and `fd_cache_entries{tier="l1"}`. Redis occupancy scanning per scrape is intentionally avoided for solo use. Evidence: S02 artifact, live proof, R041 validated. |
| AN-D traces visible to shared bearer holder | Deferred | User explicitly said not to complicate AN-D. fd is planned for solo use; no separate admin-token/multi-tenant trace scope was added. Decision: D051; R042 validated. |
| AN-E deployment-tunable policy hardcoded | Scoped / deferred | Broad config extraction was not added. For solo use, defaults remain stable; only minimal option seam needed for implemented health probes was added (`HealthOptions`). Decision: D051; R042 validated. |
| AN-F TieredCache/TEI retry/circuit options | Scoped / deferred | Broad pluggable cache/retry policy was not added. S01 added explicit invalidation primitives, and S02 added probe/metrics seams. Retry/circuit behavior remains unchanged for solo stability. Decision: D051; R042 validated. |
| AN-G dead `publicMetrics` constant | Not in requested implementation scope | Issue #8 requested focus was AN-A and AN-B/C, with AN-D/E/F scope decision. This low cleanup remains available for a later tiny cleanup if desired. |
| AN-H CORS expose embedding headers | Not in requested implementation scope | Browser-facing CORS header exposure was not part of requested M049 implementation. |
| AN-I README poll contract / single-tenant docs | Not in requested implementation scope | M049 records solo-scope decision in GSD and closure artifacts, but README docs were not changed. |

## Runtime Proof Highlights

From `benchmark-results/m049-s03-live-container-proof.md`:

```text
PASS GET /health exposes capacity and dependency context
PASS GET /metrics exposes runtime/cache gauges
PASS POST /v1/cache/flush is auth protected and works
PASS Embedding cache HIT becomes MISS after flush
PASS Embedding cache HIT becomes MISS after delete
SUMMARY passed=5 failed=0 total=5
```

## Requirement Outcomes

- R040 validated: authenticated cache invalidation for stale embeddings.
- R041 validated: health/metrics agent diagnostics.
- R042 validated: solo deployment scope boundary for AN-D and AN-E/F.

## Notes

- No GitHub issue mutation or push was performed by this milestone.
- Secret values were not recorded in artifacts or logs.
