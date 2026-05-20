# S01: Dedicated ONNX Docker packaging proof

**Goal:** Add and validate a dedicated opt-in ONNX Docker packaging path using verified local artifacts and explicit `onnx hf_tokenizers` build tags.
**Demo:** After this, there is an opt-in ONNX Docker packaging path that verifies local artifacts and preserves the default TEI Docker path.

## Must-Haves

- Packaging design chooses a safe Docker context strategy.
- Artifact verifier runs before packaging.
- Dedicated ONNX Dockerfile/script uses `onnx hf_tokenizers` tags.
- Default Docker build remains passing.
- ONNX packaging build/run either passes or records a concrete blocker.
- No binaries are tracked.

## Proof Level

- This slice proves: Verifier output, Docker build/run evidence or concrete blocker, binary hygiene.

## Integration Closure

Provides local proof that ONNX 1024 can be packaged separately from default TEI image.

## Verification

- Adds build-time artifact verification and clear runtime expectations for the ONNX image.

## Tasks

- [x] **T01: Choose ONNX Docker packaging strategy** `est:small`
  Inspect Docker context constraints and choose the dedicated ONNX packaging strategy: root-context Dockerfile, staging script, or documented blocker.
  - Files: `api/Dockerfile`, `api/.dockerignore`, `tools/verify_onnx_artifacts.py`, `docs/onnx-artifacts/README.md`
  - Verify: Strategy states context, artifact inputs, build tags, and cleanup approach.

- [x] **T02: Implement ONNX packaging path** `est:medium`
  Implement the dedicated ONNX Docker packaging path, preferably a root-context Dockerfile plus script that verifies artifacts before build and keeps binaries untracked.
  - Files: `Dockerfile.onnx`, `tools/build_onnx_image.sh`, `docs/onnx-artifacts/README.md`
  - Verify: Script shell syntax passes and artifact verifier is invoked by the script.

- [x] **T03: Run Docker packaging proof** `est:medium`
  Run the packaging proof: default Docker build, ONNX artifact verification, ONNX image build, and smoke-run `/health` if build succeeds.
  - Verify: Default Docker build passes; ONNX build/run passes or records a concrete blocker; cleanup verified.

## Files Likely Touched

- api/Dockerfile
- api/.dockerignore
- tools/verify_onnx_artifacts.py
- docs/onnx-artifacts/README.md
- Dockerfile.onnx
- tools/build_onnx_image.sh
