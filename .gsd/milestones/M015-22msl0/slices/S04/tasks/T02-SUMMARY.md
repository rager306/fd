---
id: T02
parent: S04
milestone: M015-22msl0
key_files:
  - .gsd/DECISIONS.md
  - benchmark-results/fd-legal-retrieval-m015-summary.txt
key_decisions:
  - D012 blocks ONNX packaging/tuning until long-text legal quality divergence is investigated.
duration: 
verification_result: passed
completed_at: 2026-05-20T05:06:05.246Z
blocker_discovered: false
---

# T02: Recorded the blocking legal quality decision as D012.

**Recorded the blocking legal quality decision as D012.**

## What Happened

Recorded D012: tagged ONNX failed the Russian/legal quality gate and should not proceed to packaging/tuning as the next priority. TEI remains production/default. The next recommended milestone is investigation of long-text/truncation divergence on legal corpus inputs.

## Verification

GSD decision save succeeded and regenerated `.gsd/DECISIONS.md`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_decision_save D012 quality-gate outcome` | 0 | ✅ pass — decision saved | 0ms |

## Deviations

None.

## Known Issues

Packaging/tuning remains technically possible but is deprioritized because the quality gate failed.

## Files Created/Modified

- `.gsd/DECISIONS.md`
- `benchmark-results/fd-legal-retrieval-m015-summary.txt`
