---
id: M012-3edtlz
title: "Tokenizer parity gate"
status: complete
completed_at: 2026-05-20T02:28:54.319Z
key_decisions:
  - Current `sugarme/tokenizer` path is not semantically valid for ONNX equivalence.
  - `daulet/tokenizers` with HF Rust `libtokenizers.a` is parity-correct on fixed probes.
  - Runtime integration must be gated by native packaging/build tags.
  - Do not benchmark ONNX throughput yet.
key_files:
  - tools/compare_tokenizers.py
  - benchmark-results/fd-tokenizer-baseline-m012-s01.txt
  - benchmark-results/fd-tokenizer-go-current-m012-s02.txt
  - benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt
  - .gsd/milestones/M012-3edtlz/slices/S04/S04-RESEARCH.md
lessons_learned:
  - Token-level parity is the right diagnostic layer for embedding cosine mismatches.
  - HF Rust tokenizer bindings can match exactly where pure-Go `sugarme` diverges.
  - Native dependency packaging should be treated as a separate gate from correctness.
---

# M012-3edtlz: Tokenizer parity gate

**M012 proved tokenizer parity is achievable via HF Rust tokenizers bindings, while preserving TEI default and blocking runtime integration until native packaging is designed.**

## What Happened

M012 resolved the tokenizer parity gate after M011. S01 created a safe Hugging Face baseline for fixed Russian/legal probes. S02 compared the current Go `sugarme/tokenizer` path against that baseline and proved it fails all five probes. S03 tested a Go binding to Hugging Face's Rust tokenizers implementation and proved exact parity for all probes using `daulet/tokenizers` plus `libtokenizers.a`; runtime integration was deferred because native packaging/build tags are required to avoid breaking default builds. S04 synthesized the final decision: next work should package and integrate the parity-correct tokenizer behind opt-in build tags before ONNX cosine or throughput benchmarking resumes.

## Success Criteria Results

- HF tokenizer baseline: met.
- Go tokenizer comparison: met.
- Exact parity or blocker: met; exact parity candidate found, runtime packaging blocker recorded.
- No invalid performance benchmark: met.
- TEI default: met.
- No raw probe text/large artifact commits: met.

## Definition of Done Results

- S01-S04 complete: met.
- HF baseline exists: met.
- Current Go tokenizer comparison exists: met and fails as expected.
- HF Rust binding parity exists: met and passes all probes.
- No runtime default switch: met.
- Final verification gates: met.
- GitNexus scope: low/no affected processes at final gate.

## Requirement Outcomes

- Tokenizer parity requirement: resolved in isolation, not yet integrated.
- Runtime safety: preserved.
- Benchmark validity: preserved by blocking throughput until runtime parity/cosine gates pass.
- New packaging requirement: surfaced for next milestone.

## Deviations

M012 did not integrate the parity-correct tokenizer into runtime because native packaging/build-tag design is required first. This is a deliberate safety boundary, not unfinished work.

## Follow-ups

Plan M013 for HF tokenizers native packaging and opt-in build integration: build tags, libtokenizers artifact manifest/checksum, Docker/CI support, tagged token parity tests, then rerun ONNX cosine before performance benchmarks.
