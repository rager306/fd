# M014-vjfs9f: Tagged ONNX performance benchmark

**Vision:** Measure whether the fixed-probe-correct tagged ONNX path is actually faster or operationally worthwhile compared with the current TEI+Redis runtime, without changing production defaults.

## Success Criteria

- TEI and tagged ONNX benchmark artifacts are produced under comparable, documented conditions.
- Artifacts include sanitized effective configuration and native/ONNX/ORT metadata.
- Redis cache namespace and cache effects are explicit.
- Tagged ONNX is only evaluated after correctness/cosine gate.
- Final recommendation is data-backed and does not switch production default.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, the benchmark matrix and harness changes are defined before running expensive measurements.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, there is a fresh TEI default benchmark artifact to compare against tagged ONNX.

- [ ] **S03: Tagged ONNX benchmark** `risk:high` `depends:[S02]`
  > After this: After this, tagged ONNX performance evidence exists under comparable benchmark conditions.

- [ ] **S04: Benchmark synthesis and decision** `risk:medium` `depends:[S03]`
  > After this: After this, the project has a clear data-backed recommendation for ONNX performance work.

## Boundary Map

| Area | In scope | Out of scope |
|---|---|---|
| Benchmarking | Compare TEI default vs tagged ONNX HF-tokenizer path for cold/warm/batch/cache/startup/memory | Production runtime switch |
| Runtime config | Record build tags, native artifact checksum, ONNX artifact checksum, ORT path/hash, Redis namespace | Secret/env dumping |
| Benchmark tool | Extend or wrap existing benchmark flow for tagged ONNX evidence | Rewrite benchmark framework from scratch unless necessary |
| Quality gate | Require cosine equivalence baseline before benchmark interpretation | Larger full retrieval evaluation beyond fixed probes |
| Operations | Capture startup and memory/RSS evidence for tagged ONNX local server | Full Docker/CI native artifact packaging |
