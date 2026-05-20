# S02: Opt in backend seam — UAT

**Milestone:** M011-33b7wf
**Written:** 2026-05-19T19:06:14.487Z

# S02 UAT — Opt in backend seam

## Checks

- [x] GitNexus impact was run before edits and reported LOW risk.
- [x] TEI remains default when `EMBEDDING_BACKEND` is unset.
- [x] Invalid backend names fail validation.
- [x] `EMBEDDING_BACKEND=onnx` requires `ONNX_ARTIFACT_MANIFEST`.
- [x] Manifest validation checks local file, size, SHA256, output name, dimensions, and prototype/default flags.
- [x] Explicit ONNX request does not silently fall back to TEI; it exits until S03 wires inference.
- [x] Go tests passed: 72 tests in 4 packages.
- [x] GolangCI-Lint passed: 0 issues.
- [x] Docker Compose config passed.
- [x] GitNexus detect changes reported low risk and no affected processes.

## UAT Result

Pass. The project now has a tested opt-in validation seam ready for S03 ONNX loader work.

