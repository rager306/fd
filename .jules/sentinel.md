## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-07-02 - Length-Based Timing Attack in API Key Comparison
**Vulnerability:** Direct comparison of bearer tokens and API keys using `subtle.ConstantTimeCompare` without hashing first.
**Learning:** `ConstantTimeCompare` returns early if the lengths of the two byte slices differ, which allows an attacker to deduce the length of the valid API key through a timing attack.
**Prevention:** When comparing secrets or tokens, always hash both inputs (e.g., with `crypto/sha256`) before using `subtle.ConstantTimeCompare` to guarantee they are of equal length.
