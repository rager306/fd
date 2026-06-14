---
id: T03
parent: S05
milestone: M041-4tw0w7
key_files:
  - api/middleware/ratelimit.go
  - api/middleware/ratelimit_test.go
  - api/main.go
  - benchmark-results/m041-s05-t03-go-test.txt
  - benchmark-results/m041-s05-t03-lint.txt
  - benchmark-results/m041-s05-t03-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T08:09:35.956Z
blocker_discovered: false
---

# T03: Added opt-in per-IP and per-user rate limiting with X-RateLimit headers and 429 retry envelopes.

**Added opt-in per-IP and per-user rate limiting with X-RateLimit headers and 429 retry envelopes.**

## What Happened

Implemented a concurrency-safe token bucket `RateLimiter` plus middleware wrappers. `IPRateLimitFromEnv` enforces `FD_RATE_LIMIT_IP_RPM` when `FD_RATE_LIMIT_ENABLED=true` and otherwise preserves backward compatibility. `UserRateLimitFromEnv` runs after request validation on `/v1/embeddings` so it can read the parsed `user` field and enforce a separate `FD_RATE_LIMIT_USER_RPM` bucket. Rate-limited responses set `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`, `Retry-After: 60`, and OpenAI-style `429 rate_limit_exceeded`. Wired per-IP rate limiting globally after auth and per-user rate limiting into `/v1/embeddings` after validation.

## Verification

Targeted tests passed for 101st per-IP request returning 429 with rate-limit headers and Retry-After, separate IP buckets, separate user buckets, and disabled limiter preserving compatibility. Fresh full M043 gate passed after the changes: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. GitNexus detect_changes reports LOW risk for tracked changes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./middleware -run 'Test.*RateLimit' -v` | 0 | ✅ pass: rate-limit behavior tests pass | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 4 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

The implementation uses a token bucket rather than a Redis/global distributed limiter, matching the local service scope of this milestone. Default remains disabled via `FD_RATE_LIMIT_ENABLED=false`/unset for backward compatibility.

## Known Issues

Rate limiting is in-memory per process. `api/report.json` remains an unrelated untracked generated file. `.gsd/.../S04-CONTINUE.md` remains an unrelated auto-compact artifact.

## Files Created/Modified

- `api/middleware/ratelimit.go`
- `api/middleware/ratelimit_test.go`
- `api/main.go`
- `benchmark-results/m041-s05-t03-go-test.txt`
- `benchmark-results/m041-s05-t03-lint.txt`
- `benchmark-results/m041-s05-t03-govulncheck.txt`
