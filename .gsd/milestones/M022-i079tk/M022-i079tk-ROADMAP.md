# M022-i079tk: ONNX 1024 dedicated Docker packaging

**Vision:** Prove or concretely bound a dedicated ONNX 1024 Docker packaging path using externally provisioned artifacts, while preserving the default TEI image and avoiding fake CI claims.

## Success Criteria

- Opt-in ONNX Docker packaging path is defined and locally validated or blocked with evidence.
- Default Docker build remains passing.
- Artifact verifier is part of the ONNX packaging path.
- No binary artifacts are tracked.
- CI boundary is documented without pretending local artifacts exist in CI.

## Slices

- [x] **S01: S01** `risk:high` `depends:[]`
  > After this: After this, there is an opt-in ONNX Docker packaging path that verifies local artifacts and preserves the default TEI Docker path.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, CI limitations and next automation steps are explicit, with no fake CI support claims.

## Boundary Map

Not provided.
