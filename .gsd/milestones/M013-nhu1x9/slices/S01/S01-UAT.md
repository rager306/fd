# S01: Native tokenizer artifact contract — UAT

**Milestone:** M013-nhu1x9
**Written:** 2026-05-20T03:40:50.579Z

# S01 UAT — Native tokenizer artifact contract

## Checks

- [x] Manifest exists at `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`.
- [x] Manifest records source URL, linux-amd64 platform, size, SHA256, and local path.
- [x] Local `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a` matches manifest size/SHA256.
- [x] `*.a` static libraries are gitignored.
- [x] No native binary is tracked.
- [x] Validation artifact exists at `benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt`.
- [x] Raw probe leak check passes.
- [x] Default Go tests and lint pass.

## UAT Result

Pass. S02 can use the manifest/local artifact path for opt-in build-tag feasibility.

