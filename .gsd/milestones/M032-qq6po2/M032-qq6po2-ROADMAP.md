# M032-qq6po2: ONNX model reproducibility source proof

**Vision:** Convert the remaining exact ONNX model artifact blocker into an executable local source/reproducibility proof and a precise next-step contract, without promoting ONNX or changing TEI defaults.

## Success Criteria

- Existing ONNX export contract can be verified locally with one command.
- Claims remain bounded: no byte-for-byte regeneration claim unless actually regenerated.
- TEI remains production/default and ONNX remains opt-in experimental.
- No large artifact is committed; ignored runtime artifacts remain ignored.
- No push or workflow dispatch occurs.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, the repo has an executable local verifier for the existing ONNX export contract and a clear statement of what it proves and does not prove.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, docs/manifests/outcome describe the exact ONNX model source options and remaining blocker without overclaiming.

## Boundary Map

Not provided.
