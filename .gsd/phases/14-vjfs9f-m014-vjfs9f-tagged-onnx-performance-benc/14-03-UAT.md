# S03: Tagged ONNX benchmark — UAT

**Milestone:** M014-vjfs9f
**Written:** 2026-05-20T04:31:42.576Z

# S03 UAT — Tagged ONNX benchmark

## Checks

- [x] ONNX Runtime, ONNX model, native tokenizer, tokenizer JSON, manifests, and M013 cosine artifact exist and have recorded checksums.
- [x] Tagged test passed: `go test -tags hf_tokenizers ./embed`.
- [x] Tagged ONNX server started on port 18000 and `/health` returned ok.
- [x] Startup evidence captured; RSS was about 1.69 GiB shortly after startup.
- [x] Benchmark harness supports `BENCHMARK_API_RESTART_COMMAND` for non-Compose API restart.
- [x] Final artifact exists: `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`.
- [x] Artifact includes snapshot_version 3 and tagged runtime/native/ONNX/ORT metadata.
- [x] Artifact includes isolated cache namespace `m014-onnx-hf-tokenizer`.
- [x] Raw fixed-probe text leak check returned 0.
- [x] M013 cosine correctness gate referenced and passed.
- [x] Local tagged server stopped and transient runtime files removed.

## Key ONNX Results

- Best cold latency: 10.2ms.
- Warm latency mean: 1.63ms.
- Max throughput: ~891 req/s at 4 concurrent.
- Redis L2 after API restart: 2.70ms.
- Batch L1 p95: 5.62ms.
- Batch L2 p95: 4.41ms.
- Chunk reuse warm p95: 9.00ms.

## UAT Result

Pass. S04 can compare TEI and tagged ONNX evidence.

