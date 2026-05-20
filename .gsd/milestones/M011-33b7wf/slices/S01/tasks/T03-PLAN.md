---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Document artifact validation contract

Document the artifact contract and failure expectations in S01 research: where local artifacts live, how checksum validation should behave, what to do when the file is missing, and why production runtime is unchanged.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `.gsd/milestones/M010-84qfzu/slices/S04/S04-RESEARCH.md`

## Expected Output

- `.gsd/milestones/M011-33b7wf/slices/S01/S01-RESEARCH.md`

## Verification

Research artifact exists and states missing/checksum mismatch behavior plus no production runtime change.

## Observability Impact

Gives future implementers the exact diagnostic/failure contract before runtime code changes.
