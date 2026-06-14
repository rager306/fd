# fd v2 performance validation — M041 S04

Base URL: `http://localhost:8000`
Mode: cache-hot steady-state after explicit prewarm; real cache-miss inference is diagnostic only.
Overall: PASS

## Cache-hot latency cases

| batch | count | prewarm X-Cache | p50 ms | p95 ms | p99 ms | threshold ms | errors | non-HIT | verdict |
|---:|---:|---|---:|---:|---:|---:|---:|---:|---|
| 1 | 50 | MISS | 1.280 | 2.236 | 2.514 | 50 | 0 | 0 | PASS |
| 10 | 20 | MISS | 2.828 | 3.468 | 3.468 | 200 | 0 | 0 | PASS |
| 32 | 10 | MISS | 6.894 | 7.595 | 7.595 | 1000 | 0 | 0 | PASS |

## Sequential and concurrent

- 100 sequential cache-hot zero errors: PASS (0 errors, 0 non-HIT responses)
- 4 concurrent × 8 cache-hot inputs < 2s: PASS (0.010s, 0 errors, 0 non-HIT responses)

## Cache effectiveness

- Repeated input X-Cache HIT and <5ms: PASS (first X-Cache='MISS', second X-Cache='HIT', latency=1.870ms)

## Non-blocking cache-miss diagnostics

| batch | status | X-Cache | latency ms |
|---:|---:|---|---:|
| 1 | 200 | MISS | 235.128 |
| 10 | 200 | MISS | 2106.853 |
| 32 | 200 | MISS | 6795.926 |

## Raw summary

```json
{
  "mode": "cache-hot",
  "miss_diagnostics": [
    {
      "batch": 1,
      "status": 200,
      "x_cache": "MISS",
      "latency_ms": 235.1275011897087,
      "ok": true
    },
    {
      "batch": 10,
      "status": 200,
      "x_cache": "MISS",
      "latency_ms": 2106.8529090844095,
      "ok": true
    },
    {
      "batch": 32,
      "status": 200,
      "x_cache": "MISS",
      "latency_ms": 6795.92579882592,
      "ok": true
    }
  ],
  "latency": [
    {
      "batch": 1,
      "count": 50,
      "threshold_ms": 50,
      "prewarm_status": 200,
      "prewarm_x_cache": "MISS",
      "p50": 1.2803617864847183,
      "p95": 2.2359960712492466,
      "p99": 2.5136130861938,
      "errors": 0,
      "non_hit_responses": 0,
      "passed": true,
      "error_samples": []
    },
    {
      "batch": 10,
      "count": 20,
      "threshold_ms": 200,
      "prewarm_status": 200,
      "prewarm_x_cache": "MISS",
      "p50": 2.8276038356125355,
      "p95": 3.4677847288548946,
      "p99": 3.4677847288548946,
      "errors": 0,
      "non_hit_responses": 0,
      "passed": true,
      "error_samples": []
    },
    {
      "batch": 32,
      "count": 10,
      "threshold_ms": 1000,
      "prewarm_status": 200,
      "prewarm_x_cache": "MISS",
      "p50": 6.894320249557495,
      "p95": 7.595321629196405,
      "p99": 7.595321629196405,
      "errors": 0,
      "non_hit_responses": 0,
      "passed": true,
      "error_samples": []
    }
  ],
  "sequential": {
    "count": 100,
    "prewarm_status": 200,
    "prewarm_x_cache": "MISS",
    "errors": 0,
    "non_hit_responses": 0,
    "passed": true,
    "error_samples": []
  },
  "concurrent": {
    "workers": 4,
    "batch": 8,
    "elapsed_s": 0.01014166371896863,
    "prewarm_statuses": [
      200,
      200,
      200,
      200
    ],
    "prewarm_x_cache": [
      "MISS",
      "MISS",
      "MISS",
      "MISS"
    ],
    "errors": 0,
    "non_hit_responses": 0,
    "passed": true,
    "error_samples": []
  },
  "cache": {
    "first_status": 200,
    "first_x_cache": "MISS",
    "second_status": 200,
    "x_cache": "HIT",
    "hit_latency_ms": 1.870427280664444,
    "passed": true
  }
}
```
