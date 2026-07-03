---
estimated_steps: 1
estimated_files: 5
skills_used: []
---

# T04: Updated operator docs and compose/CI surfaces to TEI-only current posture, with ONNX retained only as historical/future research context.

Update README, same-host contract, fd-v2 docs, and relevant operations docs to state TEI is the only current runtime. Mark ONNX as historical/future research, not an operator option. Remove outdated compose comments suggesting ONNX export as current optimization. Update R021/R027/R022 statuses if evidence supports it.

## Inputs

- `documents/te-perf-root-cause-m042.md`
- `documents/onnx-deactivation-inventory-m042.md`

## Expected Output

- `README.md`
- `docs/same-host-embedding-service-contract.md`
- `docs/fd-v2.md`

## Verification

Docs contain TEI-only current posture and no active ONNX operator instructions outside historical/future research notes.

## Observability Impact

Operator docs match runtime reality.
