# M023-myx9u4: ONNX 1024 packaged legal quality gate

**Vision:** Validate that the dedicated ONNX 1024 Docker image preserves Russian/legal retrieval parity in a packaged environment while TEI remains the default production runtime.

## Success Criteria

- Packaged ONNX legal quality gate is executed or blocked with evidence.
- Raw legal texts are not printed into artifacts.
- Redis namespace isolation is explicit.
- TEI remains default production runtime.
- Closure verification passes and runtime is clean.

## Slices

- [x] **S01: S01** `risk:high` `depends:[]`
  > After this: After this, the packaged ONNX Docker image has a Russian/legal retrieval gate artifact or a concrete environment blocker.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, the milestone has a decision/outcome and default guardrails are reverified.

## Boundary Map

Not provided.
