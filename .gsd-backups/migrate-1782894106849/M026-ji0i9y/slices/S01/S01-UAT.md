# S01: Startup diagnostics and health metadata — UAT

**Milestone:** M026-ji0i9y
**Written:** 2026-05-20T12:06:32.337Z

# S01 UAT — Startup diagnostics and health metadata

## Checks

- [x] Default health shape remains status/time only.
- [x] ONNX health metadata includes safe runtime fields.
- [x] Health metadata excludes path-like manifest/runtime/tokenizer fields.
- [x] Manifest validation carries validated max sequence length.
- [x] Config rejects sequence length above validated contract.
- [x] Startup logs safe ONNX preflight metadata.
- [x] Default Go tests pass.
- [x] Lint passes.
- [x] Tagged tests pass.
- [x] Default Docker build passes.

## Result

Pass. ONNX operational diagnostics are implemented without changing TEI/default behavior.

