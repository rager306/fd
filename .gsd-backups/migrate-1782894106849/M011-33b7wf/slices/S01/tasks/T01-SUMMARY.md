---
id: T01
parent: S01
milestone: M011-33b7wf
key_files:
  - .gsd/runtime/onnx/m010-s03/export-metadata.json
  - .gitignore
  - README.md
key_decisions:
  - Use a tracked JSON manifest under `docs/onnx-artifacts/` for the M010 FP32 dense artifact contract.
  - Keep the ONNX binary under ignored `.gsd/runtime/onnx/` or another externally managed artifact path; never track it in git.
  - Manifest must support missing-file and checksum-mismatch diagnostics before any ONNX runtime load.
duration: 
verification_result: mixed
completed_at: 2026-05-19T18:56:35.716Z
blocker_discovered: false
---

# T01: Defined the M011 S01 manifest requirements from M010 metadata: tracked JSON contract, ignored binary artifact, checksum-first validation, and explicit failure modes.

**Defined the M011 S01 manifest requirements from M010 metadata: tracked JSON contract, ignored binary artifact, checksum-first validation, and explicit failure modes.**

## What Happened

Inspected M010 export metadata, S04 recommendation, `.gitignore`, and README runtime notes. The M010 ONNX candidate has stable metadata suitable for a tracked manifest: model ID `deepvk/USER-bge-m3`, local source path `tei-models/deepvk--USER-bge-m3`, ONNX path `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`, size `1432482908`, SHA256 `28538a17a99302e144149732d73fb273cd7c7a0468dc59167caa5a2d5ff2a3d4`, output `dense_vecs`, shape `[batch_size,1024]`, dtype `tensor(float)`, provider `CPUExecutionProvider`, opset `17`, and dependency pin `transformers==4.51.3`. `.gitignore` already excludes `.gsd/runtime/` and `tei-models/`, so large artifacts remain untracked. README currently documents ONNX as a future measured optimization but does not yet define an artifact manifest contract.

## Verification

Read M010 export metadata, S04 recommendation, `.gitignore`, and README. Confirmed required fields and that large runtime artifacts are ignored by git.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Read `.gsd/runtime/onnx/m010-s03/export-metadata.json`; captured ONNX size/hash/output/provider/package metadata.` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Read `.gsd/milestones/M010-84qfzu/slices/S04/S04-RESEARCH.md`; confirmed opt-in prototype and no production switch boundary.` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Read `.gitignore`; confirmed `.gsd/runtime/` and `tei-models/` are ignored.` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Read `README.md`; confirmed ONNX is documented as future measured optimization but manifest contract is not yet documented.` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None. S01 T01 stayed in inspection/planning scope and did not edit runtime code.

## Known Issues

The current manifest will initially point to a local ignored artifact path. A future production path still needs artifact download/storage design beyond local development.

## Files Created/Modified

- `.gsd/runtime/onnx/m010-s03/export-metadata.json`
- `.gitignore`
- `README.md`
