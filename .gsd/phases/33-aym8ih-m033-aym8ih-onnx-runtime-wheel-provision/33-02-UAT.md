# S02: Wheel provisioning documentation and closure — UAT

**Milestone:** M033-aym8ih
**Written:** 2026-05-21T07:24:41.006Z

# UAT — M033 S02

A future operator can read `docs/onnx-artifacts/PROVISIONING.md` and `benchmark-results/fd-onnx-runtime-wheel-provisioning-m033-s02.txt` to see how `--onnx-runtime-source` handles `.whl`/`.zip` sources: only the configured `source_contract.onnx_runtime.library_member` is extracted, size/sha are verified, symlink-like members are rejected, and direct `.so` fallback still works.

The docs and outcome preserve TEI as production/default and ONNX as opt-in experimental.

