## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-07-26 - Timing Attack on API Key Validation
**Vulnerability:** The API key validation used `subtle.ConstantTimeCompare` directly on strings of potentially differing lengths, allowing an attacker to guess the API key length via a timing attack.
**Learning:** `ConstantTimeCompare` leaks length information if the slices are of different lengths.
**Prevention:** When comparing secrets of potentially different lengths, hash both inputs using a cryptographic hash (like `sha256`) before passing to `ConstantTimeCompare` to ensure both inputs have identical lengths.
