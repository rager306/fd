# M039-aexhf5: Packaged Go ONNX target runtime rerun

**Vision:** Move from local Go target-runtime proof to packaged Docker Go ONNX proof for the current artifact setup.

## Success Criteria

- Packaged Docker Go ONNX evidence exists for current artifact setup.
- Evidence is through actual packaged Go endpoint, not Python-only runtime.
- Redis namespaces and benchmark side effects are explicit.
- TEI remains production/default and ONNX remains opt-in experimental.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, a dedicated packaged ONNX Docker image has been built or verified and smoke-tested through its Go API.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, packaged ONNX legal/performance evidence is current and milestone closes cleanly with runtime SHA verification enabled.

## Boundary Map

Not provided.
