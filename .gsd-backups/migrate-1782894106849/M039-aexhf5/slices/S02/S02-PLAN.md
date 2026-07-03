# S02: Packaged ONNX closure

**Goal:** Rerun packaged legal and performance gates and close M039.
**Demo:** After this, packaged ONNX legal/performance evidence is current and milestone closes cleanly with runtime SHA verification enabled.

## Must-Haves

- Legal gate passes or blocker recorded.
- Performance benchmark passes or blocker recorded.
- Outcome matrix records passed/skipped/remaining gates.
- No external action occurred.
- Working tree clean after commit/reindex.

## Proof Level

- This slice proves: Packaged endpoint legal/perf drivers plus final guardrails.

## Integration Closure

Aligns packaged runtime evidence with M037/M038 target-runtime contract and S01 Docker smoke proof.

## Verification

- Records packaged legal/performance metrics, runtime SHA verification, namespaces, and remaining rollout blockers.

## Tasks

- [x] **T01: Run packaged legal gate** `est:medium`
  Start packaged image `fd-api:onnx1024-m039` with `ONNX_RUNTIME_SHA256` and isolated legal namespace, run legal retrieval evaluator against TEI/default and packaged ONNX endpoints, stop container, and verify artifact safety.
  - Files: `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt`
  - Verify: Legal evaluator exits 0, metrics pass, artifact contains no raw legal text/secrets/signed URLs, container stopped or prepared for next task.

- [x] **T02: Run packaged performance benchmark** `est:medium`
  Start packaged image `fd-api:onnx1024-m039` with `ONNX_RUNTIME_SHA256` and isolated benchmark namespace, run benchmark.py against packaged ONNX endpoint, stop container, and verify artifact safety.
  - Files: `benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt`
  - Verify: Benchmark exits 0, artifact metrics present, artifact contains no raw text/secrets/signed URLs, container stopped, port clean.

- [x] **T03: Record packaged acceptance matrix** `est:small`
  Write packaged ONNX acceptance matrix for M039 covering image, smoke/rerun, legal, performance, skipped gates, non-actions, and blockers.
  - Files: `benchmark-results/fd-onnx-docker-target-runtime-acceptance-m039-s02.txt`
  - Verify: Outcome artifact checks pass.

- [x] **T04: Close M039** `est:medium`
  Run final guardrails, GitNexus detect, GSD validation/completion, checkpoint, commit, reindex, and clean-state checks.
  - Verify: Final guardrails pass, milestone completes, commit created, GitNexus reindex/detect clean, working tree clean.

## Files Likely Touched

- benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt
- benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt
- benchmark-results/fd-onnx-docker-target-runtime-acceptance-m039-s02.txt
