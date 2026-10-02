## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Length-based Timing Oracle in API Key Validation
**Vulnerability:** `subtle.ConstantTimeCompare` was used directly with API key strings. Since `ConstantTimeCompare` returns immediately if slice lengths differ, it leaked the length of the expected API key, creating a timing oracle.
**Learning:** `ConstantTimeCompare` requires equal-length inputs to provide constant-time guarantees. Without length checks and dummy work balancing, attackers can deduce secret length by observing immediate returns versus timing bounds.
**Prevention:** Always check length before `ConstantTimeCompare`. If lengths mismatch, execute a dummy comparison against the secret itself to strictly balance execution time without leaking length or using computational hashing (which introduces DoS risk).
