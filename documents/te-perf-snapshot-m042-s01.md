---
milestone: M042-fjf2en
slice: S01
task: T01
captured: 2026-06-14T08:51:17Z
source_evidence: .gsd/exec/177ccc56-1ecc-434d-b12b-a0044ab13c20.stdout
---

# M042 S01 T01 — TEI cold telemetry snapshot

This snapshot captures direct TEI server-side timings for the currently running `fd_tei` container. It intentionally bypasses fd and Redis so the measurements are not hidden by fd's cache-hot behavior from M041/S04.

## Runtime and limits

`GET http://localhost:30080/info` returned:

| Field | Value |
|---|---:|
| model_id | `deepvk/USER-bge-m3` |
| model_dtype | `float32` |
| version | `1.9.3` |
| max_concurrent_requests | `512` |
| max_batch_requests | `4` |
| max_client_batch_size | `32` |
| max_batch_tokens | `16384` |
| max_input_length | `8192` |
| tokenization_workers | `11` |

The key puzzle remains visible: `max_concurrent_requests=512` is high, but batch=32 requests spend multiple seconds in TEI `queue_time` even when requests are sent sequentially.

## Existing M041 baseline context

`benchmark-results/fd-v2-baseline-before-m041-s04.md` measured fd-level cache-hot wall-clock latency, not TEI cold inference. It showed very fast p95s after cache warmth:

| fd batch size | n | p95 | note |
|---:|---:|---:|---|
| 1 | 100 | 3.7ms | cache-hot fd path |
| 10 | 20 | 3.9ms | cache-hot fd path |
| 32 | 10 | 3.5ms | cache-hot fd path |

Those values are useful for D045 cache-hot acceptance, but they do not explain the TEI cold path. The direct TEI snapshot below is the relevant evidence for M042.

## Fresh direct TEI sequential snapshot

Command evidence: `.gsd/exec/177ccc56-1ecc-434d-b12b-a0044ab13c20.stdout`.

Method:

- Sent direct `POST http://localhost:30080/v1/embeddings` requests to TEI, not through fd.
- Used `model=deepvk/USER-bge-m3`.
- Ran 10 sequential requests for each batch size: 1, 8, 32.
- Parsed `fd_tei` logs since `2026-06-14T08:51:17Z` for `total_time`, `tokenization_time`, `queue_time`, and `inference_time`.

### Summary

| Batch size | n | ok | client wall p50 | client wall p95 | TEI total p50 | TEI total p95 | queue p50 | queue p95 | inference p50 | inference p95 |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | 10 | 10 | 199.384ms | 306.682ms | 196.283ms | 303.931ms | 0.253ms | 0.300ms | 192.560ms | 299.856ms |
| 8 | 10 | 10 | 1354.058ms | 1575.901ms | 1350.442ms | 1572.372ms | 427.093ms | 496.669ms | 613.751ms | 699.894ms |
| 32 | 10 | 10 | 5342.144ms | 6034.448ms | 5337.374ms | 6027.026ms | 2434.795ms | 2477.019ms | 663.943ms | 724.354ms |

### Raw rows

| batch | i | status | client_wall_ms | tei_total_ms | token_ms | queue_ms | infer_ms |
|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | 1 | 200 | 149.340 | 128.517 | 4.476 | 0.422 | 123.515 |
| 1 | 2 | 200 | 176.519 | 173.507 | 3.120 | 0.300 | 169.983 |
| 1 | 3 | 200 | 147.749 | 145.208 | 2.360 | 0.252 | 142.517 |
| 1 | 4 | 200 | 308.667 | 305.474 | 3.390 | 0.279 | 301.723 |
| 1 | 5 | 200 | 193.852 | 191.337 | 3.390 | 0.255 | 187.580 |
| 1 | 6 | 200 | 142.761 | 139.602 | 2.919 | 0.284 | 136.302 |
| 1 | 7 | 200 | 237.685 | 234.327 | 4.291 | 0.252 | 229.693 |
| 1 | 8 | 200 | 204.916 | 201.228 | 3.349 | 0.202 | 197.541 |
| 1 | 9 | 200 | 221.935 | 219.267 | 3.935 | 0.218 | 214.978 |
| 1 | 10 | 200 | 306.682 | 303.931 | 3.775 | 0.214 | 299.856 |
| 8 | 1 | 200 | 1355.906 | 1352.418 | 0.824 | 453.367 | 565.514 |
| 8 | 2 | 200 | 1177.103 | 1173.361 | 0.263 | 312.729 | 454.102 |
| 8 | 3 | 200 | 1242.301 | 1238.833 | 0.373 | 400.819 | 466.963 |
| 8 | 4 | 200 | 1403.978 | 1400.975 | 0.627 | 458.877 | 699.894 |
| 8 | 5 | 200 | 1575.901 | 1572.372 | 0.333 | 507.793 | 589.460 |
| 8 | 6 | 200 | 1312.678 | 1309.649 | 1.178 | 288.655 | 653.884 |
| 8 | 7 | 200 | 1639.431 | 1636.241 | 0.535 | 496.669 | 817.576 |
| 8 | 8 | 200 | 1430.208 | 1426.509 | 0.541 | 483.624 | 612.847 |
| 8 | 9 | 200 | 1352.209 | 1348.465 | 0.377 | 275.551 | 673.294 |
| 8 | 10 | 200 | 1234.323 | 1230.188 | 0.359 | 317.347 | 614.654 |
| 32 | 1 | 200 | 5346.206 | 5341.182 | 1.213 | 2455.902 | 627.354 |
| 32 | 2 | 200 | 5338.083 | 5333.565 | 0.990 | 2390.142 | 666.218 |
| 32 | 3 | 200 | 5270.327 | 5264.348 | 1.924 | 2223.716 | 657.218 |
| 32 | 4 | 200 | 5250.222 | 5241.085 | 0.865 | 2463.042 | 654.793 |
| 32 | 5 | 200 | 6181.512 | 6177.634 | 1.001 | 2420.475 | 724.354 |
| 32 | 6 | 200 | 6034.448 | 6027.026 | 1.879 | 2954.085 | 752.711 |
| 32 | 7 | 200 | 5302.035 | 5297.892 | 1.611 | 2449.116 | 661.667 |
| 32 | 8 | 200 | 5623.387 | 5618.997 | 2.188 | 2477.019 | 701.865 |
| 32 | 9 | 200 | 5673.449 | 5665.359 | 0.833 | 2404.348 | 707.809 |
| 32 | 10 | 200 | 4830.548 | 4826.152 | 0.668 | 2031.324 | 602.894 |

## Correlation notes

Queue time scales with request batch size even without client-side concurrency:

- batch=1: queue p50 `0.253ms`, essentially no queueing.
- batch=8: queue p50 `427.093ms`, with total p50 `1350.442ms`.
- batch=32: queue p50 `2434.795ms`, with total p50 `5337.374ms`.

This makes the M042 symptom reproducible: queue time is not merely an overload artifact from many concurrent callers. TEI reports multi-second queue/scheduler time for large single-client batches while advertised capacity remains `max_concurrent_requests=512`, `max_batch_requests=4`, `max_client_batch_size=32`.

## T01 verdict

T01 is ready to feed T02/T03:

- The relevant bottleneck is direct TEI cold inference/scheduling, not fd HTTP overhead and not Redis cache-hit behavior.
- The strongest immediate signal is batch-size-sensitive `queue_time` growth.
- T02 should vary client concurrency and batch size to distinguish true backend serialization from TEI metric semantics or scheduler batching behavior.
