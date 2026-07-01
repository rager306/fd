---
id: S01
parent: M031-gn517a
milestone: M031-gn517a
provides:
  - Truthful artifact source matrix for S02 documentation and manifest updates.
requires:
  []
affects:
  - S02
key_files:
  - .gsd/milestones/M031-gn517a/slices/S01/S01-RESEARCH.md
key_decisions:
  - Native tokenizer candidate: pin `https://github.com/daulet/tokenizers/releases/download/v1.27.0/libtokenizers.linux-amd64.tar.gz`; archive/extracted library checksums verified.
  - Tokenizer JSON candidate: pin Hugging Face revision URL for `tokenizer.json`; checksum verified.
  - ONNX Runtime candidate: pin PyPI `onnxruntime==1.26.0` CP313 manylinux wheel; extracted shared library checksum verified.
  - ONNX model artifact remains blocked until mirrored/uploaded as an exact immutable binary or separately reproduced and revalidated.
patterns_established:
  - Use source statuses (`immutable_selected`, `candidate`, `blocked`) to avoid overclaiming rollout readiness.
observability_surfaces:
  - Source contract research artifact with source status definitions and checksum evidence references.
drill_down_paths:
  - .gsd/milestones/M031-gn517a/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M031-gn517a/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M031-gn517a/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T06:34:23.617Z
blocker_discovered: false
---

# S01: Artifact source contract research

**Defined the M031 artifact source contract research: three pinned candidates, one ONNX model blocker.**

## What Happened

S01 inventoried all required ONNX artifacts and researched immutable source candidates using tracked provenance and public metadata/checksum probes. The work selected pinned candidates for native tokenizer, tokenizer JSON, and ONNX Runtime, while preserving the ONNX model binary as a hard blocker because it exists only as a local ignored export. No external state was changed and ONNX remains experimental.

## Verification

S01 verification passed via source inventory, checksum probes, and research safety/status checks. Evidence: `.gsd/exec/78224446-6bc5-4bcb-84f1-633ade016a0c.stdout`, `.gsd/exec/24907b45-1ff7-436e-a6fa-b99bc5be5f06.stdout`, `.gsd/exec/5eb4feeb-2766-42dc-85cf-ae68c1c3a782.stdout`.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

Candidate sources have not been exercised in hosted CI. The ONNX model artifact has no immutable external source, so hosted proof remains blocked.

## Follow-ups

S02 should persist the selected candidates/blocker into docs/manifests/outcome and record a decision. It must not claim hosted workflow proof or ONNX production readiness.

## Files Created/Modified

- `.gsd/milestones/M031-gn517a/slices/S01/S01-RESEARCH.md` — Source contract research artifact for ONNX artifact source candidates and blockers.
