---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T02: Implement handler validation and tests

Introduce minimal handler dependency interfaces so tests can instantiate production handlers with mocks; add strict batch validation for dimensions and encoding_format.

## Inputs

- `S02 T01 impact notes`

## Expected Output

- `api/handlers/embeddings.go`
- `api/handlers/batch.go`
- `api/handlers/*_test.go`

## Verification

cd api && go test ./handlers -count=1

## Observability Impact

Invalid input returns clear 400 messages; production handlers become directly testable.
