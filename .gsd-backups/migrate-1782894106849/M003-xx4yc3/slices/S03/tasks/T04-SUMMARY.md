---
id: T04
parent: S03
milestone: M003-xx4yc3
key_files: []
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T08:18:51.003Z
blocker_discovered: false
---

# T04: Redis L2 cache survived API restart and served the cached key.

**Redis L2 cache survived API restart and served the cached key.**

## What Happened

Validated Redis L2 cache after API restart. After flushing Redis and priming a 1024d key, the API container was restarted. The key remained in Redis, API health recovered, and the same request after restart completed in 0.01s without increasing the TEI cache-miss log count. This proves Redis L2 serves cached embeddings after L1 is lost on process restart.

## Verification

After API restart, cached request returned 1024d embedding in 0.01s, Redis key remained, and cache-miss log count did not increase.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `redis FLUSHALL; prime key; docker compose restart api; curl same key; compare Redis keys and miss logs` | 0 | ✅ pass: keys_before=1 keys_after=1 miss count unchanged after restart | 0ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

None.
