# S05: Research Go vs C vs Rust performance tradeoffs

**Goal:** Assess whether replacing Go components with C or Rust could materially improve fd throughput, considering ecosystem maturity and current bottleneck layers.
**Demo:** After this, Go/C/Rust rewrite options are classified by likely performance gain and maintenance risk.

## Must-Haves

- Compare Go, Rust, and C for HTTP/JSON/base64, Redis client/cache path, and ONNX/native inference integration.
- Distinguish full rewrite from targeted native module/sidecar options.
- Include ecosystem maturity: libraries, ONNX Runtime bindings, Redis clients, web frameworks, build/deploy complexity, memory safety.
- Estimate likely gain ranges only with caveats and required profiling evidence.
- Define benchmark gates before any rewrite.

## Proof Level

- This slice proves: source research plus architecture assessment

## Integration Closure

Prevents a language rewrite from being recommended unless the expected gain is tied to a measured bottleneck and benchmark plan.

## Verification

- Defines profiling evidence needed before considering rewrites and what runtime metrics would prove improvement.

## Tasks

- [x] **T01: Map language-sensitive fd bottleneck layers** `est:small`
  Map fd execution layers and identify where implementation language could matter: HTTP routing, JSON/base64 encode, cache marshal/unmarshal, Redis client round trips, TEI/native inference calls, Docker/runtime overhead. Use existing code and prior benchmarks.
  - Verify: Layer map identifies measurable bottlenecks and non-language bottlenecks.

- [x] **T02: Research Rust option maturity and likely gain** `est:medium`
  Research Rust ecosystem maturity for fd-relevant paths: HTTP frameworks, Redis clients, serde/base64, ONNX Runtime/Candle/tokenizers integration, deployment complexity, and performance evidence.
  - Verify: Sources read and Rust option classified by layer, risk, and benchmarkability.

- [x] **T03: Research C option maturity and likely gain** `est:medium`
  Research C ecosystem maturity for fd-relevant paths: HTTP libraries, hiredis, cJSON/simdjson-style parsing options, ONNX Runtime C API, tokenizers constraints, memory safety risks, build/deploy complexity, and where C could realistically beat Go/Rust.
  - Verify: Sources read and C option classified by layer, risk, and benchmarkability.

- [x] **T04: Recommend language/runtime rewrite strategy** `est:small`
  Produce Go vs C vs Rust recommendation with likely gain ranges, required profiling evidence, benchmark gates, and explicit no-rewrite conditions.
  - Files: `.gsd/milestones/M008-6hnowu/slices/S05/S05-RESEARCH.md`
  - Verify: Research artifact includes recommendation, benchmark gates, and exclusions.

## Files Likely Touched

- .gsd/milestones/M008-6hnowu/slices/S05/S05-RESEARCH.md
