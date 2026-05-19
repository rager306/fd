# M008-6hnowu: Embedding runtime optimization research

**Vision:** Choose the safest measurable path to optimize embedding performance beyond the current TEI/Candle runtime, Redis cache behavior, Go implementation, and CPU inference provider stack by verifying alternatives before implementation.

## Success Criteria

- Go embedding alternatives are verified from current sources.
- Integration and benchmark risks are documented.
- Redis/cache throughput optimization opportunities are researched and benchmark-scoped.
- Go vs C vs Rust rewrite/targeted-native options are researched and benchmark-scoped.
- ONNX Runtime CPU acceleration and quantization options are researched and benchmark-scoped.
- A measured next-step plan is produced without unverified runtime migration, model replacement, provider-stack change, or language rewrite.
- GSD artifacts are complete and committed locally.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, the suggested Go embedding libraries are classified as viable, risky, or not relevant based on source evidence and Russian/legal model constraints.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, possible integration paths are mapped against the current API and benchmark setup.

- [x] **S04: S04** `risk:medium` `depends:[]`
  > After this: After this, Redis/cache throughput optimization opportunities are ranked for cached embeddings.

- [x] **S05: S05** `risk:medium` `depends:[]`
  > After this: After this, Go/C/Rust rewrite options are classified by likely performance gain and maintenance risk.

- [x] **S06: S06** `risk:medium` `depends:[]`
  > After this: After this, ONNX CPU-level optimization options are verified and benchmark-scoped for the current model.

- [x] **S03: S03** `risk:low` `depends:[]`
  > After this: After this, we have a clear next milestone recommendation for optimization work.

## Boundary Map

Not provided.
