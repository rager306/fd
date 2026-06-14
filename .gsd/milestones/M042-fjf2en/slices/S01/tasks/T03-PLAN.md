---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Wrote the M042 TEI RCA, concluding that TEI queue/startup behavior is the root target and ONNX should be deferred from the current milestone.

Write `documents/te-perf-root-cause-m042.md` with introduction, evidence from T01/T02, hypothesis tree (>=3 hypotheses with testable predictions), verdict, and recommended action. The recommendation must explicitly defer ONNX runtime implementation from M042 and focus on TEI stabilization/mitigation. Cross-reference M019/M040 as historical ONNX research only, not current implementation scope.

## Inputs

- `documents/te-perf-snapshot-m042-s01.md`
- `benchmark-results/te-concurrency-profile-m042-s01.md`
- `benchmark-results/fd-onnx-1024-performance-outcome-m019-s02.txt`
- `benchmark-results/fd-runtime-recommendation-m040-s04.md`

## Expected Output

- `documents/te-perf-root-cause-m042.md`

## Verification

Document is >=2KB and contains snapshot, hypothesis tree with >=3 hypotheses and testable predictions, verdict, recommended TEI-first action, and M019/M040 cross-references. R020 can be validated from this artifact.
