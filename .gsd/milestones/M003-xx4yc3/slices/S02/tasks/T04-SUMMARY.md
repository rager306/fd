---
id: T04
parent: S02
milestone: M003-xx4yc3
key_files: []
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T08:15:05.900Z
blocker_discovered: false
---

# T04: Live negative API tests passed with expected 400 responses.

**Live negative API tests passed with expected 400 responses.**

## What Happened

Ran live negative API tests. Invalid JSON, empty `/v1/embeddings` input, invalid dimensions, and invalid batch encoding_format all returned HTTP 400 with clear JSON error bodies.

## Verification

All four negative cases returned HTTP 400 with expected error messages.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `curl status checks for invalid_json, empty_input, bad_dim, bad_format` | 0 | ✅ pass: all returned 400 | 0ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

None.
