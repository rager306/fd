# M033-aym8ih: ONNX Runtime wheel provisioning

**Vision:** Make the already-selected ONNX Runtime wheel source candidate actually provisionable while preserving ONNX rollout blockers and TEI default.

## Success Criteria

- M031 ONNX Runtime wheel source candidate becomes usable by provisioning tooling.
- No production/default ONNX promotion occurs.
- Security hardening from M029/M030 remains intact.
- Exact ONNX model binary blocker remains explicit.
- No external state changes occur.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, `tools/provision_onnx_artifacts.py` supports a safe ONNX Runtime wheel-member extraction path and has local probes proving it.

- [x] **S02: S02** `risk:low` `depends:[]`
  > After this: After this, docs/outcome/decision describe the new wheel provisioning support and remaining blockers, then M033 is verified and committed locally.

## Boundary Map

Not provided.
