# M020-mvkq4d: ONNX 1024 artifact contract

**Vision:** Make the ONNX 1024 runtime contract explicit and auditable so packaging/CI work can proceed without confusing local feasibility with production readiness.

## Success Criteria

- 1024 runtime contract is tracked and explicit.
- Validated evidence from M018/M019 is linked.
- No binary artifacts are tracked.
- Production/default remains TEI.
- Next packaging/CI gate is explicit.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, tracked metadata clearly states the ONNX binary was exported with dynamic sequence axes and validated at runtime sequence length 1024, with evidence links and production_default=false.

- [x] **S02: S02** `risk:low` `depends:[]`
  > After this: After this, M020 records the packaging/CI follow-up gate and closes with fresh tests and artifact hygiene.

## Boundary Map

Not provided.
