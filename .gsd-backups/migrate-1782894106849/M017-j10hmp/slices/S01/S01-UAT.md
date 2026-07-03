# S01: Tagged ONNX 512 legal quality gate — UAT

**Milestone:** M017-j10hmp
**Written:** 2026-05-20T07:28:39.970Z

# S01 UAT — Tagged ONNX 512 legal quality gate

## Checks

- [x] Required ONNX/native tokenizer/corpus artifacts exist.
- [x] TEI fd API health returned ok.
- [x] Tagged Go ONNX service started on port 18000.
- [x] Service used `ONNX_MAX_SEQUENCE_LENGTH=512` and isolated namespace `m017-onnx-512-legal-quality`.
- [x] Full legal retrieval evaluator produced artifact.
- [x] Artifact hygiene check passed with `raw_legal_text_leaks=0`.
- [x] Tagged ONNX service was stopped.
- [x] No background processes remain.

## Key result

- Verdict: FAIL under strict cosine threshold.
- Minimum cross-backend cosine: `0.98982302`.
- Top-1 agreement: `1.0`.
- Mean overlap@5: `0.997701`.
- ONNX recall ratio: `1.0`.

## UAT Result

Pass as a measurement slice. The 512-token gate is a major improvement but not sufficient for strict legal equivalence; S02 should plan chunking or longer-sequence handling.

