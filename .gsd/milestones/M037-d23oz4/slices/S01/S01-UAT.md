# S01: Target runtime validation contract — UAT

**Milestone:** M037-d23oz4
**Written:** 2026-05-21T10:17:55.215Z

# UAT — M037 S01

A future operator can inspect `docs/onnx-artifacts/PROVISIONING.md`, `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`, and `benchmark-results/fd-onnx-target-runtime-validation-contract-m037-s01.txt` to see:

- Python helper checks are not production runtime proof;
- Go fd API/package gates required before ONNX promotion;
- any future Rust backend requires equivalent independent validation;
- Redis namespace isolation is mandatory for comparisons.

