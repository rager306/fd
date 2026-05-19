---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Add cache-path observability

Implement configurable log level and cache debug/warn events without changing cache semantics.

## Inputs

- `T01 findings`

## Expected Output

- `api/main.go`
- `api/cache/tiered.go`

## Verification

Go cache tests pass.

## Observability Impact

Cache path becomes inspectable at debug level.
