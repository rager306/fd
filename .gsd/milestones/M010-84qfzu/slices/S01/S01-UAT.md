# S01: ONNX artifact provenance — UAT

**Milestone:** M010-84qfzu
**Written:** 2026-05-19T18:33:34.984Z

# S01 UAT — ONNX artifact provenance

## Checks

- [x] Local TEI cache/model storage inspected.
- [x] Target model revision and local safetensors/tokenizer/config hashes recorded.
- [x] Upstream `deepvk/USER-bge-m3` model tree checked for ONNX availability.
- [x] ONNX candidate/export paths researched and ranked.
- [x] `S01-RESEARCH.md` states no production runtime change.
- [x] Downstream S02/S03 inputs are defined: comparator first, then FP32 dense-only export/load attempt.

## UAT Result

Pass. A reader can open `S01-RESEARCH.md` and determine the safe next action without relying on chat context: build a dense comparator for current TEI/API, then attempt model-preserving FP32 export from the exact `deepvk/USER-bge-m3` snapshot. No production runtime, quantization, provider stack, model replacement, sparse/ColBERT integration, Rust, or C path is introduced by this slice.

