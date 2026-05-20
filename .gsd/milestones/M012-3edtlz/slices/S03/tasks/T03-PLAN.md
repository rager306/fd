---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T03: Integrate parity path or record blocker

If candidate parity passes, integrate the binding into `api/embed/onnx.go` and add token parity tests. If it does not pass or cannot be packaged, do not integrate; write blocker summary instead.

## Inputs

- `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt`

## Expected Output

- `api/embed/onnx.go or blocker task summary`

## Verification

If code changes: Go tests pass. If blocker: blocker evidence names exact failure and runtime code remains unchanged.

## Observability Impact

Keeps runtime changes gated behind proof instead of speculative dependency churn.
