## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-11-20 - Constant Time Compare Timing Attack
**Vulnerability:** `subtle.ConstantTimeCompare` leaks length information when used to compare variable-length secrets (like API keys/tokens), acting as a timing side-channel.
**Learning:** `ConstantTimeCompare` immediately returns 0 if lengths do not match. To prevent leaking the expected secret's length, both the incoming token and the expected secret must be hashed to a fixed length (e.g., SHA-256) before comparison.
**Prevention:** Always hash secrets before comparing them with `ConstantTimeCompare` to ensure length mismatches don't expose timing side-channels.
