---
id: T01
parent: S01
milestone: M012-3edtlz
key_files:
  - tools/compare_dense_embeddings.py
  - benchmark-results/fd-dense-comparator-m010-s02.txt
key_decisions:
  - Reuse `tools/compare_dense_embeddings.py` `PROBES` as the fixed probe source so embedding and tokenizer comparisons cover the same labels without duplicating text.
  - Tokenizer baseline artifact should include token IDs and attention masks because these are not secrets and exact mismatch debugging needs token-level evidence; raw probe texts remain excluded.
  - Artifact should include package/tool metadata, tokenizer path, local revision/config hashes where available, and `raw_probe_texts_logged=false`.
duration: 
verification_result: passed
completed_at: 2026-05-20T01:56:19.157Z
blocker_discovered: false
---

# T01: Designed the sanitized tokenizer baseline format around existing fixed probes and exact token-level evidence without raw probe text.

**Designed the sanitized tokenizer baseline format around existing fixed probes and exact token-level evidence without raw probe text.**

## What Happened

Inspected the existing dense comparator and baseline artifact. The tokenizer baseline will reuse the existing `PROBES` labels/texts from source for reproducibility but will render only labels, character counts, token lengths, IDs/masks or hashes, and metadata to artifacts. Exact token IDs are allowed in the artifact because they are required to debug parity and do not expose raw text; raw probe strings must not be rendered. The artifact should have configuration JSON, a probe tokenization table, per-probe JSON blocks for exact IDs/masks or hashes, and a PASS verdict.

## Verification

Read the existing comparator source and M010 artifact, then recorded the schema decision in this task summary. The design explicitly excludes raw probe texts.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `read tools/compare_dense_embeddings.py` | 0 | ✅ pass — fixed PROBES and raw-text exclusion pattern identified | 0ms |
| 2 | `read benchmark-results/fd-dense-comparator-m010-s02.txt` | 0 | ✅ pass — existing sanitized artifact pattern identified | 0ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `tools/compare_dense_embeddings.py`
- `benchmark-results/fd-dense-comparator-m010-s02.txt`
