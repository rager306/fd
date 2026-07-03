# S01: Tagged ONNX 1024 legal quality gate — UAT

**Milestone:** M018-vq2ttb
**Written:** 2026-05-20T07:39:17.231Z

# S01 UAT — Tagged ONNX 1024 legal quality gate

## Checks

- [x] Required ONNX/native tokenizer/corpus artifacts exist.
- [x] TEI fd API health returned ok.
- [x] Tagged Go ONNX service started on port 18000.
- [x] Service used `ONNX_MAX_SEQUENCE_LENGTH=1024` and isolated namespace `m018-onnx-1024-legal-quality`.
- [x] Full legal retrieval evaluator produced artifact.
- [x] Verdict is PASS.
- [x] Artifact hygiene check passed with `raw_legal_text_leaks=0`.
- [x] Tagged ONNX service was stopped.
- [x] No background processes remain.

## Key result

- Verdict: PASS.
- Minimum cross-backend cosine: `0.99989883`.
- Top-1 agreement: `1.0`.
- Mean overlap@5: `0.997701`.
- ONNX recall ratio: `1.0`.

## UAT Result

Pass. S02 should plan performance/package validation for the 1024 path before any production promotion.

