## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-05-24 - Timing Attack in API Key Validation
**Vulnerability:** The API key validation used `subtle.ConstantTimeCompare([]byte(token), []byte(apiKey))`. Go's `subtle.ConstantTimeCompare` returns immediately if the lengths of the slices differ, leaking the expected length of the API key via timing side-channels.
**Learning:** `ConstantTimeCompare` only protects against timing attacks for inputs of the same length. Comparing a variable-length user input against a secret exposes the secret's length.
**Prevention:** Hash both the expected secret and the provided input to a fixed-length hash (e.g., SHA-256) before using `ConstantTimeCompare`, or ensure the lengths are checked beforehand (though hashing both is generally safer when dealing with string tokens of variable length).
