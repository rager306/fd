# S02: Docker CI boundary validation — UAT

**Milestone:** M021-4t2wpt
**Written:** 2026-05-20T10:26:34.724Z

# S02 UAT — Docker CI boundary validation

## Checks

- [x] Initial default Docker failure reproduced.
- [x] ONNX runtime moved behind explicit `onnx` build tag.
- [x] Default Docker build passes.
- [x] Artifact verifier passes.
- [x] Default Go tests pass.
- [x] Pinned GolangCI-Lint reports 0 issues.
- [x] Native tokenizer tagged tests pass.
- [x] ONNX+native tagged smoke tests pass.
- [x] No `.onnx` or `libtokenizers.a` files are tracked.
- [x] No background processes remain.
- [x] Port 18000 is clean.
- [x] GitNexus scope check is low.

## UAT Result

Pass. M021 can close as a packaging-contract and default-build-boundary milestone.

