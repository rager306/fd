---
id: T03
parent: S02
milestone: M014-vjfs9f
key_files:
  - benchmark-results/fd-benchmark-m014-tei-baseline.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-20T04:20:10.496Z
blocker_discovered: false
---

# T03: Verified the TEI baseline artifact is complete, parseable, and free of raw fixed-probe text leaks.

**Verified the TEI baseline artifact is complete, parseable, and free of raw fixed-probe text leaks.**

## What Happened

Verified the TEI baseline artifact. Required snapshot fields and runtime label are present, all benchmark sections are present under the actual headings, summary exists, and raw fixed probe texts from the comparator did not leak into the artifact. The first verification parser used stale heading names, which was corrected after inspecting the artifact headings.

## Verification

Artifact parser/leak checks passed after adjusting expected headings to the actual benchmark output; GitNexus detect_changes reported artifact-only low scope.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `rg '^## ' benchmark-results/fd-benchmark-m014-tei-baseline.txt` | 0 | ✅ pass — actual section headings identified | 0ms |
| 2 | `uv run --python 3.13 --with transformers --with torch --with sentencepiece python artifact hygiene check` | 0 | ✅ pass — tei_artifact_hygiene=pass; raw_probe_text_leaks=0 | 0ms |
| 3 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ✅ pass — low scope, non-code artifact changes only | 0ms |

## Deviations

Initial parser expected older/generic section headings and failed on the actual benchmark headings. The parser was corrected to match the real artifact headings and then passed.

## Known Issues

GitNexus detect shows only non-code artifact changes for S02, as expected.

## Files Created/Modified

- `benchmark-results/fd-benchmark-m014-tei-baseline.txt`
