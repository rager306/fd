---
milestone: M042-fjf2en
slice: S01
task: T02
captured: 2026-06-14T10:13:23Z
run_log: benchmark-results/te-concurrency-profile-m042-s01-run.txt
---

# M042 S01 T02 — TEI concurrency and restart profile

This artifact records the attempted TEI concurrency profile and the more important diagnostic result discovered during the run: TEI restart/recreate has a very long cold backend startup path for `deepvk/USER-bge-m3`, dominated by model backend startup/warmup and ONNX-missing fallback checks.

The task was replanned after the restart profile became operationally destructive. No additional TEI restarts are performed in this artifact.

## Service status after recovery

Current `docker compose ps tei` at collection time:

```text
NAME      IMAGE                                                   COMMAND                  SERVICE   CREATED          STATUS                    PORTS
fd_tei    ghcr.io/huggingface/text-embeddings-inference:cpu-1.9   "text-embeddings-rou…"   tei       48 minutes ago   Up 48 minutes (healthy)   0.0.0.0:30080->80/tcp, [::]:30080->80/tcp
```

TEI eventually recovered to `healthy`; the critical finding is how long it took.

## Attempted scenarios

The profiler attempted the following scenarios before failing at the restart wait gate. The run log is `benchmark-results/te-concurrency-profile-m042-s01-run.txt`.

| Scenario | Intended purpose | Status | Evidence |
|---|---|---|---|
| `sequential_batch32_control` | Four sequential direct TEI batch=32 requests as control | Attempted | Run log records start at `2026-06-14T08:56:12Z`; server logs were later lost when the container was recreated. T01 already provides 10 sequential direct TEI batch=32 datapoints. |
| `parallel_4_batch32` | Four simultaneous direct TEI batch=32 requests | Attempted | Run log records start at `2026-06-14T08:56:27Z`; metrics unavailable due later recreate/log loss. |
| `parallel_16_batch1` | Sixteen simultaneous direct TEI batch=1 requests | Attempted | Run log records start at `2026-06-14T08:56:47Z`; metrics unavailable due later recreate/log loss. |
| `idle_30s_batch32` | One batch=32 after 30s idle | Attempted | Run log records start at `2026-06-14T08:57:20Z`; metrics unavailable due later recreate/log loss. |
| `restart_then_batch32` | One batch=32 after TEI restart | Failed as diagnostic | TEI did not become healthy within 240s; later recreate showed backend startup took ~48 minutes before `Ready`. |

The missing scenario metrics are not treated as success. The run is still useful because it exposed TEI cold-start fragility and confirmed that restart profiling is not safe to repeat casually inside the milestone.

## T01 direct TEI control evidence

T01 remains the reliable sequential direct TEI timing evidence because it was captured before the destructive restart and includes parsed server-side timings.

Source: `documents/te-perf-snapshot-m042-s01.md`.

| Batch size | n | ok | client wall p50 | client wall p95 | TEI total p50 | TEI total p95 | queue p50 | queue p95 | inference p50 | inference p95 |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | 10 | 10 | 199.384ms | 306.682ms | 196.283ms | 303.931ms | 0.253ms | 0.300ms | 192.560ms | 299.856ms |
| 8 | 10 | 10 | 1354.058ms | 1575.901ms | 1350.442ms | 1572.372ms | 427.093ms | 496.669ms | 613.751ms | 699.894ms |
| 32 | 10 | 10 | 5342.144ms | 6034.448ms | 5337.374ms | 6027.026ms | 2434.795ms | 2477.019ms | 663.943ms | 724.354ms |

The key pattern is batch-size-sensitive `queue_time` growth even with sequential direct TEI requests.

## Restart and backend startup evidence

### First restart attempt

The profiler restarted TEI at `2026-06-14T08:57:26Z` and waited for `/health`:

```text
Restarting compose service tei at 2026-06-14T08:57:26+00:00
 Container fd_tei  Restarting
 Container fd_tei  Started
RuntimeError: TEI did not become healthy within 240s; last_error=ConnectionResetError: [Errno 104] Connection reset by peer
```

### Recreate/recovery evidence

A later clean compose TEI recreate reached `Ready`, but only after a very long backend path:

```text
2026-06-14T09:24:49.062231Z  INFO ... Model artifacts downloaded in 217.816µs
2026-06-14T09:24:49.151051Z  INFO ... Starting model backend
2026-06-14T10:09:53.752185Z ERROR ... Could not start ORT backend: ... /onnx/model.onnx does not exist
2026-06-14T10:09:55.334030Z  INFO ... Warming up model
2026-06-14T10:12:58.805445Z  INFO ... Ready
```

Measured from `Starting model backend` to `Ready`, TEI took approximately 48 minutes. The ONNX-missing ORT check did not fail fast; it consumed roughly 45 minutes before falling back to Candle/safetensors warmup.

## Interpretation

1. The TEI cold path problem is not just request-level queueing. Startup/backend selection is also operationally expensive and fragile.
2. The current model path attempts ORT/ONNX startup first, fails because `onnx/model.onnx` is absent, and only then warms the fallback runtime. This aligns with prior project memory: TEI logs show ONNX artifacts are missing and TEI/Candle fallback is the current measured runtime.
3. ONNX is therefore actively harmful as a current critical-path concern: it adds fallback delay and conceptual noise while not being the production-ready runtime.
4. T02 cannot honestly claim concurrency metrics for the attempted parallel scenarios because the container was later recreated and those logs were lost. The correct conclusion is to stop destructive restart profiling, preserve T01 timings, and focus the RCA on TEI backend startup/queue behavior.

## T02 verdict

- TEI is currently restored to healthy.
- Restart/recreate profiling exposed a severe cold backend startup cost.
- Further M042 work should be TEI-first: remove ONNX from the current implementation/operational path, document the fallback behavior, and evaluate any TEI mitigation without requiring repeated restarts.
- ONNX runtime implementation is deferred to future research and must not remain a M042 deliverable.
