---
id: T03
parent: S01
milestone: M010-84qfzu
key_files:
  - .gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md
  - .gsd/milestones/M010-84qfzu/slices/S01/tasks/T03-SUMMARY.md
key_decisions:
  - S02 should build the dense comparator before S03 attempts ONNX export/load.
  - S03 should use ignored local artifact storage and record ONNX/external-data/tokenizer hashes and output metadata.
  - No production runtime change is allowed from S01 findings.
duration: 
verification_result: passed
completed_at: 2026-05-19T18:32:51.671Z
blocker_discovered: false
---

# T03: Saved the S01 provenance synthesis with candidate ranking, hash/provenance checklist, dense-output risks, and S02/S03 go-forward guidance.

**Saved the S01 provenance synthesis with candidate ranking, hash/provenance checklist, dense-output risks, and S02/S03 go-forward guidance.**

## What Happened

Synthesized S01 findings from local artifact inspection and external ONNX candidate research into `.gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md`. The artifact ranks model-preserving export from the exact `deepvk/USER-bge-m3` snapshot as the only candidate suitable for M010's primary path, classifies BAAI ONNX artifacts as reference-only, defers INT8, and defines required provenance fields for future ONNX attempts. It also defines the build order for S02/S03: comparator first, then export/load proof.

## Verification

Verified with filesystem checks that `.gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md` exists and contains the required no-production-runtime-change statement plus key candidate/model identifiers.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `test -f .gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md && grep -q "No production runtime change" .gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md && grep -q "deepvk/USER-bge-m3" .gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md && grep -q "aapot/bge-m3-onnx" .gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md` | 0 | ✅ pass | 0ms |

## Deviations

None. Used GSD research artifact writer instead of manual file editing.

## Known Issues

No ready FP32 ONNX artifact for exact `deepvk/USER-bge-m3` is available; S03 may still be blocked by export/load incompatibility. This is captured as an acceptable spike risk.

## Files Created/Modified

- `.gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md`
- `.gsd/milestones/M010-84qfzu/slices/S01/tasks/T03-SUMMARY.md`
