---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Map language-sensitive fd bottleneck layers

Map fd execution layers and identify where implementation language could matter: HTTP routing, JSON/base64 encode, cache marshal/unmarshal, Redis client round trips, TEI/native inference calls, Docker/runtime overhead. Use existing code and prior benchmarks.

## Inputs

- `api/handlers/embeddings.go`
- `api/handlers/batch.go`
- `api/cache/redis.go`
- `api/cache/tiered.go`
- `benchmark.py`

## Expected Output

- `S05 T01 summary`

## Verification

Layer map identifies measurable bottlenecks and non-language bottlenecks.

## Observability Impact

Prevents generic Go/Rust/C claims from being applied to the wrong bottleneck.
