# M018-vq2ttb: ONNX 1024 legal quality gate

**Vision:** Measure whether a 1024-token tagged Go ONNX runtime passes the Russian legal strict quality gate, so the project can choose between longer-sequence remediation and deterministic chunking based on evidence.

## Success Criteria

- Tagged Go ONNX 1024 legal quality is measured or blocked with evidence.
- Outcome is compared to 128 and 512 results.
- No raw legal text appears in artifacts.
- TEI remains production/default.
- Next gate is explicit.

## Slices

- [x] **S01: S01** `risk:high` `depends:[]`
  > After this: After this, there is a measured full legal retrieval artifact for tagged Go ONNX at max sequence length 1024, with isolated cache namespace and sanitized output.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, M018 records whether 1024 is enough for quality and what the next milestone should be.

## Boundary Map

Not provided.
