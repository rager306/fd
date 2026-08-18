## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2025-02-14 - Timing Attack in API Key Validation
**Vulnerability:** The API key validation used `subtle.ConstantTimeCompare` directly on the incoming token and configured API key. If the token's length didn't match the expected key's length, `subtle.ConstantTimeCompare` would return early (not constant time), leaking the correct key's length through a timing side channel.
**Learning:** `subtle.ConstantTimeCompare` is only constant time if both inputs are the same length. It relies on the caller to ensure equal lengths.
**Prevention:** Hash both the secret and the user input with a strong hash function like SHA-256 before comparing them. This ensures both inputs are always exactly the same length.
