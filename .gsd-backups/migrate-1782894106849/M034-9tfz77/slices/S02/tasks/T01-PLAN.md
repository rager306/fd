---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Document hosted workflow inputs

Update provisioning docs and artifact README with safe manual workflow dispatch input contract: required ONNX/native sources, optional tokenizer/runtime sources, optional runtime sha override, no signed/plain secret URLs, and remaining exact model blocker.

## Inputs

- `.github/workflows/onnx-packaging.yml`
- `docs/onnx-artifacts/PROVISIONING.md`

## Expected Output

- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/README.md`

## Verification

Docs contain required/optional input policy and no overclaiming.

## Observability Impact

Makes workflow input policy discoverable.
