---
id: T01
parent: S02
milestone: M031-gn517a
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
  - docs/onnx-artifacts/PROVISIONING.md
key_decisions:
  - Persist pinned source candidates for native tokenizer, tokenizer JSON, and ONNX Runtime in docs/manifests.
  - Keep the ONNX model artifact explicitly blocked until exact binary is mirrored/uploaded or a reproducible export gate is built.
duration: 
verification_result: passed
completed_at: 2026-05-21T06:37:17.535Z
blocker_discovered: false
---

# T01: Persisted the M031 source contract in manifests and provisioning docs.

**Persisted the M031 source contract in manifests and provisioning docs.**

## What Happened

Updated the tracked ONNX model manifest with a `source_contract` that records the ONNX model binary blocker plus pinned candidates for tokenizer.json and ONNX Runtime. Updated the native tokenizer manifest to replace the mutable `latest` source with the pinned `v1.27.0` release asset and archive checksum. Updated provisioning docs with source selection status and the remaining ONNX model blocker.

## Verification

JSON parsing and docs source-status checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_exec M031 docs/manifests source contract checks` | 0 | ✅ pass — JSON valid, required source candidates present, mutable latest URL absent from provisioning doc | 54ms |

## Deviations

None.

## Known Issues

The manual hosted workflow still cannot be treated as rollout proof because the ONNX model binary has no immutable external source and no workflow has been run.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/PROVISIONING.md`
