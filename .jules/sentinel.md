## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-06-25 - Length-based Timing Attack in ConstantTimeCompare
**Vulnerability:** `subtle.ConstantTimeCompare` was used directly on `[]byte(token)` and `[]byte(apiKey)` without ensuring equal lengths.
**Learning:** `ConstantTimeCompare` returns immediately if the lengths of the two byte slices are different. This allows attackers to discover the length of the secret API key by sending requests with tokens of varying lengths and measuring the response time, which is a form of length-based timing attack.
**Prevention:** Always ensure both inputs have the same length before passing them to `ConstantTimeCompare`. A standard and secure way to do this is to hash both strings (e.g., using `crypto/sha256`) and compare the resulting fixed-size hashes.
