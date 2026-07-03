---
estimated_steps: 1
estimated_files: 6
skills_used: []
---

# T01: Confirmed issue #7 cache cleanup debt exists before S01 fixes.

Add or run static checks proving issue #7 #19/#27/#28 are present: LRUCache exists with only test/scaffold refs, duplicate short hash helpers exist, and env int parsers are duplicated. Record red evidence.

## Inputs

- `documents/issue-7-current-m048.md`

## Expected Output

- `documents/issue-7-current-m048.md`
- `api/cache/lru.go`
- `api/cache/tiered.go`
- `api/cache/redis.go`
- `api/main.go`
- `api/middleware/ratelimit.go`

## Verification

Static gsd_exec check should PASS for pre-fix presence of #19/#27/#28.

## Observability Impact

Documents exact pre-fix symbols so cleanup is measured.
