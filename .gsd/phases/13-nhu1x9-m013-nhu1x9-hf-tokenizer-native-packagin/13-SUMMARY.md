---
id: M013-nhu1x9
title: "HF tokenizer native packaging gate"
status: complete
completed_at: 2026-05-20T04:01:50.106Z
key_decisions:
  - Tagged ONNX path with HF native tokenizer is fixed-probe benchmark-ready.
  - TEI remains default and production runtime.
  - No production switch until Docker/CI packaging, pinned native artifact, larger Russian/legal quality, and benchmark evidence pass.
key_files:
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
  - api/embed/hf_tokenizer_native.go
  - api/embed/onnx.go
  - api/embed/onnx_tokenizer_default.go
  - api/embed/onnx_tokenizer_hf.go
  - benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt
  - benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt
  - benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt
  - .gsd/milestones/M013-nhu1x9/slices/S04/S04-RESEARCH.md
lessons_learned:
  - Native dependency correctness and native dependency packaging are separate gates.
  - Build tags can preserve default TEI builds while enabling experimental ONNX dependencies.
  - Once token parity and cosine pass, performance benchmarking becomes meaningful; before that it is misleading.
---

# M013-nhu1x9: HF tokenizer native packaging gate

**M013 made the parity-correct HF tokenizer available in tagged ONNX builds and proved fixed-probe cosine equivalence while preserving TEI default builds.**

## What Happened

M013 turned M012's isolated tokenizer parity finding into a safe tagged runtime path. It created a native artifact manifest for `libtokenizers.a`, added validation and ignore protections, introduced `hf_tokenizers` build-tag files so default builds remain clean, integrated the HF native tokenizer into tagged ONNX builds, and ran a tagged local ONNX API comparison. Fixed-probe TEI-vs-ONNX cosine passed at ~0.999993 for all probes with isolated Redis namespace. M013 closes with the tagged path benchmark-ready but not production-ready.

## Success Criteria Results

- Native artifact contract exists: met.
- Default builds unaffected: met.
- Opt-in build-tag path proven: met.
- Tokenizer parity and cosine gates passed: met.
- No native binary/raw text/default switch committed: met.

## Definition of Done Results

- Native artifact contract: met.
- Default build isolation: met.
- Tagged native tokenizer parity: met.
- Tagged ONNX cosine equivalence: met.
- No binary/raw text leaks: met.
- Final verification gates: met.
- GitNexus scope: met.

## Requirement Outcomes

- M012 native packaging/build-tag requirement: validated for local tagged builds.
- Benchmark validity requirement: advanced; tagged ONNX path is now valid for performance benchmarking.
- Production readiness: still deferred pending Docker/CI packaging, pinning, broader quality, and benchmark evidence.

## Deviations

M013 stopped at benchmark readiness and did not implement Docker/CI packaging or performance benchmarking. This is intentional; those are separate follow-up milestones.

## Follow-ups

Plan next milestone for tagged ONNX performance benchmarking. Include native artifact checksum/build tags in benchmark snapshots, isolate Redis namespaces, compare TEI default vs tagged ONNX for cold/warm/batch/cache/startup/memory, and do not switch production default.
