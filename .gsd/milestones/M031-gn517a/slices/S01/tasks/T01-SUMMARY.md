---
id: T01
parent: S01
milestone: M031-gn517a
key_files:
  - .gsd/milestones/M031-gn517a/slices/S01/S01-RESEARCH.md
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
  - docs/onnx-artifacts/PROVISIONING.md
key_decisions:
  - Use explicit source statuses: immutable_selected, candidate, blocked. No artifact is promoted to immutable_selected without fresh pinned-source evidence.
duration: 
verification_result: passed
completed_at: 2026-05-21T06:25:23.989Z
blocker_discovered: false
---

# T01: Inventoried ONNX artifact source requirements and current blockers.

**Inventoried ONNX artifact source requirements and current blockers.**

## What Happened

Created the M031 S01 research artifact with a source matrix for ONNX model, native tokenizer, tokenizer JSON, and ONNX Runtime. The inventory captures destination paths, exact sizes and sha256 values where known, current provenance, and source status definitions. It explicitly avoids fake/default URLs and keeps ONNX production/default blocked.

## Verification

Inventory artifact contains all four required artifacts and exact known checksums.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_summary_save M031 S01 RESEARCH draft` | 0 | ✅ pass — saved source matrix with four artifact classes | 0ms |

## Deviations

None.

## Known Issues

External candidate research remains T02; inventory currently marks ONNX model as blocked, tokenizer JSON as candidate, native tokenizer and ONNX Runtime as candidate/blocked.

## Files Created/Modified

- `.gsd/milestones/M031-gn517a/slices/S01/S01-RESEARCH.md`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/PROVISIONING.md`
