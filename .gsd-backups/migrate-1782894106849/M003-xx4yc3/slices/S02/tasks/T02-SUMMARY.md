---
id: T02
parent: S02
milestone: M003-xx4yc3
key_files: []
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T08:14:11.935Z
blocker_discovered: false
---

# T02: Live `/v1/embeddings` smoke tests passed for 1024d single and 512d array inputs.

**Live `/v1/embeddings` smoke tests passed for 1024d single and 512d array inputs.**

## What Happened

Called the real `/v1/embeddings` endpoint against live API, TEI, and Redis. A single Russian text with `dimensions:1024` returned one embedding with length 1024 and model `deepvk/USER-bge-m3`. An array request with two inputs and `dimensions:512` returned two embeddings with dimensions and lengths both 512.

## Verification

`/v1/embeddings` returned expected response shapes for 1024d and 512d requests.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `curl /v1/embeddings 1024d single + Python summary` | 0 | ✅ pass: data_len=1 dim=1024 emb_len=1024 | 0ms |
| 2 | `curl /v1/embeddings 512d array + Python summary` | 0 | ✅ pass: data_len=2 dims=[512,512] emb_lens=[512,512] | 0ms |

## Deviations

Used Python JSON summaries instead of jq to avoid depending on jq availability.

## Known Issues

Prompt token usage is approximate (`len(text)/4`) by existing implementation; not a smoke-test blocker.

## Files Created/Modified

None.
