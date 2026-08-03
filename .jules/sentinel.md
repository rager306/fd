## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-27 - Masking length mismatches in subtle.ConstantTimeCompare
**Vulnerability:** Length-based timing attack due to `subtle.ConstantTimeCompare` returning immediately in O(1) time if slice lengths differ when validating API Keys.
**Learning:** `subtle.ConstantTimeCompare` is vulnerable to length leaks because it short circuits if input lengths don't match. When hashing both inputs before comparison is too expensive (e.g., in a high-traffic HTTP auth middleware), the length mismatch leak must be masked.
**Prevention:** If lengths differ, mask the O(1) leak by comparing the expected key against itself, forcing the comparison to run in constant time relative to the expected key's length, then fail the authentication properly.
