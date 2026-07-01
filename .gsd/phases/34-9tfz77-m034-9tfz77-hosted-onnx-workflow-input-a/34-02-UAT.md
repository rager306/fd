# S02: Hosted workflow input contract documentation — UAT

**Milestone:** M034-9tfz77
**Written:** 2026-05-21T07:59:58.933Z

# UAT — M034 S02

A future operator can read `docs/onnx-artifacts/PROVISIONING.md` and `benchmark-results/fd-onnx-workflow-input-alignment-m034-s02.txt` to determine:

- required inputs: `onnx_source_url`, `native_tokenizer_source_url`;
- optional inputs: `tokenizer_json_source_url`, `onnx_runtime_source_url`, `onnx_runtime_sha256`, `image_tag`;
- `onnx_runtime_sha256` is an optional override because manifest sha is used when omitted;
- signed/plain secret URLs are not allowed;
- workflow dispatch needs explicit user approval;
- exact ONNX model binary source remains the blocking input.

