# S02: Legal retrieval evaluator — UAT

**Milestone:** M015-22msl0
**Written:** 2026-05-20T04:54:34.871Z

# S02 UAT — Legal retrieval evaluator

## Checks

- [x] `tools/evaluate_legal_retrieval.py` exists.
- [x] `python3 -m py_compile tools/evaluate_legal_retrieval.py` passed.
- [x] Dry-run artifact exists.
- [x] Dry-run artifact includes corpus SHA256 and config.
- [x] Dry-run artifact excludes sampled raw legal text.
- [x] GitNexus scope check low.

## UAT Result

Pass. S03 can run the live TEI-vs-tagged-ONNX quality gate.

