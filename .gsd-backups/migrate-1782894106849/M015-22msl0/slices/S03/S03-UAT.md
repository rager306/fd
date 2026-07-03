# S03: Run legal retrieval quality gate — UAT

**Milestone:** M015-22msl0
**Written:** 2026-05-20T05:03:22.597Z

# S03 UAT — Live legal retrieval quality gate

## Checks

- [x] TEI/default API health ok.
- [x] Tagged ONNX API started on port 18000 with namespace `m015-onnx-legal-quality`.
- [x] Evaluator ran live against both APIs.
- [x] Evaluator ID fallback bug fixed and rerun.
- [x] Artifact exists: `benchmark-results/fd-legal-retrieval-m015-s03.txt`.
- [x] Artifact excludes raw legal text.
- [x] Tagged ONNX server stopped.
- [x] No background processes remain.

## Gate Result

FAIL.

Key metrics:

- Top-1 agreement: 0.977011.
- Mean overlap@5: 0.908046.
- ONNX recall ratio: 0.991422.
- Document cross-backend cosine min: 0.369489.
- Query cross-backend cosine min: 0.656612.

## UAT Result

The gate executed successfully, but the ONNX path failed the strict legal quality threshold.

