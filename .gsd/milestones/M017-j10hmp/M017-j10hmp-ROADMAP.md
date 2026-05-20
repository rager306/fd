# M017-j10hmp: ONNX 512 legal quality gate

**Vision:** Validate the actual tagged Go ONNX runtime at 512 tokens on the Russian legal corpus before any chunking, packaging, performance, or promotion work.

## Success Criteria

- Tagged Go ONNX 512 legal quality gate is run or a concrete blocker is recorded.
- Legal quality metrics are compared to M015 128-token failure and M016 Python 512 diagnostic.
- No raw legal text appears in artifacts.
- TEI remains production/default and ONNX remains opt-in experimental.
- Next remediation step is explicit.

## Slices

- [x] **S01: S01** `risk:high` `depends:[]`
  > After this: After this, there is a measured full legal retrieval artifact for tagged Go ONNX at max sequence length 512, with isolated cache namespace and no raw text leaks.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, the project has a concrete next remediation plan based on the 512-token gate outcome.

## Boundary Map

Not provided.
