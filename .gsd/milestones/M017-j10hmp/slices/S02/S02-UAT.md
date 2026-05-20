# S02: Quality outcome decision — UAT

**Milestone:** M017-j10hmp
**Written:** 2026-05-20T07:31:34.958Z

# S02 UAT — Quality outcome decision

## Checks

- [x] Outcome assessment artifact exists.
- [x] M015, M016, and M017 metrics are compared.
- [x] Decision says 512-token ONNX is necessary but insufficient.
- [x] D015 saved in the decision register.
- [x] TEI remains production/default.
- [x] ONNX remains opt-in experimental.
- [x] Artifact hygiene passed with `raw_legal_text_leaks=0`.
- [x] Python scripts compile.
- [x] Default Go tests pass.
- [x] Pinned GolangCI-Lint reports 0 issues.
- [x] Tagged HF tokenizer tests pass.
- [x] No background processes remain.

## UAT Result

Pass. M017 can close as a measured quality-gate milestone whose outcome is that 512-token ONNX still fails strict vector equivalence and needs chunking or longer-sequence handling next.

