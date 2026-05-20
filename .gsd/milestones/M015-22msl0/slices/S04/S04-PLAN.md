# S04: Quality verdict and closure

**Goal:** Synthesize the quality gate result, record decision, and run final verification.
**Demo:** After this, the project knows whether ONNX can proceed to packaging/tuning or needs quality work.

## Must-Haves

- Gate verdict is recorded.
- Decision says continue packaging/tuning or block ONNX quality.
- Final tests/lint/tagged tests/artifact hygiene/GitNexus checks pass.
- TEI remains production/default.

## Proof Level

- This slice proves: Decision plus final verification gates.

## Integration Closure

Closes the legal quality gate milestone and updates next-step priority.

## Verification

- Captures verdict, caveats, and future dataset/technical requirements.

## Tasks

- [x] **T01: Summarize quality gate verdict** `est:small`
  Write a concise quality gate synthesis artifact that explains PASS/FAIL, key metrics, likely cause, and next recommended milestone.
  - Files: `benchmark-results/fd-legal-retrieval-m015-summary.txt`
  - Verify: Summary artifact exists and includes verdict/metrics/caveats/no raw text.

- [x] **T02: Record blocking quality decision** `est:small`
  Record a GSD decision that ONNX packaging/tuning is blocked until long-text legal quality divergence is investigated.
  - Files: `.gsd/DECISIONS.md`
  - Verify: Decision saved and references M015 evidence.

- [x] **T03: Run final M015 verification gates** `est:medium`
  Run final verification gates after last changes: py_compile evaluator, dry-run hygiene, artifact hygiene, Go tests, pinned lint, tagged tests, runtime cleanup, GitNexus detect.
  - Verify: All commands pass or expected quality-fail artifact is verified; no background process remains.

## Files Likely Touched

- benchmark-results/fd-legal-retrieval-m015-summary.txt
- .gsd/DECISIONS.md
