# S01: Artifact manifest and checksum contract

**Goal:** Define how the local ONNX artifact is referenced, checksummed, and validated without committing model files.
**Demo:** After this, ONNX artifacts have a manifest/checksum contract and missing artifact failure expectations before any runtime wiring.

## Must-Haves

- Artifact manifest schema or documented convention exists.
- M010 ONNX artifact path/hash/size are represented without committing the artifact.
- Missing artifact and checksum mismatch behavior are specified.
- No production runtime wiring occurs in this slice.

## Proof Level

- This slice proves: Tracked manifest/docs plus checksum verification against the M010 local artifact.

## Integration Closure

Provides the artifact contract consumed by runtime config and ONNX loader slices.

## Verification

- Adds explicit artifact metadata fields and failure messages needed for debugging startup/load failures.

## Tasks

- [x] **T01: Inspect artifact metadata requirements** `est:small`
  Inspect M010 export metadata, current gitignore/runtime paths, README runtime docs, and any existing artifact conventions. Determine the smallest manifest shape needed for S02/S03 without committing the 1.43GB ONNX file.
  - Files: `.gsd/runtime/onnx/m010-s03/export-metadata.json`, `.gitignore`, `README.md`
  - Verify: Required manifest fields listed; confirms no large artifact will be tracked.

- [x] **T02: Write tracked ONNX artifact manifest** `est:small`
  Create a tracked ONNX artifact manifest for the M010 FP32 dense candidate that records local expected path, size, sha256, source model revision/hash, output metadata, dependency pin, and production status. The manifest must not contain raw probe texts or secrets and must not include the ONNX binary.
  - Files: `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
  - Verify: Manifest JSON parses; path/size/sha256 match local artifact when present; git status shows ONNX binary remains ignored/untracked.

- [x] **T03: Document artifact validation contract** `est:small`
  Document the artifact contract and failure expectations in S01 research: where local artifacts live, how checksum validation should behave, what to do when the file is missing, and why production runtime is unchanged.
  - Files: `.gsd/milestones/M011-33b7wf/slices/S01/S01-RESEARCH.md`
  - Verify: Research artifact exists and states missing/checksum mismatch behavior plus no production runtime change.

## Files Likely Touched

- .gsd/runtime/onnx/m010-s03/export-metadata.json
- .gitignore
- README.md
- docs/onnx-artifacts/user-bge-m3-dense-fp32.json
- .gsd/milestones/M011-33b7wf/slices/S01/S01-RESEARCH.md
