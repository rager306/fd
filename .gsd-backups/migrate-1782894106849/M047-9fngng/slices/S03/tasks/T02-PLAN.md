---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Added bounded TEI retry and repeated-outage fast-fail behavior.

Add a small retry policy to `TEIClient`: bounded attempts, jitter/backoff injection or tiny defaults, classification for network errors and 502/503/504, and a lightweight circuit breaker that short-circuits repeated retriable failures for a cooldown. Preserve existing constructor compatibility.

## Inputs

- `api/embed/tei.go`
- `api/embed/tei_test.go`

## Expected Output

- `api/embed/tei.go`
- `api/embed/tei_test.go`

## Verification

cd api && go test ./embed && cd api && go test ./...

## Observability Impact

Dependency outage errors distinguish retry exhaustion from circuit-open fast fail.
