---
id: T03
parent: S02
milestone: M003-xx4yc3
key_files: []
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T08:14:41.704Z
blocker_discovered: false
---

# T03: Live `/embeddings/batch` smoke tests passed for base64 and float modes.

**Live `/embeddings/batch` smoke tests passed for base64 and float modes.**

## What Happened

Called the real `/embeddings/batch` endpoint in base64 and float modes. Base64 request with two 1024d inputs returned count 2, dimensions 1024, two embeddings, and the first embedding decoded as valid base64. Float request with one 512d input returned count 1, dimensions 512, and the string field decoded as a JSON array of 512 floats.

## Verification

Batch endpoint returned expected count/dimensions/payloads for base64 and float requests.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `curl /embeddings/batch base64 1024d + Python base64 decode` | 0 | ✅ pass: count=2 dimensions=1024 payload valid base64 | 0ms |
| 2 | `curl /embeddings/batch float 512d + Python JSON decode` | 0 | ✅ pass: count=1 dimensions=512 first_float_len=512 | 0ms |

## Deviations

Used Python to validate base64 decode and float payload length.

## Known Issues

`encoding_format:"float"` returns a JSON array encoded as a string field; this is existing API design and should be considered in future API improvements.

## Files Created/Modified

None.
