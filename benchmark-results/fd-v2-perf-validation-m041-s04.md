# fd v2 performance validation — M041 S04

Base URL: `http://localhost:8000`
Overall: FAIL

## Latency cases

| batch | count | p50 ms | p95 ms | p99 ms | threshold ms | errors | verdict |
|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 50 | 198.107 | 425.410 | 481.560 | 50 | 0 | FAIL |
| 10 | 20 | 1721.886 | 1993.474 | 1993.474 | 200 | 0 | FAIL |
| 32 | 10 | 5198.921 | 5670.721 | 5670.721 | 1000 | 0 | FAIL |

## Sequential and concurrent

- 100 sequential zero errors: PASS (0 errors)
- 4 concurrent × 8 input < 2s: FAIL (6.987s, 0 errors)

## Cache effectiveness

- Repeated input X-Cache HIT and <5ms: PASS (X-Cache='HIT', latency=2.039ms)

## Raw summary

```json
{
  "latency": [
    {
      "batch": 1,
      "count": 50,
      "threshold_ms": 50,
      "p50": 198.10728775337338,
      "p95": 425.40979385375977,
      "p99": 481.5601799637079,
      "errors": 0,
      "passed": false,
      "error_samples": []
    },
    {
      "batch": 10,
      "count": 20,
      "threshold_ms": 200,
      "p50": 1721.8861137516797,
      "p95": 1993.474179878831,
      "p99": 1993.474179878831,
      "errors": 0,
      "passed": false,
      "error_samples": []
    },
    {
      "batch": 32,
      "count": 10,
      "threshold_ms": 1000,
      "p50": 5198.921129107475,
      "p95": 5670.721202623099,
      "p99": 5670.721202623099,
      "errors": 0,
      "passed": false,
      "error_samples": []
    }
  ],
  "sequential": {
    "count": 100,
    "errors": 0,
    "passed": true,
    "error_samples": []
  },
  "concurrent": {
    "workers": 4,
    "batch": 8,
    "elapsed_s": 6.986875114031136,
    "errors": 0,
    "passed": false,
    "error_samples": []
  },
  "cache": {
    "first_status": 200,
    "second_status": 200,
    "x_cache": "HIT",
    "hit_latency_ms": 2.0385878160595894,
    "passed": true
  }
}
```
