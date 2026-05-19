---
id: S02
parent: M003-xx4yc3
milestone: M003-xx4yc3
provides:
  - Confirmed live endpoints for S03 cache validation and S04 benchmark.
requires:
  []
affects:
  - S03
  - S04
key_files: []
key_decisions:
  - Use Python JSON summaries instead of jq to keep smoke tests independent of jq installation.
patterns_established:
  - Smoke tests summarize response shape and payload sizes rather than printing full embeddings.
observability_surfaces:
  - Captured response shape summaries and error bodies for live HTTP calls.
drill_down_paths:
  - .gsd/milestones/M003-xx4yc3/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M003-xx4yc3/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M003-xx4yc3/slices/S02/tasks/T03-SUMMARY.md
  - .gsd/milestones/M003-xx4yc3/slices/S02/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T08:15:30.697Z
blocker_discovered: false
---

# S02: Live API smoke tests

**S02 proved live health, embeddings, batch, and negative API paths work.**

## What Happened

S02 validated live endpoints against the running stack. API, TEI, and Redis health checks passed. `/v1/embeddings` returned correct 1024d and 512d response shapes. `/embeddings/batch` returned valid base64 and float payloads. Negative validation cases returned HTTP 400 with clear error bodies.

## Verification

All live API smoke tests passed.

## Requirements Advanced

- Live API runtime validation complete. — 

## Requirements Validated

- API health ok, TEI health HTTP success, Redis PONG. — 
- /v1/embeddings 1024d single returned emb_len 1024. — 
- /v1/embeddings 512d array returned two embeddings length 512. — 
- /embeddings/batch base64 and float modes returned expected metadata and valid payloads. — 
- Negative validation cases returned 400. — 

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

TEI health endpoint returns HTTP success with empty body; verification used curl exit success.

## Known Limitations

API `/health` is shallow and does not check Redis/TEI dependencies; this is mitigated by separate explicit checks in the milestone.

## Follow-ups

Continue with S03 runtime cache validation.

## Files Created/Modified

None.
