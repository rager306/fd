## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-01 - Prevent Length-Based Timing Attacks in ConstantTimeCompare
**Vulnerability:** Length-based timing attack in API key validation when comparing raw tokens.
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the two inputs differ, exposing the length of the expected secret.
**Prevention:** Hash both inputs using a cryptographic hash function like SHA-256 before comparing them with `subtle.ConstantTimeCompare` to ensure equal lengths.
