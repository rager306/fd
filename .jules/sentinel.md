## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-05-18 - Timing Attack via ConstantTimeCompare with Differing Lengths
**Vulnerability:** The API key validation used `subtle.ConstantTimeCompare` directly on the provided token and the configured API key. If the lengths differ, `ConstantTimeCompare` returns immediately, allowing an attacker to deduce the correct API key length via timing attacks.
**Learning:** `subtle.ConstantTimeCompare` requires equal-length inputs to provide constant-time guarantees. Directly comparing variable-length user input against a secret key bypasses this protection.
**Prevention:** When comparing secrets of potentially different lengths, always hash both inputs (e.g., using `crypto/sha256`) first. The hashes will always be of equal length, ensuring a true constant-time comparison.
