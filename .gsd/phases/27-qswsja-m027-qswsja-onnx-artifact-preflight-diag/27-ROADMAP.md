# M027-qswsja: M027-qswsja: ONNX artifact preflight diagnostics

**Vision:** Close the next ONNX operational diagnostic gaps by validating tokenizer JSON, optional ONNX Runtime sha, and provider configuration at startup.

## Success Criteria

- Remaining preflight diagnostics implemented and tested.
- Default TEI runtime remains unaffected.
- ONNX remains opt-in experimental.
- Docs and decisions updated.
- No binaries tracked; verification passes.

## Slices

- [x] **S01: S01** `risk:high` `depends:[]`
  > After this: After this, ONNX startup has stricter artifact preflight for tokenizer JSON, optional runtime library sha, and provider config.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, docs/outcome/decision capture the completed preflight diagnostics and remaining security/rollout blockers.

## Boundary Map

## Boundary Map

Not provided.
