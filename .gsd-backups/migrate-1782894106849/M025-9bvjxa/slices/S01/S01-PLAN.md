# S01: Artifact provisioning contract

**Goal:** Define and implement external artifact provisioning contract for ONNX/native/runtime assets.
**Demo:** After this, there is a concrete provisioning contract and local automation for checksum-verified artifact staging without git-tracking binaries.

## Must-Haves

- Contract lists required artifacts and source responsibilities.
- Script supports dry-run and fail-fast missing URL behavior.
- Script verifies size/sha before accepting artifacts.
- Strict local verifier remains passing with existing local artifacts.
- No binaries are tracked.

## Proof Level

- This slice proves: Docs, script dry-run, verifier, binary hygiene.

## Integration Closure

Provides source interface for hosted CI/manual deploy packaging.

## Verification

- Provisioning script emits sanitized plan/status and checksum verification results.

## Tasks

- [x] **T01: Write artifact provisioning contract** `est:small`
  Write artifact provisioning contract documenting required artifact sources, destination paths, checksum policy, cache layout, and current blockers.
  - Files: `docs/onnx-artifacts/PROVISIONING.md`, `docs/onnx-artifacts/README.md`
  - Verify: Contract exists and states blockers without raw/secrets/binaries.

- [x] **T02: Implement provisioning helper** `est:medium`
  Add `tools/provision_onnx_artifacts.py` supporting dry-run and checksum-verified download/copy into manifest local paths. It should require explicit source URLs/paths and never provide fake defaults for missing external artifacts.
  - Files: `tools/provision_onnx_artifacts.py`
  - Verify: Python compile, dry-run output, missing-source failure behavior, and strict local verifier pass.

- [x] **T03: Verify provisioning contract** `est:small`
  Run S01 verification: provisioning dry-run, verifier strict/allow-missing, binary hygiene, and update README links to provisioning contract.
  - Files: `docs/onnx-artifacts/README.md`
  - Verify: All S01 checks pass and no binaries are tracked.

## Files Likely Touched

- docs/onnx-artifacts/PROVISIONING.md
- docs/onnx-artifacts/README.md
- tools/provision_onnx_artifacts.py
