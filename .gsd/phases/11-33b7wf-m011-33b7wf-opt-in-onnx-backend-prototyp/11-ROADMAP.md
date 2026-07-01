# M011-33b7wf: M011-33b7wf: Opt in ONNX backend prototype

**Vision:** Build a gated, non-default ONNX backend prototype around the proven M010 FP32 dense artifact, while preserving TEI as the production default and collecting enough evidence to decide whether deeper integration is worth pursuing.

## Success Criteria

- ONNX backend remains opt-in and TEI remains default.
- Artifact checksum/config validation exists before ONNX load.
- A dense ONNX backend prototype either works locally or is blocked with evidence.
- TEI vs ONNX comparator and performance evidence are persisted.
- No large ONNX artifacts are committed.
- No production runtime switch, model replacement, INT8, provider variants, or language rewrite occurs.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, ONNX artifacts have a manifest/checksum contract and missing artifact failure expectations before any runtime wiring.

- [x] **S02: S02** `risk:high` `depends:[]`
  > After this: After this, the Go service has a backend selection seam with TEI still default and ONNX rejected or disabled unless explicitly configured.

- [x] **S03: S03** `risk:high` `depends:[]`
  > After this: After this, an explicit ONNX backend either returns dense embeddings locally or is blocked with concrete binding/runtime evidence.

- [x] **S04: S04** `risk:medium` `depends:[]`
  > After this: After this, the project has evidence that Go ONNX runtime integration loads/runs but is blocked on tokenizer parity, plus a recommendation for the next milestone.

## Boundary Map

## Boundary Map

| Area | In scope | Out of scope |
|---|---|---|
| Runtime selection | Add explicit non-default ONNX backend selection and keep TEI default | Switch production default to ONNX |
| Artifacts | Define local/runtime artifact manifest, checksum validation, missing artifact errors | Commit ONNX model files to git |
| Go API | Introduce a small embedder seam if needed and wire ONNX behind opt-in config | Rewrite handlers/cache or change API response contract |
| ONNX execution | Load existing M010 dense FP32 artifact through ONNX Runtime CPU EP or document blocker | INT8, OpenVINO, oneDNN, ZenDNN, sparse, ColBERT |
| Benchmarks | Compare TEI default and ONNX opt-in on current fixed probes and benchmark path | Claim production superiority without larger corpus/perf evidence |
