# S05 Research: Go vs C vs Rust Performance Tradeoffs

## Scope

Assess whether fd should remain Go, move selected hot paths, use a sidecar, or rewrite the service in Rust/C for performance.

## Current bottleneck map

Current fd layers:

- Gin HTTP handlers: validation and OpenAI-compatible response shaping.
- Cache layer: L1 local cache + Redis L2 with binary dense embeddings.
- Embedding backend: currently TEI over HTTP.
- Benchmark harness: Python/uv driving API and Redis diagnostics.

Language-sensitive areas:

- JSON/base64/binary encoding and allocation under high cache-hit throughput.
- Batch cache hit loops and Redis round trips.
- Tokenization and in-process ONNX inference if added.
- Native model session/batching/threading glue.

Language-insensitive areas:

- Cold TEI inference latency while model serving is out of process.
- Redis network round trips unless batch access is changed.
- Model quality and corpus relevance.

## Rust option

Rust is credible for a future embedding sidecar or native runtime spike.

Evidence:

- `ort` crate provides Rust ONNX Runtime bindings, session/value APIs, execution provider abstractions, dynamic loading, and build info.
- Rust ecosystem has Axum/Tokio for HTTP, `tokenizers` for HuggingFace-style tokenization, Redis crates/pools, and Candle/Burn as adjacent ML inference options.

Pros:

- Memory-safe native service.
- Good fit for long-running embedding sidecar with tokenization + ORT session + batching.
- Lower operational/security risk than C.
- Can be A/B tested behind the same OpenAI-compatible API contract.

Cons:

- Rewrite cost and duplicated API/cache behavior.
- Still needs ONNX/model/tokenizer packaging.
- Does not improve TEI-dominated path unless inference moves into Rust.

Best use:

- Optional sidecar after ONNX FP32 dense quality and performance baseline exists.
- Not a full service rewrite now.

## C option

C is mature at the ORT boundary.

Evidence:

- ONNX Runtime C API is official.
- It exposes environment/session creation, tensors, execution providers, logging, allocators, graph optimization, thread pool settings, and model-from-array/file APIs.

Pros:

- Maximum ORT control.
- Lowest possible native boundary overhead.
- Good reference for what wrappers must expose.

Cons:

- Manual memory management and higher vulnerability risk.
- Harder JSON/HTTP/concurrency/error handling.
- Harder test/maintainability story for this small service.
- Native packaging remains hard.

Best use:

- Narrow audited FFI only if Go/Rust bindings cannot expose a required ORT feature and profiling proves the need.
- Not a full service rewrite.

## Go option

Go remains appropriate for the current service.

Pros:

- Existing code is small, tested, and observable.
- Good HTTP/cache/concurrency ergonomics.
- go-redis and Gin are adequate for current orchestration.
- Go ONNX path exists through `yalue/onnxruntime_go`.

Cons:

- Native ORT wrappers may lag Rust/C in ergonomics/provider surface.
- In-process tokenization/ONNX glue may need careful allocation and batch tuning.

Best use:

- Keep Go as the primary API/cache service.
- First optimize cache, benchmark comparability, and Go-native ORT adapter if practical.

## Recommendation

Do not rewrite fd in Rust or C now.

Recommended order:

1. Keep Go service.
2. Add benchmark config snapshots and Redis batch/cache metrics.
3. Add pprof/per-layer timing around handler, cache, Redis, TEI/model call, marshal/unmarshal.
4. Run ONNX FP32 dense baseline using current model quality gate.
5. If Go ONNX integration works and bottleneck is model/provider/threading, tune ORT rather than language.
6. If Go binding/native packaging is the constraint or Go allocations dominate in in-process inference, build a Rust sidecar and A/B it.
7. Use C only as narrow FFI/reference for missing ORT controls.

Stop criteria for rewrite ideas:

- If cold latency is dominated by model inference, language rewrite is not the first lever.
- If cached-hit p95 is dominated by Redis round trips, implement MGET/pipeline before rewriting.
- If pprof shows JSON/marshal allocation is material, optimize Go hot path before replacing the service.
- If Rust sidecar does not beat Go+TEI/Go+ONNX by a meaningful margin while preserving quality, keep Go.
