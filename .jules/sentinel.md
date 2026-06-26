## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2026-06-26 - Prevent Length-Based Timing Attacks in ConstantTimeCompare
**Vulnerability:** Length-based timing attack in token comparison using `subtle.ConstantTimeCompare` with inputs of potentially different lengths.
**Learning:** `subtle.ConstantTimeCompare` returns immediately if the lengths of the two byte slices are different. This allows an attacker to deduce the length of the secret key by measuring response times.
**Prevention:** When comparing secrets or tokens of unknown/variable length using `subtle.ConstantTimeCompare`, ensure both inputs are first hashed (e.g., using `crypto/sha256`) to guarantee equal length. This prevents length-based timing attacks where `ConstantTimeCompare` returns early if lengths differ.
