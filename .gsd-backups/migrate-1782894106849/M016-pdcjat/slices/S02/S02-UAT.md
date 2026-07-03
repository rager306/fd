# S02: Sequence length root cause diagnostics — UAT

**Milestone:** M016-pdcjat
**Written:** 2026-05-20T07:03:46.115Z

# S02 UAT — Sequence length root cause diagnostics

## Checks

- [x] Diagnostic script compiles.
- [x] TEI `/health` returned ok before the run.
- [x] Local ONNX artifact and tokenizer inputs existed.
- [x] Diagnostic ran against TEI API and local ONNX Runtime.
- [x] Artifact includes 128 vs 512 cosine summaries.
- [x] Raw legal text leak check passed.

## Key result

- Sequence length 128: 17/17 cases truncated, mean cosine `0.9204953`, min cosine `0.74389759`.
- Sequence length 512: 2/17 cases truncated, mean cosine `0.99885631`, min cosine `0.99026363`.

## UAT Result

Pass. S03 should plan remediation around 512-token ONNX plus chunking or longer-sequence handling for cases exceeding 512 tokens.

