---
id: T02
parent: S05
milestone: M008-6hnowu
key_files: []
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:10:45.267Z
blocker_discovered: false
---

# T02: Rust is a credible ONNX embedding sidecar option, but not a justified full rewrite without pprof/per-layer bottleneck evidence.

**Rust is a credible ONNX embedding sidecar option, but not a justified full rewrite without pprof/per-layer bottleneck evidence.**

## What Happened

Researched Rust option maturity and likely gain. Rust has a credible stack for an embedding sidecar or future service: Axum/Tokio for HTTP, `ort` for ONNX Runtime bindings, `tokenizers` for HuggingFace-compatible tokenization, Redis crates/pools for cache access, and Candle/Burn as adjacent ML inference frameworks. `ort` 2.0.0-rc.12 documents ONNX Runtime bindings, session APIs, execution provider abstractions, dynamic loading, and build/version info. This is more mature and operationally safer than C for an ML-serving spike because Rust keeps memory safety while still allowing native inference and low-overhead serialization. Likely gain is highest if fd moves inference in-process or to a Rust sidecar that owns tokenization + ORT session + batching. Likely gain is low if fd remains mostly Go orchestration calling TEI over HTTP and Redis over TCP, because the expensive work is outside Go. Recommendation: do not full-rewrite fd now; if ONNX FP32 dense benchmark passes quality gates and Go binding/packaging is painful or overhead-heavy, consider a Rust sidecar implementing the same OpenAI-compatible embedding contract for A/B comparison.

## Verification

Read `ort` docs and Rust embedding/ONNX ecosystem search results; compared against current fd bottleneck map from T01.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Fetched: https://docs.rs/ort/latest/ort/` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Read: https://ort.pyke.io/ search result summary` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Search query: Rust ONNX Runtime ort crate tokenizers redis axum embedding service performance maturity 2026` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

Rust may improve native inference/service overhead only if Go orchestration or TEI boundary is measured as bottleneck. It does not automatically improve the current TEI-dominated cold path, and a Rust rewrite would duplicate a working Go API/cache surface.

## Files Created/Modified

None.
