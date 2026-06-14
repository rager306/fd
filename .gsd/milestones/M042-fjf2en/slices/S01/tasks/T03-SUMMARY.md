---
id: T03
parent: S01
milestone: M042-fjf2en
key_files:
  - documents/te-perf-root-cause-m042.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T10:17:00.068Z
blocker_discovered: false
---

# T03: Wrote the M042 TEI RCA, concluding that TEI queue/startup behavior is the root target and ONNX should be deferred from the current milestone.

**Wrote the M042 TEI RCA, concluding that TEI queue/startup behavior is the root target and ONNX should be deferred from the current milestone.**

## What Happened

Created `documents/te-perf-root-cause-m042.md` using T01 direct TEI timings, T02 restart/startup evidence, M019 ONNX performance outcome, and M040 runtime recommendation. The RCA contains a hypothesis tree covering TEI batch scheduler queueing, queue metric semantics, missing-ONNX startup delay, async chunking uncertainty, and ONNX promotion. It concludes that fd HTTP/cache are not the root cause; TEI internal scheduling and startup/fallback behavior are. It recommends keeping M042 TEI-first, skipping/defering ONNX implementation, and treating future ONNX only as a separate research milestone with packaging/readiness gates.

## Verification

Verified `documents/te-perf-root-cause-m042.md` is 10036 bytes, contains the snapshot section, hypothesis tree with H1/H2/H3, root-cause verdict, recommended action, M019/M040 cross-references, and explicit ONNX deferral. R020 was updated to validated from this artifact.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 - <<'PY'
from pathlib import Path
p=Path('documents/te-perf-root-cause-m042.md')
text=p.read_text()
checks={
 'bytes>=2k': p.stat().st_size>=2048,
 'snapshot': 'Direct TEI sequential timing snapshot' in text,
 'hypothesis_tree': '## Hypothesis tree' in text,
 'h1': '### H1' in text and 'Prediction' in text,
 'h2': '### H2' in text,
 'h3': '### H3' in text,
 'verdict': '## Root cause verdict' in text,
 'recommended_action': '## Recommended action' in text,
 'm019_m040': 'M019' in text and 'M040' in text,
 'onnx_deferred': 'ONNX implementation is skipped/deferred' in text,
}
print('bytes', p.stat().st_size)
for k,v in checks.items(): print(k, v)
PY` | 0 | ✅ pass: RCA has required structure, evidence, verdict, and ONNX deferral | 120000ms |

## Deviations

The RCA recommendation diverges from the original M042 roadmap by rejecting ONNX as a current deliverable. This follows the user's explicit rescope and the evidence from T02 startup behavior.

## Known Issues

S02 still needs to decide whether TEI request shaping/async chunking is worthwhile after ONNX cleanup. Current TEI startup remains slow; avoid unnecessary restarts.

## Files Created/Modified

- `documents/te-perf-root-cause-m042.md`
