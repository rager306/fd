---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T06: /warmup endpoints (GET status, POST trigger)

api/handlers/warmup.go: GET /warmup — { status: ready|warming_up, progress: 0..1 }. POST /warmup — если ready, 200 { status: ready, message: already warm }; если нет, 202 { status: warming_up, message: warmup started } и trigger background warmup (sync.WaitGroup не блокирует).

## Inputs

- None specified.

## Expected Output

- `api/handlers/warmup.go`
- `api/handlers/warmup_test.go`

## Verification

Unit tests: GET /warmup → 200 status:ready после warmup, GET → 200 status:warming_up, progress:<fraction> во время. POST /warmup → 200 если ready, 202 если warming.
