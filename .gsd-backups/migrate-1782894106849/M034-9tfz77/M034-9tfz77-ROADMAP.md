# M034-9tfz77: Hosted ONNX workflow input alignment

**Vision:** Make the manual hosted ONNX workflow input contract match the local provisioning capabilities without overclaiming hosted proof or changing runtime defaults.

## Success Criteria

- M033 ONNX Runtime wheel provisioning is usable from the manual workflow input surface.
- Runtime artifacts remain checksum verified via CLI sha or manifest sha.
- Exact ONNX model source blocker remains explicit.
- No external state changes occur.
- No production/default ONNX promotion occurs.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, the workflow no longer blocks runtime wheel provisioning solely because `onnx_runtime_sha256` is omitted when manifest metadata supplies it.

- [x] **S02: S02** `risk:low` `depends:[]`
  > After this: After this, docs/outcome/decision describe safe future dispatch inputs and remaining blockers, then M034 is verified and committed locally.

## Boundary Map

Not provided.
