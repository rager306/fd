## 2024-05-24 - Unauthenticated System Endpoints
**Vulnerability:** Critical readiness and health endpoints (`/health`, `/ready`, `/v1/healthcheck`) were blocked by authentication requirements when `FD_API_KEY` was set.
**Learning:** Load balancers and orchestration systems (like Kubernetes) often probe health check endpoints without authentication. If they are blocked by a global API key requirement, the service might be incorrectly marked as unhealthy and terminated.
**Prevention:** Ensure all liveness, readiness, and health-check endpoints are explicitly excluded from global authentication middleware.

## 2026-07-11 - G304 File Inclusion Linter Warnings
**Vulnerability:** The `gosec` linter flagged `os.ReadFile` calls using dynamically constructed file paths in test fixtures as a potential file inclusion vulnerability (G304).
**Learning:** While test fixtures using `runtime.Caller` do not take user input and are inherently safe, static analysis tools cannot always differentiate them from unsafe user-supplied paths.
**Prevention:** When dynamically loading files in production, always sanitize and strictly validate user input against an allowlist before passing it to file operations. In test code, use `filepath.Clean` combined with an explicit `//nolint:gosec` directive containing a trailing explanation to silence false positives.
