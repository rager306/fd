## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-06-24 - Length-Based Timing Attack in API Key Authentication
**Vulnerability:** The API key authentication middleware used `subtle.ConstantTimeCompare` directly on the provided token and the expected API key.
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the two byte slices differ. This allows an attacker to determine the length of the secret key by sending tokens of various lengths and measuring the response time, which is a form of a length-based timing attack.
**Prevention:** Always hash both inputs (e.g., using `crypto/sha256`) before passing them to `subtle.ConstantTimeCompare` when comparing secrets or tokens of unknown or variable lengths. This guarantees equal length inputs and prevents the early return vulnerability.
