---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T02: Remove ONNX runtime selection from active API startup path

Edit active Go runtime startup/config so fd no longer accepts ONNX as a current backend selector. Remove or neutralize ONNX env parsing/config branches, update tests accordingly, and ensure invalid ONNX env usage fails closed with a clear TEI-only error or is ignored only if documented. Keep TEI behavior unchanged.

## Inputs

- `documents/onnx-deactivation-inventory-m042.md`

## Expected Output

- `api/main.go`
- `api/main_test.go`

## Verification

Targeted Go tests for runtime config pass; TEI startup config still passes; ONNX selector is absent or fails closed as TEI-only.

## Observability Impact

Runtime health should report TEI without ONNX ambiguity.
