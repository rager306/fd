# S01: Worst-case divergence profile — UAT

**Milestone:** M016-pdcjat
**Written:** 2026-05-20T05:35:58.175Z

# S01 UAT — Worst-case divergence profile

## Checks

- [x] Worst M015 IDs extracted.
- [x] IDs resolved against `tests/44-FZ-2026-articles.jsonl`.
- [x] Text hashes match M015 artifact.
- [x] Profiler compiles.
- [x] Profile artifact exists.
- [x] Raw legal text leak check passed.

## Key result

- Cases: 17.
- Resolved cases: 17.
- Hashes match: true.
- Cases truncated at 128: 17.
- Cases truncated at 512: 2.
- Max tokens with specials: 697.

## UAT Result

Pass. S02 can focus on max sequence length/truncation root-cause diagnostics.

