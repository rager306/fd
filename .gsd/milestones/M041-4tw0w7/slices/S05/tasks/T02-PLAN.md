---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: API key auth (FD_API_KEY env) и CORS

api/middleware/auth.go: если env FD_API_KEY задан, требует Authorization: Bearer <key> на всех endpoints кроме /live, /metrics, /docs, /openapi.json. На missing/wrong → 401 unauthorized. api/middleware/cors.go: Access-Control-Allow-Origin (из env FD_CORS_ORIGINS или * default), Access-Control-Allow-Methods: GET,POST,OPTIONS, Access-Control-Allow-Headers: Content-Type,Authorization,X-Request-Id. OPTIONS preflight → 204.

## Inputs

- None specified.

## Expected Output

- `api/middleware/auth.go`
- `api/middleware/cors.go`
- `api/middleware/auth_test.go`
- `api/middleware/cors_test.go`

## Verification

Unit tests: T-E-9 (с FD_API_KEY=test, без Authorization → 401 unauthorized, с правильным Bearer → 200). OPTIONS preflight → 204 с правильными CORS headers.
