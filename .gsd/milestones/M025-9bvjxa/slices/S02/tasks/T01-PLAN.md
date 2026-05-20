---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Write ONNX operations contract

Write `docs/onnx-artifacts/OPERATIONS.md` covering startup preflight, failure diagnostics, health/status checks, safe logging fields, rollout stages, rollback to TEI, and production-default safeguards.

## Inputs

- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/README.md`

## Expected Output

- `docs/onnx-artifacts/OPERATIONS.md`

## Verification

Operations doc exists and is linked from README.

## Observability Impact

Gives operators and future agents concrete diagnostic surfaces and rollback contract.
