# M024-b8pfpl: ONNX 1024 packaged performance benchmark

**Vision:** Validate whether the packaged ONNX 1024 Docker image preserves local performance viability while TEI remains the default runtime.

## Success Criteria

- Packaged ONNX performance benchmark completes or is blocked with evidence.
- Artifact contains sanitized effective config and comparable metrics.
- Restart/L2 check targets packaged ONNX container.
- Default TEI production behavior remains unchanged.
- Closure verification passes and runtime is clean.

## Slices

- [x] **S01: S01** `risk:high` `depends:[]`
  > After this: After this, there is a packaged ONNX Docker performance artifact with sanitized config and comparable metrics.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, the performance outcome and guardrails are durable and default runtime remains unchanged.

## Boundary Map

Not provided.
