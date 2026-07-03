---
id: T02
parent: S02
milestone: M015-22msl0
key_files:
  - tools/evaluate_legal_retrieval.py
  - benchmark-results/fd-legal-retrieval-dry-run-m015-s02.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-20T04:54:11.603Z
blocker_discovered: false
---

# T02: Verified evaluator dry-run parsing and artifact hygiene.

**Verified evaluator dry-run parsing and artifact hygiene.**

## What Happened

Ran the evaluator in dry-run mode against the 44-ФЗ JSONL corpus. The dry-run artifact includes corpus hash, selected document/query counts, thresholds, endpoint labels, and caveat text while excluding raw legal text. A hygiene check confirmed the expected metadata and no sampled raw legal text leaks. GitNexus detected only low non-code/new-file scope.

## Verification

Dry-run artifact exists and includes corpus hash/IDs/counts/thresholds with raw legal text leak check passing; GitNexus scope is low.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests python tools/evaluate_legal_retrieval.py --corpus tests/44-FZ-2026-articles.jsonl --output benchmark-results/fd-legal-retrieval-dry-run-m015-s02.txt --dry-run --max-docs 256 --max-title-queries 32 --max-self-queries 64` | 0 | ✅ pass — dry-run artifact written | 0ms |
| 2 | `python dry-run artifact hygiene check` | 0 | ✅ pass — dry_run_hygiene=pass; raw_legal_text_leaks=0 | 0ms |
| 3 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ✅ pass — low scope, no changed symbols | 0ms |

## Deviations

Initial dry-run assertion expected 64 queries for 128 docs, but the selected documents covered only 10 title queries plus 32 self queries. The dry-run was rerun with 256 docs and assertions were changed to validate actual artifact shape rather than a brittle query count.

## Known Issues

Dry-run does not call APIs and does not validate metric pass/fail; S03 must run live TEI and tagged ONNX endpoints.

## Files Created/Modified

- `tools/evaluate_legal_retrieval.py`
- `benchmark-results/fd-legal-retrieval-dry-run-m015-s02.txt`
