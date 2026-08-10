## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2025-10-24 - subtle.ConstantTimeCompare Timing Attack
**Vulnerability:** `subtle.ConstantTimeCompare` leaked string lengths in O(1) time because it returns early if slice lengths differ, failing to securely hide the token length.
**Learning:** Even securely-named functions like `ConstantTimeCompare` have limitations; in Go, it only guarantees constant time execution for equal-length slices. Length mismatches leak in O(1) time.
**Prevention:** If pre-hashing both inputs is too expensive (e.g. in hot paths), manually mask the length mismatch. Compare the expected value against itself if lengths differ, and bitwise AND the final result against a length-validity flag.
