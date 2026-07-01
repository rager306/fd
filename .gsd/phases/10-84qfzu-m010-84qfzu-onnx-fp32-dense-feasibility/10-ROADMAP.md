# M010-84qfzu: M010-84qfzu: ONNX FP32 dense feasibility spike

**Vision:** Verify whether a model-preserving BGE-M3 FP32 dense-only ONNX path is feasible on this host, using the M009 benchmark/comparability foundation, without changing production runtime behavior.

## Success Criteria

- A model-preserving BGE-M3 ONNX path is identified or blocked with evidence.
- A dense-output comparator baseline exists for Russian/legal probes.
- Any ONNX artifact attempt records provenance, hashes, output names/shapes, and failure modes.
- No production runtime, model, quantization, provider stack, or language rewrite is introduced.
- Next-step recommendation is evidence-based.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, candidate ONNX/export paths are ranked by provenance, artifact availability, and dense-output compatibility risk.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, fd has a repeatable dense-output comparator for current TEI baseline and future ONNX candidates.

- [x] **S03: S03** `risk:high` `depends:[]`
  > After this: After this, we know whether FP32 BGE-M3 ONNX can be exported/downloaded and loaded on CPU in this environment.

- [x] **S04: S04** `risk:low` `depends:[]`
  > After this: After this, we have a clear decision: proceed to adapter implementation, do more artifact work, or stop ONNX path.

## Boundary Map

## Boundary Map

| Area | In scope | Out of scope |
|---|---|---|
| ONNX artifacts | Identify/export/download candidate artifacts, hashes, tokenizer files, dense output name/shape | Switch production runtime |
| Quality comparator | Build dense-output comparator against current TEI/API baseline for fixed Russian/legal probes | Claim model quality equivalence without retrieval corpus |
| Runtime loading | Prototype local ONNX Runtime CPU EP load/run if artifacts are available | oneDNN/OpenVINO/ZenDNN/INT8 |
| API integration | Optional adapter design only if feasibility is proven | Wire ONNX into `/v1/embeddings` default path |
| Language | Go service remains default | Rust sidecar, C service |
