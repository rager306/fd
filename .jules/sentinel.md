## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2024-05-27 - Resolving Revive Linting Issue
**Vulnerability:** Not a vulnerability, but fixing a revive linter issue in `api/embed/tei.go` where an exported function `ObserveBatchFill` had a duplicate (differing by capitalization) function `observeBatchFill`.
**Learning:** `revive` strictly enforces `confusing-naming`. Instead of maintaining both `ObserveBatchFill` and `observeBatchFill`, consolidating to just the exported `ObserveBatchFill` resolves the confusing naming issue and cleans up unused functions.
**Prevention:** Avoid unexported and exported pairs that differ only by capitalization within the same scope.

## 2024-05-27 - Resolving Revive Linting Issue
**Vulnerability:** Not a vulnerability, but fixing a revive linter issue in `api/embed/tei.go` where an exported function `ObserveBatchFill` had a duplicate (differing by capitalization) function `observeBatchFill`.
**Learning:** `revive` strictly enforces `confusing-naming`. Instead of maintaining both `ObserveBatchFill` and `observeBatchFill`, consolidating to just the exported `ObserveBatchFill` resolves the confusing naming issue and cleans up unused functions.
**Prevention:** Avoid unexported and exported pairs that differ only by capitalization within the same scope.

## 2024-05-27 - Authentication Timing Attack
**Vulnerability:** The API key validation in `api/middleware/auth.go` uses `subtle.ConstantTimeCompare` directly on the provided token and the configured API key strings.
**Learning:** `subtle.ConstantTimeCompare` only protects against timing attacks if both byte slices are the same length. If lengths differ, it returns early.
**Prevention:** To guarantee equal length and prevent length-based timing attacks, both inputs should be hashed (e.g., using `crypto/sha256`) before comparison.
