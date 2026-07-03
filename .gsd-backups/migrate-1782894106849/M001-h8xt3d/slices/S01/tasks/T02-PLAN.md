---
estimated_steps: 1
estimated_files: 5
skills_used: []
---

# T02: Implement cache correctness fixes

Update TieredCache to use dimension-aware local/singleflight keys and update Redis marshal/set path to validate vector lengths and return errors. Add focused unit tests for dimension isolation and short-vector error behavior.

## Inputs

- `S01 T01 impact analysis`

## Expected Output

- `api/cache/tiered.go`
- `api/cache/redis.go`
- `api/cache/*_test.go`

## Verification

cd api && go test ./cache -run 'Test.*(Tiered|Binary|Redis|Local)' -count=1

## Observability Impact

Replaces panic/silent cache write failure with explicit errors.
