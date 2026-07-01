---
id: S02
parent: M012-3edtlz
milestone: M012-3edtlz
provides:
  - Current Go tokenizer mismatch artifact.
  - S03 recommendation and decision tree.
requires:
  []
affects:
  - S03
key_files:
  - tools/compare_tokenizers.py
  - benchmark-results/fd-tokenizer-go-current-m012-s02.txt
  - .gsd/milestones/M012-3edtlz/slices/S02/S02-RESEARCH.md
key_decisions:
  - Current `sugarme/tokenizer` path fails parity for all fixed probes.
  - Special-token toggles do not explain the mismatch.
  - S03 should test HF Rust tokenizers Go bindings before modifying runtime code.
patterns_established:
  - Expected-fail comparator artifacts are acceptable when they encode blocker evidence.
  - Tokenization mismatch should be proven at token level, not inferred from embedding cosine alone.
  - Dependency candidates should be isolated before runtime integration.
observability_surfaces:
  - `benchmark-results/fd-tokenizer-go-current-m012-s02.txt` records exact mismatch indexes and token windows without raw text.
  - S02 research records candidate dependency and packaging risks.
drill_down_paths:
  - .gsd/milestones/M012-3edtlz/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M012-3edtlz/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M012-3edtlz/slices/S02/tasks/T03-SUMMARY.md
  - .gsd/milestones/M012-3edtlz/slices/S02/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T02:12:24.315Z
blocker_discovered: false
---

# S02: Go tokenizer candidate comparison

**S02 proved the current Go tokenizer is not equivalent to Hugging Face and identified Rust-backed HF tokenizers bindings as the next candidate.**

## What Happened

S02 converted the M011 tokenizer suspicion into durable evidence. The tokenizer comparison tool now compares current Go `sugarme/tokenizer` output against the S01 Hugging Face baseline and writes a mismatch artifact. All five fixed probes fail parity, with token counts and IDs diverging. Research found the tokenizer JSON uses XLM-R-style Unigram/Metaspace/TemplateProcessing behavior and recommended testing Go bindings to Hugging Face's Rust tokenizers rather than hand-rolling or patching blindly.

## Verification

Python compile passed; `go-current` comparison writes expected-fail artifact; parser/leak checks passed; GitNexus affected processes are comparator-tool-only.

## Requirements Advanced

- M011-surfaced-tokenizer-parity — Retired the unknown of whether current `sugarme/tokenizer` can be used as-is; it cannot for fixed-probe parity.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S02 did not integrate a candidate tokenizer dependency; it intentionally stopped at mismatch proof and research so S03 can test the dependency in isolation first.

## Known Limitations

The recommended candidate (`daulet/tokenizers` or related bindings) likely introduces CGO/native library packaging complexity. That risk is not retired until S03.

## Follow-ups

S03 should first test `github.com/daulet/tokenizers` in isolation against the S01 baseline. If it passes parity and packaging is tractable, then integrate; otherwise block and consider sidecar/defer.

## Files Created/Modified

- `tools/compare_tokenizers.py` — Extended tokenizer comparator with go-current mode and mismatch rendering.
- `benchmark-results/fd-tokenizer-go-current-m012-s02.txt` — Current Go tokenizer mismatch artifact.
- `.gsd/milestones/M012-3edtlz/slices/S02/S02-RESEARCH.md` — Research and recommendation for S03 tokenizer path.
