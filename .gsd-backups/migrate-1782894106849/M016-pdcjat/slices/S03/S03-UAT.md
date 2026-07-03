# S03: Remediation path decision — UAT

**Milestone:** M016-pdcjat
**Written:** 2026-05-20T07:20:06.153Z

# S03 UAT — Remediation path decision

## Checks

- [x] Remediation plan artifact exists.
- [x] S01/S02 metrics are included.
- [x] Options and tradeoffs are compared.
- [x] Recommendation keeps TEI production/default and ONNX experimental.
- [x] D014 saved in the decision register.
- [x] Artifact hygiene check reports `raw_legal_text_leaks=0`.
- [x] Python diagnostic scripts compile.
- [x] Default Go tests pass.
- [x] Pinned GolangCI-Lint reports 0 issues.
- [x] Tagged HF tokenizer tests pass.

## UAT Result

Pass. M016 can close as an investigation/remediation-planning milestone. Implementation of 512-token ONNX and chunking should be a future milestone.

