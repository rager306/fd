## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-10-15 - Length-Based Timing Attacks in ConstantTimeCompare
**Vulnerability:** Length-based timing attacks due to direct use of `subtle.ConstantTimeCompare` with raw string tokens of varying lengths.
**Learning:** `subtle.ConstantTimeCompare` returns immediately if the lengths of the two byte slices differ, which can expose the length of the secret to an attacker.
**Prevention:** Always ensure both inputs are of equal length before comparing them. Hash both inputs (e.g., using `crypto/sha256`) and compare the resulting fixed-length hashes, preferably computing the expected secret's hash outside the request handler loop.
