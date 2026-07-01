# S01: ONNX Runtime wheel provisioning support — UAT

**Milestone:** M033-aym8ih
**Written:** 2026-05-21T07:16:07.521Z

# UAT — M033 S01

A future operator can pass an ONNX Runtime `.whl`/`.zip` source to `tools/provision_onnx_artifacts.py`. If `source_contract.onnx_runtime.library_member` is present, the helper extracts that exact member into `.gsd/runtime/onnxruntime/libonnxruntime.so.1.26.0`, verifies size/sha, and reports a structured `onnx_runtime` result. Missing/oversized/wrong-checksum members fail before use.

