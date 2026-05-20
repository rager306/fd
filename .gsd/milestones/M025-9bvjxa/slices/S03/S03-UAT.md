# S03: Full hosted ONNX CI skeleton — UAT

**Milestone:** M025-9bvjxa
**Written:** 2026-05-20T11:48:46.281Z

# S03 UAT — Full hosted ONNX CI skeleton

## Checks

- [x] Manual workflow exists at `.github/workflows/onnx-packaging.yml`.
- [x] Workflow uses `workflow_dispatch`, not push/PR.
- [x] Explicit artifact source inputs are required.
- [x] Workflow invokes provisioning helper, strict verifier, tagged tests, and ONNX Docker build.
- [x] Docs reference workflow and URL safety guidance.
- [x] actionlint passes.
- [x] Default guardrails pass.
- [x] Binary hygiene and cleanup pass.

## Result

Pass. Hosted ONNX CI has a safe manual skeleton, but full execution remains blocked until immutable artifact sources exist.

