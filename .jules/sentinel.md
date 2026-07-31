## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-02-14 - Constant-Time String Comparison Length Masking
**Vulnerability:** Length-based timing attack in API Key Authentication `subtle.ConstantTimeCompare`
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the slices differ, leaking length info in O(1) time. Pre-hashing to normalize length is a standard solution but can be too slow for high-throughput HTTP middleware on every request.
**Prevention:** To mitigate this without pre-hashing, always compare inputs relative to the expected secret length by masking mismatches. If the lengths differ, compare the expected secret with itself (`expected == expected`), and logically AND the result with a mismatch flag `match = 0`.
