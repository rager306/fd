# S02: CI artifact provisioning boundary — UAT

**Milestone:** M022-i079tk
**Written:** 2026-05-20T10:48:30.840Z

# S02 UAT — CI artifact provisioning boundary

## Checks

- [x] Workflow validates ONNX artifact metadata in `--allow-missing` mode.
- [x] Workflow fails if ONNX/native/runtime binaries are tracked.
- [x] Workflow triggers include ONNX packaging docs/tooling paths.
- [x] Full ONNX image CI blocker is documented.
- [x] D020 records the CI boundary decision.
- [x] actionlint passes.
- [x] Default Go tests and lint pass.
- [x] Tagged tokenizer and ONNX smoke tests pass.
- [x] Default Docker build passes.
- [x] ONNX Docker image local proof still passes.

## Result

Pass. Hosted CI now checks the safe contract and explicitly defers full ONNX image builds until artifact provisioning exists.

