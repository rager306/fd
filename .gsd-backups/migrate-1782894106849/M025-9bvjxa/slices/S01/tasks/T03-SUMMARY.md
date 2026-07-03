---
id: T03
parent: S01
milestone: M025-9bvjxa
key_files:
  - docs/onnx-artifacts/README.md
  - docs/onnx-artifacts/PROVISIONING.md
  - tools/provision_onnx_artifacts.py
key_decisions:
  - Provisioning contract is linked from README and validated with dry-run plus allow-missing verifier.
  - Binary hygiene remains enforced with the corrected `Dockerfile.onnx` exemption.
duration: 
verification_result: passed
completed_at: 2026-05-20T11:40:21.405Z
blocker_discovered: false
---

# T03: Verified the M025 artifact provisioning contract and helper.

**Verified the M025 artifact provisioning contract and helper.**

## What Happened

Ran S01 verification for the provisioning contract. The provisioning helper and verifier compile, dry-run output is valid JSON, allow-missing verifier output is valid JSON, README links the provisioning contract, the contract states hosted CI/full deployment cannot be truthful without the missing ONNX source, and binary hygiene passes.

## Verification

Provisioning tools, docs discoverability, and binary hygiene checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py tools/verify_onnx_artifacts.py && provision dry-run + verifier allow-missing` | 0 | ✅ pass — m025_provisioning_tools=pass | 15300ms |
| 2 | `grep README/PROVISIONING blocker language` | 0 | ✅ pass — provisioning_docs_discoverable=pass | 0ms |
| 3 | `git ls-files refined binary hygiene check` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0 | 0ms |

## Deviations

A first doc grep used wording that did not exactly match the contract; I reran it with the actual blocker wording and it passed.

## Known Issues

Full hosted CI remains blocked until required external artifact sources are supplied.

## Files Created/Modified

- `docs/onnx-artifacts/README.md`
- `docs/onnx-artifacts/PROVISIONING.md`
- `tools/provision_onnx_artifacts.py`
