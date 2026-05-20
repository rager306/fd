---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Run worst-case divergence profile

Run the profiler against M015 worst cases and write `benchmark-results/fd-legal-divergence-profile-m016-s01.txt`, then verify no raw legal text leaks.

## Inputs

- `tools/profile_legal_divergence.py`
- `benchmark-results/fd-legal-retrieval-m015-s03.txt`
- `tests/44-FZ-2026-articles.jsonl`

## Expected Output

- `benchmark-results/fd-legal-divergence-profile-m016-s01.txt`

## Verification

Profile artifact exists, includes token counts/truncation flags, excludes raw text, and GitNexus scope check passes.

## Observability Impact

Produces the diagnostic artifact that S02 will use for focused runtime experiments.
