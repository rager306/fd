# S01: ONNX 1024 runtime contract

**Goal:** Create or update tracked ONNX artifact metadata so validated 1024 runtime use is explicit while preserving export provenance.
**Demo:** After this, tracked metadata clearly states the ONNX binary was exported with dynamic sequence axes and validated at runtime sequence length 1024, with evidence links and production_default=false.

## Must-Haves

- Metadata JSON parses.
- Export sequence length remains documented as 128.
- Validated runtime max sequence length 1024 is documented separately.
- M018 quality and M019 performance evidence artifacts are linked.
- Production/default remains false/experimental.

## Proof Level

- This slice proves: JSON validation, evidence links, no binary tracking.

## Integration Closure

Provides a clear contract for future packaging/CI work.

## Verification

- Records validated runtime sequence length, evidence artifacts, and failure contracts in tracked metadata.

## Tasks

- [x] **T01: Choose metadata shape** `est:small`
  Inspect current ONNX manifest fields and decide whether to update the existing manifest or add a dedicated 1024 runtime contract file.
  - Files: `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
  - Verify: Chosen shape explicitly separates export provenance from validated runtime sequence length.

- [x] **T02: Update ONNX runtime contract metadata** `est:small`
  Update the tracked ONNX artifact metadata with validated 1024 runtime contract, evidence artifacts, and remaining gates while preserving prototype status.
  - Files: `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
  - Verify: `python3 -m json.tool docs/onnx-artifacts/user-bge-m3-dense-fp32.json` passes and required fields exist.

- [x] **T03: Validate contract metadata** `est:small`
  Verify metadata contract, no binary tracking, and artifact hygiene for S01 outputs.
  - Verify: JSON/field check passes and tracked binary check reports zero ONNX/native binaries.

## Files Likely Touched

- docs/onnx-artifacts/user-bge-m3-dense-fp32.json
