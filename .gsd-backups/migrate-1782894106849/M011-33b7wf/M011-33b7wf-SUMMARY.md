---
id: M011-33b7wf
title: "Opt in ONNX backend prototype"
status: complete
completed_at: 2026-05-20T01:53:30.612Z
key_decisions:
  - Close M011 as a blocked prototype with evidence, not as production-ready ONNX integration.
  - Keep TEI as default runtime.
  - Do not benchmark ONNX speed until tokenizer parity is solved.
  - Future backend comparisons must isolate Redis cache namespace.
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - api/embed/onnx_manifest.go
  - api/embed/onnx.go
  - api/embed/onnx_test.go
  - api/main.go
  - api/main_test.go
  - benchmark-results/fd-go-onnx-m011-s03.txt
  - .gsd/milestones/M011-33b7wf/slices/S04/S04-RESEARCH.md
lessons_learned:
  - A shared Redis namespace can make ONNX appear equivalent by returning TEI vectors from cache; isolate cache namespaces for backend comparisons.
  - ONNX Runtime load/run success is not semantic correctness; tokenizer parity is part of the embedding function.
  - Evidence-backed blocker closure prevents invalid performance claims.
---

# M011-33b7wf: Opt in ONNX backend prototype

**M011 added a safe opt-in Go ONNX prototype and closed it as blocked on tokenizer parity, with TEI still default and verified.**

## What Happened

M011 built the gated ONNX backend prototype around the M010 exact-model dense FP32 artifact. S01 added a tracked artifact manifest/checksum contract without committing the large ONNX binary. S02 added the opt-in backend selection seam while preserving TEI as the default. S03 implemented and wired a real Go ONNX Runtime path that can load and run the local artifact, then discovered a semantic blocker: Go tokenization via `sugarme/tokenizer` diverges from Hugging Face tokenization for Russian legal text, causing TEI-vs-Go-ONNX cosine failures under an isolated cache namespace. S04 synthesized the blocker, verified safety, and recommended tokenizer parity as the next milestone before any ONNX performance benchmark.

## Success Criteria Results

- ONNX backend remains opt-in and TEI remains default: met.
- Artifact checksum/config validation exists before ONNX load: met.
- Dense ONNX backend prototype either works locally or is blocked with evidence: met as blocked-with-evidence.
- TEI vs ONNX comparator evidence persisted: met with failed isolated-cache artifact.
- No large ONNX artifacts committed: met.
- No production runtime switch/model replacement/INT8/provider/language rewrite: met.

## Definition of Done Results

- TEI remains default: met.
- ONNX opt-in only: met via `EMBEDDING_BACKEND=onnx` plus required ONNX env vars.
- Artifact validation: met via manifest validator and checksum verification.
- Local ONNX path: met as load/run, blocked at tokenizer parity.
- Evidence persisted: met with S03 comparison artifact and S04 research.
- No large artifacts committed: met by tracked-artifact scan.
- Verification gates: Go tests, lint, Compose, health, manifest checks, GitNexus all passed or reported acceptable low-risk scope.

## Requirement Outcomes

- TEI default preservation: maintained.
- Model-preserving optimization: maintained; no model replacement.
- Benchmark comparability: advanced by documenting the cache namespace masking pitfall.
- New requirement: tokenizer parity must pass before ONNX performance benchmarking or production recommendation.

## Deviations

M011 did not produce ONNX throughput evidence because semantic equivalence failed first. This is an intentional safety deviation: benchmarking a non-equivalent embedding backend would be misleading.

## Follow-ups

Plan M012 for tokenizer parity. It should produce a Hugging Face Python tokenizer baseline, compare Go tokenizer candidates, require token-level equality for fixed Russian/legal probes, then rerun TEI-vs-Go-ONNX cosine only after token parity passes.
