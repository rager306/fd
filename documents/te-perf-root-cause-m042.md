---
milestone: M042-fjf2en
slice: S01
task: T03
captured: 2026-06-14T10:20:00Z
inputs:
  - documents/te-perf-snapshot-m042-s01.md
  - benchmark-results/te-concurrency-profile-m042-s01.md
  - benchmark-results/fd-onnx-1024-performance-outcome-m019-s02.txt
  - benchmark-results/fd-runtime-recommendation-m040-s04.md
---

# M042 S01 — TEI queue time and startup RCA

## Executive summary

M042 started with one symptom: direct TEI cold requests for `deepvk/USER-bge-m3` can take about 5–6s for batch=32, with TEI reporting about 2.4s of `queue_time` even when the client is sequential and `max_concurrent_requests=512`.

The investigation found two related bottlenecks:

1. **Batch-size-sensitive TEI scheduling/queue time.** Direct TEI requests show queue time grows from near zero at batch=1 to hundreds of ms at batch=8 and multiple seconds at batch=32.
2. **Very slow TEI backend startup after restart/recreate.** A restart/recreate path spent roughly 48 minutes between `Starting model backend` and `Ready`. The logs show a delayed missing-ONNX ORT backend failure before fallback runtime warmup.

The RCA verdict is TEI-first: stabilize and simplify the current TEI runtime path. ONNX should not remain a M042 deliverable and should not be part of the active product/runtime path right now. It remains historical research only until a future milestone explicitly reopens packaging and operational readiness.

## Evidence

### Runtime limits

From direct TEI `/info` during T01:

| Field | Value |
|---|---:|
| model_id | `deepvk/USER-bge-m3` |
| model_dtype | `float32` |
| TEI version | `1.9.3` |
| max_concurrent_requests | `512` |
| max_batch_requests | `4` |
| max_client_batch_size | `32` |
| max_batch_tokens | `16384` |
| max_input_length | `8192` |
| tokenization_workers | `11` |

The advertised request concurrency is high, but batch=32 still spends multi-second time inside TEI's reported queue/scheduler phase.

### Direct TEI sequential timing snapshot

Source: `documents/te-perf-snapshot-m042-s01.md`.

| Batch size | n | ok | client wall p50 | client wall p95 | TEI total p50 | TEI total p95 | queue p50 | queue p95 | inference p50 | inference p95 |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | 10 | 10 | 199.384ms | 306.682ms | 196.283ms | 303.931ms | 0.253ms | 0.300ms | 192.560ms | 299.856ms |
| 8 | 10 | 10 | 1354.058ms | 1575.901ms | 1350.442ms | 1572.372ms | 427.093ms | 496.669ms | 613.751ms | 699.894ms |
| 32 | 10 | 10 | 5342.144ms | 6034.448ms | 5337.374ms | 6027.026ms | 2434.795ms | 2477.019ms | 663.943ms | 724.354ms |

This proves fd HTTP overhead and Redis cache are not the cause. The measurements bypass fd and Redis and read TEI server-side timing logs.

### Restart/recreate timing

Source: `benchmark-results/te-concurrency-profile-m042-s01.md`.

The concurrency profiler attempted five scenarios, but the restart scenario became diagnostic and prevented safe completion of the original profile:

| Scenario | Result |
|---|---|
| sequential batch=32 control | Attempted; detailed comparable evidence exists in T01. |
| parallel 4 x batch=32 | Attempted; metrics lost after container recreate, not treated as success. |
| parallel 16 x batch=1 | Attempted; metrics lost after container recreate, not treated as success. |
| idle 30s then batch=32 | Attempted; metrics lost after container recreate, not treated as success. |
| restart then batch=32 | Failed health wait after 240s; later recreate eventually recovered after ~48 minutes. |

Key log sequence from the successful eventual recovery:

```text
2026-06-14T09:24:49.062231Z  Model artifacts downloaded in 217.816µs
2026-06-14T09:24:49.151051Z  Starting model backend
2026-06-14T10:09:53.752185Z  Could not start ORT backend: ... /onnx/model.onnx does not exist
2026-06-14T10:09:55.334030Z  Warming up model
2026-06-14T10:12:58.805445Z  Ready
```

The artifact download is fast. The long phase is backend startup and missing-ONNX fallback before Candle/safetensors warmup.

### Historical ONNX context

M019 proved only that a tagged ONNX 1024 path was locally performance-viable after earlier quality gates:

| Runtime | Quality status | Best cold latency | Warm latency mean | Max throughput |
|---|---|---:|---:|---:|
| TEI baseline M014 | PASS baseline | `59.0ms` | `2.25ms` | `~750 req/s` |
| Tagged ONNX M019, 1024-token path | legal quality PASS in M018 | `8.3ms` | `1.19ms` | `~858 req/s` |

But M019 explicitly did not prove Docker/CI reproducibility, native tokenizer artifact distribution, deployment safety, long-running stability, or permission to change the production runtime.

M040 then recommended: TEI remains the current/default posture; packaged ONNX is only an explicit same-host opt-in under a contract with runtime metadata, artifact/tokenizer verification, cache namespace isolation, and smoke readiness checks.

The current user direction supersedes M042's original ONNX implementation scope: ONNX should be removed or disabled from the active project path and preserved only as possible future research.

## Hypothesis tree

### H1 — TEI batch scheduler queues large client batches behind internal max batch rules

**Prediction.** Queue time should increase with client batch size even with sequential single-client requests. Inference time may not grow proportionally because the batcher is internally splitting or scheduling work around `max_batch_requests=4` / token limits.

**Evidence.** Supported. T01 shows queue p50:

- batch=1: `0.253ms`
- batch=8: `427.093ms`
- batch=32: `2434.795ms`

This happens without fd, Redis, or client concurrency.

**Verdict.** Strongly plausible. This explains the request-level queue symptom and supports S02 exploring TEI-safe chunking or request shaping, but only after startup behavior is stabilized.

### H2 — The `queue_time` metric includes internal scheduling/model-backend wait, not only external request queueing

**Prediction.** Queue time may be high even without concurrent callers. It may correlate with backend batch assembly or worker availability rather than HTTP queue depth.

**Evidence.** Supported. Sequential batch=32 requests reported ~2.4s queue p50. That cannot be explained by client-side request queueing. It likely includes internal scheduler/backend wait.

**Verdict.** Strongly plausible. Treat TEI `queue_time` as an internal scheduling/backend wait metric, not a simple proof of request overload.

### H3 — Missing ONNX artifact probing delays TEI backend startup and pollutes the active runtime path

**Prediction.** After restart/recreate, logs should spend a long time between `Starting model backend` and fallback warmup, and mention missing ONNX files before readiness.

**Evidence.** Supported. The recreate sequence spent ~45 minutes from `Starting model backend` to `Could not start ORT backend`, then warmed the fallback model and reached `Ready` roughly 3 minutes later. The missing file was `/onnx/model.onnx` under the deepvk snapshot.

**Verdict.** Strongly supported. This is the best explanation for the restart/recreate fragility observed during T02 and a concrete reason to remove ONNX from the current operational path.

### H4 — fd async chunking alone will solve all cold-path latency

**Prediction.** Splitting batch=32 or larger calls into bounded parallel chunks should reduce total fd caller latency without worsening TEI startup behavior.

**Evidence.** Not proven. T02 attempted parallel scenarios, but metrics were lost after destructive restart/recreate. T01 suggests large-batch queue time may improve if smaller chunks avoid TEI's internal scheduler path, but this remains unproven.

**Verdict.** Plausible but unproven. S02 should test this carefully and keep it opt-in. It must not rely on repeated destructive TEI restarts.

### H5 — ONNX should be promoted now because it was faster in M019

**Prediction.** ONNX artifacts, build tags, tokenizer/runtime packaging, Docker/CI, and operational smoke gates should already be ready enough to make ONNX the current mitigation.

**Evidence.** Refuted for current scope. M019 explicitly says ONNX was not production-ready. M040 kept TEI as default. The current user direction says ONNX should be removed/disabled from the active branch and left for future research.

**Verdict.** Rejected for M042.

## Root cause verdict

The immediate root cause is not fd HTTP overhead. The observed cold-path latency comes from TEI runtime behavior:

- request-time internal queue/scheduling grows with batch size;
- process startup/recreate is extremely slow because TEI attempts missing ONNX/ORT backend paths before warming the fallback runtime.

The project-level root cause is scope drift: ONNX remained present as an attractive but unready escape hatch. It complicates TEI startup, documentation, build/runtime decisions, and milestone planning, while the working service path is TEI.

## Recommended action

### For M042

1. **Keep M042 TEI-first.** S03 ONNX implementation is skipped/deferred.
2. **Use S02 to remove or disable active ONNX runtime paths** from current build/config/docs where they confuse operators or affect startup. Keep historical benchmark artifacts only.
3. **Optionally test TEI request shaping/async chunking** after the ONNX cleanup, but require fresh evidence and no destructive restart loop.
4. **Preserve D045 cache-hot M041 performance semantics.** Real TEI miss latency remains diagnostic and should be visible, but not confused with fd cache-hot acceptance.

### For future ONNX research

ONNX may return only as a separate research milestone with explicit gates:

- artifact distribution and checksum contract;
- tokenizer/runtime library packaging;
- Docker/CI build and smoke test;
- isolated cache namespace;
- quality and performance A/B evidence;
- no request-level fallback and no silent runtime switching.

## S01 closure statement

S01 retires the RCA uncertainty enough to guide implementation: do not use ONNX as the M042 mitigation. Focus on TEI operational stability and simplification first.
