## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-14 - Length-Based Timing Oracle in ConstantTimeCompare
**Vulnerability:** `subtle.ConstantTimeCompare` was used directly with un-hashed tokens. Because it returns immediately when lengths differ, it created a length-based timing oracle, leaking the secret length.
**Learning:** While `ConstantTimeCompare` prevents timing attacks on the *content* of equal-length slices, it exposes the length. If inputs are not hashed (e.g., to prevent DoS), a dummy comparison must be made when lengths differ to balance the timing.
**Prevention:** Always check if string lengths match before `ConstantTimeCompare`. If they do not, execute a dummy `ConstantTimeCompare(secret, secret)` to obscure the timing difference before rejecting the input.
