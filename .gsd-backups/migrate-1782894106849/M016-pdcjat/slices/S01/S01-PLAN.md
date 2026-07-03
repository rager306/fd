# S01: Worst-case divergence profile

**Goal:** Profile the worst M015 legal divergence cases without raw text leakage.
**Demo:** After this, the M015 failure cases are concrete and diagnosable by ID/hash/length/token counts.

## Must-Haves

- Worst M015 IDs are extracted from the artifact.
- Corpus records are resolved by ID/hash.
- Token/length/truncation diagnostics are recorded without raw text.
- No runtime switch occurs.

## Proof Level

- This slice proves: Diagnostic artifact from existing corpus and M015 output.

## Integration Closure

Turns the failed gate into a reproducible diagnostic target set.

## Verification

- Records IDs, hashes, char lengths, tokenizer token counts, truncation flags, and current ONNX max sequence assumptions.

## Tasks

- [x] **T01: Resolve worst divergence IDs** `est:small`
  Extract worst document/query IDs and hashes from `benchmark-results/fd-legal-retrieval-m015-s03.txt`, resolve them back to `tests/44-FZ-2026-articles.jsonl`, and verify IDs/hashes match without printing raw text.
  - Verify: Resolved IDs and text hashes match the M015 artifact.

- [x] **T02: Implement divergence profiler** `est:medium`
  Create `tools/profile_legal_divergence.py` to compute sanitized token/length/truncation diagnostics for the resolved worst cases using the local HF tokenizer and configurable sequence lengths.
  - Files: `tools/profile_legal_divergence.py`
  - Verify: `python3 -m py_compile tools/profile_legal_divergence.py` passes.

- [x] **T03: Run worst-case divergence profile** `est:small`
  Run the profiler against M015 worst cases and write `benchmark-results/fd-legal-divergence-profile-m016-s01.txt`, then verify no raw legal text leaks.
  - Files: `benchmark-results/fd-legal-divergence-profile-m016-s01.txt`
  - Verify: Profile artifact exists, includes token counts/truncation flags, excludes raw text, and GitNexus scope check passes.

## Files Likely Touched

- tools/profile_legal_divergence.py
- benchmark-results/fd-legal-divergence-profile-m016-s01.txt
