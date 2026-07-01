# M021-4t2wpt: M021-4t2wpt: ONNX 1024 Docker CI packaging contract

**Vision:** Make ONNX 1024 packaging safe to automate by defining how artifacts are provisioned and verified without breaking the default TEI Docker/CI path.

## Success Criteria

- ONNX 1024 Docker/CI artifact provisioning contract exists.
- Default TEI Docker/build path remains unaffected.
- No large native/ONNX binary is tracked.
- Packaging/CI next gate is explicit and evidence-backed.
- Verification passes with no background processes.

## Slices

- [x] **S01: S01** `risk:high` `depends:[]`
  > After this: After this, artifact staging/checksum validation is documented and/or scripted for ONNX 1024 and native tokenizer without tracking binaries.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, local Docker/CI checks prove the default path is not broken and the ONNX packaging path has a clear next gate or local proof.

## Boundary Map

## Boundary Map

Not provided.
