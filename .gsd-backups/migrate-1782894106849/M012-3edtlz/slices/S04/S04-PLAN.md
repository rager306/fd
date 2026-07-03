# S04: Final parity and packaging decision

**Goal:** Close M012 with the correct decision: tokenizer parity is achievable through HF Rust tokenizers bindings, but ONNX runtime integration remains gated by native packaging and build-tag design.
**Demo:** After this, M012 closes with a decision: tokenizer parity is solved in isolation, but ONNX runtime integration is gated by native packaging/build tags before performance benchmarking.

## Must-Haves

- Final recommendation states tokenizer parity is solved in isolation.
- Final recommendation states runtime ONNX integration is still blocked by native packaging/build tags.
- No ONNX throughput benchmark is recommended yet.
- Default TEI runtime is verified healthy.
- Quality/artifact/GitNexus gates pass.

## Proof Level

- This slice proves: Research synthesis plus fresh verification gates.

## Integration Closure

Provides the next milestone recommendation and prevents invalid jump to ONNX performance benchmarking.

## Verification

- Records final gate status, passing parity artifact, packaging blocker, and verification evidence.

## Tasks

- [x] **T01: Synthesize final tokenizer parity decision** `est:small`
  Write final S04 research synthesis: summarize S01 baseline, S02 current mismatch, S03 HF binding pass, and next milestone recommendation for build-tag/native packaging integration.
  - Files: `.gsd/milestones/M012-3edtlz/slices/S04/S04-RESEARCH.md`
  - Verify: Research artifact exists and names parity pass plus packaging blocker.

- [x] **T02: Verify M012 final gate state** `est:small`
  Run final M012 verification: Go tests, lint, default health, artifact parser/leak checks, and GitNexus detect_changes.
  - Verify: All final gates pass and milestone can be validated.

## Files Likely Touched

- .gsd/milestones/M012-3edtlz/slices/S04/S04-RESEARCH.md
