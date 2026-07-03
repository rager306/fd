# S01: ONNX Runtime wheel provisioning support

**Goal:** Add safe ONNX Runtime wheel extraction support to provisioning helper.
**Demo:** After this, `tools/provision_onnx_artifacts.py` supports a safe ONNX Runtime wheel-member extraction path and has local probes proving it.

## Must-Haves

- Helper can extract configured ONNX Runtime member from `.whl`/zip to approved destination.
- Member must be regular file and size-bounded.
- Destination sha256 verification remains mandatory when runtime source is provided.
- Existing ONNX/native/tokenizer behavior still works.
- No ignored native/runtime/model binaries are committed.

## Proof Level

- This slice proves: Positive/negative synthetic wheel probes plus existing provisioning/verifier checks.

## Integration Closure

Connects M031 ONNX Runtime source candidate metadata to the provisioning helper without changing runtime code.

## Verification

- Provisioning output continues to use sanitized paths and verified result records.

## Tasks

- [x] **T01: Implement ONNX Runtime wheel extraction** `est:medium`
  Update `tools/provision_onnx_artifacts.py` to read ONNX Runtime member/size/sha from `source_contract.onnx_runtime`, support safe zip/wheel member extraction, and wire it into `--onnx-runtime-source`. Preserve direct-file fallback and existing tar extraction behavior.
  - Files: `tools/provision_onnx_artifacts.py`
  - Verify: py_compile and targeted synthetic positive probe pass.

- [x] **T02: Verify wheel extraction failure modes** `est:small`
  Run positive/negative synthetic wheel probes in temporary repo roots: positive member extraction, missing member, oversized member, and checksum mismatch.
  - Verify: All synthetic probes pass/fail as expected with sanitized output.

- [x] **T03: Verify provisioning compatibility** `est:small`
  Run compatibility checks for existing dry-run and verifier behavior, then record S01 summary.
  - Verify: dry-run/verifier allow-missing and compile checks pass.

## Files Likely Touched

- tools/provision_onnx_artifacts.py
