---
milestone: M032-qq6po2
slice: S01
type: verifier-proof-boundary
status: complete
---

# M032 S01 — ONNX Export Contract Verifier Boundary

## What was added

`tools/verify_onnx_export_contract.py` verifies the existing local USER-bge-m3 ONNX export contract against:

- tracked manifest: `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`;
- local source provenance: `.gsd/runtime/onnx/m010-s03/source-provenance.json`;
- local export metadata: `.gsd/runtime/onnx/m010-s03/export-metadata.json`;
- ignored local ONNX artifact: `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`.

## What it proves

The verifier proves that the current local ignored ONNX artifact matches the tracked contract:

- `production_default=false` remains enforced;
- artifact path is under `.gsd/runtime/onnx`;
- artifact size and sha256 match the manifest;
- manifest export contract uses the pinned M010 toolchain:
  - Python `3.13.12`;
  - `torch==2.12.0`;
  - `transformers==4.51.3`;
  - `onnx==1.21.0`;
  - `onnxruntime==1.26.0`;
  - `safetensors==0.7.0`;
- manifest source files match M010 provenance checksums;
- M010 export metadata recorded the same ONNX artifact checksum;
- ONNX Runtime metadata recorded `CPUExecutionProvider` and expected 1024-dim output;
- verifier output states `claim_scope=existing_artifact_contract_verification_not_regenerated_export`.

## What it does not prove

The verifier does not regenerate the ONNX binary. It is not byte-for-byte reproducibility proof for a fresh export. It does not remove the M031 blocker that the exact ONNX model binary still needs either:

1. immutable non-secret external hosting for the exact current binary, or
2. a separate reproducible-export workflow that regenerates the binary and then reruns legal/performance/package gates.

## Failure mode evidence

Negative probes using temporary tampered metadata failed as expected:

| Probe | Expected failure label | Evidence |
|---|---|---|
| Artifact sha256 mismatch | `artifact_sha256` | `.gsd/exec/2b3f90be-849f-4c8f-af94-b87cd7aed5f1.stdout` |
| Model revision mismatch | `model_revision` | `.gsd/exec/2b3f90be-849f-4c8f-af94-b87cd7aed5f1.stdout` |
| `transformers` version mismatch | `export_metadata` | `.gsd/exec/2b3f90be-849f-4c8f-af94-b87cd7aed5f1.stdout` |

## Operational boundary

TEI remains production/default. ONNX remains opt-in experimental. No push, upload, workflow dispatch, or production switch occurred in S01.
