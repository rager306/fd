## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-07-21 - Length-Based Timing Attack in subtle.ConstantTimeCompare
**Vulnerability:** API key authentication was vulnerable to a length-based timing attack because `subtle.ConstantTimeCompare` returns early if the lengths of the two inputs differ, leaking the length of the expected secret.
**Learning:** When comparing secrets, `subtle.ConstantTimeCompare` must be used with inputs of equal length. Otherwise, the early exit allows an attacker to deduce the secret's length.
**Prevention:** Always hash both the incoming token and the expected secret using a cryptographic hash function (e.g., `crypto/sha256`) before passing them to `subtle.ConstantTimeCompare`. This guarantees both inputs are exactly the same length.
