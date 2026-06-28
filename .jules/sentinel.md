## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2025-02-18 - Length-Based Timing Attack in API Key Comparison
**Vulnerability:** API keys and tokens were being compared directly using `subtle.ConstantTimeCompare`, which can leak length information via early returns if the lengths of the two inputs differ.
**Learning:** `ConstantTimeCompare` only protects against timing attacks for the contents of the strings, not their lengths. If the lengths differ, it returns immediately, allowing attackers to infer the expected length.
**Prevention:** Always hash secrets/tokens (e.g., using `crypto/sha256`) before comparing them with `subtle.ConstantTimeCompare` to ensure both inputs are of equal length.
