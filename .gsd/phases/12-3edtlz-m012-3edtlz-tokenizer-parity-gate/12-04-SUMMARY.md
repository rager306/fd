---
id: S04
parent: M012-3edtlz
milestone: M012-3edtlz
provides:
  - Final M012 recommendation.
  - Next milestone scope: HF tokenizers native packaging and opt-in build integration.
requires:
  []
affects:
  []
key_files:
  - .gsd/milestones/M012-3edtlz/slices/S04/S04-RESEARCH.md
  - benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt
key_decisions:
  - Tokenizer parity is solved in isolation via HF Rust tokenizers binding.
  - Runtime integration is gated by native packaging/build tags.
  - ONNX throughput benchmark remains blocked until tagged integration and cosine equivalence pass.
patterns_established:
  - Separate correctness proof from packaging proof.
  - Do not let passing isolated dependency tests bypass build/CI/runtime packaging gates.
  - Keep ONNX benchmark blocked until the actual runtime path is semantically equivalent.
observability_surfaces:
  - S04 research records final gate state and next milestone recommendation.
  - T02 summary records fresh tests/lint/health/artifact/GitNexus evidence.
drill_down_paths:
  - .gsd/milestones/M012-3edtlz/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M012-3edtlz/slices/S04/tasks/T02-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T02:27:25.329Z
blocker_discovered: false
---

# S04: Final parity and packaging decision

**S04 closed M012 with the final decision: tokenizer parity is achievable, but runtime integration needs native packaging/build tags next.**

## What Happened

S04 synthesized the final M012 outcome and verified safety. The milestone answered the tokenizer question: current Go tokenizer fails parity, but `daulet/tokenizers` plus HF Rust `libtokenizers.a` passes all fixed probes. The final decision is not to benchmark performance yet; first add native packaging and build-tag integration so the parity-correct tokenizer can be used by the opt-in ONNX backend without breaking default builds.

## Verification

Fresh final verification passed: Go tests `78 passed`; lint `0 issues`; health ok; artifact/leak check passed; GitNexus low risk with no affected processes.

## Requirements Advanced

- M011-surfaced-tokenizer-parity — Resolved tokenizer parity unknown: a parity-correct Go binding exists, but integration is blocked by packaging.

## Requirements Validated

None.

## New Requirements Surfaced

- A native tokenizer artifact/package/build-tag requirement is needed before ONNX runtime integration.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

M012 closes without ONNX runtime integration because native packaging/build-tag design is outside this milestone and must precede safe integration.

## Known Limitations

No production-ready HF tokenizers packaging exists yet. No runtime code uses the parity-correct tokenizer yet.

## Follow-ups

Next milestone should package HF tokenizers native library behind opt-in build tags, keep default TEI builds unaffected, add tagged parity tests, and only then rerun ONNX cosine/performance gates.

## Files Created/Modified

- `.gsd/milestones/M012-3edtlz/slices/S04/S04-RESEARCH.md` — Final M012 decision research.
