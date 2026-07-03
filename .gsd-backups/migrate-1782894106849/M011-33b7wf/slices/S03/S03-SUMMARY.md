---
id: S03
parent: M011-33b7wf
milestone: M011-33b7wf
provides:
  - Go ONNX backend prototype code.
  - Backend wiring behind explicit config.
  - Documented tokenizer parity blocker for S04 decision.
requires:
  []
affects:
  - S04
key_files:
  - api/embed/onnx.go
  - api/main.go
  - benchmark-results/fd-go-onnx-m011-s03.txt
key_decisions:
  - Do not proceed to performance benchmarking until tokenizer parity is fixed.
  - Use isolated cache namespace for backend comparisons to avoid TEI cache masking ONNX behavior.
  - Record Go ONNX backend as technically loadable but semantically blocked.
patterns_established:
  - Isolate cache namespaces when comparing backend implementations.
  - Validate tokenizer parity before treating ONNX output differences as model/runtime issues.
  - A working native runtime load is not enough for embedding equivalence.
observability_surfaces:
  - ONNX startup logs artifact_id and dimensions.
  - `benchmark-results/fd-go-onnx-m011-s03.txt` captures failed isolated-cache comparison.
  - T04 summary captures tokenizer mismatch evidence and cache masking lesson.
drill_down_paths:
  - .gsd/milestones/M011-33b7wf/slices/S03/tasks/T01-SUMMARY.md
  - .gsd/milestones/M011-33b7wf/slices/S03/tasks/T02-SUMMARY.md
  - .gsd/milestones/M011-33b7wf/slices/S03/tasks/T03-SUMMARY.md
  - .gsd/milestones/M011-33b7wf/slices/S03/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T19:27:25.617Z
blocker_discovered: false
---

# S03: Opt in ONNX dense backend prototype

**S03 built and ran the opt-in Go ONNX backend but blocked on tokenizer parity; semantic equivalence failed despite successful ONNX Runtime execution.**

## What Happened

S03 implemented a real opt-in Go ONNX backend and wired it behind explicit configuration. Dependency feasibility was confirmed, non-live and live local ONNX embedder tests passed, and the local API started successfully with ONNX Runtime CPU EP. However, real API comparison with isolated cache failed semantic equivalence: TEI-vs-Go-ONNX cosine values ranged from `0.98266755` to `0.99713198`, below the `0.999` threshold. A tokenization probe confirmed the likely root cause: Go `sugarme/tokenizer` does not match Hugging Face tokenizer output for the Russian legal probe. This is a plan-relevant blocker for production/prototype success, but it is evidence-bearing and closes S03's allowed blocker path.

## Verification

S03 verification produced a blocker: ONNX backend started and ran, but isolated-cache comparison failed cosine threshold and tokenization probe showed Go/HF token ID mismatch.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

- Future ONNX adapter work requires tokenizer parity validation against Hugging Face token IDs before performance benchmarking.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S03 did not produce a passing Go ONNX API equivalence result. It produced structured blocker evidence: tokenizer mismatch between Go `sugarme/tokenizer` and Hugging Face Python tokenizer. Also discovered that shared Redis namespace can mask backend behavior by returning TEI cache hits; future ONNX comparisons must isolate cache namespace or flush carefully.

## Known Limitations

The Go ONNX backend is not semantically equivalent to TEI yet. `sugarme/tokenizer` produced different token IDs than Hugging Face for at least one Russian legal probe, resulting in cosine values below the threshold. The backend also depends on a uv-cache shared library path, not a stable project-managed ONNX Runtime library.

## Follow-ups

S04 should recommend resolving tokenizer parity before any performance benchmark. Candidate next steps: find a Go tokenizer that exactly matches HF tokenizers for XLM-R/SentencePiece, generate token IDs in a sidecar, or keep ONNX experimentation in Python until tokenizer parity is solved.

## Files Created/Modified

- `api/embed/onnx.go` — Go ONNX dense embedder implementation with manifest validation, tokenizer loading, ONNX Runtime dynamic session, and batch execution.
- `api/embed/onnx_test.go` — ONNX embedder tests including env-gated live local artifact test.
- `api/embed/onnx_manifest.go` — Manifest validation path resolution improvement for repo-root-relative artifact paths.
- `api/main.go` — Opt-in backend wiring in startup.
- `api/main_test.go` — Runtime config tests for ONNX env vars and defaults.
- `benchmark-results/fd-go-onnx-m011-s03.txt` — Failed isolated-cache Go ONNX API comparison artifact.
