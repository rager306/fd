## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-08-11 - Timing Attack via ConstantTimeCompare
**Vulnerability:** `subtle.ConstantTimeCompare` leaks the length of the string if the input arrays have differing lengths due to an immediate O(1) early return, leading to timing attacks that could reveal the length of API keys.
**Learning:** `ConstantTimeCompare` is only constant-time for inputs of the same length. Comparing byte slices of different lengths is vulnerable to length-revealing side channel timing attacks.
**Prevention:** If the lengths differ, reassign the user input array to the secret string array and evaluate `subtle.ConstantTimeCompare` normally (so both inputs are identical in length, running securely in constant time). Finally, bitwise AND the valid length check against the result of `subtle.ConstantTimeCompare` to guarantee correct authentication evaluation.
