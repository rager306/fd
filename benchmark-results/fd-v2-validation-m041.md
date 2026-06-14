# fd v2 contract validation — M041

Base URL: `http://localhost:8000`
Summary: 45/45 passed, 0 failed
Perf: single p95=2.008ms, batch10 p95=5.783ms

| ID | Result | Description | Detail | Latency ms |
|---|---|---|---|---:|
| T001 | PASS | GET /live returns 200 | status=200, want=200, body=b'{"status":"ok","time":"2026-06-14T08:32:14Z"}' | 26.379 |
| T002 | PASS | GET /ready returns 200 after warmup | status=200, want=200, body=b'{"status":"ready","time":"2026-06-14T08:32:14Z"}' | 2.457 |
| T003 | PASS | GET /health exposes runtime | status=200 | 1.911 |
| T004 | PASS | GET /v1/healthcheck exposes runtime | status=200 | 1.341 |
| T005 | PASS | GET /version exposes service/model | status=200 | 2.346 |
| T006 | PASS | GET /info exposes service model metadata | status=200 | 1.792 |
| T007 | PASS | GET /metrics exposes Prometheus text | status=200 | 1.511 |
| T008 | PASS | GET /warmup returns 200 | status=200, want=200, body=b'{"status":"ready","progress":1}' | 1.844 |
| T009 | PASS | POST /warmup returns 200 or 202 | status=200 | 1.387 |
| T010 | PASS | GET /openapi.json returns OpenAPI 3.1 | status=200 | 3.401 |
| T011 | PASS | GET /docs returns Swagger UI HTML | status=200 | 1.352 |
| T012 | PASS | GET /v1/traces returns JSON array | status=200 | 1.379 |
| T013 | PASS | POST /v1/embeddings string input returns one embedding | status=200 | 247.773 |
| T014 | PASS | Embedding response includes authoritative model | model=deepvk/USER-bge-m3 | 247.773 |
| T015 | PASS | Embedding response includes Server and request id headers | {'Server': 'fd/dev', 'X-Request-Id': '4df7a41f-6546-4acd-810d-9fd0bafb6709'} | 247.773 |
| T016 | PASS | Embedding response includes model/dim/cache headers | {'X-Model-Id': 'deepvk/USER-bge-m3', 'X-Dimensions': '1024', 'X-Cache': 'MISS'} | 247.773 |
| T017 | PASS | Array input returns two embeddings | status=200 | 332.281 |
| T018 | PASS | dimensions=512 returns 512 dimensions | status=200 | 338.993 |
| T019 | PASS | encoding_format=base64 returns base64 string | status=200 | 288.665 |
| T020 | PASS | priority=high is accepted | status=200, want=200, body=b'{"object":"list","data":[{"object":"embedding","embedding":[-0.045140933,0.029008793,0.018454337,-0.0035780196,-0.028233714,0.040667042,0.014471571,-0.026759688' | 263.682 |
| T021 | PASS | user field is accepted | status=200, want=200, body=b'{"object":"list","data":[{"object":"embedding","embedding":[-0.06519551,0.031328265,-0.049919073,-0.047875803,-0.035663005,0.010398183,0.011716613,-0.024761278,' | 194.622 |
| T022 | PASS | First unique embedding request is cache MISS | X-Cache=MISS | 345.996 |
| T023 | PASS | Repeated embedding request is cache HIT | X-Cache=HIT | 2.330 |
| T024 | PASS | Embedding response includes ETag and Cache-Control | {'ETag': '"8d7d494e22d887d6aa0d526d78c1b386487c41f28acb87f1a5ea6a950b43ef4b"', 'Cache-Control': 'public, max-age=86400'} | 2.330 |
| T025 | PASS | Embedding If-None-Match returns 304 | status=304, want=304, body=b'' | 2.139 |
| T026 | PASS | Info If-None-Match returns 304 | status=304, want=304, body=b'' | 2.155 |
| T027 | PASS | CORS preflight returns 204 with headers | status=204 | 1.660 |
| T028 | PASS | Malformed JSON returns invalid_json | status=400, code=invalid_json, want=400/invalid_json | 2.243 |
| T029 | PASS | Missing input returns input_required | status=400, code=input_required, want=400/input_required | 1.983 |
| T030 | PASS | Empty input returns input_required | status=400, code=input_required, want=400/input_required | 2.505 |
| T031 | PASS | Oversized input batch returns batch_too_large | status=413, code=batch_too_large, want=413/batch_too_large | 2.176 |
| T032 | PASS | Invalid dimensions returns dimensions_invalid | status=400, code=dimensions_invalid, want=400/dimensions_invalid | 1.418 |
| T033 | PASS | Invalid encoding_format returns encoding_format_invalid | status=400, code=encoding_format_invalid, want=400/encoding_format_invalid | 1.771 |
| T034 | PASS | Invalid priority returns priority_invalid | status=400, code=priority_invalid, want=400/priority_invalid | 1.632 |
| T035 | PASS | Too long input returns input_too_long | status=413, code=input_too_long, want=413/input_too_long | 1.608 |
| T036 | PASS | GET /v1/embeddings returns method_not_allowed envelope | status=405, code=method_not_allowed, want=405/method_not_allowed | 1.543 |
| T037 | PASS | Unknown route returns not_found | status=404, code=not_found, want=404/not_found | 1.356 |
| T038 | PASS | /v1/batch returns 2x4 embeddings | status=200 | 5.143 |
| T039 | PASS | /v1/batch empty batches returns input_required | status=400, code=input_required, want=400/input_required | 2.199 |
| T040 | PASS | /v1/batch oversized inner batch returns batch_too_large | status=413, code=batch_too_large, want=413/batch_too_large | 4.451 |
| T041 | PASS | Legacy /embeddings/batch returns base64 embeddings | status=200 | 2.048 |
| T042 | PASS | Traces contain recorded request entries | count=100 | 1.970 |
| T043 | PASS | OpenAPI includes /v1/embeddings, /v1/batch, /v1/traces | paths=13 | 3.401 |
| T044 | PASS | Cache-hot single input p95 < 50ms | p95=2.008ms | 2.008 |
| T045 | PASS | Cache-hot 10 input p95 < 200ms | p95=5.783ms | 5.783 |

## JSON

```json
[
  {
    "id": "T001",
    "description": "GET /live returns 200",
    "passed": true,
    "detail": "status=200, want=200, body=b'{\"status\":\"ok\",\"time\":\"2026-06-14T08:32:14Z\"}'",
    "latency_ms": 26.37916198000312
  },
  {
    "id": "T002",
    "description": "GET /ready returns 200 after warmup",
    "passed": true,
    "detail": "status=200, want=200, body=b'{\"status\":\"ready\",\"time\":\"2026-06-14T08:32:14Z\"}'",
    "latency_ms": 2.4567153304815292
  },
  {
    "id": "T003",
    "description": "GET /health exposes runtime",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 1.9109072163701057
  },
  {
    "id": "T004",
    "description": "GET /v1/healthcheck exposes runtime",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 1.3407208025455475
  },
  {
    "id": "T005",
    "description": "GET /version exposes service/model",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 2.346350345760584
  },
  {
    "id": "T006",
    "description": "GET /info exposes service model metadata",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 1.7916667275130749
  },
  {
    "id": "T007",
    "description": "GET /metrics exposes Prometheus text",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 1.51117704808712
  },
  {
    "id": "T008",
    "description": "GET /warmup returns 200",
    "passed": true,
    "detail": "status=200, want=200, body=b'{\"status\":\"ready\",\"progress\":1}'",
    "latency_ms": 1.8440471030771732
  },
  {
    "id": "T009",
    "description": "POST /warmup returns 200 or 202",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 1.3873111456632614
  },
  {
    "id": "T010",
    "description": "GET /openapi.json returns OpenAPI 3.1",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 3.4012109972536564
  },
  {
    "id": "T011",
    "description": "GET /docs returns Swagger UI HTML",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 1.3522687368094921
  },
  {
    "id": "T012",
    "description": "GET /v1/traces returns JSON array",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 1.3785180635750294
  },
  {
    "id": "T013",
    "description": "POST /v1/embeddings string input returns one embedding",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 247.77269922196865
  },
  {
    "id": "T014",
    "description": "Embedding response includes authoritative model",
    "passed": true,
    "detail": "model=deepvk/USER-bge-m3",
    "latency_ms": 247.77269922196865
  },
  {
    "id": "T015",
    "description": "Embedding response includes Server and request id headers",
    "passed": true,
    "detail": "{'Server': 'fd/dev', 'X-Request-Id': '4df7a41f-6546-4acd-810d-9fd0bafb6709'}",
    "latency_ms": 247.77269922196865
  },
  {
    "id": "T016",
    "description": "Embedding response includes model/dim/cache headers",
    "passed": true,
    "detail": "{'X-Model-Id': 'deepvk/USER-bge-m3', 'X-Dimensions': '1024', 'X-Cache': 'MISS'}",
    "latency_ms": 247.77269922196865
  },
  {
    "id": "T017",
    "description": "Array input returns two embeddings",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 332.28070195764303
  },
  {
    "id": "T018",
    "description": "dimensions=512 returns 512 dimensions",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 338.99263199418783
  },
  {
    "id": "T019",
    "description": "encoding_format=base64 returns base64 string",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 288.66526298224926
  },
  {
    "id": "T020",
    "description": "priority=high is accepted",
    "passed": true,
    "detail": "status=200, want=200, body=b'{\"object\":\"list\",\"data\":[{\"object\":\"embedding\",\"embedding\":[-0.045140933,0.029008793,0.018454337,-0.0035780196,-0.028233714,0.040667042,0.014471571,-0.026759688'",
    "latency_ms": 263.681803829968
  },
  {
    "id": "T021",
    "description": "user field is accepted",
    "passed": true,
    "detail": "status=200, want=200, body=b'{\"object\":\"list\",\"data\":[{\"object\":\"embedding\",\"embedding\":[-0.06519551,0.031328265,-0.049919073,-0.047875803,-0.035663005,0.010398183,0.011716613,-0.024761278,'",
    "latency_ms": 194.62236808612943
  },
  {
    "id": "T022",
    "description": "First unique embedding request is cache MISS",
    "passed": true,
    "detail": "X-Cache=MISS",
    "latency_ms": 345.9963980130851
  },
  {
    "id": "T023",
    "description": "Repeated embedding request is cache HIT",
    "passed": true,
    "detail": "X-Cache=HIT",
    "latency_ms": 2.3301858454942703
  },
  {
    "id": "T024",
    "description": "Embedding response includes ETag and Cache-Control",
    "passed": true,
    "detail": "{'ETag': '\"8d7d494e22d887d6aa0d526d78c1b386487c41f28acb87f1a5ea6a950b43ef4b\"', 'Cache-Control': 'public, max-age=86400'}",
    "latency_ms": 2.3301858454942703
  },
  {
    "id": "T025",
    "description": "Embedding If-None-Match returns 304",
    "passed": true,
    "detail": "status=304, want=304, body=b''",
    "latency_ms": 2.1393601782619953
  },
  {
    "id": "T026",
    "description": "Info If-None-Match returns 304",
    "passed": true,
    "detail": "status=304, want=304, body=b''",
    "latency_ms": 2.155332826077938
  },
  {
    "id": "T027",
    "description": "CORS preflight returns 204 with headers",
    "passed": true,
    "detail": "status=204",
    "latency_ms": 1.660351175814867
  },
  {
    "id": "T028",
    "description": "Malformed JSON returns invalid_json",
    "passed": true,
    "detail": "status=400, code=invalid_json, want=400/invalid_json",
    "latency_ms": 2.2427840158343315
  },
  {
    "id": "T029",
    "description": "Missing input returns input_required",
    "passed": true,
    "detail": "status=400, code=input_required, want=400/input_required",
    "latency_ms": 1.9830749370157719
  },
  {
    "id": "T030",
    "description": "Empty input returns input_required",
    "passed": true,
    "detail": "status=400, code=input_required, want=400/input_required",
    "latency_ms": 2.5047771632671356
  },
  {
    "id": "T031",
    "description": "Oversized input batch returns batch_too_large",
    "passed": true,
    "detail": "status=413, code=batch_too_large, want=413/batch_too_large",
    "latency_ms": 2.1758638322353363
  },
  {
    "id": "T032",
    "description": "Invalid dimensions returns dimensions_invalid",
    "passed": true,
    "detail": "status=400, code=dimensions_invalid, want=400/dimensions_invalid",
    "latency_ms": 1.4177169650793076
  },
  {
    "id": "T033",
    "description": "Invalid encoding_format returns encoding_format_invalid",
    "passed": true,
    "detail": "status=400, code=encoding_format_invalid, want=400/encoding_format_invalid",
    "latency_ms": 1.7712670378386974
  },
  {
    "id": "T034",
    "description": "Invalid priority returns priority_invalid",
    "passed": true,
    "detail": "status=400, code=priority_invalid, want=400/priority_invalid",
    "latency_ms": 1.6321190632879734
  },
  {
    "id": "T035",
    "description": "Too long input returns input_too_long",
    "passed": true,
    "detail": "status=413, code=input_too_long, want=413/input_too_long",
    "latency_ms": 1.6076620668172836
  },
  {
    "id": "T036",
    "description": "GET /v1/embeddings returns method_not_allowed envelope",
    "passed": true,
    "detail": "status=405, code=method_not_allowed, want=405/method_not_allowed",
    "latency_ms": 1.5425230376422405
  },
  {
    "id": "T037",
    "description": "Unknown route returns not_found",
    "passed": true,
    "detail": "status=404, code=not_found, want=404/not_found",
    "latency_ms": 1.3564946129918098
  },
  {
    "id": "T038",
    "description": "/v1/batch returns 2x4 embeddings",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 5.142655223608017
  },
  {
    "id": "T039",
    "description": "/v1/batch empty batches returns input_required",
    "passed": true,
    "detail": "status=400, code=input_required, want=400/input_required",
    "latency_ms": 2.1994998678565025
  },
  {
    "id": "T040",
    "description": "/v1/batch oversized inner batch returns batch_too_large",
    "passed": true,
    "detail": "status=413, code=batch_too_large, want=413/batch_too_large",
    "latency_ms": 4.45074588060379
  },
  {
    "id": "T041",
    "description": "Legacy /embeddings/batch returns base64 embeddings",
    "passed": true,
    "detail": "status=200",
    "latency_ms": 2.0479029044508934
  },
  {
    "id": "T042",
    "description": "Traces contain recorded request entries",
    "passed": true,
    "detail": "count=100",
    "latency_ms": 1.9702757708728313
  },
  {
    "id": "T043",
    "description": "OpenAPI includes /v1/embeddings, /v1/batch, /v1/traces",
    "passed": true,
    "detail": "paths=13",
    "latency_ms": 3.4012109972536564
  },
  {
    "id": "T044",
    "description": "Cache-hot single input p95 < 50ms",
    "passed": true,
    "detail": "p95=2.008ms",
    "latency_ms": 2.0080916583538055
  },
  {
    "id": "T045",
    "description": "Cache-hot 10 input p95 < 200ms",
    "passed": true,
    "detail": "p95=5.783ms",
    "latency_ms": 5.783475935459137
  }
]
```
