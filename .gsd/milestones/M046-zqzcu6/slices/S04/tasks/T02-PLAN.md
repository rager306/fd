---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T02: Made protected endpoints fail closed when `FD_API_KEY` is missing and protected `/metrics`.

Change auth middleware so missing API key rejects protected endpoints instead of disabling auth. Keep probe endpoints public. Remove `/metrics` from public auth carve-outs. Preserve OpenAI-style error envelopes and avoid logging secrets.

## Inputs

- `api/middleware/auth.go`
- `README.md`

## Expected Output

- `api/middleware/auth.go`
- `api/middleware/auth_test.go`
- `README.md`

## Verification

cd api && go test ./middleware -run TestAPIKeyAuth

## Observability Impact

Error responses should make auth misconfiguration diagnosable without exposing secret values.
