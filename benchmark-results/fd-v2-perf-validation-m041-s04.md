# fd v2 performance validation — M041 S04

Base URL: `http://localhost:8000`
Overall: FAIL

## Latency cases

| batch | count | p50 ms | p95 ms | p99 ms | threshold ms | errors | verdict |
|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | 50 | 175.089 | 302.094 | 381.255 | 50 | 0 | FAIL |
| 10 | 20 | 1659.292 | 1984.886 | 1984.886 | 200 | 0 | FAIL |
| 32 | 10 | 5128.126 | 5649.237 | 5649.237 | 1000 | 0 | FAIL |

## Sequential and concurrent

- 100 sequential zero errors: PASS (0 errors)
- 4 concurrent × 8 input < 2s: FAIL (4.841s, 0 errors)

## Cache effectiveness

- Repeated input X-Cache HIT and <5ms: FAIL (X-Cache='', latency=2.728ms)

## Raw summary

```json
{
  "latency": [
    {
      "batch": 1,
      "count": 50,
      "threshold_ms": 50,
      "p50": 175.0894351862371,
      "p95": 302.09397058933973,
      "p99": 381.25533098354936,
      "errors": 0,
      "passed": false,
      "error_samples": []
    },
    {
      "batch": 10,
      "count": 20,
      "threshold_ms": 200,
      "p50": 1659.2916268855333,
      "p95": 1984.8855827003717,
      "p99": 1984.8855827003717,
      "errors": 0,
      "passed": false,
      "error_samples": []
    },
    {
      "batch": 32,
      "count": 10,
      "threshold_ms": 1000,
      "p50": 5128.125797957182,
      "p95": 5649.236569181085,
      "p99": 5649.236569181085,
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
    "elapsed_s": 4.841265654191375,
    "errors": 0,
    "passed": false,
    "error_samples": []
  },
  "cache": {
    "first_status": 200,
    "second_status": 200,
    "x_cache": "",
    "hit_latency_ms": 2.7281506918370724,
    "passed": false
  }
}
```
