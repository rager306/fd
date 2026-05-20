# S01: Corpus profile and gate design

**Goal:** Profile the 44-FZ JSONL corpus and define a first-pass retrieval parity gate that does not overclaim absolute quality.
**Demo:** After this, the legal corpus and evaluation method are explicit and auditable.

## Must-Haves

- Corpus count/schema profile exists.
- Raw text is not dumped into artifacts.
- Gate metrics are defined: TEI-vs-ONNX top-k overlap/rank correlation plus synthetic known-item checks.
- Caveat about missing qrels is explicit.

## Proof Level

- This slice proves: Corpus profiling artifact plus plan.

## Integration Closure

Turns the user-provided file into a measurable quality gate contract.

## Verification

- Records corpus hash, counts, length distributions, query/doc derivation rules, and caveats.

## Tasks

- [x] **T01: Profile legal corpus** `est:small`
  Profile `tests/44-FZ-2026-articles.jsonl` for schema, counts, invalid flags, length distribution, and corpus hash without printing raw text beyond inspected samples.
  - Files: `benchmark-results/fd-legal-corpus-profile-m015-s01.txt`
  - Verify: Profile artifact exists and includes counts/hash/no raw text dump.

- [x] **T02: Define retrieval parity contract** `est:small`
  Define the quality gate metrics and acceptance thresholds for the unlabeled 44-FZ JSONL corpus.
  - Verify: Task summary names metrics, thresholds, and caveats.

## Files Likely Touched

- benchmark-results/fd-legal-corpus-profile-m015-s01.txt
