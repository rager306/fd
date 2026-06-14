---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Rate limiting (per-IP, per-user)

api/middleware/ratelimit.go: token bucket per IP (100 req/min default) и per user (1000 req/min default если user field задан). Env FD_RATE_LIMIT_IP_RPM, FD_RATE_LIMIT_USER_RPM для конфигурации. Headers X-RateLimit-Limit/Remaining/Reset на каждом response. На превышение → 429 rate_limit_exceeded + Retry-After: 60. Опционально через FD_RATE_LIMIT_ENABLED=true (default false для обратной совместимости).

## Inputs

- None specified.

## Expected Output

- `api/middleware/ratelimit.go`
- `api/middleware/ratelimit_test.go`

## Verification

Unit tests: с включённым rate limit, 101-й запрос за минуту → 429 с X-RateLimit-* headers и Retry-After: 60. Per-user limit отдельно от per-IP.
