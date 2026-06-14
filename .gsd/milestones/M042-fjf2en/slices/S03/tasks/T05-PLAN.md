---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T05: ONNX mode docs + legal quality gate reference

documents/onnx-mode-m042.md: (1) что такое ONNX mode (env flags, build command), (2) perf numbers (из S03 T03 benchmark), (3) legal quality gate deferred — explicit reference to M015/M016 findings (128-token truncation causes divergence), (4) production rollout checklist (legal quality gate close-out required as separate milestone, contact for opting in, monitoring recommendations). Обновить docs/fd-v2.md Section 5.4 с consolidated "after M042" perf table (TEI sync/async, ONNX sync/async) и known limitations.

## Inputs

- None specified.

## Expected Output

- `documents/onnx-mode-m042.md`
- `docs/fd-v2.md (Section 5.4 update)`

## Verification

Документ существует, cross-references M015/M016, M019, M041. Production rollout checklist explicit. docs/fd-v2.md Section 5.4 updated.
