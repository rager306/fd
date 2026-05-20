# M019-opzh2g: ONNX 1024 performance benchmark

**Vision:** Measure whether the legal-quality-passing ONNX 1024 runtime is fast and stable enough to justify packaging and CI work while keeping TEI as production/default.

## Success Criteria

- ONNX 1024 benchmark is run or blocked with concrete evidence.
- Performance metrics are compared to TEI baseline and prior ONNX benchmark context.
- Benchmark artifact includes sanitized config and no raw text.
- Runtime cleanup leaves no background processes.
- Next gate is explicit.

## Slices

- [x] **S01: S01** `risk:high` `depends:[]`
  > After this: After this, ONNX 1024 has a measured benchmark artifact comparable to the TEI baseline with sanitized config and isolated namespace.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, the project knows whether ONNX 1024 proceeds to packaging/CI or needs performance tuning first.

## Boundary Map

Not provided.
