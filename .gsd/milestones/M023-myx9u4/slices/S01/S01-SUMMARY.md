---
id: S01
parent: M023-myx9u4
milestone: M023-myx9u4
provides:
  - Packaged legal quality evidence for the opt-in ONNX Docker image.
requires:
  []
affects:
  - S02 outcome and closure verification
key_files:
  - benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt
  - .github/workflows/go-quality.yml
key_decisions:
  - Use packaged ONNX runtime label `packaged-onnx1024-docker`.
  - Use isolated ONNX cache namespace `m023-onnx-docker-legal`.
  - Refine binary hygiene checks to allow `Dockerfile.onnx`.
patterns_established:
  - Use isolated cache namespace for every TEI-vs-ONNX gate.
  - Treat `Dockerfile.onnx` as packaging source, not a binary artifact.
  - Use raw text leak checks before closing legal artifacts.
observability_surfaces:
  - Sanitized legal artifact with runtime/cache labels, raw-text leak check, health/smoke evidence, corrected CI binary hygiene check.
drill_down_paths:
  - .gsd/milestones/M023-myx9u4/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M023-myx9u4/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M023-myx9u4/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T10:58:52.537Z
blocker_discovered: false
---

# S01: Packaged ONNX legal quality gate

**S01 proved packaged ONNX Docker 1024 passes the Russian/legal retrieval quality gate.**

## What Happened

S01 ran the packaged ONNX legal quality gate. The TEI baseline and packaged ONNX image were healthy, the evaluator compared port 8000 against port 18000 with isolated cache namespace, and the gate passed with minimum cosine 0.99989883, top-1 agreement 1.0, mean overlap@5 0.997701, and ONNX recall ratio 1.0. The artifact excludes raw legal text. Runtime cleanup completed and the workflow binary hygiene false positive for `Dockerfile.onnx` was fixed.

## Verification

Evaluator PASS, raw text leak check, endpoint health, cleanup, and corrected CI hygiene verification passed.

## Requirements Advanced

- onnx-packaged-legal-quality — Validated packaged Docker runtime against the legal quality gate.

## Requirements Validated

- m023-packaged-legal-pass — `benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt` verdict PASS with minimum cosine 0.99989883.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Binary hygiene verification found a false positive in the CI regex: tracked `Dockerfile.onnx` matched `.onnx$`. The workflow check was corrected to exempt `Dockerfile.onnx` while still blocking actual ONNX/native/runtime binaries.

## Known Limitations

The evaluator uses an unlabeled corpus, so it proves TEI-vs-ONNX parity and synthetic known-item behavior, not absolute human relevance quality.

## Follow-ups

S02 should record the packaged legal pass outcome and rerun closure guardrails, including the corrected CI hygiene check.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt` — Packaged ONNX legal quality artifact.
- `.github/workflows/go-quality.yml` — Fixed CI binary hygiene check to distinguish Dockerfile.onnx from model artifacts.
