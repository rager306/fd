---
id: T03
parent: S03
milestone: M015-22msl0
key_files:
  - benchmark-results/fd-legal-retrieval-m015-s03.txt
  - tools/evaluate_legal_retrieval.py
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-20T05:02:32.395Z
blocker_discovered: false
---

# T03: Verified the failed legal gate artifact and cleaned up the tagged ONNX server.

**Verified the failed legal gate artifact and cleaned up the tagged ONNX server.**

## What Happened

Verified the live legal gate artifact and cleaned up the tagged ONNX runtime. The artifact includes corpus hash, isolated ONNX namespace, worst cross-backend cosine diagnostics, and verdict FAIL. Raw legal text leak checks passed. The tagged ONNX server was stopped and no background processes remain. GitNexus reports medium scope due evaluator logic changes, which are expected and verified.

## Verification

Artifact hygiene passed, ONNX server stopped, no background processes remain, and GitNexus scope was reviewed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python legal gate artifact hygiene check` | 0 | ✅ pass — legal_gate_artifact_hygiene=pass; raw_legal_text_leaks=0 | 0ms |
| 2 | `bg_shell kill eb788587` | 0 | ✅ pass — tagged ONNX process killed | 0ms |
| 3 | `curl -fsS http://localhost:18000/health after cleanup` | 0 | ✅ pass — onnx_server_stopped | 0ms |
| 4 | `bg_shell list` | 0 | ✅ pass — no background processes | 0ms |
| 5 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ⚠️ medium — expected evaluator logic changes, verified by compile/dry-run/live rerun | 0ms |

## Deviations

GitNexus risk is medium because evaluator logic changed after S02 commit; this was expected and verified via compile, dry-run, live rerun, artifact hygiene, and cleanup.

## Known Issues

Tagged ONNX quality gate failed. This should block packaging/tuning as a production-path priority until the long-text/truncation divergence is investigated.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m015-s03.txt`
- `tools/evaluate_legal_retrieval.py`
