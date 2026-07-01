# S01: Workflow runtime source alignment — UAT

**Milestone:** M034-9tfz77
**Written:** 2026-05-21T07:53:18.230Z

# UAT — M034 S01

A future operator can provide `onnx_runtime_source_url` without `onnx_runtime_sha256`. The workflow will not fail validation; it will log that provisioning uses manifest `source_contract.onnx_runtime.library_sha256`. If an override sha is supplied, the workflow passes it to provisioning.

