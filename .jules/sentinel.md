## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-03-24 - [Auth Comparison Timing Attack]
**Vulnerability:** Length-based timing attack during API key validation in auth middleware.
**Learning:** `subtle.ConstantTimeCompare` terminates early if the two byte slices have different lengths, effectively allowing an attacker to deduce the correct API key length, and potentially bruteforce the key character by character by measuring timing variations.
**Prevention:** Always hash both strings (using e.g., sha256) to fixed-length representations before passing them to `ConstantTimeCompare`, ensuring equal lengths.
