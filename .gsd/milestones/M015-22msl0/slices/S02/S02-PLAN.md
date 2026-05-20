# S02: Legal retrieval evaluator

**Goal:** Implement a minimal sanitized evaluator that embeds legal docs/queries through two API endpoints and compares retrieval rankings.
**Demo:** After this, there is a reusable local evaluator for the provided legal JSONL file.

## Must-Haves

- Evaluator reads the JSONL corpus.
- It supports TEI and ONNX API URLs and isolated namespace expectations.
- Artifacts include hashes/IDs/metrics, not raw legal text.
- Dry-run or static checks pass.

## Proof Level

- This slice proves: Script compile and small dry-run validation.

## Integration Closure

Provides repeatable quality measurement for S03.

## Verification

- Evaluator records config/corpus hash/runtime labels and avoids raw text output.

## Tasks

- [x] **T01: Implement legal retrieval evaluator** `est:medium`
  Create `tools/evaluate_legal_retrieval.py` to load the 44-ФЗ JSONL, build sanitized docs/queries, call TEI and ONNX APIs, compute top-k overlap and synthetic known-item metrics, and render a no-raw-text markdown artifact.
  - Files: `tools/evaluate_legal_retrieval.py`
  - Verify: `python3 -m py_compile tools/evaluate_legal_retrieval.py` passes.

- [x] **T02: Verify evaluator dry-run hygiene** `est:small`
  Run a non-network dry-run/profile mode of the evaluator against the corpus to verify parsing, doc/query derivation, sanitized output, and no raw text leakage.
  - Files: `benchmark-results/fd-legal-retrieval-dry-run-m015-s02.txt`
  - Verify: Dry-run artifact exists, includes IDs/counts/hash/thresholds, excludes raw legal text, and GitNexus detect_changes passes.

## Files Likely Touched

- tools/evaluate_legal_retrieval.py
- benchmark-results/fd-legal-retrieval-dry-run-m015-s02.txt
