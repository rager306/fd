# M049 S03 Live Container Proof

Captured: 2026-06-15T13:13:22Z

SUMMARY passed=5 failed=0 total=5

| Check | Result | Detail | Duration ms |
|---|---|---|---|
| GET /health exposes capacity and dependency context | PASS | capacity=0 tei_ms=0.602 redis_ms=0.293 | 35.2 |
| GET /metrics exposes runtime/cache gauges | PASS | runtime/cache gauges present | 1.9 |
| POST /v1/cache/flush is auth protected and works | PASS | unauth=401 auth=200 deleted=37 | 4.6 |
| Embedding cache HIT becomes MISS after flush | PASS | MISS->HIT->flush->MISS ms=125.1/2.3/114.3 | 244.0 |
| Embedding cache HIT becomes MISS after delete | PASS | MISS->HIT->delete->MISS deleted=1 ms=98.2/2.2/93.9 | 197.0 |

## Notes

- API container was rebuilt with `docker compose up -d --build api` before this proof.
- Secret values are intentionally omitted.
- Cache checks prove authenticated MISS/HIT invalidation behavior against the live TEI + Redis stack.
