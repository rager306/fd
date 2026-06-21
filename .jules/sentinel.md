## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-06-21 - Length-based timing attack in token comparison
**Vulnerability:** Comparing an API token and secret key directly using `subtle.ConstantTimeCompare` without hashing them first.
**Learning:** `subtle.ConstantTimeCompare` returns early if the lengths of the two inputs differ. This can allow an attacker to infer the expected length of the token via a length-based timing attack.
**Prevention:** Always ensure the inputs to `subtle.ConstantTimeCompare` have the same length. One reliable way to guarantee this is to hash both inputs (e.g., using `crypto/sha256`) before comparing the hashes.
