---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Assess Redis embedding data layout and benchmarked env retention policy

Assess Redis data layout and retention policy for embedding vectors: binary value layout, float32 bytes, metadata/key design, long configurable TTL vs no-expire research mode, explicit invalidation, model/version-aware keys, compression tradeoffs, eviction policy, persistence settings, and Redis 7/8 features where relevant. Identify which parameters should become env vars, whether they affect cache correctness or only runtime performance, and which must be recorded in benchmark config snapshots.

## Inputs

- `api/cache/redis.go`
- `api/cache/tiered.go`
- `R002`
- `R003`
- `R004`
- `D003`
- `D004`
- `D005`

## Expected Output

- `S04 T02 summary`

## Verification

Candidate data layout, retention/policy changes, env knobs, and benchmark-recorded fields are ranked with benchmark method and risks.

## Observability Impact

Maps memory/serialization decisions to cached embedding latency and capacity.
