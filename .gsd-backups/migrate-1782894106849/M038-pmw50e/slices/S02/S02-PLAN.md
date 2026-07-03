# S02: Go target runtime closure

**Goal:** Run available Go endpoint legal/performance drivers, record outcome, and close M038.
**Demo:** After this, Go target-runtime proof is either expanded to legal/performance drivers or records a truthful blocker, then milestone is closed.

## Must-Haves

- Legal/performance drivers run against actual Go endpoint when dependencies are available, or blocker is documented.
- Outcome includes sanitized config and no raw text.
- Final checks pass.
- No external action occurred.

## Proof Level

- This slice proves: Driver commands against actual endpoints plus final guardrails.

## Integration Closure

Aligns Go endpoint evidence with M037 contract and benchmark/legal tooling.

## Verification

- Records which target-runtime gates passed and which remain blocked.

## Tasks

- [x] **T01: Run Go target-runtime legal gate** `est:medium`
  Start Go ONNX API again with a fresh isolated namespace and run selected Russian/legal retrieval evaluator against TEI API on 8000 and Go ONNX API on 18000. Stop server after run.
  - Files: `benchmark-results/fd-legal-retrieval-m038-go-onnx-target-runtime.txt`
  - Verify: Legal evaluator passes or records blocker; raw legal text not logged; server stopped.

- [x] **T02: Run Go target-runtime performance gate** `est:medium`
  Start Go ONNX API with a fresh isolated namespace and run a bounded performance driver against the actual Go ONNX endpoint. Stop server after run.
  - Files: `benchmark-results/fd-benchmark-m038-go-onnx-target-runtime.txt`
  - Verify: Performance benchmark passes or records blocker; sanitized config present; server stopped.

- [x] **T03: Record Go target-runtime acceptance outcome** `est:small`
  Summarize target-runtime proof coverage and remaining blockers in outcome artifact.
  - Files: `benchmark-results/fd-onnx-go-target-runtime-acceptance-m038-s02.txt`
  - Verify: Outcome checks pass; no raw text, secrets, signed URLs, or production promotion claims.

- [x] **T04: Run final guardrails** `est:medium`
  Run final milestone guardrails: JSON/tooling checks, actionlint, Go tests/lint, tagged tests, Docker default build, leak checks, binary hygiene, port/background checks, GitNexus detect.
  - Verify: All final checks pass.

- [x] **T05: Prepare post-slice closure** `est:small`
  Record closure ordering and defer milestone validation/completion/checkpoint/commit/reindex to post-slice sequence.
  - Verify: Task records that post-slice closure will run after S02 completion.

## Files Likely Touched

- benchmark-results/fd-legal-retrieval-m038-go-onnx-target-runtime.txt
- benchmark-results/fd-benchmark-m038-go-onnx-target-runtime.txt
- benchmark-results/fd-onnx-go-target-runtime-acceptance-m038-s02.txt
