# M026-ji0i9y: ONNX operational diagnostics implementation

**Vision:** Implement the operational diagnostics promised in M025 so ONNX opt-in startup and health surfaces are debuggable without leaking sensitive data or changing the TEI default path.

## Success Criteria

- ONNX startup diagnostics are implemented and tested.
- Health metadata is safe and opt-in/runtime-aware.
- Default TEI behavior and Docker build remain passing.
- Docs capture implemented status and remaining gaps.
- No production/default switch occurs.

## Slices

- [x] **S01: S01** `risk:high` `depends:[]`
  > After this: After this, ONNX startup failures and health output have safe, actionable diagnostic metadata while TEI health remains unchanged.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, docs and closure evidence show the operational diagnostics implementation status and remaining rollout gaps.

## Boundary Map

Not provided.
