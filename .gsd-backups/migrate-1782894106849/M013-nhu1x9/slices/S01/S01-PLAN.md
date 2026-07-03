# S01: Native tokenizer artifact contract

**Goal:** Define and verify the native `libtokenizers.a` artifact contract without committing the binary.
**Demo:** After this, the native HF tokenizer artifact has a checksum/provenance contract and reproducible local setup instructions.

## Must-Haves

- Native artifact manifest exists and excludes the binary.
- Manifest records source URL, architecture, checksum, expected filename, and local runtime path.
- Local `libtokenizers.a` checksum validates against manifest.
- Docs explain setup without secrets or raw probe text.
- No native binary is tracked.

## Proof Level

- This slice proves: Manifest plus local checksum/setup verification and tracked-binary scan.

## Integration Closure

Provides S02 with a validated native library path and metadata for build-tag feasibility.

## Verification

- Adds artifact checksums, arch/version metadata, local setup path, and clear failure diagnostics for missing native library.

## Tasks

- [x] **T01: Design native tokenizer artifact contract** `est:small`
  Inspect existing artifact-manifest patterns and `.gitignore` coverage for native tokenizer artifacts. Decide the manifest path and local ignored artifact path.
  - Files: `.gitignore`, `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
  - Verify: Task summary names tracked manifest path, ignored local artifact path, and binary exclusion rule.

- [x] **T02: Write native tokenizer artifact manifest** `est:medium`
  Create the tracked native tokenizer artifact manifest and setup notes. Use the local S03 prebuilt linux-amd64 `libtokenizers.a` evidence to record checksum, size, source URL, architecture, and expected local path.
  - Files: `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`, `.gitignore`
  - Verify: Manifest JSON parses and local artifact checksum/size match when artifact exists.

- [x] **T03: Verify native tokenizer artifact contract** `est:medium`
  Add a small verifier mode or script for native tokenizer artifact manifest validation, and run it against the local temp artifact without committing the binary.
  - Files: `tools/compare_tokenizers.py`
  - Verify: Verifier exits 0 for local artifact; tracked-file scan confirms no `.a` binary is tracked.

- [x] **T04: Verify S01 artifact hygiene** `est:small`
  Run S01 verification: JSON parser, artifact checksum, default Go tests/lint if code changed, raw/binary leak checks, GitNexus detect_changes.
  - Verify: All S01 verification gates pass.

## Files Likely Touched

- .gitignore
- docs/onnx-artifacts/user-bge-m3-dense-fp32.json
- docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
- tools/compare_tokenizers.py
