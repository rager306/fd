## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-14 - Timing Attack on API Key Comparison
**Vulnerability:** `subtle.ConstantTimeCompare` returns early if the lengths of the two inputs differ, exposing the API key length to timing attacks.
**Learning:** When comparing varying length tokens (like user input vs expected API key), `subtle.ConstantTimeCompare` alone is insufficient. Both inputs must be hashed first to guarantee equal length and prevent length-based timing attacks.
**Prevention:** Always hash secrets using `crypto/sha256` before comparing them with `subtle.ConstantTimeCompare`. Compute the expected secret's hash outside the HTTP handler closure for performance.
