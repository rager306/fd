# M037-d23oz4: Target runtime validation contract

**Vision:** Prevent Python helper checks from being overinterpreted as production runtime proof by defining explicit Go and future Rust acceptance gates.

## Success Criteria

- The project no longer relies on Python checks as production runtime acceptance.
- Any regenerated or newly sourced ONNX artifact must pass Go target-runtime gates before promotion.
- Any future Rust backend must pass its own equivalent target-runtime gates.
- TEI remains production/default and ONNX remains opt-in experimental.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, the repo states exactly which Go/API/package gates are required before Python-generated/provisioned ONNX evidence can count for production acceptance.

- [x] **S02: S02** `risk:low` `depends:[]`
  > After this: After this, M037 is verified, summarized, committed locally, and reindexed.

## Boundary Map

Not provided.
