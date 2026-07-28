## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-24 - Length-Based Timing Attack in API Key Middleware
**Vulnerability:** The API key verification middleware `subtle.ConstantTimeCompare([]byte(token), []byte(apiKey))` returned immediately if the strings were of different lengths, leaking the required token length in O(1) time and enabling side-channel timing attacks.
**Learning:** `subtle.ConstantTimeCompare` requires inputs of identical length to run in constant time. In hot paths where performing cryptographic pre-hashing of inputs is too expensive, comparing raw strings directly is dangerous.
**Prevention:** Mask length mismatches. If lengths differ, substitute the incoming byte array with the expected secret's byte array, setting a valid flag to false. Perform `ConstantTimeCompare` and check both the result and the valid flag.
