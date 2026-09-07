## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.
## 2024-09-07 - Timing Attack on Variable Length Tokens
**Vulnerability:** The API key authentication middleware used `subtle.ConstantTimeCompare` directly on raw string inputs. Since `ConstantTimeCompare` returns immediately if the lengths differ, it exposes a timing side-channel that allows attackers to determine the length of the expected API key.
**Learning:** Comparing secrets of variable or unknown lengths directly with `ConstantTimeCompare` is insecure because it leaks length information.
**Prevention:** Hash both secrets (e.g., using SHA-256) before passing them to `ConstantTimeCompare` to ensure the comparison always takes constant time regardless of the original string lengths.
